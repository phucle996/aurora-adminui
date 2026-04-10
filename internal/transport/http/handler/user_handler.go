package handler

import (
	"context"
	"net/http"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/response"
)

type UserHandler struct {
	svc domainsvc.UserService
}

func NewUserHandler(svc domainsvc.UserService) *UserHandler {
	return &UserHandler{svc: svc}
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
