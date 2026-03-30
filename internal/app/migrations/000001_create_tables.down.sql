DROP INDEX IF EXISTS admin.idx_admin_sessions_expires_at;
DROP INDEX IF EXISTS admin.idx_admin_sessions_status;
DROP TABLE IF EXISTS admin.sessions;
DROP TABLE IF EXISTS admin.security_state;
DROP TABLE IF EXISTS admin.api_tokens;
DROP SCHEMA IF EXISTS admin;
