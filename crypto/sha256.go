package crypto

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
)

func RandomSha256String() string {
	return hex.EncodeToString(RandomSha256())
}

func RandomSha256() []byte {
	data := make([]byte, 64)
	for i := range data {
		data[i] = byte(rand.Intn(256))
	}

	h := sha256.New()
	h.Write(data)
	return h.Sum(nil)
}
