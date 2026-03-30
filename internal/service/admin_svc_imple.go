package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"aurora-adminui/internal/config"
	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	"aurora-adminui/internal/security"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const (
	bootstrapTTL      = 15 * time.Minute
	apiTokenRedisTTL  = 24 * time.Hour
	adminSessionTTL   = 24 * time.Hour
	adminPreauthTTL   = 150 * time.Second
	adminTOTPSetupTTL = 10 * time.Minute
)

type AdminSvcImple struct {
	repo   domainrepo.AdminRepository
	cache  domainrepo.AdminCache
	config *config.Config
}

func NewAdminService(
	repo domainrepo.AdminRepository,
	cache domainrepo.AdminCache,
	cfg *config.Config,
) domainsvc.AdminService {
	return &AdminSvcImple{repo: repo, cache: cache, config: cfg}
}

func (s *AdminSvcImple) BootstrapIfNeeded(ctx context.Context) error {
	if _, err := s.repo.GetTokenState(ctx); err == nil {
		if existing, err := s.cache.GetBootstrapPointer(ctx); err == nil && existing != "" {
			_ = s.cache.DeleteToken(ctx, existing)
		}
		_ = s.cache.ClearBootstrapPointer(ctx)
		return nil
	} else if !errors.Is(err, errorx.ErrAPITokenNotFound) {
		return err
	}

	bootstrapToken, err := s.cache.GetBootstrapPointer(ctx)
	if err != nil {
		return err
	}
	if bootstrapToken == "" {
		bootstrapToken, err = security.GenerateToken(64)
		if err != nil {
			return err
		}
		if err := s.cache.SetToken(ctx, bootstrapToken, entity.TokenTypeBootstrap, bootstrapTTL); err != nil {
			return err
		}
		if err := s.cache.SetBootstrapPointer(ctx, bootstrapToken, bootstrapTTL); err != nil {
			return err
		}
	}

	message := fmt.Sprintf(
		"Aurora bootstrap token\nToken: %s\nExpires in: %s\nUse it once to mint the admin API token.",
		bootstrapToken,
		bootstrapTTL.String(),
	)
	if err := sendBootstrapTokenToTelegram(ctx, s.config, message); err != nil {
		log.Printf("[admin-bootstrap] telegram bootstrap delivery failed, fallback to log: %s", err.Error())
		log.Printf("[admin-bootstrap] bootstrap token: %s", bootstrapToken)
	}
	return nil
}

func (s *AdminSvcImple) LoginInit(ctx context.Context, rawToken string) (*entity.LoginResult, error) {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return nil, errorx.ErrInvalidArgument
	}

	tokenType, err := s.validateFactorOne(ctx, rawToken)
	if err != nil {
		return nil, err
	}

	securityState, err := s.repo.GetSecurityState(ctx)
	if err != nil {
		return nil, err
	}
	if !securityState.TwoFactorEnabled {
		return s.finalizeLogin(ctx, rawToken, tokenType)
	}

	preauthID := uuid.NewString()
	now := time.Now().UTC()
	preauth := &entity.AdminPreauthSession{
		ID:        preauthID,
		Token:     rawToken,
		TokenType: tokenType,
		CreatedAt: now,
		ExpiresAt: now.Add(adminPreauthTTL),
	}
	if err := s.cache.SetPreauth(ctx, preauth, adminPreauthTTL); err != nil {
		return nil, err
	}
	return &entity.LoginResult{
		MFARequired:       true,
		MFAMethods:        []entity.MFAMethodType{entity.MFAMethodTOTP},
		PreauthSession:    preauthID,
		PreauthTTLSeconds: int(adminPreauthTTL.Seconds()),
	}, nil
}

func (s *AdminSvcImple) VerifySecondFactor(ctx context.Context, preauthID, code string) (*entity.LoginResult, error) {
	preauthID = strings.TrimSpace(preauthID)
	code = strings.TrimSpace(code)
	if preauthID == "" || code == "" {
		return nil, errorx.ErrInvalidArgument
	}
	preauth, err := s.cache.GetPreauth(ctx, preauthID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, errorx.ErrTokenInvalid
		}
		return nil, err
	}
	if _, err := s.validateFactorOne(ctx, preauth.Token); err != nil {
		return nil, err
	}
	if err := s.verifySecondFactorCode(ctx, code); err != nil {
		return nil, err
	}
	_ = s.cache.DeletePreauth(ctx, preauthID)
	return s.finalizeLogin(ctx, preauth.Token, preauth.TokenType)
}

