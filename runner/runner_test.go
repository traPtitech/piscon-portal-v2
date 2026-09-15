package runner_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"
	"testing/synctest"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/runner"
	"github.com/traPtitech/piscon-portal-v2/runner/benchmarker"
	benchmarkerMock "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/mock"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal/mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/sync/errgroup"
)

func Test_captureStreamOutput(t *testing.T) {
	testCases := map[string]struct {
		writeFunc func(*testing.T, io.WriteCloser, *runner.SyncStringBuilder)
		result    string
	}{
		"ok": {
			writeFunc: func(t *testing.T, w io.WriteCloser, b *runner.SyncStringBuilder) {
				t.Helper()
				for i := range 10 {
					_, err := w.Write(bytes.Repeat([]byte("a"), runner.BufSizeExported))
					require.NoError(t, err)
					time.Sleep(1 * time.Millisecond)                        // 速すぎるとテストが通らないので適当に待つ
					assert.Len(t, b.String(), runner.BufSizeExported*(i+1)) // 文字列が長いので、長さを表示するようにするため
					assert.Equal(t, string(bytes.Repeat([]byte("a"), runner.BufSizeExported*(i+1))), b.String())
				}
				w.Close()
			},
			result: strings.Repeat("a", runner.BufSizeExported*10),
		},
		"短くてもエラー無し": {
			writeFunc: func(t *testing.T, w io.WriteCloser, _ *runner.SyncStringBuilder) {
				t.Helper()
				_, err := w.Write([]byte("abc"))
				require.NoError(t, err)
				w.Close()
			},
			result: "abc",
		},
		"0文字でもエラー無し": {
			writeFunc: func(t *testing.T, w io.WriteCloser, _ *runner.SyncStringBuilder) {
				t.Helper()
				w.Close()
			},
			result: "",
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()

			// パイプを用意して、writerの方に書き込むことでテスト対象関数のreaderにデータを流す
			pr, pw := io.Pipe()
			bdr := &runner.SyncStringBuilder{}

			eg := &errgroup.Group{}
			eg.Go(func() error {
				return runner.CaptureStreamOutput(ctx, pr, bdr)
			})

			testCase.writeFunc(t, pw, bdr)

			err := eg.Wait()

			assert.NoError(t, err)
			assert.Equal(t, testCase.result, bdr.String())
		})
	}
}

