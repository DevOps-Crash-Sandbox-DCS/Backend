package sandbox

import (
	"time"

	"github.com/google/uuid"
)

const (
	SandboxStatusRunning = "running"
	SandboxStatusStopped = "stopped"
	SandboxStatusFailed  = "failed"
)

type Sandbox struct {
	ID            uuid.UUID
	SessionID     uuid.UUID
	ScenarioID    string
	ContainerName string
	Image         string
	Status        string
	StartedAt     time.Time
	StoppedAt     *time.Time
	LastSeenAt    time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CommandResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

type StartSandboxRequest struct {
	SessionID  uuid.UUID
	ScenarioID string
	Image      string
}

type ExecCommandRequest struct {
	SessionID uuid.UUID
	Command   string
}

type CleanupResult struct {
	StoppedInDB       int `json:"stoppedInDb"`
	RemovedContainers int `json:"removedContainers"`
}
