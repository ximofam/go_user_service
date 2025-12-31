package repository

import (
	"github.com/ximofam/user-service/internal/db"
	"github.com/ximofam/user-service/internal/model"
)

const (
	scanUserSelectQuery = "SELECT id, email, username, password, role, created_at, updated_at "
	userTable           = "users "
)

func scanUser(scanner db.SQLScanner) (*model.User, error) {
	user := model.User{}
	err := scanner.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
