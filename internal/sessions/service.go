package sessions

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrInvalidScenarioID = errors.New("invalid scenario id")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, req CreateSessionRequest) (*SessionResponse, error) {
	scenarioID := strings.TrimSpace(req.ScenarioID)
	if scenarioID == "" {
		return nil, ErrInvalidScenarioID
	}

	exists, err := s.repo.ScenarioExists(ctx, scenarioID)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, ErrScenarioNotFound
	}

	firstStepID, err := s.repo.GetFirstStepID(ctx, scenarioID)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	session := Session{
		ID:            uuid.New(),
		UserID:        userID,
		ScenarioID:    scenarioID,
		CurrentStepID: firstStepID,
		Status:        StatusInProgress,
		Score:         0,
		StartedAt:     now,
		FinishedAt:    nil,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(ctx, session); err != nil {
		return nil, err
	}

	response := ToSessionResponse(session)

	return &response, nil
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*SessionResponse, error) {
	session, err := s.repo.GetByIDAndUserID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	response := ToSessionResponse(*session)

	return &response, nil
}
