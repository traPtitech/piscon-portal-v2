package impl

import (
	"context"
	"fmt"
	"os/exec"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

// Example is the Example implementation of [github.com/traPtitech/piscon-portal-v2/runner/benchmarker.Benchmarker].
// It runs the benchmark script Example.sh with the target URL as an argument.
type Example struct {
	cmd    *exec.Cmd
	stdout string
	stderr string
}

func NewExample() *Example {
	return &Example{}
}

var _ benchmarker.Benchmarker = (*Example)(nil)

func (b *Example) Start(ctx context.Context, job *domain.Job) (benchmarker.Outputs, time.Time, error) {
	cmd := exec.CommandContext(ctx, "./example.sh", job.GetTargetURL())

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("get stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("start command: %w", err)
	}

	b.cmd = cmd

	startedAt := time.Now()

	return benchmarker.Outputs{
		Stdout: stdout,
		Stderr: stderr,
	}, startedAt, nil
}

func (b *Example) Wait(_ context.Context) (domain.Result, time.Time, error) {
	if err := b.cmd.Wait(); err != nil {
		return domain.ResultError, time.Now(), fmt.Errorf("wait command: %w", err)
	}
	endTime := time.Now()

	for _, line := range strings.Split(b.stdout, "\n") {
		if strings.Contains(line, "FAIL") {
			return domain.ResultFailed, endTime, nil
		}
	}

	return domain.ResultPassed, endTime, nil
}

func (b *Example) CalculateScore(_ context.Context, allStdout, allStderr string) (int, error) {
	b.stdout, b.stderr = allStdout, allStderr

	for _, line := range slices.Backward(strings.Split(allStdout, "\n")) {
		if strings.HasPrefix(line, "Score: ") {
			score, err := strconv.Atoi(strings.Split(line, " ")[1])
			if err != nil {
				return 0, fmt.Errorf("invalid score format: %w", err)
			}
			return score, nil
		}
	}

	return 0, fmt.Errorf("score not found")
}
