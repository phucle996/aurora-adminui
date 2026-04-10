package entity

import "time"

type TemplateRender struct {
	ID                   string
	Name                 string
	Description          string
	ResourceDefinitionID string
	ResourceType         string
	ModelName            string
	ResourceVersion      string
	ResourceModel        string
	StreamKey            string
	ConsumerGroup        string
	YAMLTemplate         string
	YAMLValid            bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
