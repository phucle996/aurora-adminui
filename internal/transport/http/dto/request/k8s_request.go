package request

type CreateK8sClusterRequest struct {
	Name                     string
	Description              string
	ZoneID                   string
	SupportsDBAAS            bool
	SupportsServerless       bool
	SupportsGenericWorkloads bool
}
