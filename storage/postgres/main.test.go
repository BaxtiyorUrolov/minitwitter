package postgres

import (
	"context"
	"testing"
	"twitter/config"
	"twitter/pkg/logger"
	"twitter/storage"
)

func setup(t *testing.T) (storage.IStorage, context.Context) {
	cfg := config.Load()
	log := logger.New(cfg.ServiceName)

	store, err := New(context.Background(), cfg, log)
	if err != nil {
		t.Fatalf("Error while connecting to DB: %v", err)
	}

	return store, context.Background()
}
