package agent

import (
	"context"
	"time"

	"github.com/VladKvetkin/avito-backend-trainee-assignment-2023/internal/storage"
)

func Start(ctx context.Context, storage storage.Storage) error {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	if err := storage.DeleteExpiredUserSegments(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := storage.DeleteExpiredUserSegments(ctx); err != nil {
				return err
			}
		}
	}
}
