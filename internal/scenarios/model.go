package scenarios

import "time"

type Scenario struct {
	ID               string
	Title            string
	Description      string
	Difficulty       string
	Category         string
	EstimatedMinutes int
	UserNotification string
	DesktopSymptoms  string
	TerminalSolution string
	QuickFix         string
	IsActive         bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
