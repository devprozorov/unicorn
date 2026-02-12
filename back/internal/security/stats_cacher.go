package security

import (
	"sync"
	"time"

	"unicorn-auth/internal/models"
)

// StatsCacher кэширует статистику главной страницы
type StatsCacher struct {
	mu            sync.RWMutex
	cached        *models.HomeStats
	cacheDuration time.Duration
	lastCacheTime time.Time
}

// NewStatsCacher создает новый кэширующий сервис со временем кэша по умолчанию (30 секунд)
func NewStatsCacher() *StatsCacher {
	return &StatsCacher{
		cacheDuration: 30 * time.Second,
	}
}

// NewStatsCacherWithDuration создает новый кэширующий сервис с пользовательским временем кэша
func NewStatsCacherWithDuration(duration time.Duration) *StatsCacher {
	return &StatsCacher{
		cacheDuration: duration,
	}
}

// Get возвращает закэшированную статистику, если она еще валидна, иначе nil
func (c *StatsCacher) Get() *models.HomeStats {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.cached != nil && time.Since(c.lastCacheTime) < c.cacheDuration {
		return c.cached
	}
	return nil
}

// Set кэширует новую статистику
func (c *StatsCacher) Set(stats *models.HomeStats) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cached = stats
	c.lastCacheTime = time.Now()
}

// IsCached проверяет, кэширована ли статистика и она еще валидна
func (c *StatsCacher) IsCached() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.cached != nil && time.Since(c.lastCacheTime) < c.cacheDuration
}

// Invalidate инвалидирует кэш
func (c *StatsCacher) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cached = nil
	c.lastCacheTime = time.Time{}
}
