CREATE SCHEMA IF NOT EXISTS admin;

CREATE TABLE IF NOT EXISTS admin.api_tokens (
  singleton_id BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (singleton_id = TRUE),
  current_version INT NOT NULL CHECK (current_version IN (1, 2)),
  current_token_hash TEXT NOT NULL,
  previous_version INT NULL CHECK (previous_version IS NULL OR previous_version IN (1, 2)),
  previous_token_hash TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  last_rotated_at TIMESTAMPTZ NOT NULL
);

CREATE TABLE IF NOT EXISTS admin.security_state (
  singleton_id BOOLEAN PRIMARY KEY DEFAULT TRUE CHECK (singleton_id = TRUE),
  two_factor_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  totp_secret TEXT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  totp_enabled_at TIMESTAMPTZ NULL
);

INSERT INTO admin.security_state (
  singleton_id,
  two_factor_enabled,
  totp_secret,
  created_at,
  updated_at,
  totp_enabled_at
)
VALUES (TRUE, FALSE, NULL, NOW(), NOW(), NULL)
ON CONFLICT (singleton_id) DO NOTHING;

CREATE TABLE IF NOT EXISTS admin.sessions (
  id UUID PRIMARY KEY,
  status TEXT NOT NULL CHECK (status IN ('active', 'revoked')),
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  expires_at TIMESTAMPTZ NOT NULL,
  last_seen_at TIMESTAMPTZ NULL,
  revoked_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_admin_sessions_status
  ON admin.sessions (status);

CREATE INDEX IF NOT EXISTS idx_admin_sessions_expires_at
  ON admin.sessions (expires_at);
