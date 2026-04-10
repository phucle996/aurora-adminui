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

	"github.com/gorilla/websocket"
)

type HypervisorHandler struct {
	svc      domainsvc.HypervisorService
	upgrader websocket.Upgrader
}

func NewHypervisorHandler(svc domainsvc.HypervisorService) *HypervisorHandler {
	return &HypervisorHandler{
		svc: svc,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (h *HypervisorHandler) HandleListNodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	items, err := h.svc.ListNodes(ctx)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, "failed to list hypervisor nodes", nil)
		return
	}
	response.JSON(w, http.StatusOK, "ok", httpdto.NewListHypervisorNodesResponse(items))
}

func (h *HypervisorHandler) HandleMetricsStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	nodeID := strings.TrimSpace(r.URL.Query().Get("id"))
	if nodeID == "" {
		response.JSON(w, http.StatusBadRequest, "id is required", nil)
		return
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	_ = conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(90 * time.Second))
	})

	closed := make(chan struct{})
	go h.consumeClientFrames(conn, closed)

	if err := h.writeMetricSnapshot(r.Context(), conn, nodeID); err != nil {
		_ = conn.WriteJSON(map[string]string{
			"type":    "error",
			"message": err.Error(),
		})
		return
	}

	ticker := time.NewTicker(5 * time.Second)
	pingTicker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()
	defer pingTicker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-closed:
			return
		case <-pingTicker.C:
			if err := conn.WriteControl(websocket.PingMessage, []byte("ping"), time.Now().Add(5*time.Second)); err != nil {
				return
			}
		case <-ticker.C:
			if err := h.writeMetricSnapshot(r.Context(), conn, nodeID); err != nil {
				_ = conn.WriteJSON(map[string]string{
					"type":    "error",
					"message": err.Error(),
				})
				return
			}
		}
	}
}

func (h *HypervisorHandler) HandleNodeItem(w http.ResponseWriter, r *http.Request) {
	nodeID, action, ok := parseHypervisorNodePath(r.URL.Path)
	if !ok {
		response.JSON(w, http.StatusNotFound, "hypervisor node not found", nil)
		return
	}

	switch {
	case r.Method == http.MethodGet && action == "":
		h.handleNodeDetail(w, r, nodeID)
		return
	case r.Method == http.MethodPatch && action == "name":
		h.handleNodeName(w, r, nodeID)
		return
	case r.Method == http.MethodPatch && action == "zone":
		h.handleNodeZone(w, r, nodeID)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (h *HypervisorHandler) handleNodeDetail(w http.ResponseWriter, r *http.Request, nodeID string) {
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	item, err := h.svc.GetNodeDetail(ctx, nodeID)
	if err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrHypervisorNodeNotFound):
			response.JSON(w, http.StatusNotFound, "hypervisor node not found", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to load hypervisor node", nil)
		}
		return
	}
	response.JSON(w, http.StatusOK, "ok", httpdto.NewHypervisorNodeDetailResponse(item))
}

func (h *HypervisorHandler) handleNodeName(w http.ResponseWriter, r *http.Request, nodeID string) {
	var req requestdto.UpdateHypervisorNodeNameRequest
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.UpdateNodeName(ctx, nodeID, req.Name); err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrHypervisorNodeNotFound):
			response.JSON(w, http.StatusNotFound, "hypervisor node not found", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to update node name", nil)
		}
		return
	}
	response.JSON(w, http.StatusOK, "node name updated", nil)
}

func (h *HypervisorHandler) handleNodeZone(w http.ResponseWriter, r *http.Request, nodeID string) {
	var req requestdto.AssignHypervisorNodeZoneRequest
	if err := decodeJSON(r, &req); err != nil {
		response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.svc.AssignNodeToZone(ctx, nodeID, req.ZoneID); err != nil {
		switch {
		case errors.Is(err, errorx.ErrInvalidArgument):
			response.JSON(w, http.StatusBadRequest, "invalid request", nil)
		case errors.Is(err, errorx.ErrZoneNotFound):
			response.JSON(w, http.StatusNotFound, "zone not found", nil)
		case errors.Is(err, errorx.ErrHypervisorNodeNotFound):
			response.JSON(w, http.StatusNotFound, "hypervisor node not found", nil)
		default:
			response.JSON(w, http.StatusInternalServerError, "failed to assign node to zone", nil)
		}
		return
	}
	response.JSON(w, http.StatusOK, "node assigned to zone", nil)
}

func (h *HypervisorHandler) consumeClientFrames(conn *websocket.Conn, closed chan<- struct{}) {
	defer close(closed)
	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}

func (h *HypervisorHandler) writeMetricSnapshot(parent context.Context, conn *websocket.Conn, nodeID string) error {
	ctx, cancel := context.WithTimeout(parent, 8*time.Second)
	defer cancel()

	series, err := h.svc.GetNodeMetrics(ctx, nodeID)
	if err != nil {
		return err
	}

	items := make([]httpdto.HypervisorMetricSeries, 0, len(series))
	for _, item := range series {
		points := make([]httpdto.HypervisorMetricPoint, 0, len(item.Points))
		for _, point := range item.Points {
			points = append(points, httpdto.HypervisorMetricPoint{
				Timestamp: point.Timestamp,
				Value:     point.Value,
			})
		}
		items = append(items, httpdto.HypervisorMetricSeries{
			Name:   item.Name,
			Label:  item.Label,
			Unit:   item.Unit,
			Latest: item.Latest,
			Points: points,
		})
	}

	return conn.WriteJSON(struct {
		ID          string                           `json:"id"`
		GeneratedAt time.Time                        `json:"generated_at"`
		WindowSec   int                              `json:"window_sec"`
		StepSec     int                              `json:"step_sec"`
		Series      []httpdto.HypervisorMetricSeries `json:"series"`
	}{
		ID:          nodeID,
		GeneratedAt: time.Now().UTC(),
		WindowSec:   900,
		StepSec:     30,
		Series:      items,
	})
}

func parseHypervisorNodePath(rawPath string) (nodeID string, action string, ok bool) {
	path := strings.TrimSpace(strings.TrimPrefix(rawPath, "/api/v1/admin/hypervisor/nodes/"))
	path = strings.Trim(path, "/")
	if path == "" {
		return "", "", false
	}

	parts := strings.Split(path, "/")
	if len(parts) == 1 {
		return strings.TrimSpace(parts[0]), "", parts[0] != ""
	}
	if len(parts) == 2 && parts[0] != "" && (parts[1] == "name" || parts[1] == "zone") {
		return strings.TrimSpace(parts[0]), parts[1], true
	}
	return "", "", false
}
