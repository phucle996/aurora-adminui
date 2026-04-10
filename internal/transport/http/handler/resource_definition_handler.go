package handler

import (
	"context"
	"errors"
	"strings"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	"aurora-adminui/internal/transport/http/dto/request"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/response"
	"aurora-adminui/pkg/logger"

	"github.com/gin-gonic/gin"
)

type RDHandler struct {
	svc domainsvc.RDService
}

// NewResourceDefinitionHandler builds the HTTP adapter for resource definition pages.
func NewRDHandler(svc domainsvc.RDService) *RDHandler {
	return &RDHandler{svc: svc}
}

// ListRD returns the slim catalog used by the resource definitions page.
func (h *RDHandler) ListRD(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	items, err := h.svc.ListResourceDefinitionCatalog(ctx)
	if err != nil {
		logger.HandlerError(ctx, "admin.rd.list", err, "failed to list resource definitions")
		response.RespondInternalError(c, "failed to list resource definitions")
		return
	}
	logger.HandlerInfo(ctx, "admin.rd.list", "listed %d resource definitions", len(items))
	response.RespondSuccess(c, mapResourceDefinitionCatalog(items), "ok")
}

// ListRDTemplateOptions returns only the fields the template form needs.
func (h *RDHandler) ListRDTemplateOptions(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	items, err := h.svc.ListResourceDefinitions(ctx)
	if err != nil {
		logger.HandlerError(ctx, "admin.rd.template_options", err, "failed to list resource definition template options")
		response.RespondInternalError(c, "failed to list resource definitions")
		return
	}
	response.RespondSuccess(c, mapTemplateResourceDefinitionOptions(items), "ok")
}

// CreateRD validates input and creates one exact resource definition record.
func (h *RDHandler) CreateRD(c *gin.Context) {
	var req request.CreateRDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondBadRequest(c, "invalid request")
		return
	}
	item := &entity.ResourceDefinition{
		ResourceType:    strings.TrimSpace(req.ResourceType),
		ResourceModel:   firstNonEmpty(strings.TrimSpace(req.ResourceModel), strings.TrimSpace(req.ModelName)),
		ResourceVersion: strings.TrimSpace(req.ResourceVersion),
		DisplayName:     strings.TrimSpace(req.DisplayName),
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	err := h.svc.CreateResourceDefinition(ctx, item)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			logger.HandlerWarn(ctx, "admin.rd.create", "invalid create resource definition request")
			response.RespondBadRequest(c, "invalid request")
		case errors.Is(err, errorx.ErrResourceDefinitionAlreadyExists):
			logger.HandlerWarn(ctx, "admin.rd.create", "resource definition already exists")
			response.RespondConflict(c, "resource definition already exists")
		default:
			logger.HandlerError(ctx, "admin.rd.create", err, "internal server error")
			response.RespondInternalError(c, "internal server error")
		}
		return
	}
	logger.HandlerInfo(ctx, "admin.rd.create", "created resource definition")
	response.RespondCreated(c, nil, "resource definition created")
}

func (h *RDHandler) ListRDZones(c *gin.Context) {
	definitionID := strings.TrimSpace(c.Param("id"))
	if definitionID == "" {
		response.RespondBadRequest(c, "invalid resource definition id")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	items, err := h.svc.ListResourceDefinitionZoneSupport(ctx, definitionID)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.RespondBadRequest(c, "invalid request")
		case errors.Is(err, errorx.ErrResourceDefinitionNotFound):
			response.RespondNotFound(c, "resource definition not found")
		default:
			logger.HandlerError(ctx, "admin.rd.list_zones", err, "failed to list resource definition zones id=%s", definitionID)
			response.RespondInternalError(c, "failed to list resource definition zones")
		}
		return
	}
	logger.HandlerInfo(ctx, "admin.rd.list_zones", "listed resource definition zones id=%s", definitionID)
	response.RespondSuccess(c, gin.H{"items": items}, "ok")
}

func (h *RDHandler) ReplaceRDZones(c *gin.Context) {
	definitionID := strings.TrimSpace(c.Param("id"))
	if definitionID == "" {
		response.RespondBadRequest(c, "invalid resource definition id")
		return
	}
	var req request.ReplaceResourceDefinitionZoneSupportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondBadRequest(c, "invalid request")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.ReplaceResourceDefinitionZoneSupport(ctx, definitionID, req.ZoneIDs); err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.RespondBadRequest(c, "invalid request")
		case errors.Is(err, errorx.ErrZoneNotFound):
			response.RespondNotFound(c, "zone not found")
		case errors.Is(err, errorx.ErrResourceDefinitionNotFound):
			response.RespondNotFound(c, "resource definition not found")
		case errors.Is(err, errorx.ErrResourceDefinitionNeedsZones):
			response.RespondConflict(c, "resource definition requires at least one enabled zone")
		default:
			logger.HandlerError(ctx, "admin.rd.replace_zones", err, "failed to replace resource definition zones id=%s", definitionID)
			response.RespondInternalError(c, "failed to update resource definition zones")
		}
		return
	}
	logger.HandlerInfo(ctx, "admin.rd.replace_zones", "updated resource definition zones id=%s", definitionID)
	response.RespondSuccess(c, nil, "resource definition zones updated")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

