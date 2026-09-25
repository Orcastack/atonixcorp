package ledger

import (
	"crypto/sha256"
)

// Hasher defines the interface for computing hashes over block payloads.
type Hasher interface {
	Hash(data []byte) []byte
}

// SHA256Hasher is a simple SHA-256 implementation.
type SHA256Hasher struct{}

func (h SHA256Hasher) Hash(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}
