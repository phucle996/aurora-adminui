package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	requestdto "aurora-adminui/internal/transport/http/dto/request"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/middleware"
	"aurora-adminui/internal/transport/http/response"
)

type K8sHandler struct {
	svc      domainsvc.K8sService
	adminSvc domainsvc.AdminService
}

func NewK8sHandler(svc domainsvc.K8sService, adminSvc domainsvc.AdminService) *K8sHandler {
	return &K8sHandler{svc: svc, adminSvc: adminSvc}
}

func (h *K8sHandler) RequireAdminSession(next http.HandlerFunc) http.HandlerFunc {
	return middleware.RequireAdminSession(h.adminSvc, next)
}

func (h *K8sHandler) HandleClustersCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleListClusters(w, r)
	case http.MethodPost:
		h.handleCreateCluster(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *K8sHandler) HandleClusterItem(w http.ResponseWriter, r *http.Request) {
	clusterID, action, ok := parseK8sClusterPath(r.URL.Path)
	if !ok {
		response.JSON(w, http.StatusNotFound, "cluster not found", nil)
		return
	}

	switch {
	case r.Method == http.MethodGet && action == "":
		h.handleClusterDetail(w, r, clusterID)
	case r.Method == http.MethodPost && action == "revalidate":
		h.handleRevalidateCluster(w, r, clusterID)
	case r.Method == http.MethodDelete && action == "":
		h.handleDeleteCluster(w, r, clusterID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *K8sHandler) handleListClusters(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	items, err := h.svc.ListClusters(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to list kubernetes clusters", nil)
		return
	}
	response.JSON(w, http.StatusOK, "ok", httpdto.NewListK8sClustersResponse(items))
}

func (h *K8sHandler) handleCreateCluster(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(2 << 20); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	file, _, err := r.FormFile("kubeconfig")
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "kubeconfig is required", nil)
		return
	}
	defer file.Close()

	kubeconfig, err := io.ReadAll(file)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, "failed to read kubeconfig", nil)
		return
	}

	req := requestdto.CreateK8sClusterRequest{
		Name:                     strings.TrimSpace(r.FormValue("name")),
		Description:              strings.TrimSpace(r.FormValue("description")),
		ZoneID:                   strings.TrimSpace(r.FormValue("zone_id")),
		SupportsDBAAS:            parseFormBool(r.FormValue("supports_dbaas")),
		SupportsServerless:       parseFormBool(r.FormValue("supports_serverless")),
		SupportsGenericWorkloads: parseFormBool(r.FormValue("supports_generic_workloads")),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()
	item, err := h.svc.CreateCluster(ctx, entity.K8sClusterCreateInput{
		Name:                     req.Name,
		Description:              req.Description,
		ZoneID:                   req.ZoneID,
		SupportsDBAAS:            req.SupportsDBAAS,
		SupportsServerless:       req.SupportsServerless,
		SupportsGenericWorkloads: req.SupportsGenericWorkloads,
		Kubeconfig:               kubeconfig,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrZoneNotFound):
			response.JSON(w, http.StatusNotFound, "zone not found", nil)
		case errors.Is(err, errorx.ErrK8sClusterAlreadyExists):
			response.JSON(w, http.StatusConflict, "cluster already exists", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to create kubernetes cluster", nil)
		}
		return
	}

	response.JSON(w, http.StatusCreated, "cluster created", httpdto.NewK8sClusterDetail(item))
}

func (h *K8sHandler) handleClusterDetail(w http.ResponseWriter, r *http.Request, clusterID string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	item, err := h.svc.GetClusterDetail(ctx, clusterID)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrK8sClusterNotFound):
			response.JSON(w, http.StatusNotFound, "cluster not found", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to load kubernetes cluster", nil)
		}
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewK8sClusterDetail(item))
}

func (h *K8sHandler) handleRevalidateCluster(w http.ResponseWriter, r *http.Request, clusterID string) {
	ctx, cancel := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancel()

	item, err := h.svc.RevalidateCluster(ctx, clusterID)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrK8sClusterNotFound):
			response.JSON(w, http.StatusNotFound, "cluster not found", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to revalidate kubernetes cluster", nil)
		}
		return
	}

	response.JSON(w, http.StatusOK, "cluster revalidated", httpdto.NewK8sClusterDetail(item))
}

func (h *K8sHandler) handleDeleteCluster(w http.ResponseWriter, r *http.Request, clusterID string) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.svc.DeleteCluster(ctx, clusterID); err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrK8sClusterNotFound):
			response.JSON(w, http.StatusNotFound, "cluster not found", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to delete kubernetes cluster", nil)
		}
		return
	}
	response.JSON(w, http.StatusOK, "cluster deleted", nil)
}

func parseK8sClusterPath(rawPath string) (clusterID string, action string, ok bool) {
	path := strings.TrimSpace(strings.TrimPrefix(rawPath, "/api/v1/admin/k8s/clusters/"))
	path = strings.Trim(path, "/")
	if path == "" {
		return "", "", false
	}
	parts := strings.Split(path, "/")
	switch len(parts) {
	case 1:
		return parts[0], "", true
	case 2:
		return parts[0], parts[1], true
	default:
		return "", "", false
	}
}

func parseFormBool(raw string) bool {
	value, err := strconv.ParseBool(strings.TrimSpace(raw))
	if err != nil {
		return false
	}
	return value
}
