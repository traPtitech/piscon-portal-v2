package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl"
	isucon11qualify "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/impl/isucon11-qualify"
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
	ctx := context.Background()
	// ログの収集に失敗した場合にベンチマーカーを停止させる必要があるため、
	// ポータルとの通信に使うcontextとは分けておく。
	benchCtx, cancelBench := context.WithCancel(ctx)
	defer cancelBench()

	job, err := r.portal.GetJob(ctx)
	if err != nil {
		return fmt.Errorf("get benchmark: %w", err)
	}

	out, startedAt, err := r.benchmarker.Start(benchCtx, job)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	stdoutBdr := &SyncStringBuilder{}
	stderrBdr := &SyncStringBuilder{}

	stdoutErrChan, stderrErrChan := make(chan error, 1), make(chan error, 1)

	go func() {
		stdoutErrChan <- captureStreamOutput(benchCtx, out.Stdout, stdoutBdr)
	}()
	go func() {
		stderrErrChan <- captureStreamOutput(benchCtx, out.Stderr, stderrBdr)
	}()

	streamClient, streamErr := r.portal.MakeProgressStreamClient(ctx)
	if streamErr != nil {
		streamErr = fmt.Errorf("create streaming client: %w", streamErr)
	} else {
		streamErr = r.streamJobProgress(ctx, streamClient, job, startedAt,
			stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
	}
	if streamErr != nil {
		// 出力を読み切れていないため、ベンチマーカーの正常な終了は待てない。
		log.Printf("collect data: %v", streamErr)
		cancelBench()
	}

	// ベンチマーカーのプロセス終了と最終レポートの受信を待ってから、
	// 確定したスコアとログを送信する。
	result, finishedAt, runnerErr := r.benchmarker.Wait(benchCtx)

	switch {
	case streamErr != nil:
		result = domain.ResultError
		runnerErr = errors.Join(streamErr, runnerErr)
	case streamClient != nil:
		if err := r.sendProgress(ctx, streamClient, job, startedAt, stdoutBdr, stderrBdr); err != nil {
			// ベンチマークの合否自体は確定しているので、resultは変えずにエラーだけ伝える。
			log.Printf("send final progress: %v", err)
			runnerErr = errors.Join(runnerErr, fmt.Errorf("send final progress: %w", err))
		}
	}

	if streamClient != nil {
		if err := streamClient.Close(); err != nil {
			log.Printf("close progress stream: %v", err)
			runnerErr = errors.Join(runnerErr, fmt.Errorf("close progress stream: %w", err))
		}
	}

	if err := r.portal.PostJobFinished(ctx, job.GetID(), finishedAt, result, runnerErr); err != nil {
		return fmt.Errorf("post job finished: %w", err)
	}

	return nil
}

func captureStreamOutput(_ context.Context, r io.Reader, bdr *SyncStringBuilder) error {
	for {
		n, err := io.Copy(bdr, io.LimitReader(r, bufSize))
		if err != nil {
			return fmt.Errorf("copy: %w", err)
		}
		if n == 0 {
			return nil
		}
	}
}

// sendProgress calculates the current score and sends it with the collected logs.
func (r *Runner) sendProgress(
	ctx context.Context, streamClient portal.ProgressStreamClient,
	job *domain.Job, startedAt time.Time, stdoutBdr, stderrBdr *SyncStringBuilder,
) error {
	stdout := stdoutBdr.String()
	stderr := stderrBdr.String()

	score, err := r.benchmarker.CalculateScore(ctx, stdout, stderr)
	if err != nil {
		return fmt.Errorf("calculate score: %w", err)
	}

	progress := domain.NewProgress(job.GetID(), stdout, stderr, score, startedAt)
	if err := streamClient.SendProgress(ctx, progress); err != nil {
		return fmt.Errorf("send progress: %w", err)
	}

	return nil
}

// streamJobProgress collects the benchmark job's stdout and stderr and sends the progress to the portal.
// It returns nil if, and only if both of stdout and stderr reach EOF.
//
// 最終的なスコアは [Benchmarker.Wait] がベンチマーカーの最終レポートの受信を待った後に
// 確定するため、ここでは送信しない。streamClientのCloseも呼び出し側で行う。
func (r *Runner) streamJobProgress(
	ctx context.Context, streamClient portal.ProgressStreamClient,
	job *domain.Job, startedAt time.Time,
	stdoutBdr, stderrBdr *SyncStringBuilder,
	stdoutErrChan, stderrErrChan chan error,
) error {
	// 初期状態を送信して、ポータルに開始を通知する。
	initialProgress := domain.NewProgress(job.GetID(), "", "", 0, startedAt)
	if err := streamClient.SendProgress(ctx, initialProgress); err != nil {
		return fmt.Errorf("send initial progress: %w", err)
	}

	finished := struct {
		stdout bool
		stderr bool
	}{false, false}

	ticker := time.NewTicker(sendProgressInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := r.sendProgress(ctx, streamClient, job, startedAt, stdoutBdr, stderrBdr); err != nil {
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
