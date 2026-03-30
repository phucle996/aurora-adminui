package repository

import (
	"context"
	"errors"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepoImple struct {
	db *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) domainrepo.AdminRepository {
	return &AdminRepoImple{db: db}
}

func (r *AdminRepoImple) GetTokenState(ctx context.Context) (*entity.APITokenState, error) {
	var token entity.APITokenState
	err := r.db.QueryRow(
		ctx,
		`SELECT singleton_id, current_version, current_token_hash, previous_version, previous_token_hash, created_at, last_rotated_at
		 FROM admin.api_tokens
		 WHERE singleton_id = TRUE`,
	).Scan(
		&token.SingletonID,
		&token.CurrentVersion,
		&token.CurrentTokenHash,
		&token.PreviousVersion,
		&token.PreviousTokenHash,
		&token.CreatedAt,
		&token.LastRotatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.ErrAPITokenNotFound
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *AdminRepoImple) SaveTokenState(ctx context.Context, token *entity.APITokenState) error {
	if token == nil {
		return errors.New("api token state is nil")
	}
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now().UTC()
	}
	if token.LastRotatedAt.IsZero() {
		token.LastRotatedAt = token.CreatedAt
	}
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO admin.api_tokens
		 (singleton_id, current_version, current_token_hash, previous_version, previous_token_hash, created_at, last_rotated_at)
		 VALUES (TRUE, $1, $2, $3, $4, $5, $6)
		 ON CONFLICT (singleton_id) DO UPDATE
		 SET current_version = EXCLUDED.current_version,
		     current_token_hash = EXCLUDED.current_token_hash,
		     previous_version = EXCLUDED.previous_version,
		     previous_token_hash = EXCLUDED.previous_token_hash,
		     last_rotated_at = EXCLUDED.last_rotated_at`,
		token.CurrentVersion,
		token.CurrentTokenHash,
		token.PreviousVersion,
		token.PreviousTokenHash,
		token.CreatedAt,
		token.LastRotatedAt,
	)
	return err
}

func (r *AdminRepoImple) GetSecurityState(ctx context.Context) (*entity.AdminSecurityState, error) {
	var state entity.AdminSecurityState
	err := r.db.QueryRow(
		ctx,
		`SELECT singleton_id, two_factor_enabled, totp_secret, created_at, updated_at, totp_enabled_at
		 FROM admin.security_state WHERE singleton_id = TRUE`,
	).Scan(&state.SingletonID, &state.TwoFactorEnabled, &state.TOTPSecret, &state.CreatedAt, &state.UpdatedAt, &state.TOTPEnabledAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.ErrAdminSecurityNotFound
	}
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *AdminRepoImple) EnableTwoFactor(ctx context.Context, secret string, now time.Time) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE admin.security_state
		 SET two_factor_enabled = TRUE, totp_secret = $1, updated_at = $2, totp_enabled_at = $2
		 WHERE singleton_id = TRUE`,
		secret, now,
	)
	return err
}

func (r *AdminRepoImple) DisableTwoFactor(ctx context.Context, now time.Time) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE admin.security_state
		 SET two_factor_enabled = FALSE, totp_secret = NULL, updated_at = $1, totp_enabled_at = NULL
		 WHERE singleton_id = TRUE`,
		now,
	)
	return err
}

func (r *AdminRepoImple) CreateSession(ctx context.Context, session *entity.AdminSession) error {
	if session == nil {
		return errors.New("admin session is nil")
	}
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO admin.sessions
		 (id, status, created_at, updated_at, expires_at, last_seen_at, revoked_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		session.ID, session.Status, session.CreatedAt, session.UpdatedAt, session.ExpiresAt, session.LastSeenAt, session.RevokedAt,
	)
	return err
}

func (r *AdminRepoImple) GetActiveSession(ctx context.Context, id uuid.UUID) (*entity.AdminSession, error) {
	var session entity.AdminSession
	err := r.db.QueryRow(
		ctx,
		`SELECT id, status, created_at, updated_at, expires_at, last_seen_at, revoked_at
		 FROM admin.sessions WHERE id = $1 AND status = 'active'`,
		id,
	).Scan(&session.ID, &session.Status, &session.CreatedAt, &session.UpdatedAt, &session.ExpiresAt, &session.LastSeenAt, &session.RevokedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.ErrAdminSessionNotFound
	}
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *AdminRepoImple) TouchSession(ctx context.Context, id uuid.UUID, at time.Time) error {
	if id == uuid.Nil {
		return nil
	}
	_, err := r.db.Exec(ctx, `UPDATE admin.sessions SET last_seen_at = $2, updated_at = $2 WHERE id = $1 AND status = 'active'`, id, at)
	return err
}

func (r *AdminRepoImple) RevokeSession(ctx context.Context, id uuid.UUID, at time.Time) error {
	if id == uuid.Nil {
		return nil
	}
	_, err := r.db.Exec(ctx, `UPDATE admin.sessions SET status = 'revoked', revoked_at = $2, updated_at = $2 WHERE id = $1 AND status = 'active'`, id, at)
	return err
}
