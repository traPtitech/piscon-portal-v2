package isucon11qualify

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/isucon/isucon10-portal/proto.go/isuxportal/resources"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"google.golang.org/protobuf/proto"
)

type problemConf struct {
	benchmarkerIP string
	execPath      string
}

type Isucon11Qualify struct {
	cmd  *exec.Cmd
	conf problemConf

	errCh       chan error
	passedCh    chan bool
	latestScore atomic.Int64
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

	return &Isucon11Qualify{
		errCh:    make(chan error, 1),
		passedCh: make(chan bool, 1),
		conf: problemConf{
			execPath:      path,
			benchmarkerIP: benchmarkerIP,
		},
	}, nil
}

func (b *Isucon11Qualify) Start(ctx context.Context, job *domain.Job) (benchmarker.Outputs, time.Time, error) {
	jiaServiceURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(b.conf.benchmarkerIP, "5000"),
	}
	b.cmd = exec.CommandContext(ctx, b.conf.execPath,
		"--target", job.GetTargetIPAdress(),
		// tlsを使わない場合は`--all-adresses`にtargetだけ指定すればok
		"--all-addresses", job.GetTargetIPAdress(),
		"--jia-service-url", jiaServiceURL.String())

	reportReader, reportWriter, err := os.Pipe()
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("create report pipe: %w", err)
	}
	b.cmd.Env = append(b.cmd.Environ(), "ISUXBENCH_REPORT_FD=3")
	b.cmd.ExtraFiles = []*os.File{reportWriter}

	stdout, err := b.cmd.StdoutPipe()
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("get stdout pipe: %w", err)
	}
	stderr, err := b.cmd.StderrPipe()
	if err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("get stderr pipe: %w", err)
	}

	if err := b.cmd.Start(); err != nil {
		return benchmarker.Outputs{}, time.Time{}, fmt.Errorf("start command: %w", err)
	}
	reportWriter.Close()

	go b.watchReport(reportReader)

	return benchmarker.Outputs{
		Stdout: stdout,
		Stderr: stderr,
	}, time.Now(), nil
}

func (b *Isucon11Qualify) Wait(_ context.Context) (domain.Result, time.Time, error) {
	if err := b.cmd.Wait(); err != nil {
		return domain.ResultError, time.Now(), fmt.Errorf("wait command: %w", err)
	}
	endTime := time.Now()

	err := <-b.errCh
	if err != nil {
		return domain.ResultError, endTime, err
	}

	if <-b.passedCh {
		return domain.ResultPassed, endTime, nil
	}
	return domain.ResultFailed, endTime, nil
}

func (b *Isucon11Qualify) CalculateScore(_ context.Context, _, _ string) (int, error) {
	return int(b.latestScore.Load()), nil
}

func (b *Isucon11Qualify) watchReport(report io.ReadCloser) {
	defer report.Close()
	defer close(b.errCh)
	defer close(b.passedCh)
	for {
		result, err := readResult(report)
		if err != nil {
			b.errCh <- err
			return
		}
		b.latestScore.Store(result.Score)
		if result.Finished {
			b.passedCh <- result.Passed
			return
		}
	}
}

// ISUCON11予選ではISUCON10のReporterを利用している
// https://github.com/isucon/isucon10-portal/blob/master/bench-tool.go/benchrun/reporter.go
func readResult(report io.Reader) (*resources.BenchmarkResult, error) {
	var sizeData [2]byte
	_, err := io.ReadFull(report, sizeData[:])
	if err != nil {
		return nil, fmt.Errorf("parse benchmark result: %w", err)
	}

	size := int(binary.BigEndian.Uint16(sizeData[:]))
	data := make([]byte, size)
	_, err = io.ReadFull(report, data)
	if err != nil {
		return nil, fmt.Errorf("parse benchmark result: %w", err)
	}

	var result resources.BenchmarkResult
	if err := proto.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("decode benchmark result: %w", err)
	}
	return &result, nil
}
