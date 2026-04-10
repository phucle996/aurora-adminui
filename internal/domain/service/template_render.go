package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type TemplateRenderService interface {
	ListTemplateRenderCatalog(ctx context.Context) ([]entity.TemplateRender, error)
	GetTemplateRender(ctx context.Context, id string) (*entity.TemplateRender, error)
	CreateTemplateRender(ctx context.Context, input CreateTemplateRenderInput) (*entity.TemplateRender, error)
	UpdateTemplateRender(ctx context.Context, id string, input UpdateTemplateRenderInput) (*entity.TemplateRender, error)
	DeleteTemplateRender(ctx context.Context, id string) error
}

type CreateTemplateRenderInput struct {
	ResourceDefinitionID string
	Name                 string
	Description          string
	StreamKey            string
	ConsumerGroup        string
	YAMLTemplate         string
}

type UpdateTemplateRenderInput struct {
	ResourceDefinitionID string
	Name                 string
	Description          string
	StreamKey            string
	ConsumerGroup        string
	YAMLTemplate         string
}
