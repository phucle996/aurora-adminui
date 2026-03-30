package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type HypervisorService interface {
	ListNodes(ctx context.Context) ([]entity.HypervisorNode, error)
	GetNodeDetail(ctx context.Context, nodeID string) (*entity.HypervisorNodeDetail, error)
	GetNodeMetrics(ctx context.Context, nodeID string) ([]entity.HypervisorMetricSeries, error)
	UpdateNodeName(ctx context.Context, nodeID, name string) error
	AssignNodeToZone(ctx context.Context, nodeID, zoneID string) error
}
