package response

import "aurora-adminui/internal/domain/entity"

type User struct {
	ID          string `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Status      string `json:"status"`
	Role        string `json:"role"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type ListUsersResponse struct {
	Items []User `json:"items"`
}

func NewListUsersResponse(items []entity.User) ListUsersResponse {
	out := make([]User, 0, len(items))
	for _, item := range items {
		out = append(out, User{
			ID:          item.ID,
			FirstName:   item.FirstName,
			LastName:    item.LastName,
			Username:    item.Username,
			Email:       item.Email,
			PhoneNumber: item.PhoneNumber,
			Status:      item.Status,
			Role:        item.Role,
			CreatedAt:   item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}
	return ListUsersResponse{Items: out}
}
