package response

type ResourceDefinition struct {
	ID              string `json:"id"`
	ResourceType    string `json:"resource_type"`
	ResourceModel   string `json:"resource_model"`
	ResourceVersion string `json:"resource_version"`
	DisplayName     string `json:"display_name"`
	Status          string `json:"status"`
	ResourceCount   int    `json:"resource_count"`
}

type ListResourceDefinitionsResponse struct {
	Items []ResourceDefinition `json:"items"`
}

type TemplateResourceDefinitionOption struct {
	ID              string `json:"id"`
	ResourceType    string `json:"resource_type"`
	ResourceModel   string `json:"resource_model"`
	ResourceVersion string `json:"resource_version"`
}

type ListTemplateResourceDefinitionOptionsResponse struct {
	Items []TemplateResourceDefinitionOption `json:"items"`
}
