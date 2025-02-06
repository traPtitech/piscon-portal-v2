package runner

import (
	"strings"
	"sync"
)

type Builder struct {
	b  strings.Builder
	mu sync.Mutex
}

func (b *Builder) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.Write(p)
}

func (b *Builder) WriteString(s string) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.WriteString(s)
}

func (b *Builder) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.b.String()
}
