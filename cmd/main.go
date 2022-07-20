package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitlab.com/sukharnikov.aa/mail-service-auth/internal/application"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()
	go application.Start(ctx)
	<-ctx.Done()
	application.Stop()
}
