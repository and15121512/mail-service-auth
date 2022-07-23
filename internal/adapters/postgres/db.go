package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/utils"
	"go.uber.org/zap"
)

type Database struct {
	DB     *pgxpool.Pool
	logger *zap.SugaredLogger
}

func New(ctx context.Context, logger *zap.SugaredLogger, pgconn string) (*Database, error) {
	config, err := pgxpool.ParseConfig(pgconn)
	if err != nil {
		logger.Errorf("postgres connection string parse failed: %s", err.Error())
		return nil, fmt.Errorf("postgres connection string parse failed: %s", err.Error())
	}

	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		logger.Errorf("create pool failed: %s", err.Error())
		return nil, fmt.Errorf("create pool failed: %s", err.Error())
	}

	return &Database{DB: pool, logger: logger}, nil
}

func (db *Database) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return db.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}
