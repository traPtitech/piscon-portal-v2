package isucon11final

import (
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker/internal/isuxbench"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

type problemConf struct {
	execPath       string
	benchmarkerDir string
}

type Isucon11Final struct {
	conf problemConf

	proc *isuxbench.Process
}

var _ benchmarker.Benchmarker = (*Isucon11Final)(nil)

func New(conf config.Problem) (*Isucon11Final, error) {
	path, ok := conf.Options["benchmarker-path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("invalid or missing 'benchmarker-path' option in problem configuration")
	}
	benchmarkerDir, ok := conf.Options["benchmarker-dir"].(string)
	if !ok || benchmarkerDir == "" {
		return nil, fmt.Errorf("invalid or missing 'benchmarker-dir' option in problem configuration")
	}

	return &Isucon11Final{
		conf: problemConf{
			execPath:       path,
			benchmarkerDir: benchmarkerDir,
		},
	}, nil
}

func (b *Isucon11Final) Start(ctx context.Context, job *domain.Job) (benchmarker.Outputs, time.Time, error) {
	// 本選のベンチマーカーは`--target`に`host:port`を受け取り、`http://<target>/`を組み立てる。
	// ポートを省略した場合は80番になる。
	// `--exit-status`は付けない。競技上の失敗とプロセスの異常終了を区別するため、
	// 合否は最終レポートから判定する。
	cmd := exec.CommandContext(ctx, b.conf.execPath,
		"--target", job.GetTargetIPAdress())
	cmd.Dir = b.conf.benchmarkerDir

	proc, out, startedAt, err := isuxbench.Start(cmd)
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("start benchmarker: %w", err)
	}
	b.proc = proc

	return out, startedAt, nil
}

func (b *Isucon11Final) Wait(_ context.Context) (domain.Result, time.Time, error) {
	result, finishedAt, err := b.proc.Wait()
	if err != nil {
		return result, finishedAt, fmt.Errorf("wait benchmarker: %w", err)
	}
	return result, finishedAt, nil
}

func (b *Isucon11Final) CalculateScore(_ context.Context, _, _ string) (int, error) {
	if b.proc == nil {
		return 0, nil
	}
	return b.proc.LatestScore(), nil
}
