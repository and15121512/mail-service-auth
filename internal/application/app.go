package application

import (
	"context"
	"fmt"
	"os"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/adapters/data_file"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/adapters/http"
	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/domain/auth"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

var (
	s      *http.Server
	logger *zap.Logger
)

func Start(ctx context.Context) {
	logger, _ = zap.NewProduction()

	pgconn := os.Getenv("PG_URL")

	db, err := data_file.New(ctx, pgconn) // Mock!!!
	//db, err := postgres.New(ctx, pgconn)

	if err != nil {
		logger.Sugar().Fatalf("db init failed:", err)
	}
	authS := auth.New(db)

	s, err = http.New(logger.Sugar(), authS)
	if err != nil {
		logger.Sugar().Fatalf("http server creating failed:", err)
	}

	var g errgroup.Group
	g.Go(func() error {
		return s.Start()
	})

	logger.Sugar().Info(fmt.Sprintf("app is started on port :%d", s.Port()))
	err = g.Wait()
	if err != nil {
		logger.Sugar().Fatalw("http server start failed", zap.Error(err))
	}
}

func Stop() {
	_ = s.Stop(context.Background())
	logger.Sugar().Info("app has stopped")
}
