// Sign package provides signing data tools.
package sign

import (
	"crypto/hmac"
	"crypto/sha256"
)

// SignFunc type is a function that signs data.
type SignFunc func(data []byte) []byte

func GetSignFunc(keyStr string) SignFunc {
	key := []byte(keyStr)

	return func(data []byte) []byte {
		h := hmac.New(sha256.New, key)
		h.Write(data)
		return h.Sum(nil)
	}
}