func Test_streamJobProgress(t *testing.T) {
	type testCase struct {
		useSynctest bool
		// setupMocks には初期進捗の送信以降に期待する呼び出しを書く。
		setupMocks  func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time)
		writeFunc   func(t *testing.T, stdoutBdr, stderrBdr *runner.SyncStringBuilder, stdoutErrChan, stderrErrChan chan error)
		expectedErr error
	}

	tests := map[string]testCase{
		"stdoutとstderrから何も来ずにすぐ終わる": {
			setupMocks: func(_ *mock.MockProgressStreamClient, _ *benchmarkerMock.MockBenchmarker, _ time.Time) {},
			writeFunc: func(_ *testing.T, _, _ *runner.SyncStringBuilder, stdoutErrChan, stderrErrChan chan error) {
				stdoutErrChan <- nil
				stderrErrChan <- nil
			},
			expectedErr: nil,
		},
		"データが来てもintervalをまたがなければ送信しない": {
			setupMocks: func(_ *mock.MockProgressStreamClient, _ *benchmarkerMock.MockBenchmarker, _ time.Time) {},
			writeFunc: func(t *testing.T, stdoutBdr, stderrBdr *runner.SyncStringBuilder, stdoutErrChan, stderrErrChan chan error) {
				t.Helper()
				_, err := stdoutBdr.WriteString("abc")
				require.NoError(t, err)
				_, err = stderrBdr.WriteString("def")
				require.NoError(t, err)
				stdoutErrChan <- nil
				stderrErrChan <- nil
			},
			expectedErr: nil,
		},
		"stdoutでエラー": {
			setupMocks: func(_ *mock.MockProgressStreamClient, _ *benchmarkerMock.MockBenchmarker, _ time.Time) {},
			writeFunc: func(_ *testing.T, _, _ *runner.SyncStringBuilder, stdoutErrChan, stderrErrChan chan error) {
				stdoutErrChan <- assert.AnError
				stderrErrChan <- nil
			},
			expectedErr: assert.AnError,
		},
		"stderrでエラー": {
			setupMocks: func(_ *mock.MockProgressStreamClient, _ *benchmarkerMock.MockBenchmarker, _ time.Time) {},
			writeFunc: func(_ *testing.T, _, _ *runner.SyncStringBuilder, stdoutErrChan, stderrErrChan chan error) {
				stderrErrChan <- assert.AnError
				stdoutErrChan <- nil
			},
			expectedErr: assert.AnError,
		},
		"intervalをまたいでも問題なし": {
			useSynctest: true,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				gomock.InOrder(
					bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(100, nil).Call,
					sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 100, startedAt)).Return(nil).Call,
				)
			},
			writeFunc: func(t *testing.T, stdoutBdr, stderrBdr *runner.SyncStringBuilder, stdoutErrChan, stderrErrChan chan error) {
				t.Helper()
				_, err := stdoutBdr.WriteString("abc")
				require.NoError(t, err)
				_, err = stderrBdr.WriteString("def")
				require.NoError(t, err)
				time.Sleep(runner.SendProgressIntervalExported * 3 / 2)
				stdoutErrChan <- nil
				stderrErrChan <- nil
			},
			expectedErr: nil,
		},
		"intervalのCalculateScoreでエラー": {
			useSynctest: true,
			setupMocks: func(_ *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, _ time.Time) {
				bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, assert.AnError)
			},
			writeFunc: func(t *testing.T, stdoutBdr, stderrBdr *runner.SyncStringBuilder, _, _ chan error) {
				t.Helper()
				_, err := stdoutBdr.WriteString("abc")
				require.NoError(t, err)
				_, err = stderrBdr.WriteString("def")
				require.NoError(t, err)
				time.Sleep(runner.SendProgressIntervalExported * 3 / 2)
			},
			expectedErr: assert.AnError,
		},
		"intervalのSendProgressでエラー": {
			useSynctest: true,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				gomock.InOrder(
					bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(100, nil).Call,
					sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 100, startedAt)).Return(assert.AnError).Call,
				)
			},
			writeFunc: func(t *testing.T, stdoutBdr, stderrBdr *runner.SyncStringBuilder, _, _ chan error) {
				t.Helper()
				_, err := stdoutBdr.WriteString("abc")
				require.NoError(t, err)
				_, err = stderrBdr.WriteString("def")
				require.NoError(t, err)
				time.Sleep(runner.SendProgressIntervalExported * 3 / 2)
			},
			expectedErr: assert.AnError,
		},
	}

	runFlow := func(t *testing.T, tc testCase) {
		t.Helper()
		ctrl := gomock.NewController(t)
		portal := mock.NewMockPortal(ctrl)
		bm := benchmarkerMock.NewMockBenchmarker(ctrl)
		sc := mock.NewMockProgressStreamClient(ctrl)
		r := runner.NewRunnerForTest(portal, bm)

		job := domain.NewJob("id", "target")
		startedAt := time.Now()
		stdoutBdr := &runner.SyncStringBuilder{}
		stderrBdr := &runner.SyncStringBuilder{}
		stdoutErrChan := make(chan error, 1)
		stderrErrChan := make(chan error, 1)

		// 開始時の進捗は必ず送信される。
		sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "", "", 0, startedAt)).Return(nil)
		tc.setupMocks(sc, bm, startedAt)

		done := make(chan struct{})
		go func() {
			defer close(done)
			tc.writeFunc(t, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		}()

		err := r.StreamJobProgressExported(context.Background(), sc, job, startedAt,
			stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		<-done

		if tc.expectedErr != nil {
			assert.ErrorIs(t, err, tc.expectedErr)
		} else {
			assert.NoError(t, err)
		}
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.useSynctest {
				synctest.Test(t, func(t *testing.T) {
					runFlow(t, tc)
				})
				return
			}
			t.Parallel()
			runFlow(t, tc)
		})
	}
}

// 他のコンポーネントで十分テストが書かれているため、Runは一番シンプルな、すぐに入力がおわるケースのみ書く
func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)

	portal := mock.NewMockPortal(ctrl)
	mockBenchmarker := benchmarkerMock.NewMockBenchmarker(ctrl)
	streamClient := mock.NewMockProgressStreamClient(ctrl)

	startedAt := time.Now()
	stdout := strings.Repeat("a", runner.BufSizeExported*3)
	stderr := strings.Repeat("b", runner.BufSizeExported)

	portal.EXPECT().GetJob(gomock.Any()).Return(domain.NewJob("id", "target"), nil)
	mockBenchmarker.EXPECT().Start(gomock.Any(), gomock.Any()).
		Return(benchmarker.Outputs{
			Stdout: strings.NewReader(stdout),
			Stderr: strings.NewReader(stderr),
		}, startedAt, nil)
	portal.EXPECT().MakeProgressStreamClient(gomock.Any()).Return(streamClient, nil)

	// 最終スコアはベンチマーカーの最終レポートの受信を待ってから確定するため、
	// CalculateScore と最後の SendProgress は Wait の後でなければならない。
	gomock.InOrder(
		streamClient.EXPECT().
			SendProgress(gomock.Any(), domain.NewProgress("id", "", "", 0, startedAt)).
			Return(nil).Call,
		mockBenchmarker.EXPECT().Wait(gomock.Any()).Return(domain.ResultPassed, time.Now(), nil).Call,
		mockBenchmarker.EXPECT().CalculateScore(gomock.Any(), stdout, stderr).Return(100, nil).Call,
		streamClient.EXPECT().
			SendProgress(gomock.Any(), domain.NewProgress("id", stdout, stderr, 100, startedAt)).
			Return(nil).Call,
		streamClient.EXPECT().Close().Return(nil).Call,
		portal.EXPECT().PostJobFinished(gomock.Any(), "id", gomock.Any(), domain.ResultPassed, nil).Return(nil).Call,
	)

	r := runner.NewRunnerForTest(portal, mockBenchmarker)

	err := r.Run()
	assert.NoError(t, err)
}

