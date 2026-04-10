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

type MarketplaceServiceImple struct {
	repo domainrepo.MarketplaceRepository
}

func NewMarketplaceService(repo domainrepo.MarketplaceRepository) domainsvc.MarketplaceService {
	return &MarketplaceServiceImple{repo: repo}
}

func (s *MarketplaceServiceImple) ListMarketplaceModelOptions(ctx context.Context) ([]entity.MarketplaceApp, error) {
	return s.repo.ListMarketplaceModelOptions(ctx)
}

func (s *MarketplaceServiceImple) ListMarketplaceTemplateOptions(ctx context.Context) ([]entity.MarketplaceTemplate, error) {
	return s.repo.ListMarketplaceTemplateOptions(ctx)
}

func (s *MarketplaceServiceImple) ListMarketplaceApps(ctx context.Context) ([]entity.MarketplaceApp, error) {
	return s.repo.ListMarketplaceApps(ctx)
}

func (s *MarketplaceServiceImple) GetMarketplaceApp(ctx context.Context, id string) (*entity.MarketplaceApp, error) {
	id = strings.TrimSpace(id)
	if _, err := uuid.Parse(id); err != nil {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.GetMarketplaceApp(ctx, id)
}

func (s *MarketplaceServiceImple) CreateMarketplaceApp(ctx context.Context, input domainsvc.CreateMarketplaceAppInput) (*entity.MarketplaceApp, error) {
	item := &entity.MarketplaceApp{
		ID:                   uuid.NewString(),
		Name:                 strings.TrimSpace(input.Name),
		Slug:                 strings.TrimSpace(input.Slug),
		Summary:              strings.TrimSpace(input.Summary),
		Description:          strings.TrimSpace(input.Description),
		ResourceDefinitionID: strings.TrimSpace(input.ResourceDefinitionID),
		TemplateID:           strings.TrimSpace(input.TemplateID),
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}
	if err := validateMarketplaceApp(item); err != nil {
		return nil, err
	}
	if err := s.repo.CreateMarketplaceApp(ctx, item); err != nil {
		return nil, err
	}
	return s.repo.GetMarketplaceApp(ctx, item.ID)
}

func validateMarketplaceApp(item *entity.MarketplaceApp) error {
	if item == nil {
		return errorx.ErrInvalidArgument
	}
	if item.Name == "" || item.Slug == "" || item.ResourceDefinitionID == "" || item.TemplateID == "" {
		return errorx.ErrInvalidArgument
	}
	if _, err := uuid.Parse(item.ResourceDefinitionID); err != nil {
		return errorx.ErrInvalidArgument
	}
	if _, err := uuid.Parse(item.TemplateID); err != nil {
		return errorx.ErrInvalidArgument
	}
	return nil
}
