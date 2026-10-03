package isuxbench

import (
	"encoding/binary"
	"io"
	"testing"

	"github.com/isucon/isucon10-portal/proto.go/isuxportal/resources"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func encodeReport(t *testing.T, res *resources.BenchmarkResult) []byte {
	t.Helper()
	data, err := proto.Marshal(res)
	require.NoError(t, err)
	buf := binary.BigEndian.AppendUint16(nil, uint16(len(data)))
	return append(buf, data...)
}

func TestWatchReport(t *testing.T) {
	progress := &resources.BenchmarkResult{Finished: false, Score: 100}

	tests := map[string]struct {
		reports         func(t *testing.T) []byte
		wantPassed      bool
		wantScore       int
		wantErrNoFinal  bool
		wantErrUnexpEOF bool
	}{
		"passed": {
			reports: func(t *testing.T) []byte {
				return append(encodeReport(t, progress),
					encodeReport(t, &resources.BenchmarkResult{Finished: true, Passed: true, Score: 200})...)
			},
			wantPassed: true,
			wantScore:  200,
		},
		"failed": {
			reports: func(t *testing.T) []byte {
				return encodeReport(t, &resources.BenchmarkResult{Finished: true, Passed: false, Score: 0})
			},
			wantPassed: false,
			wantScore:  0,
		},
		"no reports": {
			reports:        func(*testing.T) []byte { return nil },
			wantErrNoFinal: true,
		},
		"no final report": {
			reports:        func(t *testing.T) []byte { return encodeReport(t, progress) },
			wantScore:      100,
			wantErrNoFinal: true,
		},
		"truncated length": {
			reports:         func(*testing.T) []byte { return []byte{0x00} },
			wantErrUnexpEOF: true,
		},
		"missing body": {
			reports:         func(*testing.T) []byte { return []byte{0x00, 0x05} },
			wantErrUnexpEOF: true,
		},
		"truncated body": {
			reports: func(t *testing.T) []byte {
				b := encodeReport(t, progress)
				return b[:len(b)-1]
			},
			wantErrUnexpEOF: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			r, w := io.Pipe()
			p := &Process{resultCh: make(chan result, 1)}
			go func() {
				_, _ = w.Write(tt.reports(t))
				_ = w.Close()
			}()

			p.watchReport(r)
			res := <-p.resultCh

			assert.Equal(t, tt.wantScore, p.LatestScore())
			switch {
			case tt.wantErrNoFinal:
				assert.ErrorIs(t, res.err, ErrNoFinalReport)
			case tt.wantErrUnexpEOF:
				assert.ErrorIs(t, res.err, io.ErrUnexpectedEOF)
				assert.NotErrorIs(t, res.err, ErrNoFinalReport)
			default:
				require.NoError(t, res.err)
				assert.Equal(t, tt.wantPassed, res.passed)
			}
		})
	}
}
