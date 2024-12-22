package handler

import (
	"crypto/rand"
	"encoding/base64"
)

func generateRandomString(bytes int) string {
	b := make([]byte, bytes)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
