package terminal

type ClientMessage struct {
	Type    string `json:"type"`
	Command string `json:"command"`
}

type ServerMessage struct {
	Type          string      `json:"type"`
	Output        string      `json:"output,omitempty"`
	Stdout        string      `json:"stdout,omitempty"`
	Stderr        string      `json:"stderr,omitempty"`
	ExitCode      *int        `json:"exitCode,omitempty"`
	Error         string      `json:"error,omitempty"`
	SessionID     string      `json:"sessionId,omitempty"`
	CurrentStepID *string     `json:"currentStepId,omitempty"`
	ActionResult  interface{} `json:"actionResult,omitempty"`
}

const (
	ClientMessageTypeCommand = "command"

	ServerMessageTypeWelcome       = "welcome"
	ServerMessageTypeOutput        = "output"
	ServerMessageTypeActionResult  = "action_result"
	ServerMessageTypeError         = "error"
	ServerMessageTypeSessionClosed = "session_closed"
)
