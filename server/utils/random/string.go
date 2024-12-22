package random

import (
	"crypto/rand"
	"encoding/base64"
)

func String(bytes int) string {
	b := make([]byte, bytes)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
