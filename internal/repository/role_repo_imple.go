package repository

import (
	"context"
	"errors"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepoImple struct {
	db *pgxpool.Pool
}

func NewRoleRepo(db *pgxpool.Pool) domainrepo.RoleRepository {
	return &RoleRepoImple{db: db}
}

func (r *RoleRepoImple) ListRoles(ctx context.Context) ([]entity.Role, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
		    rl.id::text,
		    rl.name,
		    COALESCE(rl.scope::text, 'global'),
		    COALESCE(rl.description, ''),
		    COUNT(DISTINCT ur.user_id) AS user_count,
		    COUNT(DISTINCT rp.permission_id) AS permission_count
		  FROM iam.roles rl
		  LEFT JOIN iam.user_roles ur ON ur.role_id = rl.id
		  LEFT JOIN iam.role_permissions rp ON rp.role_id = rl.id
		  GROUP BY rl.id, rl.name, rl.scope, rl.description
		  ORDER BY rl.name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Role, 0)
	for rows.Next() {
		var item entity.Role
		if err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Scope,
			&item.Description,
			&item.UserCount,
			&item.PermissionCount,
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

func (r *RoleRepoImple) ListPermissions(ctx context.Context) ([]entity.Permission, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT
		    p.id::text,
		    p.name,
		    COALESCE(p.description, '')
		  FROM iam.permissions p
		  ORDER BY p.name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.Permission, 0)
	for rows.Next() {
		var item entity.Permission
		if err := rows.Scan(&item.ID, &item.Name, &item.Description); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *RoleRepoImple) CreateRole(ctx context.Context, role *entity.Role, permissionIDs []string) error {
	if role == nil {
		return errors.New("role is nil")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	now := time.Now().UTC()
	_, err = tx.Exec(
		ctx,
		`INSERT INTO iam.roles (id, name, scope, description, created_at)
		 VALUES ($1, $2, 'global', $3, $4)`,
		uuid.MustParse(role.ID),
		role.Name,
		role.Description,
		now,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errorx.ErrRoleAlreadyExists
		}
		return err
	}

	if len(permissionIDs) > 0 {
		parsedIDs := make([]uuid.UUID, 0, len(permissionIDs))
		for _, rawID := range permissionIDs {
			parsedIDs = append(parsedIDs, uuid.MustParse(rawID))
		}

		var count int
		if err := tx.QueryRow(
			ctx,
			`SELECT COUNT(1)
			   FROM iam.permissions
			  WHERE id = ANY($1::uuid[])`,
			parsedIDs,
		).Scan(&count); err != nil {
			return err
		}
		if count != len(parsedIDs) {
			return errorx.ErrPermissionNotFound
		}

		_, err = tx.Exec(
			ctx,
			`INSERT INTO iam.role_permissions (id, role_id, permission_id, created_at)
			 SELECT gen_random_uuid(), $1, permission_id, $3
			   FROM UNNEST($2::uuid[]) AS permission_id`,
			uuid.MustParse(role.ID),
			parsedIDs,
			now,
		)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
