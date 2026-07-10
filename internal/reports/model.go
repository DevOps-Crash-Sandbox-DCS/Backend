package reports

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	ScenarioID string
	Status     string
	Score      int
	StartedAt  time.Time
	FinishedAt *time.Time
}

type Scenario struct {
	ID         string
	Title      string
	Difficulty string
	Category   string
}

type ScenarioStep struct {
	ID              string
	ScenarioID      string
	Order           int
	Title           string
	Description     string
	ExpectedCommand string
}

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
