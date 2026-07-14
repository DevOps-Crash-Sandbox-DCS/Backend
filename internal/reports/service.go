package reports

import (
	"context"

	"github.com/google/uuid"
)

const pointsPerStep = 10

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetSessionReport(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (*SessionReportResponse, error) {
	session, err := s.repo.GetSessionByIDAndUserID(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}

	scenario, err := s.repo.GetScenarioByID(ctx, session.ScenarioID)
	if err != nil {
		return nil, err
	}

	steps, err := s.repo.GetStepsByScenarioID(ctx, session.ScenarioID)
	if err != nil {
		return nil, err
	}

	actions, err := s.repo.GetActionsBySessionID(ctx, session.ID)
	if err != nil {
		return nil, err
	}

	response := buildReportResponse(*session, *scenario, steps, actions)

	return &response, nil
}

func buildReportResponse(
	session Session,
	scenario Scenario,
	steps []ScenarioStep,
	actions []Action,
) SessionReportResponse {
	actionsByStepID := make(map[string][]ReportActionResponse)

	totalActions := len(actions)
	correctActions := 0
	wrongActions := 0

	for _, action := range actions {
		if action.IsCorrect {
			correctActions++
		} else {
			wrongActions++
		}

		actionResponse := ReportActionResponse{
			ID:        action.ID,
			Command:   action.Command,
			IsCorrect: action.IsCorrect,
			Points:    action.Points,
			Feedback:  action.Feedback,
			CreatedAt: action.CreatedAt,
		}

		actionsByStepID[action.StepID] = append(actionsByStepID[action.StepID], actionResponse)
	}

	stepResponses := make([]ReportStepResponse, 0, len(steps))
	completedSteps := 0

	for _, step := range steps {
		stepActions := actionsByStepID[step.ID]

		for _, action := range stepActions {
			if action.IsCorrect {
				completedSteps++
				break
			}
		}

		expectedCommand := step.ExpectedCommand
		if session.Status != "completed" {
			expectedCommand = ""
		}

		stepResponses = append(stepResponses, ReportStepResponse{
			ID:              step.ID,
			Order:           step.Order,
			Title:           step.Title,
			Description:     step.Description,
			ExpectedCommand: expectedCommand,
			Actions:         stepActions,
		})
	}

	maxScore := len(steps) * pointsPerStep
	isCompleted := session.Status == "completed"

	return SessionReportResponse{
		Session: ReportSessionResponse{
			ID:         session.ID,
			Status:     session.Status,
			Score:      session.Score,
			StartedAt:  session.StartedAt,
			FinishedAt: session.FinishedAt,
		},
		Scenario: ReportScenarioResponse{
			ID:         scenario.ID,
			Title:      scenario.Title,
			Difficulty: scenario.Difficulty,
			Category:   scenario.Category,
		},
		Summary: ReportSummaryResponse{
			TotalSteps:     len(steps),
			CompletedSteps: completedSteps,
			TotalActions:   totalActions,
			CorrectActions: correctActions,
			WrongActions:   wrongActions,
			Score:          session.Score,
			MaxScore:       maxScore,
			IsCompleted:    isCompleted,
		},
		Steps: stepResponses,
	}
}
