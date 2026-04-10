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

const (
	defaultTemplateStreamKey     = "resourceplatform.jobs"
	defaultTemplateConsumerGroup = "resourceplatform-workers"
)

type TemplateRenderServiceImple struct {
	repo domainrepo.TemplateRenderRepository
}

// NewTemplateRenderService builds the business service used by template render pages.
func NewTemplateRenderService(repo domainrepo.TemplateRenderRepository) domainsvc.TemplateRenderService {
	return &TemplateRenderServiceImple{repo: repo}
}

// ListTemplateRenderCatalog returns the lightweight template list used by the catalog page.
func (s *TemplateRenderServiceImple) ListTemplateRenderCatalog(ctx context.Context) ([]entity.TemplateRender, error) {
	return s.repo.ListTemplateRenderCatalog(ctx)
}

// GetTemplateRender validates the id and fetches one template for detail or edit pages.
func (s *TemplateRenderServiceImple) GetTemplateRender(ctx context.Context, id string) (*entity.TemplateRender, error) {
	id = strings.TrimSpace(id)
	if _, err := uuid.Parse(id); err != nil {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.GetTemplateRender(ctx, id)
}

// CreateTemplateRender normalizes defaults and persists a new template render record.
func (s *TemplateRenderServiceImple) CreateTemplateRender(ctx context.Context, input domainsvc.CreateTemplateRenderInput) (*entity.TemplateRender, error) {
	item := &entity.TemplateRender{
		ID:                   uuid.NewString(),
		Name:                 strings.TrimSpace(input.Name),
		Description:          strings.TrimSpace(input.Description),
		ResourceDefinitionID: strings.TrimSpace(input.ResourceDefinitionID),
		StreamKey:            firstNonEmpty(strings.TrimSpace(input.StreamKey), defaultTemplateStreamKey),
		ConsumerGroup:        firstNonEmpty(strings.TrimSpace(input.ConsumerGroup), defaultTemplateConsumerGroup),
		YAMLTemplate:         strings.TrimSpace(input.YAMLTemplate),
		CreatedAt:            time.Now().UTC(),
		UpdatedAt:            time.Now().UTC(),
	}
	if err := validateTemplateRender(item); err != nil {
		return nil, err
	}
	if err := s.repo.CreateTemplateRender(ctx, item); err != nil {
		return nil, err
	}
	return s.repo.GetTemplateRender(ctx, item.ID)
}

// UpdateTemplateRender overwrites a template render after validating its identifiers and payload.
func (s *TemplateRenderServiceImple) UpdateTemplateRender(ctx context.Context, id string, input domainsvc.UpdateTemplateRenderInput) (*entity.TemplateRender, error) {
	id = strings.TrimSpace(id)
	if _, err := uuid.Parse(id); err != nil {
		return nil, errorx.ErrInvalidArgument
	}
	current, err := s.repo.GetTemplateRender(ctx, id)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, nil
	}
	current.Name = strings.TrimSpace(input.Name)
	current.Description = strings.TrimSpace(input.Description)
	current.ResourceDefinitionID = strings.TrimSpace(input.ResourceDefinitionID)
	current.StreamKey = firstNonEmpty(strings.TrimSpace(input.StreamKey), defaultTemplateStreamKey)
	current.ConsumerGroup = firstNonEmpty(strings.TrimSpace(input.ConsumerGroup), defaultTemplateConsumerGroup)
	current.YAMLTemplate = strings.TrimSpace(input.YAMLTemplate)
	current.UpdatedAt = time.Now().UTC()

	if err := validateTemplateRender(current); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateTemplateRender(ctx, current); err != nil {
		return nil, err
	}
	return s.repo.GetTemplateRender(ctx, current.ID)
}

// DeleteTemplateRender validates the id before deleting the template render.
func (s *TemplateRenderServiceImple) DeleteTemplateRender(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if _, err := uuid.Parse(id); err != nil {
		return errorx.ErrInvalidArgument
	}
	return s.repo.DeleteTemplateRender(ctx, id)
}

// validateTemplateRender enforces the minimal invariants shared by create and update.
func validateTemplateRender(item *entity.TemplateRender) error {
	if item == nil {
		return errorx.ErrInvalidArgument
	}
	if item.Name == "" || item.ResourceDefinitionID == "" || item.YAMLTemplate == "" {
		return errorx.ErrInvalidArgument
	}
	if _, err := uuid.Parse(item.ResourceDefinitionID); err != nil {
		return errorx.ErrInvalidArgument
	}
	return nil
}

// firstNonEmpty returns the first non-empty value in priority order.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
