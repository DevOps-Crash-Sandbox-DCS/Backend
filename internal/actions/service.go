package actions

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidSessionState = errors.New("invalid session state")
	ErrInvalidStep         = errors.New("invalid step")
	ErrInvalidCommand      = errors.New("invalid command")
)

const (
	pointsForCorrectAction = 10

	statusInProgress = "in_progress"
	statusCompleted  = "completed"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Submit(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
	req SubmitActionRequest,
) (*SubmitActionResponse, error) {
	stepID := strings.TrimSpace(req.StepID)
	command := strings.TrimSpace(req.Command)

	if stepID == "" {
		return nil, ErrInvalidStep
	}

	if command == "" {
		return nil, ErrInvalidCommand
	}

	session, err := s.repo.GetSessionByIDAndUserID(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}

	if session.Status != statusInProgress {
		return nil, ErrInvalidSessionState
	}

	if session.CurrentStepID == nil {
		return nil, ErrInvalidSessionState
	}

	if *session.CurrentStepID != stepID {
		return nil, ErrInvalidStep
	}

	step, err := s.repo.GetStepByID(ctx, stepID)
	if err != nil {
		return nil, err
	}

	if step.ScenarioID != session.ScenarioID {
		return nil, ErrInvalidStep
	}

	acceptedCommands, err := s.repo.GetAcceptedCommandsByStepID(ctx, stepID)
	if err != nil {
		return nil, err
	}

	isCorrect := IsCommandAccepted(
		command,
		step.ExpectedCommand,
		acceptedCommands,
	)

	points := 0
	feedback := "Команда неверная. Проверьте текущий шаг и попробуйте еще раз."
	nextStepID := session.CurrentStepID
	sessionStatus := statusInProgress
	sessionScore := session.Score

	if isCorrect {
		points = pointsForCorrectAction
		feedback = "Команда выполнена верно."

		next, err := s.repo.GetNextStepID(ctx, session.ScenarioID, step.Order)
		if err != nil {
			return nil, err
		}

		nextStepID = next
		sessionScore = session.Score + points

		if next == nil {
			sessionStatus = statusCompleted
		}
	}

	action := Action{
		ID:        uuid.New(),
		SessionID: session.ID,
		StepID:    step.ID,
		Command:   command,
		IsCorrect: isCorrect,
		Points:    points,
		Feedback:  feedback,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateAction(ctx, action); err != nil {
		return nil, err
	}

	if isCorrect {
		if err := s.repo.UpdateSessionProgress(ctx, session.ID, nextStepID, sessionStatus, sessionScore); err != nil {
			return nil, err
		}
	}

	return &SubmitActionResponse{
		ID:            action.ID,
		SessionID:     action.SessionID,
		StepID:        action.StepID,
		Command:       action.Command,
		IsCorrect:     action.IsCorrect,
		Points:        action.Points,
		Feedback:      action.Feedback,
		NextStepID:    nextStepID,
		SessionStatus: sessionStatus,
		SessionScore:  sessionScore,
	}, nil
}
