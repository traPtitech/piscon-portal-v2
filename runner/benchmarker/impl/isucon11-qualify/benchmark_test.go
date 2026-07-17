package isucon11qualify

import (
	"context"
	"io"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
)

func TestStartWaitCanRunSequentially(t *testing.T) {
	truePath, err := exec.LookPath("true")
	require.NoError(t, err)

	b := &Isucon11Qualify{
		conf: problemConf{
			execPath:       truePath,
			benchmarkerIP:  "127.0.0.1",
			benchmarkerDir: "/tmp",
		},
	}
	ctx := context.Background()
	job := domain.NewJob("job-id", "192.0.2.1")

	for range 2 {
		out, _, err := b.Start(ctx, job)
		require.NoError(t, err)

		_, err = io.ReadAll(out.Stdout)
		require.NoError(t, err)
		_, err = io.ReadAll(out.Stderr)
		require.NoError(t, err)

		result, _, err := b.Wait(ctx)
		require.ErrorContains(t, err, "parse benchmark result")
		assert.Equal(t, domain.ResultError, result)
	}
}
