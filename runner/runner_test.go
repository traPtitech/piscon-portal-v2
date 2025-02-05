package runner_test

import (
	"bytes"
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/runner"
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
