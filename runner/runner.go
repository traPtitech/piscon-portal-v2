package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal"
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
	ctx := context.Background()

	job, err := r.portal.GetJob(ctx)
	if err != nil {
		return fmt.Errorf("get benchmark: %w", err)
	}

	streamClient, err := r.portal.MakeProgressStreamClient(ctx)
	if err != nil {
		return fmt.Errorf("create streaming client: %w", err)
	}
	defer streamClient.Close()

	stdoutR, stderrR, startedAt, err := r.benchmarker.Start(ctx, job)
	if err != nil {
		return fmt.Errorf("execute: %w", err)
	}

	stdoutBdr := &strings.Builder{}
	stderrBdr := &strings.Builder{}

	for {
		if _, err := io.CopyN(stdoutBdr, stdoutR, 1024); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read stdout: %w", err)
		}

		if _, err := io.CopyN(stderrBdr, stderrR, 1024); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return fmt.Errorf("read stderr: %w", err)
		}

		score, err := r.benchmarker.CalculateScore(ctx, stdoutBdr.String(), stderrBdr.String())
		if err != nil {
			return fmt.Errorf("calculate score: %w", err)
		}

		err = streamClient.SendProgress(ctx, domain.NewProgress(job.GetID(), stdoutBdr.String(), stderrBdr.String(), score, startedAt))
		if err != nil {
			return fmt.Errorf("send progress: %w", err)
		}
	}

	finishedAt := time.Now()

	result, err := r.benchmarker.Wait()
	if err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	err = r.portal.PostJobFinished(ctx, job.GetID(), finishedAt, result)
	if err != nil {
		return fmt.Errorf("post job finished: %w", err)
	}

	return nil
}
