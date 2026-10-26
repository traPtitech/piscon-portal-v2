package runner

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// 不正なUTF-8を含むログも、実行後のgzipファイルからバイト単位で復元できることを検証する。
func TestRunArchivesBothStreams(t *testing.T) {
	root := t.TempDir()
	t.Setenv("PISCON_LOG_DIR", root)
	stdout, stderr := strings.Repeat("日本語\n", 10000), "staff log\x00\xff"
	r := &Runner{portal: &recordingPortal{}, benchmarker: outputProcess{stdout: stdout, stderr: stderr}}
	require.NoError(t, r.Run())
	for name, want := range map[string]string{"stdout": stdout, "stderr": stderr} {
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(filepath.Join(root, "test-job", name+".gz"))
			require.NoError(t, err)
			defer f.Close()
			z, err := gzip.NewReader(f)
			require.NoError(t, err)
			got, err := io.ReadAll(z)
			require.NoError(t, err)
			require.NoError(t, z.Close())
			require.Equal(t, want, string(got))
		})
	}
}
