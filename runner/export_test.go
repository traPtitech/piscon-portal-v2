package runner

import (
	"context"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal"
)

const (
	BufSizeExported              = bufSize
	SendProgressIntervalExported = sendProgressInterval
)

var (
	CaptureStreamOutput = captureStreamOutput
)

// NewRunnerForTest constructs a Runner with an already created Benchmarker,
// bypassing the problem registry.
func NewRunnerForTest(p portal.Portal, b benchmarker.Benchmarker) *Runner {
	return &Runner{portal: p, benchmarker: b}
}

func (r *Runner) StreamJobProgressExported(
	ctx context.Context, streamClient portal.ProgressStreamClient,
	job *domain.Job, startedAt time.Time,
	stdoutBdr, stderrBdr *SyncStringBuilder,
	stdoutErrChan, stderrErrChan chan error,
) error {
	return r.streamJobProgress(ctx, streamClient, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
}
