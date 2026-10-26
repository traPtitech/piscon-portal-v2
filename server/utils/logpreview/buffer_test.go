package logpreview_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
	"github.com/traPtitech/piscon-portal-v2/server/utils/logpreview"
)

func TestPreviewPreservesShortLogsAndBothEnds(t *testing.T) {
	const edge = 16 * 1024
	for _, s := range []string{"", "small", strings.Repeat("a", 2*edge), strings.Repeat("日", 30000)} {
		for _, chunk := range []int{7, 1024, len(s) + 1} {
			var b logpreview.Buffer
			for i := 0; i < len(s); i += chunk {
				b.WriteString(s[i:min(i+chunk, len(s))])
			}
			got := b.String()
			require.True(t, utf8.ValidString(got))
			if len(s) <= 2*edge {
				require.Equal(t, s, got)
				continue
			}
			require.LessOrEqual(t, len(got), 2*edge+len(logpreview.Omitted))
			require.True(t, strings.HasPrefix(got, strings.Repeat("日", edge/3)))
			require.True(t, strings.HasSuffix(got, strings.Repeat("日", edge/3)))
		}
	}
}

func TestPreviewRemovesInvalidUTF8(t *testing.T) {
	for _, tc := range []struct {
		name, input, want string
	}{
		{"短い不正UTF8", "start\xffend", "startend"},
		{"上限以内の不正UTF8", strings.Repeat("\xffa", 12000),
			strings.Repeat("a", 12000)},
		{"長い不正UTF8", strings.Repeat("\xffa", 20000),
			strings.Repeat("a", 8192) + logpreview.Omitted + strings.Repeat("a", 8192)},
		{"上限が文字の途中にある", strings.Repeat("日", 12000),
			strings.Repeat("日", 5461) + logpreview.Omitted + strings.Repeat("日", 5461)},
		{"不正UTF8と日本語の混在", strings.Repeat("\xff日", 10000),
			strings.Repeat("日", 4096) + logpreview.Omitted + strings.Repeat("日", 4096)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, chunk := range []int{7, 1024, len(tc.input) + 1} {
				var b logpreview.Buffer
				for i := 0; i < len(tc.input); i += chunk {
					b.WriteString(tc.input[i:min(i+chunk, len(tc.input))])
				}
				got := b.String()
				require.True(t, utf8.ValidString(got))
				require.LessOrEqual(t, len(got), 32*1024+len(logpreview.Omitted))
				require.Equal(t, tc.want, got)
			}
		})
	}
}
