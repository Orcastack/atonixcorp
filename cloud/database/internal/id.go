package internal

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateID(prefix string) string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(b))
}
