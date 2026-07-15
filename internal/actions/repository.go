package actions

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrStepNotFound    = errors.New("step not found")
)

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
			current_step_id,
			status,
			score
		FROM sessions
		WHERE id = $1 AND user_id = $2
	`, sessionID, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.ScenarioID,
		&session.CurrentStepID,
		&session.Status,
		&session.Score,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}

		return nil, err
	}

	return &session, nil
}

func (r *Repository) GetStepByID(ctx context.Context, stepID string) (*ScenarioStep, error) {
	var step ScenarioStep

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			scenario_id,
			step_order,
			expected_command
		FROM scenario_steps
		WHERE id = $1
	`, stepID).Scan(
		&step.ID,
		&step.ScenarioID,
		&step.Order,
		&step.ExpectedCommand,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrStepNotFound
		}

		return nil, err
	}

	return &step, nil
}

func (r *Repository) GetNextStepID(ctx context.Context, scenarioID string, currentOrder int) (*string, error) {
	var nextStepID string

	err := r.db.QueryRow(ctx, `
		SELECT id
		FROM scenario_steps
		WHERE scenario_id = $1 AND step_order > $2
		ORDER BY step_order ASC
		LIMIT 1
	`, scenarioID, currentOrder).Scan(&nextStepID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &nextStepID, nil
}

func (r *Repository) CreateAction(ctx context.Context, action Action) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO actions (
			id,
			session_id,
			step_id,
			command,
			is_correct,
			points,
			feedback,
			created_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
	`,
		action.ID,
		action.SessionID,
		action.StepID,
		action.Command,
		action.IsCorrect,
		action.Points,
		action.Feedback,
		action.CreatedAt,
	)

	return err
}

func (r *Repository) UpdateSessionProgress(
	ctx context.Context,
	sessionID uuid.UUID,
	nextStepID *string,
	status string,
	score int,
) error {
	isCompleted := status == "completed"

	_, err := r.db.Exec(ctx, `
		UPDATE sessions
		SET
			current_step_id = $1,
			status = $2::varchar,
			score = $3,
			finished_at = CASE
				WHEN $5 THEN NOW()
				ELSE finished_at
			END,
			updated_at = NOW()
		WHERE id = $4
	`,
		nextStepID,
		status,
		score,
		sessionID,
		isCompleted,
	)

	return err
}

func (r *Repository) GetAcceptedCommandsByStepID(ctx context.Context, stepID string) ([]AcceptedCommand, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id::text,
			step_id,
			command,
			match_type,
			description,
			is_active
		FROM step_accepted_commands
		WHERE step_id = $1
		  AND is_active = TRUE
		ORDER BY created_at ASC
	`, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]AcceptedCommand, 0)

	for rows.Next() {
		var item AcceptedCommand

		err := rows.Scan(
			&item.ID,
			&item.StepID,
			&item.Command,
			&item.MatchType,
			&item.Description,
			&item.IsActive,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
