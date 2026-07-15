package sandbox

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrSandboxRecordNotFound = errors.New("sandbox record not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, sandbox *Sandbox) (*Sandbox, error) {
	row := r.db.QueryRow(ctx, `
		INSERT INTO session_sandboxes (
			session_id,
			scenario_id,
			container_name,
			image,
			status
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			session_id,
			scenario_id,
			container_name,
			image,
			status,
			started_at,
			stopped_at,
			last_seen_at,
			created_at,
			updated_at
	`,
		sandbox.SessionID,
		sandbox.ScenarioID,
		sandbox.ContainerName,
		sandbox.Image,
		sandbox.Status,
	)

	return scanSandbox(row)
}

func (r *Repository) GetBySessionID(ctx context.Context, sessionID uuid.UUID) (*Sandbox, error) {
	row := r.db.QueryRow(ctx, `
		SELECT
			id,
			session_id,
			scenario_id,
			container_name,
			image,
			status,
			started_at,
			stopped_at,
			last_seen_at,
			created_at,
			updated_at
		FROM session_sandboxes
		WHERE session_id = $1
	`, sessionID)

	return scanSandbox(row)
}

func (r *Repository) MarkRunning(ctx context.Context, sessionID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE session_sandboxes
		SET
			status = $2,
			stopped_at = NULL,
			last_seen_at = NOW(),
			updated_at = NOW()
		WHERE session_id = $1
	`, sessionID, SandboxStatusRunning)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrSandboxRecordNotFound
	}

	return nil
}

func (r *Repository) Touch(ctx context.Context, sessionID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE session_sandboxes
		SET
			last_seen_at = NOW(),
			updated_at = NOW()
		WHERE session_id = $1
	`, sessionID)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrSandboxRecordNotFound
	}

	return nil
}

func (r *Repository) MarkStopped(ctx context.Context, sessionID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE session_sandboxes
		SET
			status = $2,
			stopped_at = NOW(),
			updated_at = NOW()
		WHERE session_id = $1
	`, sessionID, SandboxStatusStopped)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrSandboxRecordNotFound
	}

	return nil
}

func (r *Repository) MarkFailed(ctx context.Context, sessionID uuid.UUID) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE session_sandboxes
		SET
			status = $2,
			updated_at = NOW()
		WHERE session_id = $1
	`, sessionID, SandboxStatusFailed)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrSandboxRecordNotFound
	}

	return nil
}

func (r *Repository) ListRunningOlderThan(ctx context.Context, olderThan time.Duration) ([]Sandbox, error) {
	cutoff := time.Now().Add(-olderThan)

	rows, err := r.db.Query(ctx, `
		SELECT
			id,
			session_id,
			scenario_id,
			container_name,
			image,
			status,
			started_at,
			stopped_at,
			last_seen_at,
			created_at,
			updated_at
		FROM session_sandboxes
		WHERE status = $1
		  AND last_seen_at < $2
		ORDER BY last_seen_at ASC
	`, SandboxStatusRunning, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]Sandbox, 0)

	for rows.Next() {
		item, err := scanSandbox(rows)
		if err != nil {
			return nil, err
		}

		items = append(items, *item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanSandbox(row scanner) (*Sandbox, error) {
	var item Sandbox

	err := row.Scan(
		&item.ID,
		&item.SessionID,
		&item.ScenarioID,
		&item.ContainerName,
		&item.Image,
		&item.Status,
		&item.StartedAt,
		&item.StoppedAt,
		&item.LastSeenAt,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrSandboxRecordNotFound
		}

		return nil, err
	}

	return &item, nil
}
