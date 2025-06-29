package privateisu

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

type PrivateIsu struct {
	cmd    *exec.Cmd
	stdout string
	stderr string
}

func New() *PrivateIsu {
	return &PrivateIsu{}
}

var _ benchmarker.Benchmarker = (*PrivateIsu)(nil)

const (
	benchmarkerDir         = "/home/isucon/private_isu/benchmarker"
	benchmarkerPath        = benchmarkerDir + "/bin/benchmarker"
	benchmarkerUserDataDir = benchmarkerDir + "/userdata"
)

func (b *PrivateIsu) Start(ctx context.Context, job *domain.Job) (benchmarker.Outputs, time.Time, error) {
	cmd := exec.CommandContext(ctx, benchmarkerPath,
		"--userdata", benchmarkerUserDataDir,
		"--target", job.GetTargetURL(),
	)

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

// https://github.com/catatsuy/private-isu/blob/3c5e230ca92c03965cdb6f61b92605cb41ecb259/benchmarker/cli.go#L42
type stdOut struct {
	Pass     bool     `json:"pass"`
	Score    int64    `json:"score"`
	Success  int64    `json:"success"`
	Fail     int64    `json:"fail"`
	Messages []string `json:"messages"`
}

func (b *PrivateIsu) Wait(_ context.Context) (domain.Result, time.Time, error) {
	if err := b.cmd.Wait(); err != nil {
		return domain.ResultError, time.Now(), fmt.Errorf("wait command: %w", err)
	}
	endTime := time.Now()

	var stdout stdOut
	err := json.Unmarshal([]byte(b.stdout), &b)
	if err != nil {
		return domain.ResultError, endTime, fmt.Errorf("unmarshal stdout: %w", err)
	}

	if !stdout.Pass {
		return domain.ResultFailed, endTime, nil
	}

	return domain.ResultPassed, endTime, nil
}

func (b *PrivateIsu) CalculateScore(_ context.Context, allStdout, allStderr string) (int, error) {
	b.stdout, b.stderr = allStdout, allStderr

	var stdout stdOut
	err := json.Unmarshal([]byte(b.stdout), &stdout)
	if err != nil {
		return 0, fmt.Errorf("unmarshal stdout: %w", err)
	}

	return int(stdout.Score), nil
}
