package request

type CreateRDRequest struct {
	ResourceType    string `json:"resource_type" binding:"required"`
	ResourceModel   string `json:"resource_model"`
	ModelName       string `json:"model_name"`
	ResourceVersion string `json:"resource_version" binding:"required"`
	DisplayName     string `json:"display_name" binding:"required"`
}

type UpdateResourceDefinitionStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type ReplaceResourceDefinitionZoneSupportRequest struct {
	ZoneIDs []string `json:"zone_ids"`
}
