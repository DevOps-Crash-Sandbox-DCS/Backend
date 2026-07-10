package sessions

import (
	"time"

	"github.com/google/uuid"
)

type CreateSessionRequest struct {
	ScenarioID string `json:"scenarioId" binding:"required"`
}

type SessionResponse struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"userId"`
	ScenarioID    string     `json:"scenarioId"`
	CurrentStepID *string    `json:"currentStepId"`
	Status        string     `json:"status"`
	Score         int        `json:"score"`
	StartedAt     time.Time  `json:"startedAt"`
	FinishedAt    *time.Time `json:"finishedAt"`
}

func ToSessionResponse(s Session) SessionResponse {
	return SessionResponse{
		ID:            s.ID,
		UserID:        s.UserID,
		ScenarioID:    s.ScenarioID,
		CurrentStepID: s.CurrentStepID,
		Status:        s.Status,
		Score:         s.Score,
		StartedAt:     s.StartedAt,
		FinishedAt:    s.FinishedAt,
	}
}
