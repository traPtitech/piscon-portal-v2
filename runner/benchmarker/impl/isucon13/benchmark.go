package isucon13

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

type problemConf struct {
	execPath       string
	benchmarkerDir string
	targetURL      string
	dnsPort        int
	enableSSL      bool
}

type Isucon13 struct {
	conf problemConf
	cmd  *exec.Cmd

	outputDir string
	score     int
}

var _ benchmarker.Benchmarker = (*Isucon13)(nil)

func New(conf config.Problem) (*Isucon13, error) {
	path, ok := conf.Options["benchmarker-path"].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("invalid or missing 'benchmarker-path' option in problem configuration")
	}
	dir, ok := conf.Options["benchmarker-dir"].(string)
	if !ok || dir == "" {
		return nil, fmt.Errorf("invalid or missing 'benchmarker-dir' option in problem configuration")
	}
	// matsuu/aws-isucon のAMIは u.isucon.local を使う。
	targetURL := "https://pipe.u.isucon.local"
	if v, exists := conf.Options["target-url"]; exists {
		targetURL, ok = v.(string)
		if !ok || targetURL == "" {
			return nil, fmt.Errorf("invalid 'target-url' option in problem configuration")
		}
	}
	u, err := url.Parse(targetURL)
	if err != nil || u.Hostname() == "" || (u.Scheme != "https" && u.Scheme != "http") {
		return nil, fmt.Errorf("invalid 'target-url' option in problem configuration")
	}
	dnsPort := 53
	if v, exists := conf.Options["dns-port"]; exists {
		dnsPort, ok = v.(int)
		if !ok || dnsPort <= 0 || dnsPort > 65535 {
			return nil, fmt.Errorf("invalid 'dns-port' option in problem configuration")
		}
	}
	enableSSL := true
	if v, exists := conf.Options["enable-ssl"]; exists {
		enableSSL, ok = v.(bool)
		if !ok {
			return nil, fmt.Errorf("invalid 'enable-ssl' option in problem configuration")
		}
	}
	return &Isucon13{conf: problemConf{
		execPath:       path,
		benchmarkerDir: dir,
		targetURL:      targetURL,
		dnsPort:        dnsPort,
		enableSSL:      enableSSL,
	}}, nil
}

func (b *Isucon13) Start(ctx context.Context, job *domain.Job) (benchmarker.Outputs, time.Time, error) {
	dir, err := os.MkdirTemp("", "piscon-isucon13-*")
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("create benchmark output directory: %w", err)
	}
	started := false
	defer func() {
		if !started {
			_ = os.RemoveAll(dir)
		}
	}()

	// nameserver はDNSの問い合わせ先であると同時に、接続を許可するwebappのIPになる。
	args := []string{"run",
		"--target", b.conf.targetURL,
		"--nameserver", job.GetTargetIPAdress(),
		"--dns-port", strconv.Itoa(b.conf.dnsPort),
		"--enable-ssl=" + strconv.FormatBool(b.conf.enableSSL),
		"--result-path", filepath.Join(dir, "result.json"),
		"--staff-log-path", filepath.Join(dir, "staff.log"),
		"--contestant-log-path", filepath.Join(dir, "contestant.log"),
	}
	seenIPs := map[string]struct{}{job.GetTargetIPAdress(): {}}
	for _, ip := range job.GetInstanceIPAddresses() {
		if _, seen := seenIPs[ip]; ip != "" && !seen {
			args = append(args, "--webapp", ip)
			seenIPs[ip] = struct{}{}
		}
	}
	cmd := exec.CommandContext(ctx, b.conf.execPath, args...)
	cmd.Dir = b.conf.benchmarkerDir
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("get stderr pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("start benchmarker: %w", err)
	}
	started = true
	b.cmd, b.outputDir, b.score = cmd, dir, 0
	return benchmarker.Outputs{
		Stdout: stdout,
		Stderr: &infoFilter{reader: bufio.NewReader(stderr)},
	}, time.Now(), nil
}

// infoFilter removes INFO entries from zap's tab-separated console output.
type infoFilter struct {
	reader  *bufio.Reader
	pending []byte
	err     error
}

func (r *infoFilter) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	for len(r.pending) == 0 && r.err == nil {
		r.pending, r.err = r.reader.ReadBytes('\n')
		fields := bytes.SplitN(r.pending, []byte{'\t'}, 3)
		if len(fields) == 3 && bytes.Equal(fields[1], []byte("info")) {
			r.pending = nil
		}
	}
	n := copy(p, r.pending)
	r.pending = r.pending[n:]
	if len(r.pending) > 0 {
		return n, nil
	}
	return n, r.err
}

// 成功時はstdoutに結果JSONが出ないため、--result-path のファイルを読む。
// https://github.com/isucon/isucon13/blob/main/bench/cmd/bench/bench.go
type benchResult struct {
	Pass  *bool `json:"pass"`
	Score *int  `json:"score"`
}

func (b *Isucon13) readResult() (benchResult, error) {
	data, err := os.ReadFile(filepath.Join(b.outputDir, "result.json"))
	if err != nil {
		return benchResult{}, fmt.Errorf("read benchmark result: %w", err)
	}
	var result benchResult
	if err := json.Unmarshal(data, &result); err != nil {
		return benchResult{}, fmt.Errorf("decode benchmark result: %w", err)
	}
	if result.Pass == nil || result.Score == nil {
		return benchResult{}, fmt.Errorf("benchmark result must contain pass and score")
	}
	return result, nil
}

func (b *Isucon13) Wait(_ context.Context) (domain.Result, time.Time, error) {
	waitErr := b.cmd.Wait()
	finishedAt := time.Now()
	defer func() {
		_ = os.RemoveAll(b.outputDir)
		b.outputDir = ""
	}()
	result, err := b.readResult()
	if err != nil {
		return domain.ResultError, finishedAt, errors.Join(err, waitErr)
	}
	// 最終チェック失敗時には pass=false を書き、終了コード1で終了する。
	var exitErr *exec.ExitError
	if waitErr != nil && (*result.Pass || !errors.As(waitErr, &exitErr) || exitErr.ExitCode() != 1) {
		return domain.ResultError, finishedAt, fmt.Errorf("wait benchmarker: %w", waitErr)
	}
	b.score = *result.Score
	if !*result.Pass {
		return domain.ResultFailed, finishedAt, nil
	}
	return domain.ResultPassed, finishedAt, nil
}

func (b *Isucon13) CalculateScore(_ context.Context, _, _ string) (int, error) {
	if b.outputDir != "" {
		// 実行中はファイルが未作成または書き込み途中の場合がある。
		// 終了後の厳密な検証はWaitで行う。
		if result, err := b.readResult(); err == nil {
			b.score = *result.Score
		}
	}
	return b.score, nil
}
