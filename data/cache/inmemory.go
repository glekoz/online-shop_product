package cache

import (
	"errors"
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

	//delCount  int32 // счетчик удаленных записей
	//threshold int32 // порог, при котором надо пересоздавать очередь и сам кэш
}

func (c *inMemoryCache[K, V]) Add(key K, value V, ttl time.Duration) error {
	if ttl < c.step {
		return errors.New("ttl must be more than cache check step")
	}
	t := time.Now().Truncate(c.step)
	ttl = ttl.Truncate(c.step)
	t = t.Add(ttl)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value

	keys, ok := c.queue[t]
	if !ok {
		c.queue[t] = make([]K, 0, c.queueSize)
	}
	c.queue[t] = append(keys, key)
	return nil
}

func (c inMemoryCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	value, ok := c.cache[key]
	c.mu.RUnlock()
	return value, ok
}

func (c *inMemoryCache[K, V]) Delete(key K) {
	c.mu.Lock()
	delete(c.cache, key)
	c.mu.Unlock()
}

/*
func (c *inMemoryCache[K, V]) StartGC() { // а нужен ли StopGC? Если будет нужен, то добавлю булево поле и сделаю while цикл
	go func ()  {
		for
	}()
}
*/
