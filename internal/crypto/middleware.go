package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"io"
	"net/http"
)

type middleware func(next http.Handler) http.Handler

// 1. AES-ключ (зашифрованный RSA)
// 2. Расшифровываем AES-ключ через PKCS1v15
// 3. Инициализация AES-GCM
// 4. Расшифровка тела
func GetDecryptMiddleware(privateKey *rsa.PrivateKey) middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet {
				next.ServeHTTP(w, r)
			}
			defer r.Body.Close()

			payload, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			rsaBlockSize := privateKey.Size()

			if len(payload) < rsaBlockSize {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			encryptedKey := payload[:rsaBlockSize]
			rest := payload[rsaBlockSize:]

			aesKey, err := rsa.DecryptPKCS1v15(nil, privateKey, encryptedKey)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			block, err := aes.NewCipher(aesKey)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			gcm, err := cipher.NewGCM(block)
			if err != nil {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			nonceSize := gcm.NonceSize()

			if len(rest) < nonceSize {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			nonce := rest[:nonceSize]
			ciphertext := rest[nonceSize:]

			plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(plaintext))

			next.ServeHTTP(w, r)
		})
	}
}
