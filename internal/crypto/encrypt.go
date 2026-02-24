package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"io"
)

// Сначала хотел просто по блокам шифровать тело, но чатГПТ сказал, что лучше делать так
// 1. Генерируем случайный AES-256 ключ.
// 2. Создаём AES-GCM
// 3. Шифруем AES-ключ через RSA
// 4. Формируем итоговый payload:
// [rsa_key_len(2 bytes)] [rsa_encrypted_key] [nonce] [ciphertext]
func Encrypt(body io.Reader, publicKey *rsa.PublicKey) (io.Reader, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	aesKey := make([]byte, 32)
	if _, err := rand.Read(aesKey); err != nil {
		return nil, fmt.Errorf("generate aes key: %w", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, fmt.Errorf("aes new cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("new gcm: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nil, nonce, data, nil)
	encryptedKey, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, aesKey)
	if err != nil {
		return nil, fmt.Errorf("encrypt aes key: %w", err)
	}
	var buf bytes.Buffer

	// if err := binary.Write(&buf, binary.BigEndian, uint16(len(encryptedKey))); err != nil {
	// 	return nil, fmt.Errorf("write key len: %w", err)
	// }

	buf.Write(encryptedKey)
	buf.Write(nonce)
	buf.Write(ciphertext)

	body = bytes.NewReader(buf.Bytes())
	return body, nil
}