func (s *AdminSvcImple) ValidateSession(ctx context.Context, rawSessionID string) (*entity.AdminSession, error) {
	rawSessionID = strings.TrimSpace(rawSessionID)
	if rawSessionID == "" {
		return nil, errorx.ErrTokenInvalid
	}

	sessionID, err := uuid.Parse(rawSessionID)
	if err != nil {
		return nil, errorx.ErrTokenInvalid
	}

	if ok, err := s.cache.HasSession(ctx, rawSessionID); err == nil && ok {
		session, err := s.repo.GetActiveSession(ctx, sessionID)
		if err != nil {
			return nil, err
		}
		if session.ExpiresAt.Before(time.Now().UTC()) {
			_ = s.cache.DeleteSession(ctx, rawSessionID)
			_ = s.repo.RevokeSession(ctx, sessionID, time.Now().UTC())
			return nil, errorx.ErrTokenExpired
		}
		_ = s.repo.TouchSession(ctx, sessionID, time.Now().UTC())
		return session, nil
	}

	session, err := s.repo.GetActiveSession(ctx, sessionID)
	if err != nil {
		if errors.Is(err, errorx.ErrAdminSessionNotFound) {
			return nil, errorx.ErrTokenInvalid
		}
		return nil, err
	}
	if session.ExpiresAt.Before(time.Now().UTC()) {
		_ = s.repo.RevokeSession(ctx, sessionID, time.Now().UTC())
		return nil, errorx.ErrTokenExpired
	}
	if err := s.cache.SetSession(ctx, rawSessionID, time.Until(session.ExpiresAt)); err != nil {
		return nil, err
	}
	_ = s.repo.TouchSession(ctx, sessionID, time.Now().UTC())
	return session, nil
}

func (s *AdminSvcImple) Logout(ctx context.Context, rawSessionID string) error {
	rawSessionID = strings.TrimSpace(rawSessionID)
	if rawSessionID == "" {
		return nil
	}
	sessionID, err := uuid.Parse(rawSessionID)
	if err != nil {
		return nil
	}
	_ = s.cache.DeleteSession(ctx, rawSessionID)
	return s.repo.RevokeSession(ctx, sessionID, time.Now().UTC())
}

func (s *AdminSvcImple) GetTwoFactorStatus(ctx context.Context) (*entity.AdminSecurityState, error) {
	return s.repo.GetSecurityState(ctx)
}

func (s *AdminSvcImple) BeginTOTPSetup(ctx context.Context) (*entity.TOTPSetupBeginResult, error) {
	state, err := s.repo.GetSecurityState(ctx)
	if err != nil {
		return nil, err
	}
	if state.TwoFactorEnabled {
		return &entity.TOTPSetupBeginResult{AlreadyEnabled: true}, nil
	}
	secret, err := security.GenerateTOTPSecret()
	if err != nil {
		return nil, err
	}
	setupID := uuid.NewString()
	now := time.Now().UTC()
	setup := &entity.AdminTOTPSetupSession{
		ID:        setupID,
		Secret:    secret,
		CreatedAt: now,
		ExpiresAt: now.Add(adminTOTPSetupTTL),
	}
	if err := s.cache.SetTOTPSetup(ctx, setup, adminTOTPSetupTTL); err != nil {
		return nil, err
	}
	return &entity.TOTPSetupBeginResult{
		SetupSession: setupID,
		Secret:       secret,
		OTPAuthURL:   security.BuildOTPAuthURL("aurora-admin", secret),
		TTLSeconds:   int(adminTOTPSetupTTL.Seconds()),
	}, nil
}

