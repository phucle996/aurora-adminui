package request

type CreateK8sClusterRequest struct {
	Name        string
	Description string
	ZoneID      string
}

type UpdateK8sClusterRequest struct {
	ZoneID string `json:"zone_id"`
}
