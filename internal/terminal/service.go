package terminal

import (
	"context"
	"strings"
)

const (
	OutputMatchTypeExact    = "exact"
	OutputMatchTypeContains = "contains"
)

type CommandResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetCommandOutput(ctx context.Context, stepID string, command string) (*CommandResult, error) {
	outputs, err := s.repo.GetCommandOutputsByStepID(ctx, stepID)
	if err != nil {
		return nil, err
	}

	normalizedCommand := normalizeTerminalCommand(command)

	for _, output := range outputs {
		normalizedPattern := normalizeTerminalCommand(output.CommandPattern)

		switch output.MatchType {
		case OutputMatchTypeExact:
			if normalizedCommand == normalizedPattern {
				return &CommandResult{
					Stdout:   output.Stdout,
					Stderr:   output.Stderr,
					ExitCode: output.ExitCode,
				}, nil
			}

		case OutputMatchTypeContains:
			if normalizedPattern != "" && strings.Contains(normalizedCommand, normalizedPattern) {
				return &CommandResult{
					Stdout:   output.Stdout,
					Stderr:   output.Stderr,
					ExitCode: output.ExitCode,
				}, nil
			}

		default:
			if normalizedCommand == normalizedPattern {
				return &CommandResult{
					Stdout:   output.Stdout,
					Stderr:   output.Stderr,
					ExitCode: output.ExitCode,
				}, nil
			}
		}
	}

	// Fallback, чтобы старые команды не ломались.
	return &CommandResult{
		Stdout:   SimulateCommandOutput(command),
		Stderr:   "",
		ExitCode: 0,
	}, nil
}

func normalizeTerminalCommand(command string) string {
	fields := strings.Fields(strings.TrimSpace(command))

	return strings.ToLower(strings.Join(fields, " "))
}
