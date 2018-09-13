package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

type HMAC struct {
	hmac hash.Hash
}

func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hmac: h,
	}
}

// Hash will hash the provided input string using HMAC with
// the secret key provided when the HMAC object was created
func (h HMAC) Hash(input string) string {
	h.hmac.Reset()
	h.hmac.Write([]byte(input))
	hashedData := h.hmac.Sum(nil)
	return base64.URLEncoding.EncodeToString(hashedData)
}
