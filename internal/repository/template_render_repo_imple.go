package repository

import (
	"context"
	"strings"
	"time"

	"aurora-adminui/internal/domain/entity"
	domainrepo "aurora-adminui/internal/domain/repository"
	controlplanegrpc "aurora-adminui/internal/transport/grpc"
)

type TemplateRenderRepoImple struct {
	client *controlplanegrpc.Client
}

// NewTemplateRenderRepo wires template render storage to the controlplane gRPC catalog.
func NewTemplateRenderRepo(client *controlplanegrpc.Client) domainrepo.TemplateRenderRepository {
	return &TemplateRenderRepoImple{client: client}
}

// ListTemplateRenderCatalog fetches the slim template rows used by the list page.
func (r *TemplateRenderRepoImple) ListTemplateRenderCatalog(ctx context.Context) ([]entity.TemplateRender, error) {
	items := make([]entity.TemplateRender, 0)
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.ListResourceTemplateCatalog(ctx, &controlplanegrpc.ListResourceTemplateCatalogRequest{})
		if err != nil {
			return err
		}
		items = make([]entity.TemplateRender, 0, len(resp.Items))
		for _, item := range resp.Items {
			items = append(items, entity.TemplateRender{
				ID:            item.Id,
				Name:          item.Name,
				Description:   item.Description,
				ResourceType:  item.ResourceType,
				ResourceModel: item.ResourceModel,
				StreamKey:     item.StreamKey,
				ConsumerGroup: item.ConsumerGroup,
				YAMLValid:     item.YamlValid,
				UpdatedAt:     parseTemplateTime(item.UpdatedAt),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetTemplateRender fetches a single template by id.
func (r *TemplateRenderRepoImple) GetTemplateRender(ctx context.Context, id string) (*entity.TemplateRender, error) {
	var out *entity.TemplateRender
	err := r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.GetResourceTemplate(ctx, &controlplanegrpc.GetResourceTemplateRequest{
			TemplateId: id,
		})
		if err != nil {
			return err
		}
		if resp == nil || strings.TrimSpace(resp.Id) == "" {
			out = nil
			return nil
		}
		item := mapTemplateRender(resp)
		out = &item
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CreateTemplateRender sends a create request to controlplane and refreshes the local entity with the result.
func (r *TemplateRenderRepoImple) CreateTemplateRender(ctx context.Context, item *entity.TemplateRender) error {
	return r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resp, err := client.CreateResourceTemplate(ctx, &controlplanegrpc.CreateResourceTemplateRequest{
			ResourceDefinitionId: item.ResourceDefinitionID,
			Name:                 item.Name,
			Description:          item.Description,
			StreamKey:            item.StreamKey,
			ConsumerGroup:        item.ConsumerGroup,
			TemplateBody:         item.YAMLTemplate,
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return nil
		}
		*item = mapTemplateRender(resp)
		return nil
	})
}

// UpdateTemplateRender sends a patch request to controlplane and refreshes the local entity with the result.
func (r *TemplateRenderRepoImple) UpdateTemplateRender(ctx context.Context, item *entity.TemplateRender) error {
	return r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		resourceDefinitionID := item.ResourceDefinitionID
		name := item.Name
		description := item.Description
		streamKey := item.StreamKey
		consumerGroup := item.ConsumerGroup
		templateBody := item.YAMLTemplate

		resp, err := client.UpdateResourceTemplate(ctx, &controlplanegrpc.UpdateResourceTemplateRequest{
			TemplateId:           item.ID,
			ResourceDefinitionId: &resourceDefinitionID,
			Name:                 &name,
			Description:          &description,
			StreamKey:            &streamKey,
			ConsumerGroup:        &consumerGroup,
			TemplateBody:         &templateBody,
		})
		if err != nil {
			return err
		}
		if resp == nil {
			return nil
		}
		*item = mapTemplateRender(resp)
		return nil
	})
}

// DeleteTemplateRender deletes one template by id through the controlplane catalog gRPC service.
func (r *TemplateRenderRepoImple) DeleteTemplateRender(ctx context.Context, id string) error {
	return r.client.Invoke(ctx, func(client controlplanegrpc.AdminCatalogServiceClient) error {
		_, err := client.DeleteResourceTemplate(ctx, &controlplanegrpc.DeleteResourceTemplateRequest{
			TemplateId: id,
		})
		return err
	})
}

// mapTemplateRender translates the gRPC wire shape into the adminui entity model.
func mapTemplateRender(item *controlplanegrpc.AdminResourceTemplate) entity.TemplateRender {
	if item == nil {
		return entity.TemplateRender{}
	}
	return entity.TemplateRender{
		ID:                   item.Id,
		Name:                 item.Name,
		Description:          item.Description,
		ResourceDefinitionID: item.ResourceDefinitionId,
		ResourceType:         item.ResourceType,
		ModelName:            item.ModelName,
		ResourceVersion:      item.ResourceVersion,
		ResourceModel:        item.ResourceModel,
		StreamKey:            item.StreamKey,
		ConsumerGroup:        item.ConsumerGroup,
		YAMLTemplate:         item.TemplateBody,
		YAMLValid:            item.YamlValid,
		CreatedAt:            parseTemplateTime(item.CreatedAt),
		UpdatedAt:            parseTemplateTime(item.UpdatedAt),
	}
}

// parseTemplateTime keeps catalog parsing resilient if an upstream timestamp is absent or malformed.
func parseTemplateTime(value string) time.Time {
	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}
	}
	return parsed
}
