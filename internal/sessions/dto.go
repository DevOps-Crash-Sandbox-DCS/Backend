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

type SessionHistoryItemResponse struct {
	ID            uuid.UUID  `json:"id"`
	ScenarioID    string     `json:"scenarioId"`
	ScenarioTitle string     `json:"scenarioTitle"`
	Difficulty    string     `json:"difficulty"`
	Category      string     `json:"category"`
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

func ToSessionHistoryItemResponse(item SessionHistoryItem) SessionHistoryItemResponse {
	return SessionHistoryItemResponse{
		ID:            item.ID,
		ScenarioID:    item.ScenarioID,
		ScenarioTitle: item.ScenarioTitle,
		Difficulty:    item.Difficulty,
		Category:      item.Category,
		Status:        item.Status,
		Score:         item.Score,
		StartedAt:     item.StartedAt,
		FinishedAt:    item.FinishedAt,
	}
}
