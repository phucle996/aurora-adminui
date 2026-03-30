package handler

import (
	"context"
	"net/http"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/transport/http/middleware"
	"aurora-adminui/internal/transport/http/response"
)

type AdminSecurityHandler struct {
	svc domainsvc.AdminService
}

func NewAdminSecurityHandler(svc domainsvc.AdminService) *AdminSecurityHandler {
	return &AdminSecurityHandler{svc: svc}
}

func (h *AdminSecurityHandler) RequireAdminSession(next http.HandlerFunc) http.HandlerFunc {
	return middleware.RequireAdminSession(h.svc, next)
}

func (h *AdminSecurityHandler) HandleTwoFactorStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	securityState, err := h.svc.GetTwoFactorStatus(ctx)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, "ok", map[string]interface{}{
		"two_factor_enabled": securityState.TwoFactorEnabled,
		"totp_enabled_at":    securityState.TOTPEnabledAt,
	})
}

func (h *AdminSecurityHandler) HandleBeginTOTPSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	result, err := h.svc.BeginTOTPSetup(ctx)
	if err != nil {
		writeAdminError(w, err)
		return
	}
	if result.AlreadyEnabled {
		response.JSON(w, http.StatusConflict, "2fa already enabled", nil)
		return
	}
	response.JSON(w, http.StatusOK, "totp setup created", map[string]interface{}{
		"setup_session":     result.SetupSession,
		"secret":            result.Secret,
		"otpauth_url":       result.OTPAuthURL,
		"setup_ttl_seconds": result.TTLSeconds,
	})
}

func (h *AdminSecurityHandler) HandleConfirmTOTPSetup(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		SetupSession string `json:"setup_session"`
		Code         string `json:"code"`
	}
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.ConfirmTOTPSetup(ctx, req.SetupSession, req.Code); err != nil {
		writeAdminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, "2fa enabled", nil)
}

func (h *AdminSecurityHandler) HandleDisableTwoFactor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Code string `json:"code"`
	}
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.DisableTwoFactor(ctx, req.Code); err != nil {
		writeAdminError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, "2fa disabled", nil)
}
