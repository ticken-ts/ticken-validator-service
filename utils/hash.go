package utils

import (
	"crypto/sha256"
)

func HashSHA256(msg string) []byte {
	h := sha256.New()
	h.Write([]byte(msg))
	return h.Sum(nil)
}
