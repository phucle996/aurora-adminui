package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	domainsvc "aurora-adminui/internal/domain/service"
)

type UserSvcImple struct {
	repo domainrepo.UserRepository
}

func NewUserService(repo domainrepo.UserRepository) domainsvc.UserService {
	return &UserSvcImple{repo: repo}
}

func (s *UserSvcImple) ListUsers(ctx context.Context) ([]entity.User, error) {
	return s.repo.ListUsers(ctx)
}
