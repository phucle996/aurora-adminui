package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type RoleRepository interface {
	ListRoles(ctx context.Context) ([]entity.Role, error)
	ListPermissions(ctx context.Context) ([]entity.Permission, error)
	CreateRole(ctx context.Context, role *entity.Role, permissionIDs []string) error
}
