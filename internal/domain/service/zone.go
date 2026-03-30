package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type ZoneService interface {
	ListZones(ctx context.Context) ([]entity.Zone, error)
	CreateZone(ctx context.Context, name, description string) (*entity.Zone, error)
	DeleteZone(ctx context.Context, rawID string) error
}
