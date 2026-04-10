package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type RDRepoInterface interface {
	ListResourceDefinitions(ctx context.Context) ([]entity.ResourceDefinition, error)
	ListResourceDefinitionCatalog(ctx context.Context) ([]entity.ResourceDefinitionCatalogItem, error)
	CreateResourceDefinition(ctx context.Context, item *entity.ResourceDefinition) error
	ListResourceDefinitionZoneSupport(ctx context.Context, definitionID string) ([]entity.ResourceDefinitionZoneSupport, error)
	ReplaceResourceDefinitionZoneSupport(ctx context.Context, definitionID string, zoneIDs []string) error
	UpdateResourceDefinitionStatus(ctx context.Context, definitionID, status string) (*entity.ResourceDefinition, error)
	DeleteResourceDefinition(ctx context.Context, definitionID string) error
}
