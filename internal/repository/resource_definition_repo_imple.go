package repository

import (
	"context"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	"aurora-adminui/internal/errorx"
	controlplanegrpc "aurora-adminui/internal/transport/grpc"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResourceDefinitionRepoImple struct {
	client *controlplanegrpc.Client
}

func NewResourceDefinitionRepo(client *controlplanegrpc.Client) domainrepo.RDRepoInterface {
	return &ResourceDefinitionRepoImple{client: client}
}

func (r *ResourceDefinitionRepoImple) ListResourceDefinitions(ctx context.Context) ([]entity.ResourceDefinition, error) {
	items := make([]entity.ResourceDefinition, 0)
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.ListResourceDefinitions(ctx, &controlplanegrpc.ListResourceDefinitionsRequest{})
		if err != nil {
			return err
		}
		items = make([]entity.ResourceDefinition, 0, len(resp.Items))
		for _, item := range resp.Items {
			items = append(items, entity.ResourceDefinition{
				ID:              parseUUID(item.GetId()),
				ResourceType:    item.GetResourceType(),
				ResourceModel:   item.GetResourceModel(),
				ResourceVersion: item.GetResourceVersion(),
				DisplayName:     item.GetDisplayName(),
				Status:          item.GetStatus(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, mapRDGRPCError(err)
	}
	return items, nil
}

func (r *ResourceDefinitionRepoImple) ListResourceDefinitionCatalog(ctx context.Context) ([]entity.ResourceDefinitionCatalogItem, error) {
	items := make([]entity.ResourceDefinitionCatalogItem, 0)
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.ListResourceDefinitions(ctx, &controlplanegrpc.ListResourceDefinitionsRequest{})
		if err != nil {
			return err
		}
		items = make([]entity.ResourceDefinitionCatalogItem, 0, len(resp.Items))
		for _, item := range resp.Items {
			items = append(items, entity.ResourceDefinitionCatalogItem{
				ResourceDefinition: entity.ResourceDefinition{
					ID:              parseUUID(item.GetId()),
					ResourceType:    item.GetResourceType(),
					ResourceModel:   item.GetResourceModel(),
					ResourceVersion: item.GetResourceVersion(),
					DisplayName:     item.GetDisplayName(),
					Status:          item.GetStatus(),
				},
				ResourceCount: int(item.GetResourceCount()),
			})
		}
		return nil
	})
	if err != nil {
		return nil, mapRDGRPCError(err)
	}
	return items, nil
}

func (r *ResourceDefinitionRepoImple) CreateResourceDefinition(ctx context.Context, item *entity.ResourceDefinition) error {
	if item == nil {
		return errorx.ErrInvalidArgument
	}
	return mapRDGRPCError(r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		_, err := client.CreateResourceDefinition(ctx, &controlplanegrpc.CreateResourceDefinitionRequest{
			ResourceType:    item.ResourceType,
			ResourceModel:   item.ResourceModel,
			ResourceVersion: item.ResourceVersion,
			DisplayName:     item.DisplayName,
		})
		return err
	}))
}

func (r *ResourceDefinitionRepoImple) ListResourceDefinitionZoneSupport(ctx context.Context, definitionID string) ([]entity.ResourceDefinitionZoneSupport, error) {
	items := make([]entity.ResourceDefinitionZoneSupport, 0)
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.ListResourceDefinitionZoneSupport(ctx, &controlplanegrpc.ListResourceDefinitionZoneSupportRequest{
			DefinitionId: definitionID,
		})
		if err != nil {
			return err
		}
		items = make([]entity.ResourceDefinitionZoneSupport, 0, len(resp.Items))
		for _, item := range resp.Items {
			items = append(items, entity.ResourceDefinitionZoneSupport{
				ZoneID:   parseUUID(item.GetZoneId()),
				ZoneName: item.GetZoneName(),
				Enabled:  item.GetEnabled(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, mapRDGRPCError(err)
	}
	return items, nil
}

func (r *ResourceDefinitionRepoImple) ReplaceResourceDefinitionZoneSupport(ctx context.Context, definitionID string, zoneIDs []string) error {
	return mapRDGRPCError(r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		_, err := client.ReplaceResourceDefinitionZoneSupport(ctx, &controlplanegrpc.ReplaceResourceDefinitionZoneSupportRequest{
			DefinitionId: definitionID,
			ZoneIds:      zoneIDs,
		})
		return err
	}))
}

func (r *ResourceDefinitionRepoImple) UpdateResourceDefinitionStatus(ctx context.Context, definitionID, status string) (*entity.ResourceDefinition, error) {
	var out *entity.ResourceDefinition
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.UpdateResourceDefinitionStatus(ctx, &controlplanegrpc.UpdateResourceDefinitionStatusRequest{
			DefinitionId: definitionID,
			Status:       status,
		})
		if err != nil {
			return err
		}
		out = &entity.ResourceDefinition{
			ID:              parseUUID(resp.GetId()),
			ResourceType:    resp.GetResourceType(),
			ResourceModel:   resp.GetResourceModel(),
			ResourceVersion: resp.GetResourceVersion(),
			DisplayName:     resp.GetDisplayName(),
			Status:          resp.GetStatus(),
		}
		return nil
	})
	if err != nil {
		return nil, mapRDGRPCError(err)
	}
	return out, nil
}

func (r *ResourceDefinitionRepoImple) DeleteResourceDefinition(ctx context.Context, definitionID string) error {
	return mapRDGRPCError(r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		_, err := client.DeleteResourceDefinition(ctx, &controlplanegrpc.DeleteResourceDefinitionRequest{
			DefinitionId: definitionID,
		})
		return err
	}))
}

func mapRDGRPCError(err error) error {
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
	case codes.AlreadyExists:
		return errorx.ErrResourceDefinitionAlreadyExists
	case codes.NotFound:
		if st.Message() == errorx.ErrZoneNotFound.Error() {
			return errorx.ErrZoneNotFound
		}
		return errorx.ErrResourceDefinitionNotFound
	case codes.FailedPrecondition:
		switch st.Message() {
		case errorx.ErrResourceDefinitionNeedsZones.Error():
			return errorx.ErrResourceDefinitionNeedsZones
		default:
			return errorx.ErrResourceDefinitionHasResources
		}
	default:
		return err
	}
}

func parseUUID(value string) uuid.UUID {
	id, err := uuid.Parse(value)
	if err != nil {
		return uuid.UUID{}
	}
	return id
}
