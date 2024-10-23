package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"twitter/api"
	"twitter/api/handler"
	"twitter/config"
	"twitter/pkg/logger"
	"twitter/service"
	"twitter/storage/postgres"
)

func main() {

	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	store, err := postgres.New(context.Background(), cfg, log)
	if err != nil {
		log.Error("Error while connecting to DB: %v", logger.Error(err))
		return
	}
	defer store.Close()

	services := service.New(store, log)

	h := handler.New(services, log)

	server := api.New(cfg, services, store, log, h)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		if err := server.Run(); err != nil {
			log.Error("Error run http server", zap.Error(err))
		}
	}()

	<-ctx.Done()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("failed http graceful shutdown", zap.Error(err))
	}
}
