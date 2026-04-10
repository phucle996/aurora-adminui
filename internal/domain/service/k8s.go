package service

import (
	"context"

	"aurora-adminui/internal/domain/entity"
)

type K8sService interface {
	ListClusters(ctx context.Context) ([]entity.K8sCluster, error)
	GetClusterDetail(ctx context.Context, id string) (*entity.K8sCluster, error)
	ListClusterNodes(ctx context.Context, id string) ([]entity.K8sClusterNode, error)
	CreateCluster(ctx context.Context, input entity.K8sClusterCreateInput) (*entity.K8sCluster, error)
	UpdateCluster(ctx context.Context, id string, input entity.K8sClusterUpdateInput) (*entity.K8sCluster, error)
	RevalidateCluster(ctx context.Context, id string) (*entity.K8sCluster, error)
	DeleteCluster(ctx context.Context, id string) error
}
