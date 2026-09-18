// Package isuxbench runs a benchmarker child process that reports its results
// through the ISUXBENCH_REPORT_FD protocol of the isuxportal supervisor.
//
// ISUCON11の予選・本選のベンチマーカーはどちらもISUCON10のReporterを利用しており、
// 結果の通知方法が共通である。
// https://github.com/isucon/isucon10-portal/blob/master/bench-tool.go/benchrun/reporter.go
package isuxbench

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/isucon/isucon10-portal/proto.go/isuxportal/resources"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"google.golang.org/protobuf/proto"
)

// reportFD is the file descriptor number that the child process writes reports to.
// [exec.Cmd.ExtraFiles] の先頭がfd 3になる。
const reportFD = 3

// Process is a running benchmarker child process.
type Process struct {
	cmd *exec.Cmd
	// resultCh receives exactly one value once the report stream ends.
	resultCh    chan result
	latestScore atomic.Int64
}

type result struct {
	passed bool
	err    error
}

// Start attaches a report pipe to cmd as [reportFD], sets ISUXBENCH_REPORT_FD,
// starts cmd and begins reading reports in the background.
// cmd.ExtraFiles, cmd.Stdout and cmd.Stderr must not be set by the caller.
// cmd.Env may be set; ISUXBENCH_REPORT_FD is appended to it.
//
// After a successful call, [Process.Wait] must be called in order to release
// associated system resources.
func Start(cmd *exec.Cmd) (*Process, benchmarker.Outputs, time.Time, error) {
	reportReader, reportWriter, err := os.Pipe()
	if err != nil {
		return nil, benchmarker.Outputs{}, time.Time{}, fmt.Errorf("create report pipe: %w", err)
	}
	// 書き込み側は子プロセスに渡した後、親では不要になる。
	// 起動に失敗した場合は読み取り側も解放する。
	defer reportWriter.Close()
	started := false
	defer func() {
		if !started {
			reportReader.Close()
		}
	}()

	cmd.Env = append(cmd.Environ(), fmt.Sprintf("ISUXBENCH_REPORT_FD=%d", reportFD))
	cmd.ExtraFiles = []*os.File{reportWriter}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, benchmarker.Outputs{}, time.Time{}, fmt.Errorf("get stdout pipe: %w", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, benchmarker.Outputs{}, time.Time{}, fmt.Errorf("get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, benchmarker.Outputs{}, time.Time{}, fmt.Errorf("start command: %w", err)
	}
	started = true

	p := &Process{
		cmd:      cmd,
		resultCh: make(chan result, 1),
	}
	go p.watchReport(reportReader)

	return p, benchmarker.Outputs{Stdout: stdout, Stderr: stderr}, time.Now(), nil
}

// LatestScore returns the score of the latest report received so far.
func (p *Process) LatestScore() int {
	return int(p.latestScore.Load())
}

// Wait waits for the process to exit and for the report stream to end, then
// determines the result from the last report.
// The caller must finish reading stdout and stderr before calling Wait.
//
// 呼び出しが返った時点で最終レポートの内容が [Process.LatestScore] に反映されている。
func (p *Process) Wait() (domain.Result, time.Time, error) {
	waitErr := p.cmd.Wait()
	finishedAt := time.Now()

	// 子プロセスが終了すればパイプの書き込み側が閉じるため、watchReportも必ず終わる。
	res := <-p.resultCh

	if waitErr != nil {
		return domain.ResultError, finishedAt, fmt.Errorf("wait command: %w", waitErr)
	}
	if res.err != nil {
		return domain.ResultError, finishedAt, res.err
	}
	if res.passed {
		return domain.ResultPassed, finishedAt, nil
	}
	return domain.ResultFailed, finishedAt, nil
}

func (p *Process) watchReport(report io.ReadCloser) {
	defer report.Close()
	for {
		res, err := readResult(report)
		if err != nil {
			// 最終レポートを受け取る前にストリームが終わった場合もここに来る。
			p.resultCh <- result{err: err}
			return
		}
		p.latestScore.Store(res.GetScore())
		if res.GetFinished() {
			p.resultCh <- result{passed: res.GetPassed()}
			return
		}
	}
}

func readResult(report io.Reader) (*resources.BenchmarkResult, error) {
	var sizeData [2]byte
	if _, err := io.ReadFull(report, sizeData[:]); err != nil {
		return nil, fmt.Errorf("parse benchmark result: %w", err)
	}

	size := int(binary.BigEndian.Uint16(sizeData[:]))
	data := make([]byte, size)
	if _, err := io.ReadFull(report, data); err != nil {
		return nil, fmt.Errorf("parse benchmark result: %w", err)
	}

	var res resources.BenchmarkResult
	if err := proto.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("decode benchmark result: %w", err)
	}
	return &res, nil
}
