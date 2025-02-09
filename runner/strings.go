package runner

import (
	"strings"
	"sync"
)

type SyncStringBuilder struct {
	b  strings.Builder
	mu sync.Mutex
}

func (b *SyncStringBuilder) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Write(p)
}

func (b *SyncStringBuilder) WriteString(s string) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteString(s)
}

func (b *SyncStringBuilder) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}
