package sessions

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSessionNotFound  = errors.New("session not found")
	ErrScenarioNotFound = errors.New("scenario not found")
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) ScenarioExists(ctx context.Context, scenarioID string) (bool, error) {
	var exists bool

	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM scenarios
			WHERE id = $1 AND is_active = TRUE
		)
	`, scenarioID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Repository) GetFirstStepID(ctx context.Context, scenarioID string) (*string, error) {
	var stepID string

	err := r.db.QueryRow(ctx, `
		SELECT id
		FROM scenario_steps
		WHERE scenario_id = $1
		ORDER BY step_order ASC
		LIMIT 1
	`, scenarioID).Scan(&stepID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &stepID, nil
}

func (r *Repository) Create(ctx context.Context, session Session) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO sessions (
			id,
			user_id,
			scenario_id,
			current_step_id,
			status,
			score,
			started_at,
			finished_at,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
		)
	`,
		session.ID,
		session.UserID,
		session.ScenarioID,
		session.CurrentStepID,
		session.Status,
		session.Score,
		session.StartedAt,
		session.FinishedAt,
		session.CreatedAt,
		session.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetByIDAndUserID(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*Session, error) {
	var session Session

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			user_id,
			scenario_id,
			current_step_id,
			status,
			score,
			started_at,
			finished_at,
			created_at,
			updated_at
		FROM sessions
		WHERE id = $1 AND user_id = $2
	`, id, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.ScenarioID,
		&session.CurrentStepID,
		&session.Status,
		&session.Score,
		&session.StartedAt,
		&session.FinishedAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSessionNotFound
		}

		return nil, err
	}

	return &session, nil
}

func (r *Repository) GetHistoryByUserID(ctx context.Context, userID uuid.UUID) ([]SessionHistoryItem, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			s.id,
			s.scenario_id,
			sc.title,
			sc.difficulty,
			sc.category,
			s.status,
			s.score,
			s.started_at,
			s.finished_at
		FROM sessions s
		INNER JOIN scenarios sc ON sc.id = s.scenario_id
		WHERE s.user_id = $1
		ORDER BY s.started_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]SessionHistoryItem, 0)

	for rows.Next() {
		var item SessionHistoryItem

		err := rows.Scan(
			&item.ID,
			&item.ScenarioID,
			&item.ScenarioTitle,
			&item.Difficulty,
			&item.Category,
			&item.Status,
			&item.Score,
			&item.StartedAt,
			&item.FinishedAt,
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
