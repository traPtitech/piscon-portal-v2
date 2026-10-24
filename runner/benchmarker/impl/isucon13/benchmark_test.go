package isucon13

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/runner/config"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

func TestStartWait(t *testing.T) {
	for _, tt := range []struct {
		name   string
		json   string
		exit   int
		result domain.Result
		score  int
		err    bool
	}{
		{"passed", `{"pass":true,"score":1234}`, 0, domain.ResultPassed, 1234, false},
		{"initialize failed", `{"pass":false,"score":0}`, 0, domain.ResultFailed, 0, false},
		{"finalcheck failed", `{"pass":false,"score":0}`, 1, domain.ResultFailed, 0, false},
		{"crashed", `{"pass":false,"score":0}`, 2, domain.ResultError, 0, true},
		{"passed but abnormal exit", `{"pass":true,"score":1234}`, 1, domain.ResultError, 0, true},
		{"missing result", "", 0, domain.ResultError, 0, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			path := filepath.Join(dir, "bench")
			// 実際のCLIと同様、成功時は結果をstdoutへ出さずファイルに書く。
			script := fmt.Sprintf(`#!/bin/sh
while [ "$#" -gt 0 ]; do
  case "$1" in
    --result-path) shift; result_path="$1" ;;
  esac
  shift
done
if [ -n '%s' ]; then
  printf '%%s' '%s' > "$result_path"
fi
printf '2026-10-24T00:00:00Z\tinfo\tstdout progress\n'
printf '2026-10-24T00:00:00Z\twarn\tstdout warning\n'
printf '2026-10-24T00:00:00Z\tinfo\tstderr progress\n' >&2
printf '2026-10-24T00:00:00Z\terror\tstderr error\n' >&2
exit %d
`, tt.json, tt.json, tt.exit)
			require.NoError(t, os.WriteFile(path, []byte(script), 0o755))
			b, err := New(config.Problem{Options: map[string]any{
				"benchmarker-path": path, "benchmarker-dir": dir,
			}})
			require.NoError(t, err)
			ctx := context.Background()
			job := domain.NewJob("job-id", "192.0.2.1")
			var previousOutputDir string
			for range 2 {
				out, _, err := b.Start(ctx, job)
				require.NoError(t, err)
				assert.Zero(t, b.score)
				outputDir := b.outputDir
				assert.NotEqual(t, previousOutputDir, outputDir)
				previousOutputDir = outputDir
				stdout, err := io.ReadAll(out.Stdout)
				require.NoError(t, err)
				assert.Equal(t, "2026-10-24T00:00:00Z\twarn\tstdout warning\n", string(stdout))
				stderr, err := io.ReadAll(out.Stderr)
				require.NoError(t, err)
				assert.Equal(t, "2026-10-24T00:00:00Z\terror\tstderr error\n", string(stderr))
				// Runnerはstdout/stderrがEOFになった後、Waitより先に最終スコアを送る。
				if !tt.err {
					score, err := b.CalculateScore(ctx, "", "")
					require.NoError(t, err)
					assert.Equal(t, tt.score, score)
				}
				result, _, err := b.Wait(ctx)
				if tt.err {
					require.Error(t, err)
				} else {
					require.NoError(t, err)
					score, err := b.CalculateScore(ctx, "", "")
					require.NoError(t, err)
					assert.Equal(t, tt.score, score)
				}
				assert.Equal(t, tt.result, result)
				_, err = os.Stat(outputDir)
				require.True(t, os.IsNotExist(err))
			}
		})
	}
}

func TestInfoFilter(t *testing.T) {
	const timestamp = "2026-10-24T00:00:00Z\t"
	longWarning := timestamp + "warn\t" + strings.Repeat("x", 100_000) + "\n"
	for _, tt := range []struct {
		name, input, want string
	}{
		{"empty", "", ""},
		{"info only", timestamp + "info\tprogress\n" + timestamp + "info\tlast", ""},
		{"mixed", timestamp + "info\tprogress\n" + timestamp + "warn\twarning\n" + timestamp + "error\terror\n", timestamp + "warn\twarning\n" + timestamp + "error\terror\n"},
		{"other output", "panic: failed\n{\"pass\":false,\"score\":0}\n", "panic: failed\n{\"pass\":false,\"score\":0}\n"},
		{"info in message", timestamp + "warn\tmessage\tinfo\tvalue\n", timestamp + "warn\tmessage\tinfo\tvalue\n"},
		{"debug", timestamp + "debug\tdetails\n", timestamp + "debug\tdetails\n"},
		{"long line", longWarning, longWarning},
		{"no final newline", timestamp + "error\tfailed", timestamp + "error\tfailed"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			filter := &infoFilter{reader: bufio.NewReader(iotest.OneByteReader(strings.NewReader(tt.input)))}
			require.NoError(t, iotest.TestReader(filter, []byte(tt.want)))
		})
	}
	filter := &infoFilter{reader: bufio.NewReader(iotest.ErrReader(io.ErrUnexpectedEOF))}
	_, err := io.ReadAll(filter)
	require.ErrorIs(t, err, io.ErrUnexpectedEOF)
}
