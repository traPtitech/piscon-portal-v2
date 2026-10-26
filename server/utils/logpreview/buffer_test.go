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

func TestPreviewBoundsAfterUTF8Replacement(t *testing.T) {
	for _, tc := range []struct {
		name, input, want string
	}{
		{"短い不正UTF8", "start\xffend", "start�end"},
		{"置換後に初めて上限を超える", strings.Repeat("\xffa", 12000),
			strings.Repeat("�a", 4096) + logpreview.Omitted + strings.Repeat("�a", 4096)},
		{"長い不正UTF8", strings.Repeat("\xffa", 20000),
			strings.Repeat("�a", 4096) + logpreview.Omitted + strings.Repeat("�a", 4096)},
		{"上限が文字の途中にある", strings.Repeat("\xff日", 10000),
			strings.Repeat("�日", 2730) + "�" + logpreview.Omitted + "日" + strings.Repeat("�日", 2730)},
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
