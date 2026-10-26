package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl"
	isucon11final "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl/isucon11-final"
	isucon11qualify "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl/isucon11-qualify"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl/isucon13"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl/isucon14"
	privateisu "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl/private_isu"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal"
)

const (
	bufSize              = 1024
	sendProgressInterval = 5 * time.Second
)

const (
	problemExample         string = "example"
	problemPrivateIsu      string = "private_isu"
	problemIsucon11Qualify string = "isucon11-qualify"
	problemIsucon11Final   string = "isucon11-final"
	problemIsucon13        string = "isucon13"
	problemIsucon14        string = "isucon14"
)

var (
	problemBenchmarks = map[string]func(conf config.Problem) (benchmarker.Benchmarker, error){
		problemExample: func(_ config.Problem) (benchmarker.Benchmarker, error) {
			return impl.NewExample(), nil
		},
		problemPrivateIsu: func(conf config.Problem) (benchmarker.Benchmarker, error) {
			return privateisu.New(conf)
		},
		problemIsucon11Qualify: func(conf config.Problem) (benchmarker.Benchmarker, error) {
			return isucon11qualify.New(conf)
		},
		problemIsucon11Final: func(conf config.Problem) (benchmarker.Benchmarker, error) {
			return isucon11final.New(conf)
		},
		problemIsucon13: func(conf config.Problem) (benchmarker.Benchmarker, error) {
			return isucon13.New(conf)
		},
		problemIsucon14: func(conf config.Problem) (benchmarker.Benchmarker, error) {
			return isucon14.New(conf)
		},
	}
)

type Runner struct {
	portal      portal.Portal
	benchmarker benchmarker.Benchmarker
}

func Prepare(portal portal.Portal, problemConfig config.Problem) (*Runner, error) {
	newBenchmarkerFn, ok := problemBenchmarks[problemConfig.Name]
	if !ok {
		return nil, fmt.Errorf("unknown problem: %q", problemConfig.Name)
	}

	benchmarker, err := newBenchmarkerFn(problemConfig)
	if err != nil {
		return nil, fmt.Errorf("create benchmarker: %w", err)
	}

	return &Runner{
		portal:      portal,
		benchmarker: benchmarker,
	}, nil
}

func (r *Runner) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	job, err := r.portal.GetJob(ctx)
	if err != nil {
		return fmt.Errorf("get benchmark: %w", err)
	}

	stdoutArchive, stderrArchive, err := openLogArchive(job.GetID())
	if err != nil {
		archiveErr := fmt.Errorf("create log archive: %w", err)
		return errors.Join(archiveErr, r.portal.PostJobFinished(ctx, job.GetID(), time.Now(), domain.ResultError, archiveErr))
	}
	out, startedAt, err := r.benchmarker.Start(ctx, job)
	if err != nil {
		return errors.Join(fmt.Errorf("execute: %w", err), stdoutArchive.Close(), stderrArchive.Close())
	}
	stdoutR := out.Stdout
	stderrR := out.Stderr

	stdoutBdr := &SyncStringBuilder{}
	stderrBdr := &SyncStringBuilder{}

	stdoutErrChan, stderrErrChan := make(chan error, 1), make(chan error, 1)
	var readers sync.WaitGroup

	readers.Go(func() {
		err := captureStreamOutput(ctx, stdoutR, io.MultiWriter(stdoutArchive, stdoutBdr))
		stdoutErrChan <- errors.Join(err, stdoutArchive.Close())
	})
	readers.Go(func() {
		err := captureStreamOutput(ctx, stderrR, io.MultiWriter(stderrArchive, stderrBdr))
		stderrErrChan <- errors.Join(err, stderrArchive.Close())
	})

	progressErr := r.streamJobProgress(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
	if progressErr != nil {
		cancel()
	}
	// パイプの読み取りとgzipの書き込み完了後に、子プロセスを回収して結果を報告する。
	readers.Wait()
	result, finishedAt, waitErr := r.benchmarker.Wait(ctx)
	runErr := errors.Join(progressErr, waitErr)
	if runErr != nil {
		result = domain.ResultError
	}
	// 送信失敗時は実行用のコンテキストがキャンセル済みなので、完了報告用に作り直す。
	reportCtx, reportCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer reportCancel()
	if err := r.portal.PostJobFinished(reportCtx, job.GetID(), finishedAt, result, runErr); err != nil {
		return errors.Join(runErr, fmt.Errorf("post job finished: %w", err))
	}
	return runErr
}

func captureStreamOutput(_ context.Context, r io.Reader, w io.Writer) error {
	for {
		n, err := io.Copy(w, io.LimitReader(r, bufSize))
		if err != nil {
			return fmt.Errorf("copy: %w", err)
		}
		if n == 0 {
			return nil
		}
	}
}

// streamJobProgress はstdout/stderrを収集し、終了時に最終進捗を送信する。
// 読み取り・送信・portal側の保存に失敗した場合はエラーを返す。
func (r *Runner) streamJobProgress(
	ctx context.Context, job *domain.Job, startedAt time.Time,
	stdoutBdr, stderrBdr *SyncStringBuilder,
	stdoutErrChan, stderrErrChan chan error,
) (err error) {
	streamClient, err := r.portal.MakeProgressStreamClient(ctx)
	if err != nil {
		return fmt.Errorf("create streaming client: %w", err)
	}
	defer func() { err = errors.Join(err, streamClient.Close()) }()

	// 初期状態を送信して、ポータルに開始を通知する。
	initialProgress := domain.NewProgress(job.GetID(), "", "", 0, startedAt)
	if err := streamClient.SendProgress(ctx, initialProgress); err != nil {
		return fmt.Errorf("send initial progress: %w", err)
	}

	sendFailed := false
	calcAndSendProgress := func() error {
		stdout := stdoutBdr.String()
		stderr := stderrBdr.String()
		score, err := r.benchmarker.CalculateScore(ctx, stdout, stderr)
		if err != nil {
			return fmt.Errorf("calculate score: %w", err)
		}

		progress := domain.NewProgress(job.GetID(), stdout, stderr, score, startedAt)
		err = streamClient.SendProgress(ctx, progress)
		if err != nil {
			sendFailed = true
			return fmt.Errorf("send progress: %w", err)
		}
		return nil
	}

	// 送信失敗後のストリームには再送せず、同じエラーを重複して返さない。
	defer func() {
		if sendFailed {
			return
		}
		if calcErr := calcAndSendProgress(); calcErr != nil {
			err = errors.Join(err, fmt.Errorf("send final progress: %w", calcErr))
		}
	}()

	finished := struct {
		stdout bool
		stderr bool
	}{false, false}

	ticker := time.NewTicker(sendProgressInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := calcAndSendProgress(); err != nil {
				return fmt.Errorf("calc and send progress in tick: %w", err)
			}
		case err := <-stdoutErrChan:
			finished.stdout = true
			if err != nil {
				return fmt.Errorf("read stdout: %w", err)
			}
			if finished.stderr {
				return nil
			}

		case err := <-stderrErrChan:
			finished.stderr = true
			if err != nil {
				return fmt.Errorf("read stderr: %w", err)
			}
			if finished.stdout {
				return nil
			}
		}
	}
}