func (s *AdminSvcImple) ConfirmTOTPSetup(ctx context.Context, setupSessionID, code string) error {
	setupSessionID = strings.TrimSpace(setupSessionID)
	code = strings.TrimSpace(code)
	if setupSessionID == "" || code == "" {
		return errorx.ErrInvalidArgument
	}
	state, err := s.repo.GetSecurityState(ctx)
	if err != nil {
		return err
	}
	if state.TwoFactorEnabled {
		return errorx.ErrMFAMethodAlreadyEnabled
	}
	setup, err := s.cache.GetTOTPSetup(ctx, setupSessionID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return errorx.ErrTokenInvalid
		}
		return err
	}
	if !security.ValidateTOTP(setup.Secret, code, 1) {
		return errorx.ErrMFACodeInvalid
	}
	now := time.Now().UTC()
	if err := s.repo.EnableTwoFactor(ctx, setup.Secret, now); err != nil {
		return err
	}
	_ = s.cache.DeleteTOTPSetup(ctx, setupSessionID)
	return nil
}

func (s *AdminSvcImple) DisableTwoFactor(ctx context.Context, code string) error {
	if err := s.verifySecondFactorCode(ctx, code); err != nil {
		return err
	}
	return s.repo.DisableTwoFactor(ctx, time.Now().UTC())
}

func (s *AdminSvcImple) validateFactorOne(ctx context.Context, rawToken string) (entity.TokenType, error) {
	tokenType, found, err := s.cache.GetTokenType(ctx, rawToken)
	if err != nil {
		return "", err
	}
	if found && tokenType == entity.TokenTypeBootstrap {
		return entity.TokenTypeBootstrap, nil
	}

	tokenState, err := s.repo.GetTokenState(ctx)
	if errors.Is(err, errorx.ErrAPITokenNotFound) {
		return "", errorx.ErrTokenInvalid
	}
	if err != nil {
		return "", err
	}

	version, err := matchAPITokenVersion(tokenState, rawToken)
	if err != nil {
		if found && tokenType == entity.TokenTypeAPIToken {
			_ = s.cache.DeleteToken(ctx, rawToken)
		}
		return "", err
	}
	if err := s.cache.SetToken(ctx, rawToken, entity.TokenTypeAPIToken, apiTokenRedisTTL); err != nil {
		return "", err
	}
	_ = version
	return entity.TokenTypeAPIToken, nil
}

func (s *AdminSvcImple) finalizeLogin(ctx context.Context, rawToken string, tokenType entity.TokenType) (*entity.LoginResult, error) {
	result := &entity.LoginResult{TokenType: tokenType}
	switch tokenType {
	case entity.TokenTypeAPIToken:
		tokenState, err := s.repo.GetTokenState(ctx)
		if err != nil {
			return nil, err
		}
		version, err := matchAPITokenVersion(tokenState, rawToken)
		if err != nil {
			return nil, err
		}
		result.TokenVersion = version
	case entity.TokenTypeBootstrap:
		newToken, err := s.exchangeBootstrapToken(ctx, rawToken)
		if err != nil {
			return nil, err
		}
		result.TokenType = entity.TokenTypeAPIToken
		result.PlaintextToken = newToken
		result.BootstrapExchanged = true
		tokenState, err := s.repo.GetTokenState(ctx)
		if err != nil {
			return nil, err
		}
		result.TokenVersion = tokenState.CurrentVersion
	default:
		return nil, errorx.ErrTokenInvalid
	}

	sessionID, err := s.issueSession(ctx)
	if err != nil {
		return nil, err
	}
	result.SessionID = sessionID
	return result, nil
}

func (s *AdminSvcImple) exchangeBootstrapToken(ctx context.Context, rawToken string) (string, error) {
	tokenType, found, err := s.cache.GetTokenType(ctx, rawToken)
	if err != nil {
		return "", err
	}
	if !found || tokenType != entity.TokenTypeBootstrap {
		return "", errorx.ErrTokenInvalid
	}
	_ = s.cache.DeleteToken(ctx, rawToken)
	_ = s.cache.ClearBootstrapPointer(ctx)

	newRawToken, err := security.GenerateToken(64)
	if err != nil {
		return "", err
	}
	tokenHash, err := security.HashPassword(newRawToken)
	if err != nil {
		return "", err
	}
	now := time.Now().UTC()
	if err := s.repo.SaveTokenState(ctx, &entity.APITokenState{
		SingletonID:      true,
		CurrentVersion:   1,
		CurrentTokenHash: tokenHash,
		CreatedAt:        now,
		LastRotatedAt:    now,
	}); err != nil {
		return "", err
	}
	if err := s.cache.SetToken(ctx, newRawToken, entity.TokenTypeAPIToken, apiTokenRedisTTL); err != nil {
		return "", err
	}
	return newRawToken, nil
}

