package request

type UpdateHypervisorNodeNameRequest struct {
	Name string `json:"name"`
}

type AssignHypervisorNodeZoneRequest struct {
	ZoneID string `json:"zone_id"`
}
