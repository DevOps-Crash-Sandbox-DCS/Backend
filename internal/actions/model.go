package actions

import (
	"time"

	"github.com/google/uuid"
)

type Action struct {
	ID        uuid.UUID
	SessionID uuid.UUID
	StepID    string
	Command   string
	IsCorrect bool
	Points    int
	Feedback  string
	CreatedAt time.Time
}

type Session struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	ScenarioID    string
	CurrentStepID *string
	Status        string
	Score         int
}

type ScenarioStep struct {
	ID              string
	ScenarioID      string
	Order           int
	ExpectedCommand string
}
