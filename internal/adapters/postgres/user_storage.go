package postgres

import (
	"context"
	"fmt"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/errors"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/ports"
)

var _ ports.UserStorage = (*Database)(nil)

func (db *Database) Get(ctx context.Context, login string) (*models.User, error) {
	logger := db.annotatedLogger(ctx)
	var user models.User

	rows, err := db.DB.Query(ctx, "SELECT users.login as login, users.Password as Password FROM users WHERE users.login = $1", login)
	if err != nil {
		logger.Errorf("query exec failed: %s", err)
		return nil, fmt.Errorf("query exec failed: %s", err)
	}

	if !rows.Next() {
		return nil, errors.ErrNotFound
	}

	err = rows.Scan(&user.Login, &user.PasswordHash)
	if err != nil {
		logger.Errorf("scan exec failed: %s", err)
		return nil, fmt.Errorf("scan exec failed: %s", err)
	}

	return &user, nil
}
