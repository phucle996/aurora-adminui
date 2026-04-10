package request

type CreateMarketplaceAppRequest struct {
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	Summary              string `json:"summary"`
	Description          string `json:"description"`
	ResourceDefinitionID string `json:"resource_definition_id"`
	TemplateID           string `json:"template_id"`
}
