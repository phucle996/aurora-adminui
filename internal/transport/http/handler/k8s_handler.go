package handler

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
	requestdto "aurora-adminui/internal/transport/http/dto/request"
	httpdto "aurora-adminui/internal/transport/http/dto/response"
	"aurora-adminui/internal/transport/http/response"
)

type K8sHandler struct {
	svc     domainsvc.K8sService
	zoneSvc domainsvc.ZoneService
}

// NewK8sHandler wires the HTTP adapter for kubernetes cluster management pages.
func NewK8sHandler(svc domainsvc.K8sService, zoneSvc domainsvc.ZoneService) *K8sHandler {
	return &K8sHandler{svc: svc, zoneSvc: zoneSvc}
}

// HandleClustersCollection serves list and create operations for the kubernetes catalog.
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

// HandleClusterItem serves the page-data, update, revalidate, and delete flows for one cluster.
func (h *K8sHandler) HandleClusterItem(w http.ResponseWriter, r *http.Request) {
	clusterID, action, ok := parseK8sClusterPath(r.URL.Path)
	if !ok {
		response.JSON(w, http.StatusNotFound, "cluster not found", nil)
		return
	}

	switch {
	case r.Method == http.MethodGet && action == "":
		h.handleClusterDetail(w, r, clusterID)
	case r.Method == http.MethodGet && action == "page-data":
		h.handleClusterDetailPageData(w, r, clusterID)
	case r.Method == http.MethodPatch && action == "":
		h.handleUpdateCluster(w, r, clusterID)
	case r.Method == http.MethodPost && action == "revalidate":
		h.handleRevalidateCluster(w, r, clusterID)
	case r.Method == http.MethodDelete && action == "":
		h.handleDeleteCluster(w, r, clusterID)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *K8sHandler) handleClusterDetailPageData(w http.ResponseWriter, r *http.Request, clusterID string) {
	ctx, cancel := context.WithTimeout(r.Context(), 20*time.Second)
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

	nodes, _ := h.svc.ListClusterNodes(ctx, clusterID)
	zones, err := h.zoneSvc.ListZones(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to load zones", nil)
		return
	}

	response.JSON(w, http.StatusOK, "ok", httpdto.NewK8sClusterDetailPageData(item, nodes, zones))
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
		Name:        strings.TrimSpace(r.FormValue("name")),
		Description: strings.TrimSpace(r.FormValue("description")),
		ZoneID:      strings.TrimSpace(r.FormValue("zone_id")),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	item, err := h.svc.CreateCluster(ctx, entity.K8sClusterCreateInput{
		Name:        req.Name,
		Description: req.Description,
		ZoneID:      req.ZoneID,
		Kubeconfig:  kubeconfig,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrZoneNotFound):
			response.JSON(w, http.StatusNotFound, "zone not found", nil)
		case errors.Is(err, errorx.ErrK8sClusterAlreadyExists):
			response.JSON(w, http.StatusConflict, "cluster already exists", nil)
		case errors.Is(err, errorx.ErrNoHealthyDataplane):
			response.JSON(w, http.StatusServiceUnavailable, "no healthy dataplane available", nil)
		case errors.Is(err, errorx.ErrDataplaneValidationFailed):
			response.JSON(w, http.StatusFailedDependency, "failed to validate kubernetes cluster", nil)
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
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	item, err := h.svc.RevalidateCluster(ctx, clusterID)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrK8sClusterNotFound):
			response.JSON(w, http.StatusNotFound, "cluster not found", nil)
		case errors.Is(err, errorx.ErrNoHealthyDataplane):
			response.JSON(w, http.StatusServiceUnavailable, "no healthy dataplane available", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to revalidate kubernetes cluster", nil)
		}
		return
	}

	response.JSON(w, http.StatusOK, "cluster revalidated", httpdto.NewK8sClusterDetail(item))
}

func (h *K8sHandler) handleUpdateCluster(w http.ResponseWriter, r *http.Request, clusterID string) {
	req := requestdto.UpdateK8sClusterRequest{}
	var kubeconfig []byte

	if strings.HasPrefix(strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type"))), "multipart/form-data") {
		if err := r.ParseMultipartForm(2 << 20); err != nil {
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
			return
		}
		req.ZoneID = strings.TrimSpace(r.FormValue("zone_id"))

		file, _, err := r.FormFile("kubeconfig")
		switch {
		case err == nil:
			defer file.Close()
			kubeconfig, err = io.ReadAll(file)
			if err != nil {
				response.JSON(w, http.StatusBadRequest, "failed to read kubeconfig", nil)
				return
			}
		case !errors.Is(err, http.ErrMissingFile):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
			return
		}
	} else {
		if err := decodeJSON(r, &req); err != nil {
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
			return
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()

	item, err := h.svc.UpdateCluster(ctx, clusterID, entity.K8sClusterUpdateInput{
		ZoneID:     strings.TrimSpace(req.ZoneID),
		Kubeconfig: kubeconfig,
	})
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrZoneNotFound):
			response.JSON(w, http.StatusNotFound, "zone not found", nil)
		case errors.Is(err, errorx.ErrK8sClusterNotFound):
			response.JSON(w, http.StatusNotFound, "cluster not found", nil)
		case errors.Is(err, errorx.ErrNoHealthyDataplane):
			response.JSON(w, http.StatusServiceUnavailable, "no healthy dataplane available", nil)
		case errors.Is(err, errorx.ErrDataplaneValidationFailed):
			response.JSON(w, http.StatusFailedDependency, "failed to validate kubernetes cluster", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to update kubernetes cluster", nil)
		}
		return
	}

	response.JSON(w, http.StatusOK, "cluster updated", httpdto.NewK8sClusterDetail(item))
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
		case errors.Is(err, errorx.ErrK8sClusterHasResources):
			response.JSON(w, http.StatusConflict, "cluster still has resources", nil)
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
