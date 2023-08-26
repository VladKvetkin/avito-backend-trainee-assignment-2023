package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"

	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/config"
	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/server"
	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/services/agent"
	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/storage"
	"go.uber.org/zap"
)

func main() {
	os.Exit(start())
}

func start() int {
	config, err := config.NewConfig()
	if err != nil {
		zap.L().Info("error create config", zap.Error(err))
		return 1
	}

	defer zap.L().Sync()

	db, err := sqlx.Connect("postgres", config.DatabaseDSN)
	if err != nil {
		zap.L().Info("error failed to connect to db: %w", zap.Error(err))
		return 1
	}

	defer db.Close()

	postgresStorage := storage.NewPostgresStorage(db)

	server := server.NewServer(config, postgresStorage)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := server.Start(); err != nil {
			zap.L().Info("error starting server", zap.Error(err))
			return err
		}

		return nil
	})

	eg.Go(func() error {
		if err := agent.Start(ctx, postgresStorage); err != nil {
			zap.L().Info("error start agent", zap.Error(err))
			return err
		}

		return nil
	})

	<-ctx.Done()

	eg.Go(func() error {
		if err := server.Stop(); err != nil {
			zap.L().Info("error stopping server", zap.Error(err))
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return 1
	}

	return 0
}
