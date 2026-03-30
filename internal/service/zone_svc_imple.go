package service

import (
	"context"
	"strings"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
)

type ZoneSvcImple struct {
	repo domainrepo.ZoneRepository
}

func NewZoneService(repo domainrepo.ZoneRepository) domainsvc.ZoneService {
	return &ZoneSvcImple{repo: repo}
}

func (s *ZoneSvcImple) ListZones(ctx context.Context) ([]entity.Zone, error) {
	return s.repo.ListZones(ctx)
}

func (s *ZoneSvcImple) CreateZone(ctx context.Context, name, description string) (*entity.Zone, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return nil, errorx.ErrInvalidArgument
	}

	zone := &entity.Zone{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.repo.CreateZone(ctx, zone); err != nil {
		return nil, err
	}
	return zone, nil
}

func (s *ZoneSvcImple) DeleteZone(ctx context.Context, rawID string) error {
	id, err := uuid.Parse(strings.TrimSpace(rawID))
	if err != nil {
		return errorx.ErrInvalidArgument
	}
	zone, err := s.repo.GetZoneByID(ctx, id)
	if err != nil {
		return err
	}
	count, err := s.repo.CountZoneObjects(ctx, zone.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errorx.ErrZoneHasResources
	}
	return s.repo.DeleteZone(ctx, id)
}
