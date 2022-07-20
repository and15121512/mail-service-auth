package ports

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/models"
)

type UserStorage interface {
	Get(ctx context.Context, login string) (*models.User, error)
}
