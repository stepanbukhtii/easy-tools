package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"hash"
)

// Const's declaration for common.go operations
const (
	HashSHA1 = iota
	HashSHA256
	HashSHA512
	HashSHA512384
	HashMD5
)

// GetHMAC returns a keyed-hash message authentication code using the desired hashtype
func GetHMAC(hashType int, input, key []byte) []byte {
	var h func() hash.Hash

	switch hashType {
	case HashSHA1:
		h = sha1.New
	case HashSHA256:
		h = sha256.New
	case HashSHA512:
		h = sha512.New
	case HashSHA512384:
		h = sha512.New384
	case HashMD5:
		h = md5.New
	}

	hm := hmac.New(h, key)
	hm.Write(input)
	return hm.Sum(nil)
}
