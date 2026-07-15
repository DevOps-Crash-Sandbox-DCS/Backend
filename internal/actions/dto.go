package actions

import "github.com/google/uuid"

type SubmitActionRequest struct {
	StepID  string `json:"stepId" binding:"required"`
	Command string `json:"command" binding:"required"`
}

type SubmitActionResponse struct {
	ID            uuid.UUID `json:"id"`
	SessionID     uuid.UUID `json:"sessionId"`
	StepID        string    `json:"stepId"`
	Command       string    `json:"command"`
	IsCorrect     bool      `json:"isCorrect"`
	Points        int       `json:"points"`
	Feedback      string    `json:"feedback"`
	NextStepID    *string   `json:"nextStepId"`
	SessionStatus string    `json:"sessionStatus"`
	SessionScore  int       `json:"sessionScore"`
}
