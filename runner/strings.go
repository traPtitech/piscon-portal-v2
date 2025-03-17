package runner

import (
	"fmt"
	"strings"
	"sync"
)

type SyncStringBuilder struct {
	b  strings.Builder
	mu sync.Mutex
}

func (b *SyncStringBuilder) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	n, err := b.b.Write(p)
	if err != nil {
		return 0, fmt.Errorf("write: %w", err)
	}
	return n, nil
}

func (b *SyncStringBuilder) WriteString(s string) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	n, err := b.b.WriteString(s)
	if err != nil {
		return 0, fmt.Errorf("write string: %w", err)
	}
	return n, nil
}

func (b *SyncStringBuilder) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}
