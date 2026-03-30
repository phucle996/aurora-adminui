package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type UserRepository interface {
	ListUsers(ctx context.Context) ([]entity.User, error)
}
