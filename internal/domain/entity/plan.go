package entity

import (
	"time"

	"github.com/google/uuid"
)

type ResourceType string

const (
	ResourceTypeVPS ResourceType = "vps"
)

type PlanStatus string

const (
	PlanStatusActive  PlanStatus = "active"
	PlanStatusRetired PlanStatus = "retired"
)

type Plan struct {
	ID           uuid.UUID
	ResourceType ResourceType
	Code         string
	Name         string
	Description  string
	Status       PlanStatus
	VCPU         int
	RAMGB        int
	DiskGB       int
	CreatedAt    time.Time
}
