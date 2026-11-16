package repository

import (
	"context"
	"database/sql"

	"github.com/dotuan0601/go-backend-template/internal/model"
)

type UserRepository interface {
	GetByID(
		ctx context.Context,
		id int64,
	) (*model.User, error)

	GetAll(
		ctx context.Context,
	) ([]*model.User, error)

	Create(
		ctx context.Context,
		user *model.User,
	) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(
	db *sql.DB,
) UserRepository {

	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByID(
	ctx context.Context,
	id int64,
) (*model.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			created_at,
			updated_at
		FROM users
		WHERE id = ?
	`

	user := &model.User{}

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepository) GetAll(
	ctx context.Context,
) ([]*model.User, error) {

	query := `
		SELECT
			id,
			name,
			email,
			created_at,
			updated_at
		FROM users
		ORDER BY id DESC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make(
		[]*model.User,
		0,
	)

	for rows.Next() {

		user := &model.User{}

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(
			users,
			user,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) Create(
	ctx context.Context,
	user *model.User,
) error {

	query := `
		INSERT INTO users (
			name,
			email
		)
		VALUES (?, ?)
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		user.Name,
		user.Email,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	user.ID = id

	return nil
}
