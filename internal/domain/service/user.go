package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type UserService interface {
	ListUsers(ctx context.Context) ([]entity.User, error)
}
