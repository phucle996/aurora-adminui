package repository

import (
	"context"
	"errors"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PlanRepoImple struct {
	db *pgxpool.Pool
}

func NewPlanRepo(db *pgxpool.Pool) domainrepo.PlanRepository {
	return &PlanRepoImple{db: db}
}

func (r *PlanRepoImple) ListPlans(ctx context.Context) ([]entity.Plan, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
		    rp.id,
		    rp.resource_type,
		    COALESCE(rp.resource_model, ''),
		    rp.code,
		    rp.name,
		    COALESCE(rp.description, ''),
		    rp.status,
		    COALESCE(vp.vcpu, 0),
		    COALESCE(vp.ram_gb, 0),
		    COALESCE(vp.disk_gb, 0),
		    rp.created_at
		  FROM plan.resource_packages rp
		  LEFT JOIN plan.vps_packages vp ON vp.package_id = rp.id
		  ORDER BY rp.created_at DESC, rp.name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Plan, 0)
	for rows.Next() {
		var item entity.Plan
		if err := rows.Scan(
			&item.ID,
			&item.ResourceType,
			&item.ResourceModel,
			&item.Code,
			&item.Name,
			&item.Description,
			&item.Status,
			&item.VCPU,
			&item.RAMGB,
			&item.DiskGB,
			&item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PlanRepoImple) CreateVPSPlan(ctx context.Context, item *entity.Plan) error {
	if item == nil {
		return errors.New("plan is nil")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	_, err = tx.Exec(
		ctx,
		`INSERT INTO plan.resource_packages (
		    id, resource_type, resource_model, code, name, description, status, created_at, retired_at
		  ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NULL)`,
		item.ID,
		item.ResourceType,
		item.ResourceModel,
		item.Code,
		item.Name,
		item.Description,
		item.Status,
		item.CreatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errorx.ErrPlanAlreadyExists
		}
		return err
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO plan.vps_packages (
		    package_id, vcpu, ram_gb, disk_gb
		  ) VALUES ($1, $2, $3, $4)`,
		item.ID,
		item.VCPU,
		item.RAMGB,
		item.DiskGB,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
