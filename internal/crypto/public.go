package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func GetPublicKey(path string) (*rsa.PublicKey, error) {
	publicKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ReadFile: %w", err)
	}
	publicKeyPemBlock, _ := pem.Decode(publicKeyBytes)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyPemBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("ParsePKCS1PublicKey: %w", err)
	}
	return publicKey, nil
}
