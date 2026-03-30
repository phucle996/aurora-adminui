package request

type CreatePlanRequest struct {
	ResourceType string `json:"resourceType"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	VCPU         int    `json:"vcpu"`
	RAMGB        int    `json:"ramGb"`
	DiskGB       int    `json:"diskGb"`
}
