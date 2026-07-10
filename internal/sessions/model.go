package sessions

import (
	"time"

	"github.com/google/uuid"
)

const (
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusFailed     = "failed"
)

type Session struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ScenarioID    string
	CurrentStepID *string
	Status        string
	Score         int
	StartedAt     time.Time
	FinishedAt    *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
