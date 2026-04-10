package service

import (
	"context"
	"strings"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	domainsvc "aurora-adminui/internal/domain/service"
	"aurora-adminui/internal/errorx"
)

type K8sSvcImple struct {
	repo domainrepo.K8sRepository
}

// NewK8sService builds the kubernetes service as a thin adminui -> controlplane orchestrator.
func NewK8sService(repo domainrepo.K8sRepository) domainsvc.K8sService {
	return &K8sSvcImple{repo: repo}
}

func (s *K8sSvcImple) ListClusters(ctx context.Context) ([]entity.K8sCluster, error) {
	return s.repo.ListClusters(ctx)
}

func (s *K8sSvcImple) GetClusterDetail(ctx context.Context, id string) (*entity.K8sCluster, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.GetClusterByID(ctx, id)
}

func (s *K8sSvcImple) ListClusterNodes(ctx context.Context, id string) ([]entity.K8sClusterNode, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.ListClusterNodes(ctx, id)
}

func (s *K8sSvcImple) CreateCluster(ctx context.Context, input entity.K8sClusterCreateInput) (*entity.K8sCluster, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Description = strings.TrimSpace(input.Description)
	input.ZoneID = strings.TrimSpace(input.ZoneID)
	if input.Name == "" || len(input.Kubeconfig) == 0 {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.CreateCluster(ctx, input)
}

func (s *K8sSvcImple) UpdateCluster(ctx context.Context, id string, input entity.K8sClusterUpdateInput) (*entity.K8sCluster, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errorx.ErrInvalidArgument
	}
	input.ZoneID = strings.TrimSpace(input.ZoneID)
	if input.ZoneID == "" && len(input.Kubeconfig) == 0 {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.UpdateCluster(ctx, id, input)
}

func (s *K8sSvcImple) RevalidateCluster(ctx context.Context, id string) (*entity.K8sCluster, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errorx.ErrInvalidArgument
	}
	return s.repo.RevalidateCluster(ctx, id)
}

func (s *K8sSvcImple) DeleteCluster(ctx context.Context, id string) error {
	id = strings.TrimSpace(id)
	if id == "" {
		return errorx.ErrInvalidArgument
	}
	return s.repo.DeleteCluster(ctx, id)
}
