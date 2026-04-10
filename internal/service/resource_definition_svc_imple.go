package service

import (
	"context"
	"strings"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
)

type RDSvcImple struct {
	repo domainrepo.RDRepoInterface
}

func NewResourceDefinitionService(repo domainrepo.RDRepoInterface) domainsvc.RDService {
	return &RDSvcImple{repo: repo}
}

func (s *RDSvcImple) ListResourceDefinitions(ctx context.Context) ([]entity.ResourceDefinition, error) {
	return s.repo.ListResourceDefinitions(ctx)
}

func (s *RDSvcImple) ListResourceDefinitionCatalog(ctx context.Context) ([]entity.ResourceDefinitionCatalogItem, error) {
	return s.repo.ListResourceDefinitionCatalog(ctx)
}

func (s *RDSvcImple) CreateResourceDefinition(ctx context.Context, item *entity.ResourceDefinition) error {
	if item == nil {
		return errorx.ErrInvalidArgument
	}
	next := &entity.ResourceDefinition{
		ResourceType:    strings.TrimSpace(item.ResourceType),
		ResourceModel:   strings.TrimSpace(item.ResourceModel),
		ResourceVersion: strings.TrimSpace(item.ResourceVersion),
		DisplayName:     strings.TrimSpace(item.DisplayName),
	}
	if next.ResourceType == "" || next.ResourceModel == "" || next.ResourceVersion == "" || next.DisplayName == "" {
		return errorx.ErrInvalidArgument
	}
	return s.repo.CreateResourceDefinition(ctx, next)
}

func (s *RDSvcImple) ListResourceDefinitionZoneSupport(ctx context.Context, definitionID string) ([]entity.ResourceDefinitionZoneSupport, error) {
	definitionID = strings.TrimSpace(definitionID)
	if definitionID == "" {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.ListResourceDefinitionZoneSupport(ctx, definitionID)
}

func (s *RDSvcImple) ReplaceResourceDefinitionZoneSupport(ctx context.Context, definitionID string, zoneIDs []string) error {
	definitionID = strings.TrimSpace(definitionID)
	if definitionID == "" {
		return errorx.ErrInvalidArgument
	}
	normalized := make([]string, 0, len(zoneIDs))
	for _, rawID := range zoneIDs {
		value := strings.TrimSpace(rawID)
		if value == "" {
			continue
		}
		normalized = append(normalized, value)
	}
	return s.repo.ReplaceResourceDefinitionZoneSupport(ctx, definitionID, normalized)
}

func (s *RDSvcImple) UpdateResourceDefinitionStatus(ctx context.Context, definitionID, status string) (*entity.ResourceDefinition, error) {
	definitionID = strings.TrimSpace(definitionID)
	nextStatus := strings.TrimSpace(strings.ToLower(status))
	if definitionID == "" {
		return nil, errorx.ErrInvalidArgument
	}
	switch nextStatus {
	case "draft", "ready", "maintain", "disabled":
	default:
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.UpdateResourceDefinitionStatus(ctx, definitionID, nextStatus)
}

func (s *RDSvcImple) DeleteResourceDefinition(ctx context.Context, definitionID string) error {
	definitionID = strings.TrimSpace(definitionID)
	if definitionID == "" {
		return errorx.ErrInvalidArgument
	}
	return s.repo.DeleteResourceDefinition(ctx, definitionID)
}
