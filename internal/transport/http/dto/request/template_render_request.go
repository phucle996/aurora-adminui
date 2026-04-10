package request

type CreateTemplateRenderRequest struct {
	ResourceDefinitionID string `json:"resource_definition_id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	StreamKey            string `json:"stream_key"`
	ConsumerGroup        string `json:"consumer_group"`
	YAMLTemplate         string `json:"yaml_template"`
}

type UpdateTemplateRenderRequest struct {
	ResourceDefinitionID string `json:"resource_definition_id"`
	Name                 string `json:"name"`
	Description          string `json:"description"`
	StreamKey            string `json:"stream_key"`
	ConsumerGroup        string `json:"consumer_group"`
	YAMLTemplate         string `json:"yaml_template"`
}
