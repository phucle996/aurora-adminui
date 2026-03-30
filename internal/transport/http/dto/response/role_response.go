package response

import "aurora-adminui/internal/domain/entity"

type Role struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Scope           string `json:"scope"`
	Description     string `json:"description"`
	UserCount       int    `json:"userCount"`
	PermissionCount int    `json:"permissionCount"`
}

type ListRolesResponse struct {
	Items []Role `json:"items"`
}

type Permission struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ListPermissionsResponse struct {
	Items []Permission `json:"items"`
}

func NewListRolesResponse(items []entity.Role) ListRolesResponse {
	out := make([]Role, 0, len(items))
	for _, item := range items {
		out = append(out, Role{
			ID:              item.ID,
			Name:            item.Name,
			Scope:           item.Scope,
			Description:     item.Description,
			UserCount:       item.UserCount,
			PermissionCount: item.PermissionCount,
		})
	}
	return ListRolesResponse{Items: out}
}

func NewListPermissionsResponse(items []entity.Permission) ListPermissionsResponse {
	out := make([]Permission, 0, len(items))
	for _, item := range items {
		out = append(out, Permission{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
		})
	}
	return ListPermissionsResponse{Items: out}
}
