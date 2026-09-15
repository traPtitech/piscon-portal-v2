package isucon11qualify

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker/internal/isuxbench"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

type problemConf struct {
	benchmarkerIP  string
	execPath       string
	benchmarkerDir string
}

type Isucon11Qualify struct {
	conf problemConf

	proc *isuxbench.Process
}

var _ benchmarker.Benchmarker = (*Isucon11Qualify)(nil)

func New(conf config.Problem) (*Isucon11Qualify, error) {
	path, ok := conf.Options["benchmarker-path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("invalid or missing 'benchmarker-path' option in problem configuration")
	}
	benchmarkerIP, ok := conf.Options["benchmarker-ip"].(string)
	if !ok || benchmarkerIP == "" {
		return nil, fmt.Errorf("invalid or missing 'benchmarker-ip' option in problem configuration")
	}
	benchmarkerDir, ok := conf.Options["benchmarker-dir"].(string)
	if !ok || benchmarkerDir == "" {
		return nil, fmt.Errorf("invalid or missing 'benchmarker-dir' option in problem configuration")
	}

	return &Isucon11Qualify{
		conf: problemConf{
			execPath:       path,
			benchmarkerIP:  benchmarkerIP,
			benchmarkerDir: benchmarkerDir,
		},
	}, nil
}

func (b *Isucon11Qualify) Start(ctx context.Context, job *domain.Job) (benchmarker.Outputs, time.Time, error) {
	jiaServiceURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(b.conf.benchmarkerIP, "4999"),
	}
	cmd := exec.CommandContext(ctx, b.conf.execPath,
		"--target", job.GetTargetIPAdress(),
		// tlsを使わない場合は`--all-adresses`にtargetだけ指定すればok
		"--all-addresses", job.GetTargetIPAdress(),
		"--jia-service-url", jiaServiceURL.String())
	cmd.Dir = b.conf.benchmarkerDir

	proc, out, startedAt, err := isuxbench.Start(cmd)
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("start benchmarker: %w", err)
	}
	b.proc = proc

	return out, startedAt, nil
}

func (b *Isucon11Qualify) Wait(_ context.Context) (domain.Result, time.Time, error) {
	result, finishedAt, err := b.proc.Wait()
	if err != nil {
		return result, finishedAt, fmt.Errorf("wait benchmarker: %w", err)
	}
	return result, finishedAt, nil
}

func (b *Isucon11Qualify) CalculateScore(_ context.Context, _, _ string) (int, error) {
	if b.proc == nil {
		return 0, nil
	}
	return b.proc.LatestScore(), nil
}
