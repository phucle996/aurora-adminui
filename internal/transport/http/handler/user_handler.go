package handler

import (
	"context"
	"net/http"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/middleware"
	"aurora-adminui/internal/transport/http/response"
)

type UserHandler struct {
	svc      domainsvc.UserService
	adminSvc domainsvc.AdminService
}

func NewUserHandler(svc domainsvc.UserService, adminSvc domainsvc.AdminService) *UserHandler {
	return &UserHandler{svc: svc, adminSvc: adminSvc}
}

func (h *UserHandler) RequireAdminSession(next http.HandlerFunc) http.HandlerFunc {
	return middleware.RequireAdminSession(h.adminSvc, next)
}

func (h *UserHandler) HandleListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.ListUsers(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to list users", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewListUsersResponse(items))
}
