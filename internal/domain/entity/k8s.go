package entity

import (
	"time"

	"github.com/google/uuid"
)

type K8sClusterImportMode string

const (
	K8sClusterImportModeKubeconfig K8sClusterImportMode = "kubeconfig"
)

type K8sClusterValidationStatus string

const (
	K8sClusterValidationStatusPending     K8sClusterValidationStatus = "pending"
	K8sClusterValidationStatusValid       K8sClusterValidationStatus = "valid"
	K8sClusterValidationStatusInvalid     K8sClusterValidationStatus = "invalid"
	K8sClusterValidationStatusUnreachable K8sClusterValidationStatus = "unreachable"
)

type K8sCluster struct {
	ID                       uuid.UUID
	Name                     string
	Description              string
	ImportMode               K8sClusterImportMode
	KubeconfigCiphertext     string
	APIServerURL             string
	CurrentContext           string
	KubernetesVersion        string
	ValidationStatus         K8sClusterValidationStatus
	LastValidatedAt          *time.Time
	LastValidationError      string
	SupportsDBAAS            bool
	SupportsServerless       bool
	SupportsGenericWorkloads bool
	CreatedAt                time.Time
	ZoneID                   *uuid.UUID
	ZoneName                 string
}

type K8sClusterCreateInput struct {
	Name                     string
	Description              string
	ZoneID                   string
	SupportsDBAAS            bool
	SupportsServerless       bool
	SupportsGenericWorkloads bool
	Kubeconfig               []byte
}