func (s *AdminSvcImple) RotateAPIToken(ctx context.Context, code string) (*entity.TokenRotationResult, error) {
	state, err := s.repo.GetTokenState(ctx)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, errorx.ErrAPITokenNotFound
	}

	securityState, err := s.repo.GetSecurityState(ctx)
	if err != nil {
		return nil, err
	}
	if securityState.TwoFactorEnabled {
		if err := s.verifySecondFactorCode(ctx, code); err != nil {
			return nil, err
		}
	}

	newRawToken, err := security.GenerateToken(64)
	if err != nil {
		return nil, err
	}
	newHash, err := security.HashPassword(newRawToken)
	if err != nil {
		return nil, err
	}

	nextVersion := nextAPITokenVersion(state.CurrentVersion)
	now := time.Now().UTC()
	previousVersion := state.CurrentVersion
	previousHash := state.CurrentTokenHash

	nextState := &entity.APITokenState{
		SingletonID:       true,
		CurrentVersion:    nextVersion,
		CurrentTokenHash:  newHash,
		PreviousVersion:   &previousVersion,
		PreviousTokenHash: &previousHash,
		CreatedAt:         state.CreatedAt,
		LastRotatedAt:     now,
	}
	if err := s.repo.SaveTokenState(ctx, nextState); err != nil {
		return nil, err
	}
	if err := s.cache.SetToken(ctx, newRawToken, entity.TokenTypeAPIToken, apiTokenRedisTTL); err != nil {
		return nil, err
	}

	message := fmt.Sprintf(
		"Aurora admin API token rotated\nVersion: v%d\nToken: %s\nPrevious version: v%d",
		nextVersion,
		newRawToken,
		previousVersion,
	)
	telegramSent := true
	if err := sendBootstrapTokenToTelegram(ctx, s.config, message); err != nil {
		telegramSent = false
		log.Printf("[admin-rotate] telegram delivery failed, fallback to log: %s", err.Error())
		log.Printf("[admin-rotate] api token v%d: %s", nextVersion, newRawToken)
	}

	return &entity.TokenRotationResult{
		Version:      nextVersion,
		TelegramSent: telegramSent,
	}, nil
}

func (s *AdminSvcImple) issueSession(ctx context.Context) (string, error) {
	now := time.Now().UTC()
	session := &entity.AdminSession{
		ID:         uuid.New(),
		Status:     entity.AdminSessionStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
		ExpiresAt:  now.Add(adminSessionTTL),
		LastSeenAt: &now,
	}
	if err := s.repo.CreateSession(ctx, session); err != nil {
		return "", err
	}
	if err := s.cache.SetSession(ctx, session.ID.String(), adminSessionTTL); err != nil {
		return "", err
	}
	return session.ID.String(), nil
}

func (s *AdminSvcImple) verifySecondFactorCode(ctx context.Context, code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return errorx.ErrInvalidArgument
	}
	state, err := s.repo.GetSecurityState(ctx)
	if err != nil {
		return err
	}
	if !state.TwoFactorEnabled {
		return errorx.ErrMFAMethodNotFound
	}
	if state.TOTPSecret == nil || strings.TrimSpace(*state.TOTPSecret) == "" {
		return errorx.ErrMFAMethodNotFound
	}
	if !security.ValidateTOTP(*state.TOTPSecret, code, 1) {
		return errorx.ErrMFACodeInvalid
	}
	return nil
}

func matchAPITokenVersion(state *entity.APITokenState, rawToken string) (int, error) {
	if state == nil {
		return 0, errorx.ErrTokenInvalid
	}
	if err := security.ComparePassword(state.CurrentTokenHash, rawToken); err == nil {
		return state.CurrentVersion, nil
	}
	if state.PreviousTokenHash != nil && strings.TrimSpace(*state.PreviousTokenHash) != "" {
		if err := security.ComparePassword(*state.PreviousTokenHash, rawToken); err == nil {
			if state.PreviousVersion != nil {
				return *state.PreviousVersion, nil
			}
			return 0, nil
		}
	}
	return 0, errorx.ErrTokenInvalid
}

func nextAPITokenVersion(current int) int {
	if current == 1 {
		return 2
	}
	return 1
}
