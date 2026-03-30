package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"

	"github.com/google/uuid"
)

type HypervisorRepository interface {
	ListNodes(ctx context.Context) ([]entity.HypervisorNode, error)
	GetNodeDetail(ctx context.Context, nodeID string) (*entity.HypervisorNodeDetail, error)
	UpdateNodeName(ctx context.Context, nodeID, name string) error
	AssignNodeToZone(ctx context.Context, nodeID string, zoneID uuid.UUID) error
}
