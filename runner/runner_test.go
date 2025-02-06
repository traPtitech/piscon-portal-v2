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
	benchmarkerMock "github.com/traPtitech/piscon-portal-v2/runner/benchmarker/mock"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal/mock"
	"go.uber.org/mock/gomock"
	"golang.org/x/sync/errgroup"
)

func Test_captureStreamOutput(t *testing.T) {
	testCases := map[string]struct {
		writeFunc func(*testing.T, io.WriteCloser, *runner.Builder)
		result    string
	}{
		"ok": {
			writeFunc: func(t *testing.T, w io.WriteCloser, b *runner.Builder) {
				t.Helper()
				for i := range 10 {
					_, err := w.Write(bytes.Repeat([]byte("a"), runner.BufSizeExported))
					require.NoError(t, err)
					assert.Equal(t, string(bytes.Repeat([]byte("a"), runner.BufSizeExported*i)), b.String())
				}
				w.Close()
			},
			result: strings.Repeat("a", runner.BufSizeExported*10),
		},
		"短くてもエラー無し": {
			writeFunc: func(t *testing.T, w io.WriteCloser, _ *runner.Builder) {
				t.Helper()
				_, err := w.Write([]byte("abc"))
				require.NoError(t, err)
				w.Close()
			},
			result: "abc",
		},
		"0文字でもエラー無し": {
			writeFunc: func(t *testing.T, w io.WriteCloser, _ *runner.Builder) {
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
			bdr := &runner.Builder{}

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
	ctrl := gomock.NewController(t)

	setupRunner := func(t *testing.T) (*runner.Runner, *mock.MockPortal, *mock.MockProgressStreamClient, *benchmarkerMock.MockBenchmarker) {
		t.Helper()
		portal := mock.NewMockPortal(ctrl)
		benchmarker := benchmarkerMock.NewMockBenchmarker(ctrl)
		streamClient := mock.NewMockProgressStreamClient(ctrl)
		r := runner.Prepare(portal, benchmarker)

		portal.EXPECT().MakeProgressStreamClient(gomock.Any()).
			Return(streamClient, nil)
		streamClient.EXPECT().Close().Return(nil)

		return r, portal, streamClient, benchmarker
	}

	setupArgs := func(t *testing.T) (*domain.Job, time.Time, *runner.Builder, *runner.Builder, chan error, chan error) {
		t.Helper()
		job := domain.NewJob("id", "target")
		startedAt := time.Now()
		stdoutBdr := &runner.Builder{}
		stderrBdr := &runner.Builder{}
		stdoutErrChan := make(chan error, 1)
		stderrErrChan := make(chan error, 1)
		return job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan
	}

	type testCase struct {
		useSynctest bool
		setupMocks  func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time)
		writeFunc   func(stdoutBdr, stderrBdr *runner.Builder, stdoutErrChan, stderrErrChan chan error)
		expectedErr error
	}

	tests := map[string]testCase{
		"stdoutとstderrから何も来ずにすぐ終わる": {
			useSynctest: false,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				bm.EXPECT().CalculateScore(gomock.Any(), "", "").Return(0, nil)
				sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "", "", 0, startedAt)).Return(nil)
			},
			writeFunc: func(_, _ *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					stdoutErrChan <- nil
					stderrErrChan <- nil
				}()
			},
			expectedErr: nil,
		},
		"stdoutにデータが来てすぐ終わる": {
			useSynctest: false,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				bm.EXPECT().CalculateScore(gomock.Any(), "abc", "").Return(0, nil)
				sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "", 0, startedAt)).Return(nil)
			},
			writeFunc: func(stdoutBdr, _ *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					stdoutErrChan <- nil
					stderrErrChan <- nil
				}()
			},
			expectedErr: nil,
		},
		"stderrにデータが来てすぐ終わる": {
			useSynctest: false,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				bm.EXPECT().CalculateScore(gomock.Any(), "", "abc").Return(0, nil)
				sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "", "abc", 0, startedAt)).Return(nil)
			},
			writeFunc: func(_, stderrBdr *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					_, err := stderrBdr.WriteString("abc")
					require.NoError(t, err)
					stdoutErrChan <- nil
					stderrErrChan <- nil
				}()
			},
			expectedErr: nil,
		},
		"stdoutとstderrにデータが来てすぐ終わる": {
			useSynctest: false,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil)
				sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil)
			},
			writeFunc: func(stdoutBdr, stderrBdr *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("def")
					require.NoError(t, err)
					stdoutErrChan <- nil
					stderrErrChan <- nil
				}()
			},
			expectedErr: nil,
		},
		"stdoutでエラー": {
			useSynctest: false,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil)
				sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil)
			},
			writeFunc: func(stdoutBdr, stderrBdr *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("def")
					require.NoError(t, err)
					stdoutErrChan <- assert.AnError
					stderrErrChan <- nil
				}()
			},
			expectedErr: assert.AnError,
		},
		"stderrでエラー": {
			useSynctest: false,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil)
				sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil)
			},
			writeFunc: func(stdoutBdr, stderrBdr *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("def")
					require.NoError(t, err)
					stdoutErrChan <- nil
					stderrErrChan <- assert.AnError
				}()
			},
			expectedErr: assert.AnError,
		},
		"intervalをまたいでも問題なし": {
			useSynctest: true,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				gomock.InOrder(
					bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil).Call,
					sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil).Call,
					bm.EXPECT().CalculateScore(gomock.Any(), "abcdef", "defghi").Return(100, nil).Call,
					sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abcdef", "defghi", 100, startedAt)).Return(nil).Call,
				)
			},
			writeFunc: func(stdoutBdr, stderrBdr *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("def")
					require.NoError(t, err)
					time.Sleep(runner.SendProgressIntervalExported * 3 / 2)
					_, err = stdoutBdr.WriteString("def")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("ghi")
					require.NoError(t, err)
					stdoutErrChan <- nil
					stderrErrChan <- nil
				}()
			},
			expectedErr: nil,
		},
		"intervalをまたいで最後のcalculateScoreでエラー": {
			useSynctest: true,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				gomock.InOrder(
					bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil).Call,
					sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil).Call,
					bm.EXPECT().CalculateScore(gomock.Any(), "abcdef", "defghi").Return(100, assert.AnError).Call,
				)
			},
			writeFunc: func(stdoutBdr, stderrBdr *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("def")
					require.NoError(t, err)
					time.Sleep(runner.SendProgressIntervalExported * 3 / 2)
					_, err = stdoutBdr.WriteString("def")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("ghi")
					require.NoError(t, err)
					stdoutErrChan <- nil
					stderrErrChan <- nil
				}()
			},
			expectedErr: assert.AnError,
		},
		"intervalをまたいで最後のsendProgressでエラー": {
			useSynctest: true,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				gomock.InOrder(
					bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil).Call,
					sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil).Call,
					bm.EXPECT().CalculateScore(gomock.Any(), "abcdef", "defghi").Return(100, nil).Call,
					sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abcdef", "defghi", 100, startedAt)).Return(assert.AnError).Call,
				)
			},
			writeFunc: func(stdoutBdr, stderrBdr *runner.Builder, stdoutErrChan, stderrErrChan chan error) {
				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("def")
					require.NoError(t, err)
					time.Sleep(runner.SendProgressIntervalExported * 3 / 2)
					_, err = stdoutBdr.WriteString("def")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("ghi")
					require.NoError(t, err)
					stdoutErrChan <- nil
					stderrErrChan <- nil
				}()
			},
			expectedErr: assert.AnError,
		},
		"intervalでデータを読んでCalculateScoreでエラー": {
			useSynctest: true,
			setupMocks: func(sc *mock.MockProgressStreamClient, bm *benchmarkerMock.MockBenchmarker, startedAt time.Time) {
				gomock.InOrder(
					bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, assert.AnError).Call,
					// SendProgress は呼ばれない想定
					bm.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(100, nil).Call,
					sc.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 100, startedAt)).Return(nil).Call,
				)
			},
			writeFunc: func(stdoutBdr, stderrBdr *runner.Builder, _ chan error, _ chan error) {
				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("def")
					require.NoError(t, err)
					time.Sleep(runner.SendProgressIntervalExported * 3 / 2)
				}()
			},
			expectedErr: assert.AnError,
		},
	}

	runFlow := func(tc testCase, t *testing.T) {
		r, _, sc, bm := setupRunner(t)
		ctx := context.Background()
		job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

		tc.setupMocks(sc, bm, startedAt)
		tc.writeFunc(stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)

		err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		if tc.expectedErr != nil {
			assert.ErrorIs(t, err, tc.expectedErr)
		} else {
			assert.NoError(t, err)
		}
	}

	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if tc.useSynctest {
				synctest.Run(func() {
					runFlow(tc, t)
				})
			} else {
				runFlow(tc, t)
			}
		})
	}
}

