package data_file

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
)

func (db *DataFile) Get(ctx context.Context, login string) (*models.User, error) {
	logger := db.annotatedLogger(ctx)

	user := &models.User{
		Login:        config.GetConfig(logger).Auth.Login,
		PasswordHash: config.GetConfig(logger).Auth.PasswordHash,
	}
	return user, nil
}
