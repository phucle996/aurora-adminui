package entity

import (
	"time"

	"github.com/google/uuid"
)

type K8sClusterValidationStatus string

const (
	K8sClusterValidationStatusPending     K8sClusterValidationStatus = "pending"
	K8sClusterValidationStatusValid       K8sClusterValidationStatus = "valid"
	K8sClusterValidationStatusInvalid     K8sClusterValidationStatus = "invalid"
	K8sClusterValidationStatusUnreachable K8sClusterValidationStatus = "unreachable"
)

type K8sCluster struct {
	ID                   uuid.UUID
	Name                 string
	Description          string
	KubeconfigCiphertext string
	APIServerURL         string
	CurrentContext       string
	KubernetesVersion    string
	ValidationStatus     K8sClusterValidationStatus
	LastValidatedAt      *time.Time
	CreatedAt            time.Time
	ZoneID               *uuid.UUID
	ZoneName             string
	Nodes                []K8sClusterNode
}

type K8sClusterCreateInput struct {
	Name        string
	Description string
	ZoneID      string
	Kubeconfig  []byte
}

type K8sClusterNode struct {
	Name             string
	Roles            []string
	KubeletVersion   string
	ContainerRuntime string
	OSImage          string
	KernelVersion    string
	Ready            bool
}

type K8sClusterUpdateInput struct {
	ZoneID     string
	Kubeconfig []byte
}
