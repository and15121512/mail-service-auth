package data_file

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/config"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
)

func (*DataFile) Get(ctx context.Context, login string) (*models.User, error) {
	user := &models.User{
		Login:        config.GetConfig().Auth.Login,
		PasswordHash: config.GetConfig().Auth.PasswordHash,
	}
	return user, nil
}
