package repository

import (
	"context"
	"strings"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepoImple struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) domainrepo.UserRepository {
	return &UserRepoImple{db: db}
}

func (r *UserRepoImple) ListUsers(ctx context.Context) ([]entity.User, error) {
	rows, err := r.db.Query(
		ctx,
		`WITH primary_roles AS (
		    SELECT
		      ur.user_id,
		      rl.name,
		      ROW_NUMBER() OVER (
		        PARTITION BY ur.user_id
		        ORDER BY
		          CASE rl.name
		            WHEN 'root' THEN 0
		            WHEN 'admin' THEN 1
		            WHEN 'user' THEN 2
		            ELSE 9
		          END,
		          rl.name ASC
		      ) AS rn
		    FROM iam.user_roles ur
		    JOIN iam.roles rl ON rl.id = ur.role_id
		  )
		  SELECT
		    u.id::text,
		    COALESCE(p.full_name, ''),
		    u.username,
		    u.email,
		    COALESCE(p.phone, ''),
		    COALESCE(u.status::text, 'pending'),
		    COALESCE(pr.name, 'user'),
		    COALESCE(u.created_at, now()),
		    COALESCE(u.updated_at, now())
		  FROM iam.users u
		  LEFT JOIN iam.profiles p ON p.user_id = u.id
		  LEFT JOIN primary_roles pr ON pr.user_id = u.id AND pr.rn = 1
		  ORDER BY COALESCE(u.updated_at, u.created_at, now()) DESC, u.username ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entity.User, 0)
	for rows.Next() {
		var (
			item     entity.User
			fullName string
		)
		if err := rows.Scan(
			&item.ID,
			&fullName,
			&item.Username,
			&item.Email,
			&item.PhoneNumber,
			&item.Status,
			&item.Role,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.FirstName, item.LastName = splitFullName(fullName, item.Username)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func splitFullName(fullName, fallback string) (string, string) {
	trimmed := strings.TrimSpace(fullName)
	if trimmed == "" {
		name := strings.TrimSpace(fallback)
		if name == "" {
			return "Unknown", ""
		}
		return name, ""
	}
	parts := strings.Fields(trimmed)
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], strings.Join(parts[1:], " ")
}
