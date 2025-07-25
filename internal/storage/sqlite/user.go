package sqlite

import (
	"blog/internal/model"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) Create(ctx context.Context, u *model.User) (uint64, error) {
	const op = "create"
	var wrap = wrapUnique
	query := `INSERT INTO users(email, username, password) 
		VALUES (?, ?, ?) RETURNING id`
	row := r.db.QueryRowContext(ctx, query, u.Email, u.Username, u.Password)
	var id uint64
	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, wrap(err))
	}
	return id, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint64) (uint64, error) {
	return 0, nil
}

func (r *UserRepository) FindByID(ctx context.Context, id uint64) (*model.User, error) {
	const op = "find by id"
	var wrap = wrapNoRows
	u := model.User{}
	var d *string
	query := `SELECT id, email, username, description, password, avatar FROM users
		WHERE (id == ?)`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&u.ID, &u.Email, &u.Username, &d, &u.Password, &u.Avatar)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	if d == nil {
		u.Description = "no description"
	} else {
		u.Description = *d
	}
	return &u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	const op = "find by email"
	var wrap = wrapNoRows
	u := model.User{}
	var d *string
	query := `SELECT id, email, username, description, password, avatar FROM users
		WHERE (email == ?)`
	row := r.db.QueryRowContext(ctx, query, email)
	err := row.Scan(&u.ID, &u.Email, &u.Username, &d, &u.Password, &u.Avatar)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, wrap(err))
	}
	if d == nil {
		u.Description = "no description"
	} else {
		u.Description = *d
	}
	return &u, nil
}

// Update обновляет информацию пользователя
func (r *UserRepository) Update(ctx context.Context, u *model.User) error {
	const op = "update"
	var wrap = wrapNotNull
	var description *string
	if strings.TrimSpace(u.Description) == "" {
		description = nil
	} else {
		description = &u.Description
	}
	query := `UPDATE users SET username = ?, description = ?, avatar = ? WHERE id == ?`
	_, err := r.db.ExecContext(ctx, query, u.Username, description, u.Avatar, u.ID)

	if err != nil {
		return fmt.Errorf("%s: %w", op, wrap(err))
	}
	return nil
}