// 他のコンポーネントで十分テストが書かれているため、Runは一番シンプルな、すぐに入力がおわるケースのみ書く
func TestRun(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	portal := mock.NewMockPortal(ctrl)
	benchmarker := benchmarkerMock.NewMockBenchmarker(ctrl)
	streamClient := mock.NewMockProgressStreamClient(ctrl)

	portal.EXPECT().MakeProgressStreamClient(gomock.Any()).Return(streamClient, nil)
	portal.EXPECT().GetJob(gomock.Any()).Return(domain.NewJob("id", "target"), nil)

	startedAt := time.Now()
	stdout := strings.Repeat("a", runner.BufSizeExported*3)
	stderr := strings.Repeat("b", runner.BufSizeExported)
	benchmarker.EXPECT().Start(gomock.Any(), gomock.Any()).
		Return(io.NopCloser(strings.NewReader(stdout)), io.NopCloser(strings.NewReader(stderr)), startedAt, nil)

	streamClient.EXPECT().Close().Return(nil)

	benchmarker.EXPECT().CalculateScore(gomock.Any(), stdout, stderr).Return(100, nil)

	streamClient.EXPECT().
		SendProgress(gomock.Any(),
			domain.NewProgress("id", stdout, stderr, 100, startedAt)).
		Return(nil)

	benchmarker.EXPECT().Wait(gomock.Any()).Return(domain.ResultPassed, time.Now(), nil)

	portal.EXPECT().PostJobFinished(gomock.Any(), "id", gomock.Any(), domain.ResultPassed, nil).Return(nil)

	// テスト対象の関数を用意
	r := runner.Prepare(portal, benchmarker)

	// テスト対象の関数を呼び出す
	err := r.Run()
	assert.NoError(t, err)
}
