package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type PlanService interface {
	ListPlans(ctx context.Context) ([]entity.Plan, error)
	CreatePlan(ctx context.Context, resourceType, code, name, description string, vcpu, ramGB, diskGB int) (*entity.Plan, error)
}
