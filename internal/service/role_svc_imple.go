package service

import (
	"context"
	"strings"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"

	"github.com/google/uuid"
)

type RoleServiceImple struct {
	repo domainrepo.RoleRepository
}

func NewRoleService(repo domainrepo.RoleRepository) domainsvc.RoleService {
	return &RoleServiceImple{repo: repo}
}

func (s *RoleServiceImple) ListRoles(ctx context.Context) ([]entity.Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *RoleServiceImple) ListPermissions(ctx context.Context) ([]entity.Permission, error) {
	return s.repo.ListPermissions(ctx)
}

func (s *RoleServiceImple) CreateRole(ctx context.Context, name, description string, permissionIDs []string) (*entity.Role, error) {
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)
	if name == "" {
		return nil, errorx.ErrInvalidArgument
	}

	seen := make(map[string]struct{}, len(permissionIDs))
	dedupedIDs := make([]string, 0, len(permissionIDs))
	for _, rawID := range permissionIDs {
		id := strings.TrimSpace(rawID)
		if id == "" {
			continue
		}
		if _, err := uuid.Parse(id); err != nil {
			return nil, errorx.ErrInvalidArgument
		}
		if _, exists := seen[id]; exists {
			continue
		}
		seen[id] = struct{}{}
		dedupedIDs = append(dedupedIDs, id)
	}

	role := &entity.Role{
		ID:              uuid.NewString(),
		Name:            name,
		Scope:           "global",
		Description:     description,
		UserCount:       0,
		PermissionCount: len(dedupedIDs),
	}
	if err := s.repo.CreateRole(ctx, role, dedupedIDs); err != nil {
		return nil, err
	}
	return role, nil
}
