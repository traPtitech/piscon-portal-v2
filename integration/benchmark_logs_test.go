package integration_test

import (
	"context"
	"io"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	portalv1 "github.com/traPtitech/piscon-portal-v2/gen/portal/v1"
	rd "github.com/traPtitech/piscon-portal-v2/runner/domain"
	client "github.com/traPtitech/piscon-portal-v2/runner/portal/grpc"
	"github.com/traPtitech/piscon-portal-v2/server/domain"
	"github.com/traPtitech/piscon-portal-v2/server/server"
	"github.com/traPtitech/piscon-portal-v2/server/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/test/bufconn"
)

// 保存境界で結果を記録し、内部処理の呼び出し順には依存しない。
type storedProgress struct {
	id    uuid.UUID
	log   domain.BenchmarkLog
	score int64
}
type progressStore struct {
	usecase.BenchmarkUseCase
	mu    sync.Mutex
	saved []storedProgress
	err   error
}

func (s *progressStore) SaveBenchmarkProgress(_ context.Context, id uuid.UUID, log domain.BenchmarkLog, score int64, _ time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.saved = append(s.saved, storedProgress{id, log, score})
	return s.err
}

type compressionObserver struct{ headers chan string }

func (s compressionObserver) TagRPC(ctx context.Context, _ *stats.RPCTagInfo) context.Context {
	return ctx
}
func (s compressionObserver) HandleRPC(_ context.Context, event stats.RPCStats) {
	if h, ok := event.(*stats.InHeader); ok {
		s.headers <- h.Compression
	}
}
func (s compressionObserver) TagConn(ctx context.Context, _ *stats.ConnTagInfo) context.Context {
	return ctx
}
func (s compressionObserver) HandleConn(context.Context, stats.ConnStats) {}

func logPortal(t *testing.T, store *progressStore) (*client.Portal, <-chan string) {
	t.Helper()
	listener := bufconn.Listen(1024 * 1024)
	observer := compressionObserver{make(chan string, 10)}
	s := grpc.NewServer(grpc.StatsHandler(observer))
	portalv1.RegisterBenchmarkServiceServer(s, server.NewBenchmarkService(store))
	serveErr := make(chan error, 1)
	go func() { serveErr <- s.Serve(listener) }()
	t.Cleanup(func() {
		s.Stop()
		require.NoError(t, <-serveErr)
	})
	conn, err := grpc.NewClient("passthrough:///bufnet", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { conn.Close() })
	return client.NewPortal(portalv1.NewBenchmarkServiceClient(conn), 0), observer.headers
}

// 実際のgRPC通信を使い、差分の復元・圧縮・大量ログの保存制限を検証する。
func TestLogTransportAndPersistence(t *testing.T) {
	store := &progressStore{}
	p, headers := logPortal(t, store)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := p.MakeProgressStreamClient(ctx)
	require.NoError(t, err)
	id, started := uuid.New(), time.Now()
	head, tail := strings.Repeat("H", 16*1024), strings.Repeat("T", 16*1024)
	large := head + strings.Repeat("middle日", 700000) + tail
	snapshots := []struct{ out, errOut string }{
		{"abc", "staff"}, {"abcdef", "staff more"},
		{"abcdef" + large, "staff more" + large},
		{"abcdef" + large + "END", "staff more" + large + "END"},
	}
	for i, s := range snapshots {
		require.NoError(t, stream.SendProgress(ctx, rd.NewProgress(id.String(), s.out, s.errOut, i, started)))
	}
	require.NoError(t, stream.Close())
	t.Run("gzipで通信する", func(t *testing.T) {
		require.Equal(t, "gzip", <-headers)
	})
	store.mu.Lock()
	defer store.mu.Unlock()
	if !t.Run("分割数によらず進捗ごとに1回保存する", func(t *testing.T) {
		require.Len(t, store.saved, 4)
	}) {
		return
	}
	t.Run("差分を重複欠落なく復元する", func(t *testing.T) {
		require.Equal(t, domain.BenchmarkLog{UserLog: "abc", AdminLog: "staff"}, store.saved[0].log)
		require.Equal(t, domain.BenchmarkLog{UserLog: "abcdef", AdminLog: "staff more"}, store.saved[1].log)
	})
	t.Run("ベンチIDとスコアを維持する", func(t *testing.T) {
		for i, entry := range store.saved {
			require.Equal(t, id, entry.id)
			require.Equal(t, int64(i), entry.score)
		}
	})
	t.Run("大量ログの先頭と最新の末尾だけを保存する", func(t *testing.T) {
		for i, entry := range store.saved[2:] {
			for _, log := range []string{entry.log.UserLog, entry.log.AdminLog} {
				require.LessOrEqual(t, len(log), 32*1024+256)
				require.NotContains(t, log, "middle")
			}
			require.True(t, strings.HasPrefix(entry.log.UserLog, "abcdef"+head[:len(head)-6]))
			require.True(t, strings.HasPrefix(entry.log.AdminLog, "staff more"+head[:len(head)-10]))
			ending := tail
			if i == 1 {
				ending = tail[3:] + "END"
			}
			require.True(t, strings.HasSuffix(entry.log.UserLog, ending))
			require.True(t, strings.HasSuffix(entry.log.AdminLog, ending))
		}
	})
}

func TestLogPersistenceFailureIsReturnedToRunner(t *testing.T) {
	p, _ := logPortal(t, &progressStore{err: io.ErrUnexpectedEOF})
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := p.MakeProgressStreamClient(ctx)
	require.NoError(t, err)
	_ = stream.SendProgress(ctx, rd.NewProgress(uuid.NewString(), "log", "", 0, time.Now()))
	require.Error(t, stream.Close(), "a successful send does not imply durable storage")
}
