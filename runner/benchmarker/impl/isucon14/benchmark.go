package isucon14

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"strconv"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker/internal/isuxbench"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

const (
	// 競技環境のnginxは Host ヘッダーが `*.xiv.isucon.net`宛てのリクエストにしかアプリを返さない。
	defaultTargetURL = "https://isuride.xiv.isucon.net"
	// 問題のAMIでは`isuride-payment_mock`が12345番を使っているため、それと被らないようにする。
	defaultPaymentBindPort = 12346
	targetPort             = "443"
)

type problemConf struct {
	execPath              string
	benchmarkerIP         string
	benchmarkerDir        string
	targetURL             string
	paymentBindPort       int
	skipStaticSanityCheck bool
}

type Isucon14 struct {
	conf problemConf

	proc *isuxbench.Process
}

var _ benchmarker.Benchmarker = (*Isucon14)(nil)

func New(conf config.Problem) (*Isucon14, error) {
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

	targetURL := defaultTargetURL
	if v, exists := conf.Options["target-url"]; exists {
		targetURL, ok = v.(string)
		if !ok || targetURL == "" {
			return nil, fmt.Errorf("invalid 'target-url' option in problem configuration")
		}
	}

	paymentBindPort := defaultPaymentBindPort
	if v, exists := conf.Options["payment-bind-port"]; exists {
		paymentBindPort, ok = v.(int)
		if !ok || paymentBindPort <= 0 || paymentBindPort > 65535 {
			return nil, fmt.Errorf("invalid 'payment-bind-port' option in problem configuration")
		}
	}

	skipStaticSanityCheck := false
	if v, exists := conf.Options["skip-static-sanity-check"]; exists {
		skipStaticSanityCheck, ok = v.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid 'skip-static-sanity-check' option in problem configuration")
		}
	}

	return &Isucon14{
		conf: problemConf{
			execPath:              path,
			benchmarkerIP:         benchmarkerIP,
			benchmarkerDir:        benchmarkerDir,
			targetURL:             targetURL,
			paymentBindPort:       paymentBindPort,
			skipStaticSanityCheck: skipStaticSanityCheck,
		},
	}, nil
}

func (b *Isucon14) Start(ctx context.Context, job *domain.Job) (benchmarker.Outputs, time.Time, error) {
	port := strconv.Itoa(b.conf.paymentBindPort)
	// ターゲットのアプリは`--payment-url`に決済のリクエストを送るので、ターゲットから届くアドレスにする。
	paymentURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(b.conf.benchmarkerIP, port),
	}
	args := []string{"run",
		// `--target`はHostヘッダーとSNIに使われ、実際の接続先は`--addr`になる。
		"--target", b.conf.targetURL,
		"--addr", net.JoinHostPort(job.GetTargetIPAdress(), targetPort),
		"--payment-url", paymentURL.String(),
		"--payment-bind-port", port,
	}
	if b.conf.skipStaticSanityCheck {
		args = append(args, "--skip-static-sanity-check")
	}
	cmd := exec.CommandContext(ctx, b.conf.execPath, args...)
	cmd.Dir = b.conf.benchmarkerDir

	proc, out, startedAt, err := isuxbench.Start(cmd)
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("start benchmarker: %w", err)
	}
	b.proc = proc

	return out, startedAt, nil
}

func (b *Isucon14) Wait(_ context.Context) (domain.Result, time.Time, error) {
	result, finishedAt, err := b.proc.Wait()
	// Prepare(initializeや静的ファイルの検証など)で失敗すると、
	// ベンチマーカーは最終レポートを送らずに正常終了する。これは競技上の失敗として扱う。
	if errors.Is(err, isuxbench.ErrNoFinalReport) {
		return domain.ResultFailed, finishedAt, nil
	}
	if err != nil {
		return result, finishedAt, fmt.Errorf("wait benchmarker: %w", err)
	}
	return result, finishedAt, nil
}

func (b *Isucon14) CalculateScore(_ context.Context, _, _ string) (int, error) {
	if b.proc == nil {
		return 0, nil
	}
	return b.proc.LatestScore(), nil
}
