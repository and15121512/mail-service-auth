package data_file

import (
	"context"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/utils"
	"go.uber.org/zap"
)

type DataFile struct {
	logger *zap.SugaredLogger
}

func New(ctx context.Context, logger *zap.SugaredLogger, pgconn string) (*DataFile, error) {
	return &DataFile{logger}, nil
}

func (db *DataFile) annotatedLogger(ctx context.Context) *zap.SugaredLogger {
	request_id, _ := ctx.Value(utils.CtxKeyRequestIDGet()).(string)
	method, _ := ctx.Value(utils.CtxKeyMethodGet()).(string)
	url, _ := ctx.Value(utils.CtxKeyURLGet()).(string)

	return db.logger.With(
		"request_id", request_id,
		"method", method,
		"url", url,
	)
}
