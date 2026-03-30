package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
)

type SymmetricCipher struct {
	aead cipher.AEAD
}

func NewSymmetricCipher(rawKey string) (*SymmetricCipher, error) {
	key, err := normalizeSymmetricKey(rawKey)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &SymmetricCipher{aead: aead}, nil
}

func (c *SymmetricCipher) Encrypt(plaintext []byte) (string, error) {
	if len(plaintext) == 0 {
		return "", fmt.Errorf("plaintext is empty")
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	ciphertext := c.aead.Seal(nil, nonce, plaintext, nil)
	payload := append(nonce, ciphertext...)
	return base64.RawStdEncoding.EncodeToString(payload), nil
}

func (c *SymmetricCipher) Decrypt(encoded string) ([]byte, error) {
	encoded = strings.TrimSpace(encoded)
	if encoded == "" {
		return nil, fmt.Errorf("ciphertext is empty")
	}
	payload, err := base64.RawStdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	nonceSize := c.aead.NonceSize()
	if len(payload) <= nonceSize {
		return nil, fmt.Errorf("ciphertext payload is invalid")
	}
	nonce := payload[:nonceSize]
	ciphertext := payload[nonceSize:]
	return c.aead.Open(nil, nonce, ciphertext, nil)
}

func normalizeSymmetricKey(raw string) ([]byte, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, fmt.Errorf("symmetric key is empty")
	}
	if decoded, err := base64.RawStdEncoding.DecodeString(trimmed); err == nil && len(decoded) == 32 {
		return decoded, nil
	}
	if decoded, err := base64.StdEncoding.DecodeString(trimmed); err == nil && len(decoded) == 32 {
		return decoded, nil
	}
	if len(trimmed) == 32 {
		return []byte(trimmed), nil
	}
	return nil, fmt.Errorf("symmetric key must be 32 bytes or base64-encoded 32 bytes")
}
