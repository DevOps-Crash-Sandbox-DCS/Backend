package hints

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
)

var ErrInvalidHintLevel = errors.New("invalid hint level")

type Service struct {
	repo   *Repository
	client *Client
}

func NewService(repo *Repository, client *Client) *Service {
	return &Service{
		repo:   repo,
		client: client,
	}
}

func (s *Service) CreateHint(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
	req CreateHintRequest,
) (*HintResponse, error) {
	hintLevel := normalizeHintLevel(req.Level)

	contextData, err := s.repo.GetSessionContext(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}

	mlReq := buildMLHintRequest(contextData, hintLevel, req.RecentTerminalOutput)

	mlResp, err := s.client.GetHint(ctx, mlReq)
	if err != nil {
		mlResp = buildFallbackHint(contextData, err)
	}

	if err := s.repo.SaveHint(
		ctx,
		sessionID,
		userID,
		contextData.Session.ScenarioID,
		contextData.Session.CurrentStepID,
		hintLevel,
		mlReq,
		*mlResp,
	); err != nil {
		return nil, err
	}

	return &HintResponse{
		Hint:          mlResp.Hint,
		Confidence:    mlResp.Confidence,
		Source:        mlResp.Source,
		CurrentStepID: contextData.Session.CurrentStepID,
	}, nil
}

func normalizeHintLevel(level string) string {
	value := strings.ToLower(strings.TrimSpace(level))

	switch value {
	case string(HintLevelBasic), string(HintLevelDetailed), string(HintLevelDirect):
		return value
	case "":
		return string(HintLevelBasic)
	default:
		return string(HintLevelBasic)
	}
}

func buildMLHintRequest(contextData *SessionContext, hintLevel string, recentTerminalOutput string) MLHintRequest {
	history := make([]MLActionEntry, 0, len(contextData.History))

	for _, item := range contextData.History {
		history = append(history, MLActionEntry{
			StepID:    item.StepID,
			Command:   item.Command,
			IsCorrect: item.IsCorrect,
			Points:    item.Points,
			Feedback:  item.Feedback,
			CreatedAt: item.CreatedAt,
		})
	}

	var currentStep *MLCurrentStep

	if contextData.Step != nil {
		currentStep = &MLCurrentStep{
			ID:              contextData.Step.ID,
			Title:           contextData.Step.Title,
			Description:     contextData.Step.Description,
			Hint:            contextData.Step.Hint,
			ExpectedCommand: contextData.Step.ExpectedCommand,
			ExpectedResult:  contextData.Step.ExpectedResult,
		}
	}

	return MLHintRequest{
		UserID:               contextData.Session.UserID,
		SessionID:            contextData.Session.ID,
		ScenarioID:           contextData.Session.ScenarioID,
		CurrentStepID:        contextData.Session.CurrentStepID,
		SessionStatus:        contextData.Session.Status,
		Score:                contextData.Session.Score,
		HintLevel:            hintLevel,
		CurrentStep:          currentStep,
		History:              history,
		RecentTerminalOutput: truncateString(strings.TrimSpace(recentTerminalOutput), 4000),
	}
}

func buildFallbackHint(contextData *SessionContext, cause error) *MLHintResponse {
	source := "fallback"

	confidence := 0.30

	if contextData.Step != nil && strings.TrimSpace(contextData.Step.Hint) != "" {
		return &MLHintResponse{
			Hint:       contextData.Step.Hint,
			Confidence: &confidence,
			Source:     source,
			Reasoning:  "ML service unavailable: " + cause.Error(),
		}
	}

	return &MLHintResponse{
		Hint:       "Посмотри на текущий шаг, историю команд и попробуй определить следующий диагностический шаг.",
		Confidence: &confidence,
		Source:     source,
		Reasoning:  "ML service unavailable: " + cause.Error(),
	}
}

func truncateString(value string, limit int) string {
	if limit <= 0 {
		return ""
	}

	if len(value) <= limit {
		return value
	}

	return value[len(value)-limit:]
}
