package grpc

import (
	"context"
	"fmt"

	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	"github.com/traPtitech/piscon-portal-v2/runner/domain"
	"github.com/traPtitech/piscon-portal-v2/runner/portal"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ProgressStreamClient struct {
	stdoutOffset, stderrOffset int
	cl                         grpc.ClientStreamingClient[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse]
}

func NewProgressStreamClient(cl grpc.ClientStreamingClient[portalv1.SendBenchmarkProgressRequest, portalv1.SendBenchmarkProgressResponse]) *ProgressStreamClient {
	return &ProgressStreamClient{cl: cl}
}

var _ portal.ProgressStreamClient = (*ProgressStreamClient)(nil)

// 圧縮が効かないログでも、各メッセージをgRPCの受信上限より十分小さくする。
const logChunkSize = 64 * 1024

func (c *ProgressStreamClient) SendProgress(ctx context.Context, progress *domain.Progress) error {
	stdout, stderr := progress.GetStdout(), progress.GetStderr()
	if len(stdout) < c.stdoutOffset || len(stderr) < c.stderrOffset {
		return fmt.Errorf("benchmark logs must be append-only")
	}
	for {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("send progress: %w", err)
		}
		outN := min(len(stdout)-c.stdoutOffset, logChunkSize)
		errN := min(len(stderr)-c.stderrOffset, logChunkSize)
		req := &portalv1.SendBenchmarkProgressRequest{
			BenchmarkId:   progress.GetBenchmarkID(),
			StartedAt:     timestamppb.New(progress.GetStartedAt()),
			StdoutDelta:   []byte(stdout[c.stdoutOffset : c.stdoutOffset+outN]),
			StderrDelta:   []byte(stderr[c.stderrOffset : c.stderrOffset+errN]),
			Score:         int64(progress.GetScore()),
			LogsAreDelta:  true,
			MoreLogChunks: c.stdoutOffset+outN < len(stdout) || c.stderrOffset+errN < len(stderr),
		}
		if err := c.cl.Send(req); err != nil {
			return fmt.Errorf("send progress: %w", err)
		}
		c.stdoutOffset += outN
		c.stderrOffset += errN
		if c.stdoutOffset == len(stdout) && c.stderrOffset == len(stderr) {
			return nil
		}
	}
}

func (c *ProgressStreamClient) Close() error {
	_, err := c.cl.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close progress stream: %w", err)
	}

	return nil
}
