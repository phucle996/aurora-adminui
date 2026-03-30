package repository

import (
	"context"
	"time"

	"aurora-adminui/internal/domain/entity"

	"github.com/google/uuid"
)

type AdminRepository interface {
	GetTokenState(ctx context.Context) (*entity.APITokenState, error)
	SaveTokenState(ctx context.Context, token *entity.APITokenState) error
	GetSecurityState(ctx context.Context) (*entity.AdminSecurityState, error)
	EnableTwoFactor(ctx context.Context, secret string, at time.Time) error
	DisableTwoFactor(ctx context.Context, at time.Time) error
	CreateSession(ctx context.Context, session *entity.AdminSession) error
	GetActiveSession(ctx context.Context, id uuid.UUID) (*entity.AdminSession, error)
	TouchSession(ctx context.Context, id uuid.UUID, at time.Time) error
	RevokeSession(ctx context.Context, id uuid.UUID, at time.Time) error
}

type AdminCache interface {
	GetTokenType(ctx context.Context, token string) (entity.TokenType, bool, error)
	SetToken(ctx context.Context, token string, tokenType entity.TokenType, ttl time.Duration) error
	DeleteToken(ctx context.Context, token string) error
	SetBootstrapPointer(ctx context.Context, token string, ttl time.Duration) error
	GetBootstrapPointer(ctx context.Context) (string, error)
	ClearBootstrapPointer(ctx context.Context) error
	SetPreauth(ctx context.Context, preauth *entity.AdminPreauthSession, ttl time.Duration) error
	GetPreauth(ctx context.Context, id string) (*entity.AdminPreauthSession, error)
	DeletePreauth(ctx context.Context, id string) error
	SetSession(ctx context.Context, sessionID string, ttl time.Duration) error
	HasSession(ctx context.Context, sessionID string) (bool, error)
	DeleteSession(ctx context.Context, sessionID string) error
	SetTOTPSetup(ctx context.Context, session *entity.AdminTOTPSetupSession, ttl time.Duration) error
	GetTOTPSetup(ctx context.Context, id string) (*entity.AdminTOTPSetupSession, error)
	DeleteTOTPSetup(ctx context.Context, id string) error
}
