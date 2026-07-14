package reports

import (
	"time"

	"github.com/google/uuid"
)

type ReportSessionResponse struct {
	ID         uuid.UUID  `json:"id"`
	Status     string     `json:"status"`
	Score      int        `json:"score"`
	StartedAt  time.Time  `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`
}

type ReportScenarioResponse struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Category   string `json:"category"`
}

type ReportSummaryResponse struct {
	TotalSteps     int  `json:"totalSteps"`
	CompletedSteps int  `json:"completedSteps"`
	TotalActions   int  `json:"totalActions"`
	CorrectActions int  `json:"correctActions"`
	WrongActions   int  `json:"wrongActions"`
	Score          int  `json:"score"`
	MaxScore       int  `json:"maxScore"`
	IsCompleted    bool `json:"isCompleted"`
}

type ReportActionResponse struct {
	ID        uuid.UUID `json:"id"`
	Command   string    `json:"command"`
	IsCorrect bool      `json:"isCorrect"`
	Points    int       `json:"points"`
	Feedback  string    `json:"feedback"`
	CreatedAt time.Time `json:"createdAt"`
}

type ReportStepResponse struct {
	ID              string                 `json:"id"`
	Order           int                    `json:"order"`
	Title           string                 `json:"title"`
	Description     string                 `json:"description"`
	ExpectedCommand string                 `json:"expectedCommand"`
	Actions         []ReportActionResponse `json:"actions"`
}

type SessionReportResponse struct {
	Session  ReportSessionResponse  `json:"session"`
	Scenario ReportScenarioResponse `json:"scenario"`
	Summary  ReportSummaryResponse  `json:"summary"`
	Steps    []ReportStepResponse   `json:"steps"`
}
