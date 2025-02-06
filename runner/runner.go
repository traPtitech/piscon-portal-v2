package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal"
)

const (
	bufSize              = 1024
	sendProgressInterval = 5 * time.Second
)

type Runner struct {
	portal      portal.Portal
	benchmarker benchmarker.Benchmarker
}

func Prepare(portal portal.Portal, benchmarker benchmarker.Benchmarker) *Runner {
	return &Runner{
		portal:      portal,
		benchmarker: benchmarker,
	}
}

func (r *Runner) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	job, err := r.portal.GetJob(ctx)
	if err != nil {
		return fmt.Errorf("get benchmark: %w", err)
	}

	stdoutR, stderrR, startedAt, err := r.benchmarker.Start(ctx, job)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	stdoutBdr := &strings.Builder{}
	stderrBdr := &strings.Builder{}

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

func captureStreamOutput(_ context.Context, r io.Reader, bdr *strings.Builder) error {
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
	stdoutBdr, stderrBdr *strings.Builder,
	stdoutErrChan, stderrErrChan chan error,
) (err error) {
	streamClient, err := r.portal.MakeProgressStreamClient(ctx)
	if err != nil {
		return fmt.Errorf("create streaming client: %w", err)
	}
	defer streamClient.Close()

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
				return fmt.Errorf("calc and send progress: %w", err)
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
