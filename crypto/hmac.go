package crypto

import (
	"crypto/hmac"
	"hash"
)

// GetHMAC returns a keyed-hash message authentication code using the desired hashtype
// sha1.New, sha256.New, sha512.New, sha512.New384, md5.New
func GetHMAC(h func() hash.Hash, key, input []byte) []byte {
	hm := hmac.New(h, key)
	hm.Write(input)
	return hm.Sum(nil)
}
