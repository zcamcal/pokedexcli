package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries map[string]cacheEntry
	mutex   *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{make(map[string]cacheEntry), &sync.Mutex{}}

	go cache.reapLoop(interval)

	return cache
}

func (cache *Cache) Add(key string, value []byte) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	now := time.Now()
	cache.entries[key] = cacheEntry{
		createdAt: now,
		val:       value,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	val, ok := cache.entries[key]
	if !ok {
		return nil, false
	}

	return val.val, true
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case now := <-ticker.C:
			go cache.cleaner(now)
		}
	}
}

func (cache *Cache) cleaner(time time.Time) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	for key, value := range cache.entries {
		if time.Before(value.createdAt) {
			continue
		}
		delete(cache.entries, key)
	}
}
