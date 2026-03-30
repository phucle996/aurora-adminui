package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"

	"github.com/google/uuid"
)

type ZoneRepository interface {
	ListZones(ctx context.Context) ([]entity.Zone, error)
	CreateZone(ctx context.Context, zone *entity.Zone) error
	GetZoneByID(ctx context.Context, id uuid.UUID) (*entity.Zone, error)
	CountZoneObjects(ctx context.Context, zoneID uuid.UUID) (int, error)
	DeleteZone(ctx context.Context, id uuid.UUID) error
}
