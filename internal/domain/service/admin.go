package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type AdminService interface {
	BootstrapIfNeeded(ctx context.Context) error
	LoginInit(ctx context.Context, rawToken string) (*entity.LoginResult, error)
	VerifySecondFactor(ctx context.Context, preauthID, code string) (*entity.LoginResult, error)
	ValidateSession(ctx context.Context, rawSessionID string) (*entity.AdminSession, error)
	Logout(ctx context.Context, rawSessionID string) error
	GetTwoFactorStatus(ctx context.Context) (*entity.AdminSecurityState, error)
	BeginTOTPSetup(ctx context.Context) (*entity.TOTPSetupBeginResult, error)
	ConfirmTOTPSetup(ctx context.Context, setupSessionID, code string) error
	DisableTwoFactor(ctx context.Context, code string) error
	RotateAPIToken(ctx context.Context, code string) (*entity.TokenRotationResult, error)
}
