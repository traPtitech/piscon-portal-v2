// Package logpreview は追記されるログの先頭と末尾だけを保持する。
package logpreview

import (
	"strings"
	"unicode/utf8"
)

const EdgeBytes = 16 * 1024
const Omitted = "\n...（ログの中間部分を省略しています）...\n"

type Buffer struct {
	head, tail string
	total      int64
}

func (b *Buffer) WriteString(s string) {
	b.total += int64(len(s))
	n := min(EdgeBytes-len(b.head), len(s))
	if n > 0 {
		if b.head == "" {
			// 初回は切り出し元のログ全体を参照し続けないようコピーする。
			b.head = strings.Clone(s[:n])
		} else {
			b.head += s[:n]
		}
	}
	s = s[n:]
	if len(s) >= EdgeBytes {
		b.tail = strings.Clone(s[len(s)-EdgeBytes:])
	} else {
		b.tail += s
	}
	if len(b.tail) > EdgeBytes {
		b.tail = b.tail[len(b.tail)-EdgeBytes:]
	}
}

func (b *Buffer) String() string {
	if b.total <= 2*EdgeBytes {
		s := strings.ToValidUTF8(b.head+b.tail, "�")
		if len(s) <= 2*EdgeBytes {
			return s
		}
		return truncateEdge(s, false) + Omitted + truncateEdge(s, true)
	}
	return truncateEdge(strings.ToValidUTF8(b.head, "�"), false) + Omitted +
		truncateEdge(strings.ToValidUTF8(b.tail, "�"), true)
}

// 不正UTF-8の置換で増えた分も含め、文字を分断せず各端部を上限内に収める。
func truncateEdge(s string, fromEnd bool) string {
	if len(s) <= EdgeBytes {
		return s
	}
	if fromEnd {
		i := len(s) - EdgeBytes
		for !utf8.RuneStart(s[i]) {
			i++
		}
		return s[i:]
	}
	i := EdgeBytes
	for !utf8.RuneStart(s[i]) {
		i--
	}
	return s[:i]
}