// ログの収集に失敗した場合も、完了通知は1回だけ送信される。
func TestRunStreamError(t *testing.T) {
	ctrl := gomock.NewController(t)

	portal := mock.NewMockPortal(ctrl)
	mockBenchmarker := benchmarkerMock.NewMockBenchmarker(ctrl)

	startedAt := time.Now()

	portal.EXPECT().GetJob(gomock.Any()).Return(domain.NewJob("id", "target"), nil)
	mockBenchmarker.EXPECT().Start(gomock.Any(), gomock.Any()).
		Return(benchmarker.Outputs{
			Stdout: strings.NewReader(""),
			Stderr: strings.NewReader(""),
		}, startedAt, nil)
	portal.EXPECT().MakeProgressStreamClient(gomock.Any()).Return(nil, assert.AnError)

	// 進捗を送れないので、ベンチマーカーは停止させて回収するだけ。
	mockBenchmarker.EXPECT().Wait(gomock.Any()).Return(domain.ResultError, time.Now(), nil)

	portal.EXPECT().
		PostJobFinished(gomock.Any(), "id", gomock.Any(), domain.ResultError, gomock.Not(gomock.Nil())).
		Return(nil).
		Times(1)

	r := runner.NewRunnerForTest(portal, mockBenchmarker)

	err := r.Run()
	assert.NoError(t, err)
}

// 最終スコアの送信に失敗しても、ベンチマークの合否自体は確定しているため、
// resultは変えずにエラーだけ完了通知に載せる。
func TestRunFinalProgressError(t *testing.T) {
	tests := map[string]struct {
		calcErr  error
		sendErr  error
		closeErr error
	}{
		"最終のCalculateScoreでエラー": {calcErr: assert.AnError},
		"最終のSendProgressでエラー":   {sendErr: assert.AnError},
		"進捗ストリームのCloseでエラー":     {closeErr: assert.AnError},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			portal := mock.NewMockPortal(ctrl)
			mockBenchmarker := benchmarkerMock.NewMockBenchmarker(ctrl)
			streamClient := mock.NewMockProgressStreamClient(ctrl)

			startedAt := time.Now()

			portal.EXPECT().GetJob(gomock.Any()).Return(domain.NewJob("id", "target"), nil)
			mockBenchmarker.EXPECT().Start(gomock.Any(), gomock.Any()).
				Return(benchmarker.Outputs{
					Stdout: strings.NewReader("out"),
					Stderr: strings.NewReader("err"),
				}, startedAt, nil)
			portal.EXPECT().MakeProgressStreamClient(gomock.Any()).Return(streamClient, nil)
			streamClient.EXPECT().
				SendProgress(gomock.Any(), domain.NewProgress("id", "", "", 0, startedAt)).
				Return(nil)
			mockBenchmarker.EXPECT().Wait(gomock.Any()).Return(domain.ResultPassed, time.Now(), nil)

			mockBenchmarker.EXPECT().CalculateScore(gomock.Any(), "out", "err").Return(100, tc.calcErr)
			if tc.calcErr == nil {
				streamClient.EXPECT().
					SendProgress(gomock.Any(), domain.NewProgress("id", "out", "err", 100, startedAt)).
					Return(tc.sendErr)
			}
			streamClient.EXPECT().Close().Return(tc.closeErr)

			portal.EXPECT().
				PostJobFinished(gomock.Any(), "id", gomock.Any(), domain.ResultPassed, gomock.Not(gomock.Nil())).
				Return(nil).
				Times(1)

			r := runner.NewRunnerForTest(portal, mockBenchmarker)

			assert.NoError(t, r.Run())
		})
	}
}
