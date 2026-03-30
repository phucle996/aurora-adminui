package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type K8sRepository interface {
	ListClusters(ctx context.Context) ([]entity.K8sCluster, error)
	GetClusterByID(ctx context.Context, id string) (*entity.K8sCluster, error)
	CreateCluster(ctx context.Context, cluster *entity.K8sCluster) error
	UpdateClusterValidation(ctx context.Context, cluster *entity.K8sCluster) error
	DeleteCluster(ctx context.Context, id string) error
}
