package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type RoleService interface {
	ListRoles(ctx context.Context) ([]entity.Role, error)
	ListPermissions(ctx context.Context) ([]entity.Permission, error)
	CreateRole(ctx context.Context, name, description string, permissionIDs []string) (*entity.Role, error)
}
