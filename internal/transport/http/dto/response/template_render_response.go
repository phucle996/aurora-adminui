package response

import (
	"aurora-adminui/internal/domain/entity"
	"time"
)

type TemplateRender struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	ResourceDefinitionID string `json:"resource_definition_id"`
	ResourceType         string `json:"resource_type"`
	ResourceModel        string `json:"resource_model"`
	StreamKey            string `json:"stream_key"`
	ConsumerGroup        string `json:"consumer_group"`
	YAMLTemplate         string `json:"yaml_template"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

type TemplateRenderCatalogItem struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ResourceType  string `json:"resource_type"`
	ResourceModel string `json:"resource_model"`
	StreamKey     string `json:"stream_key"`
	ConsumerGroup string `json:"consumer_group"`
	YAMLValid     bool   `json:"yaml_valid"`
	UpdatedAt     string `json:"updated_at"`
}

type ListTemplateRenderCatalogResponse struct {
	Items []TemplateRenderCatalogItem `json:"items"`
}

func NewTemplateRender(item entity.TemplateRender) TemplateRender {
	return TemplateRender{
		ID:                   item.ID,
		Name:                 item.Name,
		Description:          item.Description,
		ResourceDefinitionID: item.ResourceDefinitionID,
		ResourceType:         item.ResourceType,
		ResourceModel:        item.ResourceModel,
		StreamKey:            item.StreamKey,
		ConsumerGroup:        item.ConsumerGroup,
		YAMLTemplate:         item.YAMLTemplate,
		CreatedAt:            item.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            item.UpdatedAt.Format(time.RFC3339),
	}
}

func NewListTemplateRenderCatalogResponse(items []entity.TemplateRender) ListTemplateRenderCatalogResponse {
	out := make([]TemplateRenderCatalogItem, 0, len(items))
	for _, item := range items {
		out = append(out, NewTemplateRenderCatalogItem(item))
	}
	return ListTemplateRenderCatalogResponse{Items: out}
}

func NewTemplateRenderCatalogItem(item entity.TemplateRender) TemplateRenderCatalogItem {
	return TemplateRenderCatalogItem{
		ID:            item.ID,
		Name:          item.Name,
		Description:   item.Description,
		ResourceType:  item.ResourceType,
		ResourceModel: item.ResourceModel,
		StreamKey:     item.StreamKey,
		ConsumerGroup: item.ConsumerGroup,
		YAMLValid:     item.YAMLValid,
		UpdatedAt:     item.UpdatedAt.Format(time.RFC3339),
	}
}
