package response

import (
	"time"

	"aurora-adminui/internal/domain/entity"
)

type K8sClusterListItem struct {
	ID                string     `json:"id"`
	Name              string     `json:"name"`
	Description       string     `json:"description"`
	APIServerURL      string     `json:"api_server_url"`
	KubernetesVersion string     `json:"kubernetes_version"`
	ValidationStatus  string     `json:"validation_status"`
	LastValidatedAt   *time.Time `json:"last_validated_at,omitempty"`
	ZoneName          string     `json:"zone_name"`
}

type ListK8sClustersResponse struct {
	Items []K8sClusterListItem `json:"items"`
}

type K8sClusterDetail struct {
	ID                string           `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	APIServerURL      string           `json:"api_server_url"`
	CurrentContext    string           `json:"current_context"`
	KubernetesVersion string           `json:"kubernetes_version"`
	ValidationStatus  string           `json:"validation_status"`
	LastValidatedAt   *time.Time       `json:"last_validated_at,omitempty"`
	CreatedAt         time.Time        `json:"created_at"`
	ZoneID            string           `json:"zone_id,omitempty"`
	ZoneName          string           `json:"zone_name"`
	Nodes             []K8sClusterNode `json:"nodes"`
}

type K8sClusterZoneOption struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type K8sClusterDetailPageData struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	Description       string               `json:"description"`
	APIServerURL      string               `json:"api_server_url"`
	CurrentContext    string               `json:"current_context"`
	KubernetesVersion string               `json:"kubernetes_version"`
	ValidationStatus  string               `json:"validation_status"`
	LastValidatedAt   *time.Time           `json:"last_validated_at,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
	ZoneID            string               `json:"zone_id,omitempty"`
	ZoneName          string               `json:"zone_name"`
	Nodes             []K8sClusterNode     `json:"nodes"`
	ZoneOptions       []K8sClusterZoneOption `json:"zone_options"`
}

type K8sClusterNode struct {
	Name             string   `json:"name"`
	Roles            []string `json:"roles"`
	KubeletVersion   string   `json:"kubelet_version"`
	ContainerRuntime string   `json:"container_runtime"`
	OSImage          string   `json:"os_image"`
	KernelVersion    string   `json:"kernel_version"`
	Ready            bool     `json:"ready"`
}

func NewListK8sClustersResponse(items []entity.K8sCluster) ListK8sClustersResponse {
	out := make([]K8sClusterListItem, 0, len(items))
	for _, item := range items {
		out = append(out, K8sClusterListItem{
			ID:                item.ID.String(),
			Name:              item.Name,
			Description:       item.Description,
			APIServerURL:      item.APIServerURL,
			KubernetesVersion: item.KubernetesVersion,
			ValidationStatus:  string(item.ValidationStatus),
			LastValidatedAt:   item.LastValidatedAt,
			ZoneName:          item.ZoneName,
		})
	}
	return ListK8sClustersResponse{Items: out}
}

func NewK8sClusterDetail(item *entity.K8sCluster) K8sClusterDetail {
	if item == nil {
		return K8sClusterDetail{}
	}
	detail := K8sClusterDetail{
		ID:                item.ID.String(),
		Name:              item.Name,
		Description:       item.Description,
		APIServerURL:      item.APIServerURL,
		CurrentContext:    item.CurrentContext,
		KubernetesVersion: item.KubernetesVersion,
		ValidationStatus:  string(item.ValidationStatus),
		LastValidatedAt:   item.LastValidatedAt,
		CreatedAt:         item.CreatedAt,
		ZoneName:          item.ZoneName,
		Nodes:             make([]K8sClusterNode, 0, len(item.Nodes)),
	}
	if item.ZoneID != nil {
		detail.ZoneID = item.ZoneID.String()
	}
	for _, node := range item.Nodes {
		detail.Nodes = append(detail.Nodes, K8sClusterNode{
			Name:             node.Name,
			Roles:            append([]string(nil), node.Roles...),
			KubeletVersion:   node.KubeletVersion,
			ContainerRuntime: node.ContainerRuntime,
			OSImage:          node.OSImage,
			KernelVersion:    node.KernelVersion,
			Ready:            node.Ready,
		})
	}
	return detail
}

func NewK8sClusterDetailPageData(item *entity.K8sCluster, nodes []entity.K8sClusterNode, zones []entity.Zone) K8sClusterDetailPageData {
	detail := K8sClusterDetailPageData{
		ID:                item.ID.String(),
		Name:              item.Name,
		Description:       item.Description,
		APIServerURL:      item.APIServerURL,
		CurrentContext:    item.CurrentContext,
		KubernetesVersion: item.KubernetesVersion,
		ValidationStatus:  string(item.ValidationStatus),
		LastValidatedAt:   item.LastValidatedAt,
		CreatedAt:         item.CreatedAt,
		ZoneName:          item.ZoneName,
		Nodes:             make([]K8sClusterNode, 0, len(nodes)),
		ZoneOptions:       make([]K8sClusterZoneOption, 0, len(zones)),
	}
	if item.ZoneID != nil {
		detail.ZoneID = item.ZoneID.String()
	}
	for _, node := range nodes {
		detail.Nodes = append(detail.Nodes, K8sClusterNode{
			Name:             node.Name,
			Roles:            append([]string(nil), node.Roles...),
			KubeletVersion:   node.KubeletVersion,
			ContainerRuntime: node.ContainerRuntime,
			OSImage:          node.OSImage,
			KernelVersion:    node.KernelVersion,
			Ready:            node.Ready,
		})
	}
	for _, zone := range zones {
		detail.ZoneOptions = append(detail.ZoneOptions, K8sClusterZoneOption{
			ID:   zone.ID.String(),
			Name: zone.Name,
		})
	}
	return detail
}
