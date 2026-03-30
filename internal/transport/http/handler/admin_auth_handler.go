package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	"aurora-adminui/internal/transport/http/middleware"
	"aurora-adminui/internal/transport/http/response"
)

type AdminAuthHandler struct {
	svc domainsvc.AdminService
}

func NewAdminAuthHandler(svc domainsvc.AdminService) *AdminAuthHandler {
	return &AdminAuthHandler{svc: svc}
}

func (h *AdminAuthHandler) BootstrapIfNeeded(ctx context.Context) error {
	return h.svc.BootstrapIfNeeded(ctx)
}

func (h *AdminAuthHandler) RequireAdminSession(next http.HandlerFunc) http.HandlerFunc {
	return middleware.RequireAdminSession(h.svc, next)
}

func (h *AdminAuthHandler) HandleTokenLoginInit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Token string `json:"token"`
	}
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	result, err := h.svc.LoginInit(ctx, req.Token)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	if result.MFARequired {
		response.JSON(w, http.StatusOK, "mfa required", map[string]interface{}{
			"mfa_required":        true,
			"methods":             result.MFAMethods,
			"preauth_session":     result.PreauthSession,
			"preauth_ttl_seconds": result.PreauthTTLSeconds,
		})
		return
	}
	middleware.SetSessionCookie(w, result.SessionID, int((24 * time.Hour).Seconds()), middleware.RequestSecure(r))
	payload := map[string]interface{}{
		"token_type":    result.TokenType,
		"token_version": result.TokenVersion,
	}
	if result.BootstrapExchanged {
		payload["api_token"] = result.PlaintextToken
	}
	response.JSON(w, http.StatusOK, "admin login successful", payload)
}

func (h *AdminAuthHandler) HandleTokenLoginVerify2FA(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		PreauthSession string `json:"preauth_session"`
		Code           string `json:"code"`
	}
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	result, err := h.svc.VerifySecondFactor(ctx, req.PreauthSession, req.Code)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	middleware.SetSessionCookie(w, result.SessionID, int((24 * time.Hour).Seconds()), middleware.RequestSecure(r))
	payload := map[string]interface{}{
		"token_type":    result.TokenType,
		"token_version": result.TokenVersion,
	}
	if result.BootstrapExchanged {
		payload["api_token"] = result.PlaintextToken
	}
	response.JSON(w, http.StatusOK, "admin mfa verified", payload)
}

func (h *AdminAuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	sessionID := middleware.ReadAdminSessionCookie(r)
	_ = h.svc.Logout(ctx, sessionID)
	middleware.ClearSessionCookie(w, middleware.RequestSecure(r))
	response.JSON(w, http.StatusOK, "admin logout successful", nil)
}

func (h *AdminAuthHandler) HandleSessionStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	sessionID := middleware.ReadAdminSessionCookie(r)
	session, err := h.svc.ValidateSession(ctx, sessionID)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	securityState, err := h.svc.GetTwoFactorStatus(ctx)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, "ok", map[string]interface{}{
		"authenticated":      true,
		"session_id":         session.ID.String(),
		"session_expires_at": session.ExpiresAt,
		"two_factor_enabled": securityState.TwoFactorEnabled,
	})
}

func (h *AdminAuthHandler) HandleRotateToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Code string `json:"code"`
	}
	if err := decodeJSON(r, &req); err != nil && !errors.Is(err, io.EOF) {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	result, err := h.svc.RotateAPIToken(ctx, req.Code)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, "admin api token rotated", map[string]interface{}{
		"version":       result.Version,
		"telegram_sent": result.TelegramSent,
	})
}

func writeAdminError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errorx.ErrInvalidArgument):
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
	case errors.Is(err, errorx.ErrTokenInvalid),
		errors.Is(err, errorx.ErrTokenExpired),
		errors.Is(err, errorx.ErrMFACodeInvalid),
		errors.Is(err, errorx.ErrMFAMethodNotFound):
		response.JSON(w, http.StatusUnauthorized, "unauthorized", nil)
	case errors.Is(err, errorx.ErrMFAMethodAlreadyEnabled):
		response.JSON(w, http.StatusConflict, "2fa already enabled", nil)
	default:
		response.JSON(w, http.StatusInternalServerError, "internal server error", nil)
	}
}
