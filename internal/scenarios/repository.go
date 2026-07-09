package scenarios

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrScenarioNotFound = errors.New("scenario not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetAllActive(ctx context.Context) ([]Scenario, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			title,
			description,
			difficulty,
			category,
			estimated_minutes,
			user_notification,
			desktop_symptoms,
			terminal_solution,
			quick_fix,
			is_active,
			created_at,
			updated_at
		FROM scenarios
		WHERE is_active = TRUE
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Scenario, 0)

	for rows.Next() {
		var scenario Scenario

		err := rows.Scan(
			&scenario.ID,
			&scenario.Title,
			&scenario.Description,
			&scenario.Difficulty,
			&scenario.Category,
			&scenario.EstimatedMinutes,
			&scenario.UserNotification,
			&scenario.DesktopSymptoms,
			&scenario.TerminalSolution,
			&scenario.QuickFix,
			&scenario.IsActive,
			&scenario.CreatedAt,
			&scenario.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, scenario)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) GetByID(ctx context.Context, id string) (*Scenario, error) {
	var scenario Scenario

	err := r.db.QueryRow(ctx, `
		SELECT
			id,
			title,
			description,
			difficulty,
			category,
			estimated_minutes,
			user_notification,
			desktop_symptoms,
			terminal_solution,
			quick_fix,
			is_active,
			created_at,
			updated_at
		FROM scenarios
		WHERE id = $1 AND is_active = TRUE
	`, id).Scan(
		&scenario.ID,
		&scenario.Title,
		&scenario.Description,
		&scenario.Difficulty,
		&scenario.Category,
		&scenario.EstimatedMinutes,
		&scenario.UserNotification,
		&scenario.DesktopSymptoms,
		&scenario.TerminalSolution,
		&scenario.QuickFix,
		&scenario.IsActive,
		&scenario.CreatedAt,
		&scenario.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrScenarioNotFound
		}

		return nil, err
	}

	return &scenario, nil
}
