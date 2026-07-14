package auth

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"DCS/internal/users"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, user users.User) error {
	query := `
		INSERT INTO users (id, email, name, password_hash, role)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(
		ctx,
		query,
		user.ID,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Role,
	)

	return err
}

func (r *Repository) GetUserByEmail(ctx context.Context, email string) (*users.User, error) {
	query := `
		SELECT id, email, name, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user users.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return &user, nil
}
