package scenarios

type ScenarioResponse struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Difficulty       string `json:"difficulty"`
	Category         string `json:"category"`
	EstimatedMinutes int    `json:"estimatedMinutes"`
	UserNotification string `json:"userNotification"`
	DesktopSymptoms  string `json:"desktopSymptoms"`
	TerminalSolution string `json:"terminalSolution"`
	QuickFix         string `json:"quickFix"`
}

type ScenarioStepResponse struct {
	ID              string `json:"id"`
	Order           int    `json:"order"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	Hint            string `json:"hint"`
	ExpectedCommand string `json:"expectedCommand"`
	ExpectedResult  string `json:"expectedResult"`
}

type ScenarioDetailsResponse struct {
	ID               string                 `json:"id"`
	Title            string                 `json:"title"`
	Description      string                 `json:"description"`
	Difficulty       string                 `json:"difficulty"`
	Category         string                 `json:"category"`
	EstimatedMinutes int                    `json:"estimatedMinutes"`
	UserNotification string                 `json:"userNotification"`
	DesktopSymptoms  string                 `json:"desktopSymptoms"`
	TerminalSolution string                 `json:"terminalSolution"`
	QuickFix         string                 `json:"quickFix"`
	Steps            []ScenarioStepResponse `json:"steps"`
}

func ToScenarioResponse(s Scenario) ScenarioResponse {
	return ScenarioResponse{
		ID:               s.ID,
		Title:            s.Title,
		Description:      s.Description,
		Difficulty:       s.Difficulty,
		Category:         s.Category,
		EstimatedMinutes: s.EstimatedMinutes,
		UserNotification: s.UserNotification,
		DesktopSymptoms:  s.DesktopSymptoms,
		TerminalSolution: s.TerminalSolution,
		QuickFix:         s.QuickFix,
	}
}

func ToScenarioStepResponse(s ScenarioStep) ScenarioStepResponse {
	return ScenarioStepResponse{
		ID:              s.ID,
		Order:           s.Order,
		Title:           s.Title,
		Description:     s.Description,
		Hint:            s.Hint,
		ExpectedCommand: s.ExpectedCommand,
		ExpectedResult:  s.ExpectedResult,
	}
}

func ToScenarioDetailsResponse(s Scenario, steps []ScenarioStep) ScenarioDetailsResponse {
	stepResponses := make([]ScenarioStepResponse, 0, len(steps))

	for _, step := range steps {
		stepResponses = append(stepResponses, ToScenarioStepResponse(step))
	}

	return ScenarioDetailsResponse{
		ID:               s.ID,
		Title:            s.Title,
		Description:      s.Description,
		Difficulty:       s.Difficulty,
		Category:         s.Category,
		EstimatedMinutes: s.EstimatedMinutes,
		UserNotification: s.UserNotification,
		DesktopSymptoms:  s.DesktopSymptoms,
		TerminalSolution: s.TerminalSolution,
		QuickFix:         s.QuickFix,
		Steps:            stepResponses,
	}
}
