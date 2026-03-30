package response

import "aurora-adminui/internal/domain/entity"

type Zone struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	ResourceCount int    `json:"resource_count"`
	CanDelete     bool   `json:"can_delete"`
}

type ListZonesResponse struct {
	Items []Zone `json:"items"`
}

func NewListZonesResponse(items []entity.Zone) ListZonesResponse {
	out := make([]Zone, 0, len(items))
	for _, item := range items {
		out = append(out, Zone{
			ID:            item.ID.String(),
			Name:          item.Name,
			Description:   item.Description,
			ResourceCount: item.ResourceCount,
			CanDelete:     item.ResourceCount == 0,
		})
	}
	return ListZonesResponse{Items: out}
}

func NewZone(item *entity.Zone) Zone {
	if item == nil {
		return Zone{}
	}
	return Zone{
		ID:            item.ID.String(),
		Name:          item.Name,
		Description:   item.Description,
		ResourceCount: item.ResourceCount,
		CanDelete:     item.ResourceCount == 0,
	}
}
