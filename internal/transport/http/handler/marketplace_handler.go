package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	requestdto "aurora-adminui/internal/transport/http/dto/request"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/response"
)

type MarketplaceHandler struct {
	svc domainsvc.MarketplaceService
}

func NewMarketplaceHandler(svc domainsvc.MarketplaceService) *MarketplaceHandler {
	return &MarketplaceHandler{svc: svc}
}

func (h *MarketplaceHandler) HandleMarketplaceCollection(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/api/v1/admin/marketplace/model-options" && r.Method == http.MethodGet {
		h.handleListMarketplaceModelOptions(w, r)
		return
	}
	if r.URL.Path == "/api/v1/admin/marketplace/template-options" && r.Method == http.MethodGet {
		h.handleListMarketplaceTemplateOptions(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		h.handleListMarketplaceApps(w, r)
	case http.MethodPost:
		h.handleCreateMarketplaceApp(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MarketplaceHandler) handleListMarketplaceModelOptions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.ListMarketplaceModelOptions(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to load marketplace model options", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewListMarketplaceModelOptionsResponse(items))
}

func (h *MarketplaceHandler) handleListMarketplaceTemplateOptions(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.ListMarketplaceTemplateOptions(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to load marketplace template options", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewListMarketplaceTemplateOptionsResponse(items))
}

func (h *MarketplaceHandler) HandleMarketplaceItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetMarketplaceApp(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *MarketplaceHandler) handleListMarketplaceApps(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.ListMarketplaceApps(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to load marketplace apps", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewListMarketplaceAppsResponse(items))
}

func (h *MarketplaceHandler) handleGetMarketplaceApp(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/marketplace/")
	id = strings.TrimSuffix(id, "/")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	item, err := h.svc.GetMarketplaceApp(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid marketplace app id", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to load marketplace app", nil)
		}
		return
	}
	if item == nil {
		response.JSON(w, http.StatusNotFound, "marketplace app not found", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewMarketplaceApp(*item))
}

func (h *MarketplaceHandler) handleCreateMarketplaceApp(w http.ResponseWriter, r *http.Request) {
	var req requestdto.CreateMarketplaceAppRequest
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	item, err := h.svc.CreateMarketplaceApp(ctx, domainsvc.CreateMarketplaceAppInput{
		Name:                 req.Name,
		Slug:                 req.Slug,
		Summary:              req.Summary,
		Description:          req.Description,
		ResourceDefinitionID: req.ResourceDefinitionID,
		TemplateID:           req.TemplateID,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrMarketplaceAppAlreadyExists):
			response.JSON(w, http.StatusConflict, "marketplace app already exists", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to create marketplace app", nil)
		}
		return
	}

	response.JSON(w, http.StatusCreated, "marketplace app created", httpdto.NewMarketplaceApp(*item))
}
