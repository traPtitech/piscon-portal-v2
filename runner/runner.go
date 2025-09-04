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
	problemExample    string = "example"
	problemPrivateIsu string = "private_isu"
)

var (
	problemBenchmarks = map[string]func(conf config.Problem) (benchmarker.Benchmarker, error){
		problemExample: func(_ config.Problem) (benchmarker.Benchmarker, error) {
			return impl.NewExample(), nil
		},
		problemPrivateIsu: func(conf config.Problem) (benchmarker.Benchmarker, error) {
			return privateisu.New(conf)
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

	out, startedAt, err := r.benchmarker.Start(ctx, job)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}
	stdoutR := out.Stdout
	stderrR := out.Stderr

	stdoutBdr := &SyncStringBuilder{}
	stderrBdr := &SyncStringBuilder{}

	stdoutErrChan, stderrErrChan := make(chan error), make(chan error)

	go func() {
		stdoutErrChan <- captureStreamOutput(ctx, stdoutR, stdoutBdr)
	}()
	go func() {
		stderrErrChan <- captureStreamOutput(ctx, stderrR, stderrBdr)
	}()

	err = r.streamJobProgress(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
	if err != nil {
		log.Printf("collect data: %v", err)
		err := r.portal.PostJobFinished(ctx, job.GetID(), time.Now(), domain.ResultError, err)
		if err != nil {
			return fmt.Errorf("post job finished: %w", err)
		}
	}

	result, finishedAt, err := r.benchmarker.Wait(ctx)
	err = r.portal.PostJobFinished(ctx, job.GetID(), finishedAt, result, err)
	if err != nil {
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

// streamJobProgress collects the benchmark job's stdout and stderr and sends the progress to the portal.
// It returns nil if, and only if both of stdout and stderr reach EOF.
func (r *Runner) streamJobProgress(
	ctx context.Context, job *domain.Job, startedAt time.Time,
	stdoutBdr, stderrBdr *SyncStringBuilder,
	stdoutErrChan, stderrErrChan chan error,
) (err error) {
	streamClient, err := r.portal.MakeProgressStreamClient(ctx)
	if err != nil {
		return fmt.Errorf("create streaming client: %w", err)
	}
	defer streamClient.Close()

	// 初期状態を送信して、ポータルに開始を通知する。
	initialProgress := domain.NewProgress(job.GetID(), "", "", 0, startedAt)
	if err := streamClient.SendProgress(ctx, initialProgress); err != nil {
		return fmt.Errorf("send initial progress: %w", err)
	}

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
			return fmt.Errorf("send progress: %w", err)
		}
		return nil
	}

	// 最後に必ず結果を計算して送信するようにする
	defer func() {
		if calcErr := calcAndSendProgress(); calcErr != nil {
			err = errors.Join(err, fmt.Errorf("defer: calcAndSendProgress%w", calcErr))
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
