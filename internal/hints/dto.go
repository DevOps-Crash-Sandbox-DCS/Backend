package hints

import "time"

type HintLevel string

const (
	HintLevelBasic    HintLevel = "basic"
	HintLevelDetailed HintLevel = "detailed"
	HintLevelDirect   HintLevel = "direct"
)

type CreateHintRequest struct {
	Level string `json:"level"`
}

type HintResponse struct {
	Hint          string   `json:"hint"`
	Confidence    *float64 `json:"confidence,omitempty"`
	Source        string   `json:"source"`
	CurrentStepID *string  `json:"currentStepId,omitempty"`
}

type MLHintRequest struct {
	UserID               string          `json:"userId"`
	SessionID            string          `json:"sessionId"`
	ScenarioID           string          `json:"scenarioId"`
	CurrentStepID        *string         `json:"currentStepId"`
	SessionStatus        string          `json:"sessionStatus"`
	Score                int             `json:"score"`
	HintLevel            string          `json:"hintLevel"`
	CurrentStep          *MLCurrentStep  `json:"currentStep,omitempty"`
	History              []MLActionEntry `json:"history"`
	RecentTerminalOutput string          `json:"recentTerminalOutput"`
}

type MLCurrentStep struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Hint            string `json:"hint"`
	ExpectedCommand string `json:"expectedCommand"`
	ExpectedResult  string `json:"expectedResult"`
}

type MLActionEntry struct {
	StepID    string    `json:"stepId"`
	Command   string    `json:"command"`
	IsCorrect bool      `json:"isCorrect"`
	Points    int       `json:"points"`
	Feedback  string    `json:"feedback"`
	CreatedAt time.Time `json:"createdAt"`
}

type MLHintResponse struct {
	Hint       string   `json:"hint"`
	Confidence *float64 `json:"confidence,omitempty"`
	Source     string   `json:"source"`
	Reasoning  string   `json:"reasoning,omitempty"`
}

type SessionContext struct {
	Session SessionInfo
	Step    *StepInfo
	History []ActionHistoryItem
}

type SessionInfo struct {
	ID            string
	UserID        string
	ScenarioID    string
	CurrentStepID *string
	Status        string
	Score         int
}

type StepInfo struct {
	ID              string
	Title           string
	Description     string
	Hint            string
	ExpectedCommand string
	ExpectedResult  string
}

type ActionHistoryItem struct {
	StepID    string
	Command   string
	IsCorrect bool
	Points    int
	Feedback  string
	CreatedAt time.Time
}
