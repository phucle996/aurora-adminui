package entity

import (
	"time"

	"github.com/google/uuid"
)

type Zone struct {
	ID            uuid.UUID
	Name          string
	Description   string
	ResourceCount int
	CreatedAt     time.Time
}
