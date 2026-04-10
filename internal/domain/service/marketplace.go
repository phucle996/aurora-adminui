package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type MarketplaceService interface {
	ListMarketplaceModelOptions(ctx context.Context) ([]entity.MarketplaceApp, error)
	ListMarketplaceTemplateOptions(ctx context.Context) ([]entity.MarketplaceTemplate, error)
	ListMarketplaceApps(ctx context.Context) ([]entity.MarketplaceApp, error)
	GetMarketplaceApp(ctx context.Context, id string) (*entity.MarketplaceApp, error)
	CreateMarketplaceApp(ctx context.Context, input CreateMarketplaceAppInput) (*entity.MarketplaceApp, error)
}

type CreateMarketplaceAppInput struct {
	Name                 string
	Slug                 string
	Summary              string
	Description          string
	ResourceDefinitionID string
	TemplateID           string
}
