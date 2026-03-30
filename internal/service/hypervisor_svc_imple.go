package service

import (
	"context"
	"strings"
	"time"

	"aurora-adminui/infra/victoria"
	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type HypervisorSvcImple struct {
	repo     domainrepo.HypervisorRepository
	victoria *victoria.Client
}

func NewHypervisorService(repo domainrepo.HypervisorRepository, victoriaClient *victoria.Client) domainsvc.HypervisorService {
	return &HypervisorSvcImple{
		repo:     repo,
		victoria: victoriaClient,
	}
}

func (s *HypervisorSvcImple) ListNodes(ctx context.Context) ([]entity.HypervisorNode, error) {
	return s.repo.ListNodes(ctx)
}

func (s *HypervisorSvcImple) GetNodeDetail(ctx context.Context, nodeID string) (*entity.HypervisorNodeDetail, error) {
	nodeID = strings.TrimSpace(nodeID)
	if nodeID == "" {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.GetNodeDetail(ctx, nodeID)
}

func (s *HypervisorSvcImple) GetNodeMetrics(ctx context.Context, nodeID string) ([]entity.HypervisorMetricSeries, error) {
	nodeID = strings.TrimSpace(nodeID)
	if nodeID == "" {
		return nil, errorx.ErrInvalidArgument
	}
	return s.loadNodeMetrics(ctx, nodeID), nil
}

func (s *HypervisorSvcImple) UpdateNodeName(ctx context.Context, nodeID, name string) error {
	nodeID = strings.TrimSpace(nodeID)
	name = strings.TrimSpace(name)
	if nodeID == "" || name == "" {
		return errorx.ErrInvalidArgument
	}
	return s.repo.UpdateNodeName(ctx, nodeID, name)
}

func (s *HypervisorSvcImple) AssignNodeToZone(ctx context.Context, nodeID, zoneID string) error {
	nodeID = strings.TrimSpace(nodeID)
	zoneID = strings.TrimSpace(zoneID)
	if nodeID == "" || zoneID == "" {
		return errorx.ErrInvalidArgument
	}
	parsedZoneID, err := uuid.Parse(zoneID)
	if err != nil {
		return errorx.ErrInvalidArgument
	}
	return s.repo.AssignNodeToZone(ctx, nodeID, parsedZoneID)
}

type metricQuerySpec struct {
	Name  string
	Unit  string
	Query string
}

func (s *HypervisorSvcImple) loadNodeMetrics(ctx context.Context, nodeID string) []entity.HypervisorMetricSeries {
	if s.victoria == nil || !s.victoria.Configured() {
		return nil
	}

	step := 30 * time.Second
	now := time.Now().UTC().Truncate(step)
	start := now.Add(-15 * time.Minute)
	specs := []metricQuerySpec{
		{Name: "cpu_usage_percent", Unit: "percent", Query: `max(hypervisor_host_cpu_usage_percent{agent_node_id="` + nodeID + `"})`},
		{Name: "ram_used_bytes", Unit: "bytes", Query: `max(hypervisor_host_ram_used_bytes{agent_node_id="` + nodeID + `"})`},
		{Name: "disk_io_bytes_per_sec", Unit: "bytes_per_sec", Query: `sum(hypervisor_host_disk_read_bytes_per_sec{agent_node_id="` + nodeID + `"}) + sum(hypervisor_host_disk_write_bytes_per_sec{agent_node_id="` + nodeID + `"})`},
		{Name: "network_bytes_per_sec", Unit: "bytes_per_sec", Query: `sum(hypervisor_host_network_rx_bytes_per_sec{agent_node_id="` + nodeID + `"}) + sum(hypervisor_host_network_tx_bytes_per_sec{agent_node_id="` + nodeID + `"})`},
		{Name: "gpu_usage_percent", Unit: "percent", Query: `avg(hypervisor_host_gpu_utilization_percent{agent_node_id="` + nodeID + `"})`},
	}

	series := make([]entity.HypervisorMetricSeries, len(specs))
	g, groupCtx := errgroup.WithContext(ctx)
	for i, spec := range specs {
		i := i
		spec := spec
		g.Go(func() error {
			latest, err := s.victoria.Query(groupCtx, spec.Query, now)
			if err != nil {
				return err
			}

			points, err := s.victoria.QueryRange(groupCtx, spec.Query, start, now, step)
			if err != nil {
				return err
			}
			metricPoints := make([]entity.HypervisorMetricPoint, 0, len(points))
			for _, point := range points {
				metricPoints = append(metricPoints, entity.HypervisorMetricPoint{
					Timestamp: point.Timestamp,
					Value:     point.Value,
				})
			}
			series[i] = entity.HypervisorMetricSeries{
				Name:   spec.Name,
				Unit:   spec.Unit,
				Latest: latest,
				Points: metricPoints,
			}
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		return nil
	}

	perGPUSeries, err := s.victoria.QueryRangeSeries(
		ctx,
		`hypervisor_host_gpu_utilization_percent{agent_node_id="`+nodeID+`"}`,
		start,
		now,
		step,
	)
	if err != nil {
		return series
	}

	for _, rawSeries := range perGPUSeries {
		label := strings.TrimSpace(rawSeries.Labels["gpu_model"])
		if label == "" {
			label = strings.TrimSpace(rawSeries.Labels["gpu_pci_address"])
		}
		if label == "" {
			label = strings.TrimSpace(rawSeries.Labels["gpu_uuid"])
		}
		if label == "" {
			label = "Unknown GPU"
		}

		points := make([]entity.HypervisorMetricPoint, 0, len(rawSeries.Points))
		latest := 0.0
		for _, point := range rawSeries.Points {
			points = append(points, entity.HypervisorMetricPoint{
				Timestamp: point.Timestamp,
				Value:     point.Value,
			})
			latest = point.Value
		}
		series = append(series, entity.HypervisorMetricSeries{
			Name:   "gpu_usage_percent_device",
			Label:  label,
			Unit:   "percent",
			Latest: latest,
			Points: points,
		})
	}

	return series
}