// DeleteRD deletes one exact resource definition record.
func (h *RDHandler) DeleteRD(c *gin.Context) {
	definitionID := strings.TrimSpace(c.Param("id"))
	if definitionID == "" {
		response.RespondBadRequest(c, "invalid resource definition id")
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.DeleteResourceDefinition(ctx, definitionID); err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.RespondBadRequest(c, "invalid request")
		case errors.Is(err, errorx.ErrResourceDefinitionHasResources):
			response.RespondConflict(c, "resource definition still has resources")
		case errors.Is(err, errorx.ErrResourceDefinitionNotFound):
			response.RespondNotFound(c, "resource definition not found")
		default:
			logger.HandlerError(ctx, "admin.rd.delete", err, "failed to delete resource definition id=%s", definitionID)
			response.RespondInternalError(c, "failed to delete resource definition")
		}
		return
	}
	logger.HandlerInfo(ctx, "admin.rd.delete", "deleted resource definition id=%s", definitionID)
	response.RespondSuccess(c, gin.H{"id": definitionID}, "resource definition deleted")
}

// UpdateRDStatus changes one definition lifecycle state after validation.
func (h *RDHandler) UpdateRDStatus(c *gin.Context) {
	definitionID := strings.TrimSpace(c.Param("id"))
	if definitionID == "" {
		response.RespondBadRequest(c, "invalid resource definition id")
		return
	}
	var req request.UpdateResourceDefinitionStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.RespondBadRequest(c, "invalid request")
		return
	}
	req.Status = strings.TrimSpace(strings.ToLower(req.Status))
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	item, err := h.svc.UpdateResourceDefinitionStatus(ctx, definitionID, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.RespondBadRequest(c, "invalid request")
		case errors.Is(err, errorx.ErrResourceDefinitionNeedsZones):
			response.RespondConflict(c, "resource definition requires at least one enabled zone")
		case errors.Is(err, errorx.ErrResourceDefinitionNotFound):
			response.RespondNotFound(c, "resource definition not found")
		default:
			logger.HandlerError(ctx, "admin.rd.update_status", err, "failed to update resource definition status id=%s", definitionID)
			response.RespondInternalError(c, "failed to update resource definition status")
		}
		return
	}
	logger.HandlerInfo(ctx, "admin.rd.update_status", "updated resource definition id=%s status=%s", definitionID, item.Status)
	response.RespondSuccess(c, mapResourceDefinition(*item, 0), "resource definition status updated")
}

func mapResourceDefinitionCatalog(items []entity.ResourceDefinitionCatalogItem) httpdto.ListResourceDefinitionsResponse {
	out := make([]httpdto.ResourceDefinition, 0, len(items))
	for _, item := range items {
		out = append(out, mapResourceDefinition(item.ResourceDefinition, item.ResourceCount))
	}
	return httpdto.ListResourceDefinitionsResponse{Items: out}
}

func mapResourceDefinition(item entity.ResourceDefinition, resourceCount int) httpdto.ResourceDefinition {
	return httpdto.ResourceDefinition{
		ID:              item.ID.String(),
		ResourceType:    item.ResourceType,
		ResourceModel:   item.ResourceModel,
		ResourceVersion: item.ResourceVersion,
		DisplayName:     item.DisplayName,
		Status:          item.Status,
		ResourceCount:   resourceCount,
	}
}

func mapTemplateResourceDefinitionOptions(items []entity.ResourceDefinition) httpdto.ListTemplateResourceDefinitionOptionsResponse {
	out := make([]httpdto.TemplateResourceDefinitionOption, 0, len(items))
	for _, item := range items {
		out = append(out, httpdto.TemplateResourceDefinitionOption{
			ID:              item.ID.String(),
			ResourceType:    item.ResourceType,
			ResourceModel:   item.ResourceModel,
			ResourceVersion: item.ResourceVersion,
		})
	}
	return httpdto.ListTemplateResourceDefinitionOptionsResponse{Items: out}
}
