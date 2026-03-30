package entity

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type MFAMethodType string

const (
	MFAMethodTOTP MFAMethodType = "totp"
)

type TokenType string

const (
	TokenTypeBootstrap TokenType = "bootstrap"
	TokenTypeAPIToken  TokenType = "apitoken"
)

type APITokenState struct {
	SingletonID       bool
	CurrentVersion    int
	CurrentTokenHash  string
	PreviousVersion   *int
	PreviousTokenHash *string
	CreatedAt         time.Time
	LastRotatedAt     time.Time
}

type AdminSecurityState struct {
	SingletonID      bool
	TwoFactorEnabled bool
	TOTPSecret       *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	TOTPEnabledAt    *time.Time
}

type AdminSessionStatus string

const (
	AdminSessionStatusActive  AdminSessionStatus = "active"
	AdminSessionStatusRevoked AdminSessionStatus = "revoked"
)

type AdminSession struct {
	ID         uuid.UUID
	Status     AdminSessionStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ExpiresAt  time.Time
	LastSeenAt *time.Time
	RevokedAt  *time.Time
}

type AdminPreauthSession struct {
	ID        string
	Token     string
	TokenType TokenType
	CreatedAt time.Time
	ExpiresAt time.Time
}

type AdminTOTPSetupSession struct {
	ID        string
	Secret    string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type LoginResult struct {
	TokenType          TokenType
	TokenVersion       int
	PlaintextToken     string
	BootstrapExchanged bool
	MFARequired        bool
	MFAMethods         []MFAMethodType
	PreauthSession     string
	PreauthTTLSeconds  int
	SessionID          string
}

func (r LoginResult) PlaintextTokenOrFallback(fallback string) string {
	if strings.TrimSpace(r.PlaintextToken) != "" {
		return r.PlaintextToken
	}
	return strings.TrimSpace(fallback)
}

type TOTPSetupBeginResult struct {
	SetupSession   string
	Secret         string
	OTPAuthURL     string
	TTLSeconds     int
	AlreadyEnabled bool
}

type TokenRotationResult struct {
	Version      int
	TelegramSent bool
}
