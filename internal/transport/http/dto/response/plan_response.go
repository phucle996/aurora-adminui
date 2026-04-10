package response

import (
	"aurora-adminui/internal/domain/entity"
)

type Plan struct {
	ID            string `json:"id"`
	ResourceType  string `json:"resource_type"`
	ResourceModel string `json:"resource_model"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	VCPU          int    `json:"vcpu"`
	RAMGB         int    `json:"ram_gb"`
	DiskGB        int    `json:"disk_gb"`
}

type ListPlansResponse struct {
	Items []Plan `json:"items"`
}

func NewListPlansResponse(items []entity.Plan) ListPlansResponse {
	out := make([]Plan, 0, len(items))
	for _, item := range items {
		out = append(out, NewPlan(item))
	}
	return ListPlansResponse{Items: out}
}

func NewPlan(item entity.Plan) Plan {
	return Plan{
		ID:            item.ID.String(),
		ResourceType:  string(item.ResourceType),
		ResourceModel: item.ResourceModel,
		Code:          item.Code,
		Name:          item.Name,
		Description:   item.Description,
		Status:        string(item.Status),
		VCPU:          item.VCPU,
		RAMGB:         item.RAMGB,
		DiskGB:        item.DiskGB,
	}
}
