package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	domainrequest "aurora-adminui/internal/transport/http/dto/request"
	httpdto "aurora-adminui/internal/transport/http/dto/response"

	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	"aurora-adminui/internal/transport/http/response"
)

type ZoneHandler struct {
	svc domainsvc.ZoneService
}

func NewZoneHandler(svc domainsvc.ZoneService) *ZoneHandler {
	return &ZoneHandler{svc: svc}
}

func (h *ZoneHandler) HandleZonesCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListZones(w, r)
	case http.MethodPost:
		h.handleCreateZone(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ZoneHandler) HandleZoneItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		h.handleDeleteZone(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *ZoneHandler) handleListZones(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	items, err := h.svc.ListZones(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to list zones", nil)
		return
	}
	response.JSON(w, http.StatusOK, "ok", httpdto.NewListZonesResponse(items))
}

func (h *ZoneHandler) handleCreateZone(w http.ResponseWriter, r *http.Request) {
	var req domainrequest.CreateZoneRequest
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	item, err := h.svc.CreateZone(ctx, req.Name, req.Description)
	if err != nil {
		h.writeZoneError(w, err)
		return
	}
	response.JSON(w, http.StatusCreated, "zone created", httpdto.NewZone(item))
}

func (h *ZoneHandler) handleDeleteZone(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/zones/")
	if strings.TrimSpace(id) == "" || strings.Contains(id, "/") {
		response.JSON(w, http.StatusNotFound, "zone not found", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.DeleteZone(ctx, id); err != nil {
		h.writeZoneError(w, err)
		return
	}
	response.JSON(w, http.StatusOK, "zone deleted", nil)
}

func (h *ZoneHandler) writeZoneError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, errorx.ErrInvalidArgument):
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
	case errors.Is(err, errorx.ErrZoneAlreadyExists):
		response.JSON(w, http.StatusConflict, "zone already exists", nil)
	case errors.Is(err, errorx.ErrZoneHasResources):
		response.JSON(w, http.StatusConflict, "zone still has attached nodes", nil)
	case errors.Is(err, errorx.ErrZoneNotFound):
		response.JSON(w, http.StatusNotFound, "zone not found", nil)
	default:
		response.JSON(w, http.StatusInternalServerError, "internal server error", nil)
	}
}
