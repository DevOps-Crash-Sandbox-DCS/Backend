package terminal

type CommandOutput struct {
	ID             string
	StepID         string
	CommandPattern string
	MatchType      string
	Stdout         string
	Stderr         string
	ExitCode       int
	Description    string
	Priority       int
	IsActive       bool
}
