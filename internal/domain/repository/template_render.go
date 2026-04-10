package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type TemplateRenderRepository interface {
	ListTemplateRenderCatalog(ctx context.Context) ([]entity.TemplateRender, error)
	GetTemplateRender(ctx context.Context, id string) (*entity.TemplateRender, error)
	CreateTemplateRender(ctx context.Context, item *entity.TemplateRender) error
	UpdateTemplateRender(ctx context.Context, item *entity.TemplateRender) error
	DeleteTemplateRender(ctx context.Context, id string) error
}
