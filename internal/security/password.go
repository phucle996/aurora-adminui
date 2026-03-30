package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"aurora-adminui/internal/errorx"

	"golang.org/x/crypto/argon2"
)

type argonParams struct {
	Memory      uint32
	Time        uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

func defaultArgonParams() argonParams {
	return argonParams{
		Memory:      64 * 1024,
		Time:        3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}
}

func HashPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password is empty")
	}
	p := defaultArgonParams()
	salt := make([]byte, p.SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Parallelism, p.KeyLength)
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.Memory,
		p.Time,
		p.Parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key),
	), nil
}

func ComparePassword(storedHash, password string) error {
	if storedHash == "" || password == "" {
		return errorx.ErrTokenInvalid
	}
	parts := strings.Split(storedHash, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" {
		return errorx.ErrTokenInvalid
	}

	var p argonParams
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Time, &p.Parallelism); err != nil {
		return errorx.ErrTokenInvalid
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil || len(salt) == 0 {
		return errorx.ErrTokenInvalid
	}
	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil || len(hash) == 0 {
		return errorx.ErrTokenInvalid
	}
	other := argon2.IDKey([]byte(password), salt, p.Time, p.Memory, p.Parallelism, uint32(len(hash)))
	if subtle.ConstantTimeCompare(hash, other) != 1 {
		return errorx.ErrTokenInvalid
	}
	return nil
}
