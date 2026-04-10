package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	domainrequest "aurora-adminui/internal/transport/http/dto/request"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/response"
)

type RoleHandler struct {
	svc domainsvc.RoleService
}

func NewRoleHandler(svc domainsvc.RoleService) *RoleHandler {
	return &RoleHandler{svc: svc}
}

func (h *RoleHandler) HandleListRoles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.ListRoles(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to list roles", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewListRolesResponse(items))
}

func (h *RoleHandler) HandleListPermissions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.ListPermissions(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to list permissions", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewListPermissionsResponse(items))
}

func (h *RoleHandler) HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req domainrequest.CreateRoleRequest
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	item, err := h.svc.CreateRole(ctx, req.Name, req.Description, req.PermissionIDs)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrRoleAlreadyExists):
			response.JSON(w, http.StatusConflict, "role already exists", nil)
		case errors.Is(err, errorx.ErrPermissionNotFound):
			response.JSON(w, http.StatusBadRequest, "permission not found", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to create role", nil)
		}
		return
	}

	response.JSON(w, http.StatusCreated, "role created", httpdto.Role{
		ID:              item.ID,
		Name:            item.Name,
		Scope:           item.Scope,
		Description:     item.Description,
		UserCount:       item.UserCount,
		PermissionCount: item.PermissionCount,
	})
}
