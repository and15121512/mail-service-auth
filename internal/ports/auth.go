package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
)

type Auth interface {
	// Info(ctx context.Context, login string) (*models.User, error)
	Validate(ctx context.Context, access_token string) (*models.User, error)
	Login(ctx context.Context, login, password string) (models.TokenPair, error)
}
