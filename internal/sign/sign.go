package sign

import (
	"crypto/hmac"
	"crypto/sha256"
)

type signFunc func(data []byte) []byte

func GetSignFunc(keyStr string) signFunc {
	key := []byte(keyStr)

	return func(data []byte) []byte {
		h := hmac.New(sha256.New, key)
		h.Write(data)
		return h.Sum(nil)
	}
}
