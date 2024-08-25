package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

type PasswordHasher interface {
	Hash(password string) string
}

type SHA256Hasher struct {
}

func NewSHA256Hasher() *SHA256Hasher {
	return &SHA256Hasher{}
}

func (h *SHA256Hasher) Hash(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}
