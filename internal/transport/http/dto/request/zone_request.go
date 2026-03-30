package request

type CreateZoneRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
