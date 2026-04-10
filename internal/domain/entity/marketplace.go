package entity

import "time"

type MarketplaceApp struct {
	ID                   string
	Name                 string
	Slug                 string
	Summary              string
	Description          string
	ResourceDefinitionID string
	TemplateID           string
	TemplateName         string
	ResourceType         string
	ResourceModel        string
	Versions             []string
	Templates            []MarketplaceTemplate
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type MarketplaceTemplate struct {
	ID            string
	Name          string
	ResourceType  string
	ResourceModel string
	Version       string
}
