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
