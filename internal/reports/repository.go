package reports

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSessionNotFound = errors.New("session not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetSessionByIDAndUserID(ctx context.Context, sessionID uuid.UUID, userID uuid.UUID) (*Session, error) {
	var session Session

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			scenario_id,
			status,
			score,
			started_at,
			finished_at
		FROM sessions
		WHERE id = $1 AND user_id = $2
	`, sessionID, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.ScenarioID,
		&session.Status,
		&session.Score,
		&session.StartedAt,
		&session.FinishedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}

		return nil, err
	}

	return &session, nil
}

func (r *Repository) GetScenarioByID(ctx context.Context, scenarioID string) (*Scenario, error) {
	var scenario Scenario

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			title,
			difficulty,
			category
		FROM scenarios
		WHERE id = $1
	`, scenarioID).Scan(
		&scenario.ID,
		&scenario.Title,
		&scenario.Difficulty,
		&scenario.Category,
	)

	if err != nil {
		return nil, err
	}

	return &scenario, nil
}

func (r *Repository) GetStepsByScenarioID(ctx context.Context, scenarioID string) ([]ScenarioStep, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			scenario_id,
			step_order,
			title,
			description,
			expected_command
		FROM scenario_steps
		WHERE scenario_id = $1
		ORDER BY step_order ASC
	`, scenarioID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := make([]ScenarioStep, 0)

	for rows.Next() {
		var step ScenarioStep

		err := rows.Scan(
			&step.ID,
			&step.ScenarioID,
			&step.Order,
			&step.Title,
			&step.Description,
			&step.ExpectedCommand,
		)
		if err != nil {
			return nil, err
		}

		steps = append(steps, step)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return steps, nil
}

func (r *Repository) GetActionsBySessionID(ctx context.Context, sessionID uuid.UUID) ([]Action, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			session_id,
			step_id,
			command,
			is_correct,
			points,
			feedback,
			created_at
		FROM actions
		WHERE session_id = $1
		ORDER BY created_at ASC
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	actions := make([]Action, 0)

	for rows.Next() {
		var action Action

		err := rows.Scan(
			&action.ID,
			&action.SessionID,
			&action.StepID,
			&action.Command,
			&action.IsCorrect,
			&action.Points,
			&action.Feedback,
			&action.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		actions = append(actions, action)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actions, nil
}
