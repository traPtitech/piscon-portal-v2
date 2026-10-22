package isucon13

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

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
				_, err = io.Copy(io.Discard, out.Stdout)
				require.NoError(t, err)
				_, err = io.Copy(io.Discard, out.Stderr)
				require.NoError(t, err)
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
