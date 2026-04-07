package main

import (
	"fmt"
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

type Cache struct {
	mu   sync.Mutex
	data map[string]entry
}

func NewCache() *Cache {
	return &Cache{data: make(map[string]entry)}
}

// Get возвращает значение по ключу. Не возвращает expired записи.
func (c *Cache) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	e, ok := c.data[key]
	if !ok {
		return "", false
	}
	if time.Now().After(e.expiresAt) {
		c.deleteLocked(key)
		return "", false
	}
	return e.value, true
}

// Set устанавливает значение с TTL.
func (c *Cache) Set(key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	fmt.Printf("Set: %s=%s\n", key, value)
}

// Delete удаляет ключ из кэша.
func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// deleteLocked удаляет ключ из кэша (вспомогательная функция)
func (c *Cache) deleteLocked(key string) {
	delete(c.data, key)
}


// Keys возвращает все ключи (включая expired).
func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()

	keys := make([]string, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}
	return keys
}

// Cleanup удаляет все expired записи, возвращает кол-во удалённых.
func (c *Cache) Cleanup() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	deleted := 0
	now := time.Now()
	for key, e := range c.data {
		if now.After(e.expiresAt) {
			c.deleteLocked(key)
			deleted++
		}
	}
	return deleted
}

func main() {
	cache := NewCache()

	cache.Set("a", "1", 50*time.Millisecond)  // expires fast
	cache.Set("b", "2", 10*time.Second)        // lives long
	cache.Set("c", "3", 50*time.Millisecond)   // expires fast

	time.Sleep(100 * time.Millisecond) // wait for a, c to expire

	fmt.Println("Keys before cleanup:", cache.Keys())

	deleted := cache.Cleanup()
	fmt.Printf("Cleanup: deleted %d expired keys\n", deleted)

	fmt.Println("Keys after cleanup:", cache.Keys())

	val, found := cache.Get("b")
	fmt.Printf("Get b: %s, found: %v\n", val, found)

	val, found = cache.Get("a")
	fmt.Printf("Get a: %s, found: %v\n", val, found)
}
