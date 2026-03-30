package response

import (
	"time"

	"aurora-adminui/internal/domain/entity"
)

type K8sClusterListItem struct {
	ID                       string     `json:"id"`
	Name                     string     `json:"name"`
	Description              string     `json:"description"`
	APIServerURL             string     `json:"api_server_url"`
	KubernetesVersion        string     `json:"kubernetes_version"`
	ValidationStatus         string     `json:"validation_status"`
	LastValidatedAt          *time.Time `json:"last_validated_at,omitempty"`
	SupportsDBAAS            bool       `json:"supports_dbaas"`
	SupportsServerless       bool       `json:"supports_serverless"`
	SupportsGenericWorkloads bool       `json:"supports_generic_workloads"`
	ZoneName                 string     `json:"zone_name"`
}

type ListK8sClustersResponse struct {
	Items []K8sClusterListItem `json:"items"`
}

type K8sClusterDetail struct {
	ID                       string     `json:"id"`
	Name                     string     `json:"name"`
	Description              string     `json:"description"`
	ImportMode               string     `json:"import_mode"`
	APIServerURL             string     `json:"api_server_url"`
	CurrentContext           string     `json:"current_context"`
	KubernetesVersion        string     `json:"kubernetes_version"`
	ValidationStatus         string     `json:"validation_status"`
	LastValidatedAt          *time.Time `json:"last_validated_at,omitempty"`
	LastValidationError      string     `json:"last_validation_error,omitempty"`
	SupportsDBAAS            bool       `json:"supports_dbaas"`
	SupportsServerless       bool       `json:"supports_serverless"`
	SupportsGenericWorkloads bool       `json:"supports_generic_workloads"`
	CreatedAt                time.Time  `json:"created_at"`
	ZoneID                   string     `json:"zone_id,omitempty"`
	ZoneName                 string     `json:"zone_name"`
}

func NewListK8sClustersResponse(items []entity.K8sCluster) ListK8sClustersResponse {
	out := make([]K8sClusterListItem, 0, len(items))
	for _, item := range items {
		out = append(out, K8sClusterListItem{
			ID:                       item.ID.String(),
			Name:                     item.Name,
			Description:              item.Description,
			APIServerURL:             item.APIServerURL,
			KubernetesVersion:        item.KubernetesVersion,
			ValidationStatus:         string(item.ValidationStatus),
			LastValidatedAt:          item.LastValidatedAt,
			SupportsDBAAS:            item.SupportsDBAAS,
			SupportsServerless:       item.SupportsServerless,
			SupportsGenericWorkloads: item.SupportsGenericWorkloads,
			ZoneName:                 item.ZoneName,
		})
	}
	return ListK8sClustersResponse{Items: out}
}

func NewK8sClusterDetail(item *entity.K8sCluster) K8sClusterDetail {
	if item == nil {
		return K8sClusterDetail{}
	}
	detail := K8sClusterDetail{
		ID:                       item.ID.String(),
		Name:                     item.Name,
		Description:              item.Description,
		ImportMode:               string(item.ImportMode),
		APIServerURL:             item.APIServerURL,
		CurrentContext:           item.CurrentContext,
		KubernetesVersion:        item.KubernetesVersion,
		ValidationStatus:         string(item.ValidationStatus),
		LastValidatedAt:          item.LastValidatedAt,
		LastValidationError:      item.LastValidationError,
		SupportsDBAAS:            item.SupportsDBAAS,
		SupportsServerless:       item.SupportsServerless,
		SupportsGenericWorkloads: item.SupportsGenericWorkloads,
		CreatedAt:                item.CreatedAt,
		ZoneName:                 item.ZoneName,
	}
	if item.ZoneID != nil {
		detail.ZoneID = item.ZoneID.String()
	}
	return detail
}
