package response

import "aurora-adminui/internal/domain/entity"

type MarketplaceCatalogItem struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Slug          string   `json:"slug"`
	Summary       string   `json:"summary"`
	TemplateID    string   `json:"template_id"`
	TemplateName  string   `json:"template_name"`
	ResourceType  string   `json:"resource_type"`
	ResourceModel string   `json:"resource_model"`
	Versions      []string `json:"versions"`
}

type ListMarketplaceAppsResponse struct {
	Items []MarketplaceCatalogItem `json:"items"`
}

type MarketplaceModelOption struct {
	ResourceDefinitionID string   `json:"resource_definition_id"`
	ResourceType         string   `json:"resource_type"`
	ResourceModel        string   `json:"resource_model"`
	Versions             []string `json:"versions"`
}

type ListMarketplaceModelOptionsResponse struct {
	Items []MarketplaceModelOption `json:"items"`
}

type MarketplaceTemplateOption struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	ResourceType  string `json:"resource_type"`
	ResourceModel string `json:"resource_model"`
	Version       string `json:"version"`
}

type ListMarketplaceTemplateOptionsResponse struct {
	Items []MarketplaceTemplateOption `json:"items"`
}

type MarketplaceApp struct {
	ID                   string   `json:"id"`
	Name                 string   `json:"name"`
	Slug                 string   `json:"slug"`
	Summary              string   `json:"summary"`
	Description          string   `json:"description"`
	ResourceDefinitionID string   `json:"resource_definition_id"`
	TemplateID           string   `json:"template_id"`
	TemplateName         string   `json:"template_name"`
	ResourceType         string   `json:"resource_type"`
	ResourceModel        string   `json:"resource_model"`
	Versions             []string `json:"versions"`
}

func NewListMarketplaceAppsResponse(items []entity.MarketplaceApp) ListMarketplaceAppsResponse {
	out := make([]MarketplaceCatalogItem, 0, len(items))
	for _, item := range items {
		out = append(out, NewMarketplaceCatalogItem(item))
	}
	return ListMarketplaceAppsResponse{Items: out}
}

func NewMarketplaceCatalogItem(item entity.MarketplaceApp) MarketplaceCatalogItem {
	return MarketplaceCatalogItem{
		ID:            item.ID,
		Name:          item.Name,
		Slug:          item.Slug,
		Summary:       item.Summary,
		TemplateID:    item.TemplateID,
		TemplateName:  item.TemplateName,
		ResourceType:  item.ResourceType,
		ResourceModel: item.ResourceModel,
		Versions:      append([]string(nil), item.Versions...),
	}
}

func NewListMarketplaceModelOptionsResponse(items []entity.MarketplaceApp) ListMarketplaceModelOptionsResponse {
	out := make([]MarketplaceModelOption, 0, len(items))
	for _, item := range items {
		out = append(out, MarketplaceModelOption{
			ResourceDefinitionID: item.ResourceDefinitionID,
			ResourceType:         item.ResourceType,
			ResourceModel:        item.ResourceModel,
			Versions:             append([]string(nil), item.Versions...),
		})
	}
	return ListMarketplaceModelOptionsResponse{Items: out}
}

func NewListMarketplaceTemplateOptionsResponse(items []entity.MarketplaceTemplate) ListMarketplaceTemplateOptionsResponse {
	out := make([]MarketplaceTemplateOption, 0, len(items))
	for _, item := range items {
		out = append(out, MarketplaceTemplateOption{
			ID:            item.ID,
			Name:          item.Name,
			ResourceType:  item.ResourceType,
			ResourceModel: item.ResourceModel,
			Version:       item.Version,
		})
	}
	return ListMarketplaceTemplateOptionsResponse{Items: out}
}

func NewMarketplaceApp(item entity.MarketplaceApp) MarketplaceApp {
	return MarketplaceApp{
		ID:                   item.ID,
		Name:                 item.Name,
		Slug:                 item.Slug,
		Summary:              item.Summary,
		Description:          item.Description,
		ResourceDefinitionID: item.ResourceDefinitionID,
		TemplateID:           item.TemplateID,
		TemplateName:         item.TemplateName,
		ResourceType:         item.ResourceType,
		ResourceModel:        item.ResourceModel,
		Versions:             append([]string(nil), item.Versions...),
	}
}
