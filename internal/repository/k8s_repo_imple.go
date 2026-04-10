package repository

import (
	"context"
	"strings"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"
	controlplanegrpc "aurora-adminui/internal/transport/grpc"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type K8sRepoImple struct {
	client *controlplanegrpc.Client
}

func NewK8sRepo(client *controlplanegrpc.Client) domainrepo.K8sRepository {
	return &K8sRepoImple{client: client}
}

func (r *K8sRepoImple) ListClusters(ctx context.Context) ([]entity.K8sCluster, error) {
	items := make([]entity.K8sCluster, 0)
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.ListK8SClusters(ctx, &controlplanegrpc.ListK8SClustersRequest{})
		if err != nil {
			return err
		}
		items = make([]entity.K8sCluster, 0, len(resp.Items))
		for _, item := range resp.Items {
			items = append(items, entity.K8sCluster{
				ID:                parseUUID(item.GetId()),
				Name:              item.GetName(),
				Description:       item.GetDescription(),
				APIServerURL:      item.GetApiServerUrl(),
				KubernetesVersion: item.GetKubernetesVersion(),
				ValidationStatus:  entity.K8sClusterValidationStatus(item.GetValidationStatus()),
				LastValidatedAt:   parseOptionalTime(item.GetLastValidatedAt()),
				CreatedAt:         parseTime(item.GetCreatedAt()),
				ZoneName:          item.GetZoneName(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, mapK8sGRPCError(err)
	}
	return items, nil
}

func (r *K8sRepoImple) GetClusterByID(ctx context.Context, id string) (*entity.K8sCluster, error) {
	var item *entity.K8sCluster
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.GetK8SCluster(ctx, &controlplanegrpc.GetK8SClusterRequest{
			ClusterId: id,
		})
		if err != nil {
			return err
		}
		item = mapClusterDetail(resp)
		return nil
	})
	if err != nil {
		return nil, mapK8sGRPCError(err)
	}
	return item, nil
}

func (r *K8sRepoImple) ListClusterNodes(ctx context.Context, id string) ([]entity.K8sClusterNode, error) {
	items := make([]entity.K8sClusterNode, 0)
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.GetClusterNodesAdmin(ctx, &controlplanegrpc.GetClusterNodesAdminRequest{
			ClusterId: id,
		})
		if err != nil {
			return err
		}
		items = make([]entity.K8sClusterNode, 0, len(resp.Items))
		for _, item := range resp.Items {
			items = append(items, entity.K8sClusterNode{
				Name:             item.GetName(),
				Roles:            append([]string(nil), item.GetRoles()...),
				KubeletVersion:   item.GetKubeletVersion(),
				ContainerRuntime: item.GetContainerRuntime(),
				OSImage:          item.GetOsImage(),
				KernelVersion:    item.GetKernelVersion(),
				Ready:            item.GetReady(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, mapK8sGRPCError(err)
	}
	return items, nil
}

func (r *K8sRepoImple) CreateCluster(ctx context.Context, input entity.K8sClusterCreateInput) (*entity.K8sCluster, error) {
	var item *entity.K8sCluster
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.CreateK8SCluster(ctx, &controlplanegrpc.CreateK8SClusterRequest{
			Name:        input.Name,
			Description: input.Description,
			ZoneId:      input.ZoneID,
			Kubeconfig:  input.Kubeconfig,
		})
		if err != nil {
			return err
		}
		item = mapClusterDetail(resp)
		return nil
	})
	if err != nil {
		return nil, mapK8sGRPCError(err)
	}
	return item, nil
}

func (r *K8sRepoImple) UpdateCluster(ctx context.Context, id string, input entity.K8sClusterUpdateInput) (*entity.K8sCluster, error) {
	var item *entity.K8sCluster
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.UpdateK8SCluster(ctx, &controlplanegrpc.UpdateK8SClusterRequest{
			ClusterId:  id,
			ZoneId:     input.ZoneID,
			Kubeconfig: input.Kubeconfig,
		})
		if err != nil {
			return err
		}
		item = mapClusterDetail(resp)
		return nil
	})
	if err != nil {
		return nil, mapK8sGRPCError(err)
	}
	return item, nil
}

func (r *K8sRepoImple) RevalidateCluster(ctx context.Context, id string) (*entity.K8sCluster, error) {
	var item *entity.K8sCluster
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.RevalidateK8SCluster(ctx, &controlplanegrpc.RevalidateK8SClusterRequest{
			ClusterId: id,
		})
		if err != nil {
			return err
		}
		item = mapClusterDetail(resp)
		return nil
	})
	if err != nil {
		return nil, mapK8sGRPCError(err)
	}
	return item, nil
}

func (r *K8sRepoImple) DeleteCluster(ctx context.Context, id string) error {
	return mapK8sGRPCError(r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		_, err := client.DeleteK8SCluster(ctx, &controlplanegrpc.DeleteK8SClusterRequest{
			ClusterId: id,
		})
		return err
	}))
}

func mapClusterDetail(item *controlplanegrpc.AdminK8SClusterDetail) *entity.K8sCluster {
	if item == nil {
		return nil
	}
	cluster := &entity.K8sCluster{
		ID:                parseUUID(item.GetId()),
		Name:              item.GetName(),
		Description:       item.GetDescription(),
		APIServerURL:      item.GetApiServerUrl(),
		CurrentContext:    item.GetCurrentContext(),
		KubernetesVersion: item.GetKubernetesVersion(),
		ValidationStatus:  entity.K8sClusterValidationStatus(item.GetValidationStatus()),
		LastValidatedAt:   parseOptionalTime(item.GetLastValidatedAt()),
		CreatedAt:         parseTime(item.GetCreatedAt()),
		ZoneName:          item.GetZoneName(),
	}
	if id := parseUUID(item.GetZoneId()); id != uuid.Nil {
		cluster.ZoneID = &id
	}
	return cluster
}

func mapK8sGRPCError(err error) error {
	if err == nil {
		return nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return err
	}
	switch st.Code() {
	case codes.InvalidArgument:
		return errorx.ErrInvalidArgument
	case codes.NotFound:
		if st.Message() == errorx.ErrZoneNotFound.Error() {
			return errorx.ErrZoneNotFound
		}
		return errorx.ErrK8sClusterNotFound
	case codes.AlreadyExists:
		return errorx.ErrK8sClusterAlreadyExists
	case codes.FailedPrecondition:
		switch st.Message() {
		case errorx.ErrDataplaneValidationFailed.Error():
			return errorx.ErrDataplaneValidationFailed
		case errorx.ErrK8sClusterHasResources.Error():
			return errorx.ErrK8sClusterHasResources
		default:
			return errorx.ErrDataplaneValidationFailed
		}
	case codes.Unavailable:
		return errorx.ErrNoHealthyDataplane
	default:
		return err
	}
}

func parseOptionalTime(value string) *time.Time {
	value = trimSpace(value)
	if value == "" {
		return nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}
	return &parsed
}

func parseTime(value string) time.Time {
	value = trimSpace(value)
	if value == "" {
		return time.Time{}
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}
	}
	return parsed
}

func trimSpace(value string) string {
	return strings.TrimSpace(value)
}
