package hints

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrHintSessionNotFound = errors.New("hint session not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetSessionContext(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
) (*SessionContext, error) {
	session, err := r.getSession(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}

	var step *StepInfo

	if session.CurrentStepID != nil {
		step, err = r.getStep(ctx, *session.CurrentStepID)
		if err != nil {
			return nil, err
		}
	}

	history, err := r.listActionHistory(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return &SessionContext{
		Session: *session,
		Step:    step,
		History: history,
	}, nil
}

func (r *Repository) SaveHint(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
	scenarioID string,
	stepID *string,
	hintLevel string,
	requestPayload MLHintRequest,
	responsePayload MLHintResponse,
) error {
	requestJSON, err := json.Marshal(requestPayload)
	if err != nil {
		return err
	}

	responseJSON, err := json.Marshal(responsePayload)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, `
		INSERT INTO session_hints (
			session_id,
			user_id,
			scenario_id,
			step_id,
			hint_level,
			request_payload,
			response_payload,
			hint,
			confidence,
			source
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		sessionID,
		userID,
		scenarioID,
		stepID,
		hintLevel,
		requestJSON,
		responseJSON,
		responsePayload.Hint,
		responsePayload.Confidence,
		responsePayload.Source,
	)

	return err
}

func (r *Repository) getSession(
	ctx context.Context,
	sessionID uuid.UUID,
	userID uuid.UUID,
) (*SessionInfo, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id::text,
			user_id::text,
			scenario_id,
			current_step_id,
			status,
			score
		FROM sessions
		WHERE id = $1
		  AND user_id = $2
	`, sessionID, userID)

	var item SessionInfo
	var currentStepID sql.NullString

	err := row.Scan(
		&item.ID,
		&item.UserID,
		&item.ScenarioID,
		&currentStepID,
		&item.Status,
		&item.Score,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrHintSessionNotFound
		}

		return nil, err
	}

	if currentStepID.Valid {
		item.CurrentStepID = &currentStepID.String
	}

	return &item, nil
}

func (r *Repository) getStep(ctx context.Context, stepID string) (*StepInfo, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id,
			title,
			description,
			COALESCE(hint, ''),
			COALESCE(expected_command, ''),
			COALESCE(expected_result, '')
		FROM scenario_steps
		WHERE id = $1
	`, stepID)

	var item StepInfo

	err := row.Scan(
		&item.ID,
		&item.Title,
		&item.Description,
		&item.Hint,
		&item.ExpectedCommand,
		&item.ExpectedResult,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) listActionHistory(
	ctx context.Context,
	sessionID uuid.UUID,
) ([]ActionHistoryItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			step_id,
			command,
			is_correct,
			points,
			COALESCE(feedback, ''),
			created_at
		FROM actions
		WHERE session_id = $1
		ORDER BY created_at ASC
	`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]ActionHistoryItem, 0)

	for rows.Next() {
		var item ActionHistoryItem

		err := rows.Scan(
			&item.StepID,
			&item.Command,
			&item.IsCorrect,
			&item.Points,
			&item.Feedback,
			&item.CreatedAt,
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
