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

type PlanSvcImple struct {
	repo domainrepo.PlanRepository
}

func NewPlanService(repo domainrepo.PlanRepository) domainsvc.PlanService {
	return &PlanSvcImple{repo: repo}
}

func (s *PlanSvcImple) ListPlans(ctx context.Context) ([]entity.Plan, error) {
	return s.repo.ListPlans(ctx)
}

func (s *PlanSvcImple) CreatePlan(ctx context.Context, resourceType, code, name, description string, vcpu, ramGB, diskGB int) (*entity.Plan, error) {
	resourceType = strings.ToLower(strings.TrimSpace(resourceType))
	code = strings.TrimSpace(code)
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	if resourceType == "" || code == "" || name == "" {
		return nil, errorx.ErrInvalidArgument
	}
	if vcpu <= 0 || ramGB <= 0 || diskGB <= 0 {
		return nil, errorx.ErrInvalidArgument
	}

	switch entity.ResourceType(resourceType) {
	case entity.ResourceTypeVPS:
	default:
		return nil, errorx.ErrInvalidArgument
	}

	item := &entity.Plan{
		ID:           uuid.New(),
		ResourceType: entity.ResourceType(resourceType),
		Code:         code,
		Name:         name,
		Description:  description,
		Status:       entity.PlanStatusActive,
		VCPU:         vcpu,
		RAMGB:        ramGB,
		DiskGB:       diskGB,
		CreatedAt:    time.Now().UTC(),
	}
	if err := s.repo.CreateVPSPlan(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}
