package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/ximofam/user-service/internal/db"
	"github.com/ximofam/user-service/internal/dto"
	"github.com/ximofam/user-service/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	UpdatePassword(ctx context.Context, userID uint, newPassword string) error
	IsExists(ctx context.Context, field, value string) error
	GetAll(ctx context.Context, input *dto.ListUsersInput) ([]model.User, *dto.PagingOutput, error)
	SoftDelete(ctx context.Context, id uint) error
}

type userRepo struct {
	*db.Database
}

func NewUserRepository(database *db.Database) *userRepo {
	return &userRepo{
		Database: database,
	}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	query := "INSERT INTO users (email, username, password) VALUES(?,?,?)"
	if _, err := r.Querier(ctx).ExecContext(ctx, query, user.Email, user.Username, user.Password); err != nil {
		return err
	}

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uint) (*model.User, error) {
	query := fmt.Sprintf("%s FROM %s WHERE id = ? AND deleted_at IS NULL", scanUserSelectQuery, userTable)

	user, err := scanUser(r.Querier(ctx).QueryRowContext(ctx, query, id))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := fmt.Sprintf("%s FROM %s WHERE email = ? AND deleted_at IS NULL", scanUserSelectQuery, userTable)

	user, err := scanUser(r.Querier(ctx).QueryRowContext(ctx, query, email))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, userID uint, newPassword string) error {
	query := fmt.Sprintf("UPDATE %s SET password = ? WHERE id = ?", userTable)
	res, err := r.Querier(ctx).ExecContext(ctx, query, newPassword, userID)
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

func (r *userRepo) IsExists(ctx context.Context, field, value string) error {
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE %s = ? LIMIT 1", userTable, field)
	count := 0
	if err := r.Querier(ctx).QueryRowContext(ctx, query, value).Scan(&count); err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("%s:%s does not exists", field, value)
	}

	return nil
}

func (r *userRepo) GetAll(ctx context.Context, input *dto.ListUsersInput) ([]model.User, *dto.PagingOutput, error) {
	query := strings.Builder{}
	args := []any{}

	query.WriteString(fmt.Sprintf("FROM %s WHERE 1 ", userTable))
	if input.Username != "" {
		query.WriteString("AND username LIKE ? ")
		args = append(args, "%"+input.Username+"%")
	}

	query.WriteString("AND deleted_at IS NULL ")

	query.WriteString(fmt.Sprintf("ORDER BY %s %s ", input.GetSortColumn(), input.GetSortDirection()))

	query.WriteString("LIMIT ? OFFSET ? ")
	args = append(args, input.GetLimit(), input.GetOffset())

	querier := r.Querier(ctx)

	selectQuery := scanUserSelectQuery + query.String()
	log.Println(selectQuery)
	rows, err := querier.QueryContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	res := []model.User{}
	for rows.Next() {
		user, err := scanUser(rows)
		if err != nil {
			return nil, nil, err
		}

		res = append(res, *user)
	}

	count := 0
	countQuery := "SELECT COUNT(*) " + query.String()
	err = querier.QueryRowContext(ctx, countQuery, args...).Scan(&count)
	if err != nil {
		return nil, nil, err
	}

	totalPages := math.Ceil(float64(count) / float64(input.GetLimit()))

	return res,
		&dto.PagingOutput{
			CurrentPage: input.GetPage(),
			TotalPages:  int(totalPages),
			Limit:       input.GetLimit(),
			Total:       int64(count),
		},
		nil
}

func (r *userRepo) SoftDelete(ctx context.Context, id uint) error {
	query := fmt.Sprintf("UPDATE %s SET deleted_at = NOW() WHERE id = ?", userTable)
	res, err := r.Querier(ctx).ExecContext(ctx, query, id)
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
