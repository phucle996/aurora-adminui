package entity

import "github.com/google/uuid"

type ResourceDefinition struct {
	ID              uuid.UUID
	ResourceType    string
	ResourceModel   string
	ResourceVersion string
	DisplayName     string
	Status          string
}

type ResourceDefinitionCatalogItem struct {
	ResourceDefinition
	ResourceCount int
}

type ResourceDefinitionZoneSupport struct {
	ZoneID   uuid.UUID
	ZoneName string
	Enabled  bool
}
