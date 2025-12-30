package repository

import (
	"context"
	"database/sql"

	"github.com/ximofam/user-service/internal/db"
	"github.com/ximofam/user-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UpdatePassword(ctx context.Context, userID uint, newPassword string) error
}

type userRepo struct {
	*db.DBProvider
}

func NewUserRepository(dbProvider *db.DBProvider) *userRepo {
	return &userRepo{
		DBProvider: dbProvider,
	}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (email, password) VALUES(?,?)"
	if _, err := r.GetQuerier(ctx).ExecContext(ctx, query, user.Email, user.Password); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*model.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE id = ? AND deleted_at IS NULL
	`
	user := model.User{}
	err := r.GetQuerier(ctx).QueryRowContext(ctx, query, id).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = ? AND deleted_at IS NULL
	`
	user := model.User{}
	err := r.GetQuerier(ctx).QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, userID uint, newPassword string) error {
	query := "UPDATE users SET password = ? WHERE id = ?"
	res, err := r.GetQuerier(ctx).ExecContext(ctx, query, newPassword, userID)
	if err != nil {
		return err
	}

	row, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if row == 0 {
		return sql.ErrNoRows
	}

	return nil
}
