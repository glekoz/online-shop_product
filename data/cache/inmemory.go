package cache

import (
	"sync"
	"time"
)

// увеличиваю потреблении памяти, но уменьшаю потребление CPU при очистке памяти
// минимальный шаг - 1 секунда
/*
type Queue[K comparable] struct {
	ttl  time.Time
	keys []K
}
*/
type inMemoryCache[K comparable, V any] struct {
	mu    sync.RWMutex
	cache map[K]V // тут тип string, а не uuid.UUID, потому что в Redis'e нет такого типа данных

	queue     map[time.Time][]K // надо иногда создавать новую - иначе утечка памяти
	step      time.Duration     // раз в какой период проверять устаревший кэш (соответственно, и шаг добавления в кэш) и не пора ли пересоздать мапу и очередь
	queueSize int               // размер []K - зависить от use case

	delCount  int32 // счетчик удаленных записей
	threshold int32 // порог, при котором надо пересоздавать очередь и сам кэш
}

func (c *inMemoryCache[K, V]) Add(key K, value V, ttl time.Duration) {
	//t := time.Now().Truncate(c.step)
	if ttl < c.step {
		ttl = c.step
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value

	//if c.queue
}
