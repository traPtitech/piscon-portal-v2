package runner

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal"
)

// 外部プロセスとportalだけを置き換え、runnerの実行・完了処理は実装をそのまま使う。
type outputProcess struct{ stdout, stderr string }

func (b outputProcess) Start(context.Context, *domain.Job) (benchmarker.Outputs, time.Time, error) {
	return benchmarker.Outputs{Stdout: strings.NewReader(b.stdout), Stderr: strings.NewReader(b.stderr)}, time.Now(), nil
}
func (outputProcess) CalculateScore(context.Context, string, string) (int, error) { return 100, nil }
func (outputProcess) Wait(context.Context) (domain.Result, time.Time, error) {
	return domain.ResultPassed, time.Now(), nil
}

type recordingPortal struct {
	closeErr  error
	sendErr   error
	result    domain.Result
	runnerErr error
	finished  int
}

func (*recordingPortal) GetJob(context.Context) (*domain.Job, error) {
	return domain.NewJob("test-job", "target"), nil
}
func (p *recordingPortal) MakeProgressStreamClient(context.Context) (portal.ProgressStreamClient, error) {
	return p, nil
}
func (p *recordingPortal) SendProgress(_ context.Context, progress *domain.Progress) error {
	if progress.GetScore() != 0 {
		return p.sendErr
	}
	return nil
}
func (p *recordingPortal) Close() error { return p.closeErr }
func (p *recordingPortal) PostJobFinished(_ context.Context, _ string, _ time.Time, result domain.Result, err error) error {
	p.result, p.runnerErr = result, err
	p.finished++
	return nil
}

// portalからの終了応答が失敗した場合は、成功として報告せず、一度だけ完了を通知する。
func TestRunReportsOutcomeOnce(t *testing.T) {
	tests := map[string]struct{ closeErr error }{
		"success":                   {},
		"portal_rejects_completion": {closeErr: io.ErrUnexpectedEOF},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := &recordingPortal{closeErr: tc.closeErr}
			r := &Runner{portal: p, benchmarker: outputProcess{stdout: "log", stderr: "staff log"}}
			err := r.Run()
			if tc.closeErr == nil {
				require.NoError(t, err)
				require.NoError(t, p.runnerErr)
				require.Equal(t, domain.ResultPassed, p.result)
			} else {
				require.ErrorIs(t, err, tc.closeErr)
				require.ErrorIs(t, p.runnerErr, tc.closeErr)
				require.Equal(t, domain.ResultError, p.result)
			}
			require.Equal(t, 1, p.finished)
		})
	}
}

// キャンセルされるまで出力を開いたままにして、実行中の進捗送信を発生させる。
type waitingProcess struct{ outputProcess }

func (waitingProcess) Start(ctx context.Context, _ *domain.Job) (benchmarker.Outputs, time.Time, error) {
	reader, writer := io.Pipe()
	go func() { <-ctx.Done(); writer.Close() }()
	return benchmarker.Outputs{Stdout: reader, Stderr: strings.NewReader("")}, time.Now(), nil
}

func TestRunReportsProgressSendFailureOnce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		sendErr := errors.New("progress unavailable")
		p := &recordingPortal{sendErr: sendErr}
		runner := &Runner{portal: p, benchmarker: waitingProcess{}}
		err := runner.Run()
		require.ErrorIs(t, err, sendErr)
		require.Equal(t, 1, strings.Count(err.Error(), sendErr.Error()))
		require.ErrorIs(t, p.runnerErr, sendErr)
		require.Equal(t, 1, strings.Count(p.runnerErr.Error(), sendErr.Error()))
		require.Equal(t, domain.ResultError, p.result)
		require.Equal(t, 1, p.finished)
	})
}
