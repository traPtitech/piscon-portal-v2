package runner

import (
	"context"
	"time"

	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

const (
	BufSizeExported              = bufSize
	SendProgressIntervalExported = sendProgressInterval
)

var (
	CaptureStreamOutput = captureStreamOutput
)

func (r *Runner) StreamJobProgressExported(
	ctx context.Context, job *domain.Job, startedAt time.Time,
	stdoutBdr, stderrBdr *Builder,
	stdoutErrChan, stderrErrChan chan error,
) error {
	return r.streamJobProgress(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
}
