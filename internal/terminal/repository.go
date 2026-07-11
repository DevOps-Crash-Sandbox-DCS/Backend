package terminal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetCommandOutputsByStepID(ctx context.Context, stepID string) ([]CommandOutput, error) {
	rows, err := r.db.Query(ctx, `
		SELECT
			id::text,
			step_id,
			command_pattern,
			match_type,
			stdout,
			stderr,
			exit_code,
			description,
			priority,
			is_active
		FROM step_command_outputs
		WHERE step_id = $1
		  AND is_active = TRUE
		ORDER BY priority ASC, created_at ASC
	`, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]CommandOutput, 0)

	for rows.Next() {
		var item CommandOutput

		err := rows.Scan(
			&item.ID,
			&item.StepID,
			&item.CommandPattern,
			&item.MatchType,
			&item.Stdout,
			&item.Stderr,
			&item.ExitCode,
			&item.Description,
			&item.Priority,
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
