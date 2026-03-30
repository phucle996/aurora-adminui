package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	"aurora-adminui/internal/transport/http/response"
)

func RequireAdminSession(svc domainsvc.AdminService, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := ValidateAdminSession(r.Context(), r, svc); err != nil {
			response.JSON(w, http.StatusUnauthorized, "unauthorized", nil)
			return
		}
		next(w, r)
	}
}

func ValidateAdminSession(parent context.Context, r *http.Request, svc domainsvc.AdminService) error {
	sessionID := ReadAdminSessionCookie(r)
	if sessionID == "" {
		return errorx.ErrTokenInvalid
	}
	ctx, cancel := context.WithTimeout(parent, 5*time.Second)
	defer cancel()
	_, err := svc.ValidateSession(ctx, sessionID)
	return err
}

func ReadAdminSessionCookie(r *http.Request) string {
	if r == nil {
		return ""
	}
	cookie, err := r.Cookie("admin_session")
	if err != nil {
		return ""
	}
	return cookie.Value
}

func SetSessionCookie(w http.ResponseWriter, value string, maxAge int, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    value,
		Path:     "/",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func ClearSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     "admin_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	})
}

func RequestSecure(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	return strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}
