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
		writeFunc func(*testing.T, io.WriteCloser, *strings.Builder)
		result    string
	}{
		"ok": {
			writeFunc: func(t *testing.T, w io.WriteCloser, b *strings.Builder) {
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
			writeFunc: func(t *testing.T, w io.WriteCloser, _ *strings.Builder) {
				t.Helper()
				_, err := w.Write([]byte("abc"))
				require.NoError(t, err)
				w.Close()
			},
			result: "abc",
		},
		"0文字でもエラー無し": {
			writeFunc: func(t *testing.T, w io.WriteCloser, _ *strings.Builder) {
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
			bdr := &strings.Builder{}

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

	setupArgs := func(t *testing.T) (*domain.Job, time.Time, *strings.Builder, *strings.Builder, chan error, chan error) {
		t.Helper()

		job := domain.NewJob("id", "target")
		startedAt := time.Now()
		stdoutBdr := &strings.Builder{}
		stderrBdr := &strings.Builder{}
		stdoutErrChan := make(chan error, 1)
		stderrErrChan := make(chan error, 1)
		return job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan
	}

	// mockの設定とかが複雑すぎるので、TDTを諦める。

	t.Run("stdoutとstderrから何も来ずにすぐ終わる", func(t *testing.T) {

		t.Parallel()
		r, _, streamClient, benchmarker := setupRunner(t)

		ctx := context.Background()
		job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

		benchmarker.EXPECT().CalculateScore(gomock.Any(), "", "").Return(0, nil)
		streamClient.EXPECT().
			SendProgress(gomock.Any(), domain.NewProgress("id", "", "", 0, startedAt)).Return(nil)

		go func() {
			stdoutErrChan <- nil
			stderrErrChan <- nil
		}()

		err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		assert.NoError(t, err)
	})

	t.Run("stdoutにデータが来てすぐ終わる", func(t *testing.T) {
		t.Parallel()
		r, _, streamClient, benchmarker := setupRunner(t)

		ctx := context.Background()
		job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

		benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "").Return(0, nil)
		streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "", 0, startedAt)).Return(nil)

		go func() {
			_, err := stdoutBdr.WriteString("abc")
			require.NoError(t, err)

			stdoutErrChan <- nil
			stderrErrChan <- nil
		}()

		err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		assert.NoError(t, err)
	})

	t.Run("stderrにデータが来てすぐ終わる", func(t *testing.T) {
		t.Parallel()
		r, _, streamClient, benchmarker := setupRunner(t)

		ctx := context.Background()
		job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

		benchmarker.EXPECT().CalculateScore(gomock.Any(), "", "abc").Return(0, nil)
		streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "", "abc", 0, startedAt)).Return(nil)

		go func() {
			_, err := stderrBdr.WriteString("abc")
			require.NoError(t, err)

			stdoutErrChan <- nil
			stderrErrChan <- nil
		}()

		err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		assert.NoError(t, err)
	})

	t.Run("stdoutとstderrにデータが来てすぐ終わる", func(t *testing.T) {
		t.Parallel()
		r, _, streamClient, benchmarker := setupRunner(t)

		ctx := context.Background()
		job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

		benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil)
		streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil)

		go func() {
			_, err := stdoutBdr.WriteString("abc")
			require.NoError(t, err)
			_, err = stderrBdr.WriteString("def")
			require.NoError(t, err)

			stdoutErrChan <- nil
			stderrErrChan <- nil
		}()

		err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		assert.NoError(t, err)
	})

	t.Run("stdoutでエラー", func(t *testing.T) {
		t.Parallel()
		r, _, streamClient, benchmarker := setupRunner(t)

		ctx := context.Background()
		job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

		benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil)
		streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil)

		go func() {
			_, err := stdoutBdr.WriteString("abc")
			require.NoError(t, err)
			_, err = stderrBdr.WriteString("def")
			require.NoError(t, err)

			stdoutErrChan <- assert.AnError
			stderrErrChan <- nil
		}()

		err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("stderrでエラー", func(t *testing.T) {
		t.Parallel()
		r, _, streamClient, benchmarker := setupRunner(t)

		ctx := context.Background()
		job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

		benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil)
		streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil)

		go func() {
			_, err := stdoutBdr.WriteString("abc")
			require.NoError(t, err)
			_, err = stderrBdr.WriteString("def")
			require.NoError(t, err)

			stdoutErrChan <- nil
			stderrErrChan <- assert.AnError
		}()

		err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("intervalをまたいでも問題なし", func(t *testing.T) {
		t.Parallel()
		synctest.Run(
			func() {
				r, _, streamClient, benchmarker := setupRunner(t)

				ctx := context.Background()
				job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

				gomock.InOrder(
					benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil).Call,
					streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil).Call,
					benchmarker.EXPECT().CalculateScore(gomock.Any(), "abcdef", "defghi").Return(100, nil).Call,
					streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abcdef", "defghi", 100, startedAt)).Return(nil).Call,
				)

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

				err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
				assert.NoError(t, err)
			})
	})

	t.Run("intervalをまたいで最後のcalculateScoreでエラー", func(t *testing.T) {
		t.Parallel()
		synctest.Run(
			func() {
				r, _, streamClient, benchmarker := setupRunner(t)

				ctx := context.Background()
				job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

				gomock.InOrder(
					benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil).Call,
					streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil).Call,
					benchmarker.EXPECT().CalculateScore(gomock.Any(), "abcdef", "defghi").Return(100, assert.AnError).Call,
				)

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

				err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
				assert.ErrorIs(t, err, assert.AnError)
			})
	})

	t.Run("intervalをまたいで最後のsendProgressでエラー", func(t *testing.T) {
		t.Parallel()
		synctest.Run(
			func() {
				r, _, streamClient, benchmarker := setupRunner(t)

				ctx := context.Background()
				job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

				gomock.InOrder(
					benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, nil).Call,
					streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil).Call,
					benchmarker.EXPECT().CalculateScore(gomock.Any(), "abcdef", "defghi").Return(100, nil).Call,
					streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abcdef", "defghi", 100, startedAt)).Return(assert.AnError).Call,
				)

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

				err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
				assert.ErrorIs(t, err, assert.AnError)
			})
	})

	t.Run("intervalでデータを読んでCalculateScoreでエラー", func(t *testing.T) {
		t.Parallel()
		synctest.Run(
			func() {
				r, _, streamClient, benchmarker := setupRunner(t)

				ctx := context.Background()
				job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan := setupArgs(t)

				gomock.InOrder(
					benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(0, assert.AnError).Call,
					// streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 0, startedAt)).Return(nil).Call,
					benchmarker.EXPECT().CalculateScore(gomock.Any(), "abc", "def").Return(100, nil).Call,
					streamClient.EXPECT().SendProgress(gomock.Any(), domain.NewProgress("id", "abc", "def", 100, startedAt)).Return(nil).Call,
				)

				go func() {
					_, err := stdoutBdr.WriteString("abc")
					require.NoError(t, err)
					_, err = stderrBdr.WriteString("def")
					require.NoError(t, err)

					time.Sleep(runner.SendProgressIntervalExported * 3 / 2)
				}()

				err := r.StreamJobProgressExported(ctx, job, startedAt, stdoutBdr, stderrBdr, stdoutErrChan, stderrErrChan)
				assert.ErrorIs(t, err, assert.AnError)
			})
	})
}
