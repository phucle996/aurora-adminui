package repository

import (
	"context"
	"errors"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ZoneRepoImple struct {
	db *pgxpool.Pool
}

func NewZoneRepo(db *pgxpool.Pool) domainrepo.ZoneRepository {
	return &ZoneRepoImple{db: db}
}

func (r *ZoneRepoImple) ListZones(ctx context.Context) ([]entity.Zone, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT z.id,
		        z.name,
		        z.description,
		        COALESCE(COUNT(zo.object_id), 0) AS resource_count,
		        z.created_at
		   FROM zone.zones z
		   LEFT JOIN zone.zone_objects zo ON zo.zone_id = z.id
		  GROUP BY z.id, z.name, z.description, z.created_at
		  ORDER BY z.name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Zone, 0)
	for rows.Next() {
		var item entity.Zone
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.ResourceCount,
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

func (r *ZoneRepoImple) CreateZone(ctx context.Context, zone *entity.Zone) error {
	if zone == nil {
		return errors.New("zone is nil")
	}
	_, err := r.db.Exec(
		ctx,
		`INSERT INTO zone.zones (id, name, description, created_at)
		 VALUES ($1, $2, $3, $4)`,
		zone.ID,
		zone.Name,
		zone.Description,
		zone.CreatedAt,
	)
	if err == nil {
		return nil
	}
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return errorx.ErrZoneAlreadyExists
	}
	return err
}

func (r *ZoneRepoImple) GetZoneByID(ctx context.Context, id uuid.UUID) (*entity.Zone, error) {
	var item entity.Zone
	err := r.db.QueryRow(
		ctx,
		`SELECT id, name, description, created_at
		   FROM zone.zones
		  WHERE id = $1`,
		id,
	).Scan(&item.ID, &item.Name, &item.Description, &item.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errorx.ErrZoneNotFound
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ZoneRepoImple) CountZoneObjects(ctx context.Context, zoneID uuid.UUID) (int, error) {
	var count int
	if err := r.db.QueryRow(
		ctx,
		`SELECT COUNT(1)
		   FROM zone.zone_objects
		  WHERE zone_id = $1`,
		zoneID,
	).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ZoneRepoImple) DeleteZone(ctx context.Context, id uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM zone.zones WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errorx.ErrZoneNotFound
	}
	return nil
}
