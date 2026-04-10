package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type MarketplaceRepository interface {
	ListMarketplaceModelOptions(ctx context.Context) ([]entity.MarketplaceApp, error)
	ListMarketplaceTemplateOptions(ctx context.Context) ([]entity.MarketplaceTemplate, error)
	ListMarketplaceApps(ctx context.Context) ([]entity.MarketplaceApp, error)
	GetMarketplaceApp(ctx context.Context, id string) (*entity.MarketplaceApp, error)
	CreateMarketplaceApp(ctx context.Context, item *entity.MarketplaceApp) error
}
