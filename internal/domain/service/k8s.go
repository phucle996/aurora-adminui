package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type K8sService interface {
	ListClusters(ctx context.Context) ([]entity.K8sCluster, error)
	GetClusterDetail(ctx context.Context, id string) (*entity.K8sCluster, error)
	CreateCluster(ctx context.Context, input entity.K8sClusterCreateInput) (*entity.K8sCluster, error)
	RevalidateCluster(ctx context.Context, id string) (*entity.K8sCluster, error)
	DeleteCluster(ctx context.Context, id string) error
}
