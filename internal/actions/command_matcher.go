package actions

import (
	"strings"
)

const (
	CommandMatchTypeExact    = "exact"
	CommandMatchTypeContains = "contains"
)

func IsCommandAccepted(input string, expectedCommand string, acceptedCommands []AcceptedCommand) bool {
	normalizedInput := NormalizeCommand(input)

	for _, accepted := range acceptedCommands {
		normalizedAccepted := NormalizeCommand(accepted.Command)

		switch accepted.MatchType {
		case CommandMatchTypeExact:
			if normalizedInput == normalizedAccepted {
				return true
			}

		case CommandMatchTypeContains:
			if normalizedAccepted != "" && strings.Contains(normalizedInput, normalizedAccepted) {
				return true
			}

		default:
			if normalizedInput == normalizedAccepted {
				return true
			}
		}
	}

	normalizedExpected := NormalizeCommand(expectedCommand)

	return normalizedExpected != "" && normalizedInput == normalizedExpected
}

func NormalizeCommand(command string) string {
	normalized := strings.TrimSpace(command)

	normalized = strings.ReplaceAll(normalized, "\t", " ")

	for strings.Contains(normalized, "  ") {
		normalized = strings.ReplaceAll(normalized, "  ", " ")
	}

	return normalized
}
