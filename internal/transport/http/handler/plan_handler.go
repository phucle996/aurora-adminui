package handler

import (
	"context"
	"errors"
	"net/http"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	requestdto "aurora-adminui/internal/transport/http/dto/request"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/middleware"
	"aurora-adminui/internal/transport/http/response"
)

type PlanHandler struct {
	svc      domainsvc.PlanService
	adminSvc domainsvc.AdminService
}

func NewPlanHandler(svc domainsvc.PlanService, adminSvc domainsvc.AdminService) *PlanHandler {
	return &PlanHandler{svc: svc, adminSvc: adminSvc}
}

func (h *PlanHandler) RequireAdminSession(next http.HandlerFunc) http.HandlerFunc {
	return middleware.RequireAdminSession(h.adminSvc, next)
}

func (h *PlanHandler) HandlePlansCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListPlans(w, r)
	case http.MethodPost:
		h.handleCreatePlan(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *PlanHandler) handleListPlans(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.ListPlans(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to list plans", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewListPlansResponse(items))
}

func (h *PlanHandler) handleCreatePlan(w http.ResponseWriter, r *http.Request) {
	var req requestdto.CreatePlanRequest
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	item, err := h.svc.CreatePlan(
		ctx,
		req.ResourceType,
		req.Code,
		req.Name,
		req.Description,
		req.VCPU,
		req.RAMGB,
		req.DiskGB,
	)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrPlanAlreadyExists):
			response.JSON(w, http.StatusConflict, "plan already exists", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to create plan", nil)
		}
		return
	}

	response.JSON(w, http.StatusCreated, "plan created", httpdto.NewPlan(*item))
}
