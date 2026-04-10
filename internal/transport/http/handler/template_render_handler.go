package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	domainrequest "aurora-adminui/internal/transport/http/dto/request"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/response"

	"gopkg.in/yaml.v3"
)

var (
	templateRenderParamPattern   = regexp.MustCompile(`\{\{[a-zA-Z_][a-zA-Z0-9_]*\}\}`)
	templateRenderBuiltinPattern = regexp.MustCompile(`\{\$random_(12_number|9_char)\}`)
)

type TemplateRenderHandler struct {
	svc domainsvc.TemplateRenderService
}

// NewTemplateRenderHandler builds the HTTP adapter for template render pages.
func NewTemplateRenderHandler(svc domainsvc.TemplateRenderService) *TemplateRenderHandler {
	return &TemplateRenderHandler{svc: svc}
}

// HandleTemplateRenderCollection serves create operations for template renders.
func (h *TemplateRenderHandler) HandleTemplateRenderCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.HandleCreateTemplateRender(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleTemplateRenderItem serves detail and update operations for one template render.
func (h *TemplateRenderHandler) HandleTemplateRenderItem(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.HandleGetTemplateRender(w, r)
	case http.MethodPatch:
		h.HandleUpdateTemplateRender(w, r)
	case http.MethodDelete:
		h.HandleDeleteTemplateRender(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// HandleListTemplateRenderCatalog returns the minimal catalog payload used by the list page.
func (h *TemplateRenderHandler) HandleListTemplateRenderCatalog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	items, err := h.svc.ListTemplateRenderCatalog(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to list template renders", nil)
		return
	}
	response.JSON(w, http.StatusOK, "ok", httpdto.NewListTemplateRenderCatalogResponse(items))
}

// HandleGetTemplateRender loads one template render for the detail page.
func (h *TemplateRenderHandler) HandleGetTemplateRender(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/resource-templates/")
	id = strings.TrimSuffix(id, "/")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	item, err := h.svc.GetTemplateRender(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid template render id", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to load template render", nil)
		}
		return
	}
	if item == nil {
		response.JSON(w, http.StatusNotFound, "template render not found", nil)
		return
	}
	response.JSON(w, http.StatusOK, "ok", httpdto.NewTemplateRender(*item))
}

// HandleCreateTemplateRender validates and creates a template render.
func (h *TemplateRenderHandler) HandleCreateTemplateRender(w http.ResponseWriter, r *http.Request) {
	var req domainrequest.CreateTemplateRenderRequest
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := validateTemplateRenderYAML(req.YAMLTemplate); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid yaml template", nil)
		return
	}
	item, err := h.svc.CreateTemplateRender(ctx, domainsvc.CreateTemplateRenderInput{
		ResourceDefinitionID: req.ResourceDefinitionID,
		Name:                 req.Name,
		Description:          req.Description,
		StreamKey:            req.StreamKey,
		ConsumerGroup:        req.ConsumerGroup,
		YAMLTemplate:         req.YAMLTemplate,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to create template render", nil)
		}
		return
	}
	response.JSON(w, http.StatusCreated, "template render created", httpdto.NewTemplateRender(*item))
}

// HandleUpdateTemplateRender validates and updates one template render.
func (h *TemplateRenderHandler) HandleUpdateTemplateRender(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/resource-templates/")
	id = strings.TrimSuffix(id, "/")
	var req domainrequest.UpdateTemplateRenderRequest
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := validateTemplateRenderYAML(req.YAMLTemplate); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid yaml template", nil)
		return
	}
	item, err := h.svc.UpdateTemplateRender(ctx, id, domainsvc.UpdateTemplateRenderInput{
		ResourceDefinitionID: req.ResourceDefinitionID,
		Name:                 req.Name,
		Description:          req.Description,
		StreamKey:            req.StreamKey,
		ConsumerGroup:        req.ConsumerGroup,
		YAMLTemplate:         req.YAMLTemplate,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to update template render", nil)
		}
		return
	}
	if item == nil {
		response.JSON(w, http.StatusNotFound, "template render not found", nil)
		return
	}
	response.JSON(w, http.StatusOK, "template render updated", httpdto.NewTemplateRender(*item))
}

// HandleDeleteTemplateRender deletes one template render record.
func (h *TemplateRenderHandler) HandleDeleteTemplateRender(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/admin/resource-templates/")
	id = strings.TrimSuffix(id, "/")
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.DeleteTemplateRender(ctx, id); err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid template render id", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to delete template render", nil)
		}
		return
	}
	response.JSON(w, http.StatusOK, "template render deleted", map[string]any{"id": id})
}

func validateTemplateRenderYAML(body string) error {
	trimmed := strings.TrimSpace(body)
	sanitized := templateRenderParamPattern.ReplaceAllString(trimmed, `"param"`)
	sanitized = templateRenderBuiltinPattern.ReplaceAllString(sanitized, `"generated"`)
	if strings.Contains(sanitized, "{{") || strings.Contains(sanitized, "}}") || strings.Contains(sanitized, "{$") {
		return errors.New("invalid template placeholder syntax")
	}
	decoder := yaml.NewDecoder(strings.NewReader(sanitized))
	hasDocument := false
	for {
		var doc map[string]any
		err := decoder.Decode(&doc)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		if len(doc) > 0 {
			hasDocument = true
		}
	}
	if !hasDocument {
		return errors.New("empty yaml")
	}
	return nil
}
