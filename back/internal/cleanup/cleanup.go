package cleanup

import (
	"context"
	"log"
	"time"

	"unicorn-auth/internal/repo"
)

// Cleaner выполняет периодическую очистку старых данных
type Cleaner struct {
	apps *repo.ApplicationRepo
}

func NewCleaner(apps *repo.ApplicationRepo) *Cleaner {
	return &Cleaner{apps: apps}
}

// Start запускает периодическую очистку
func (c *Cleaner) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Первая очистка сразу при запуске
	c.cleanup(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.cleanup(ctx)
		}
	}
}

func (c *Cleaner) cleanup(ctx context.Context) {
	// Удаляем отклики, скрытые более месяца назад обеими сторонами
	deleted, err := c.apps.DeleteOldHidden(ctx)
	if err != nil {
		log.Printf("Error cleaning old hidden applications: %v", err)
		return
	}

	if deleted > 0 {
		log.Printf("Cleaned up %d old hidden applications", deleted)
	}
}
