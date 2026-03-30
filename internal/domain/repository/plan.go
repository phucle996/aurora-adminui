package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type PlanRepository interface {
	ListPlans(ctx context.Context) ([]entity.Plan, error)
	CreateVPSPlan(ctx context.Context, item *entity.Plan) error
}
