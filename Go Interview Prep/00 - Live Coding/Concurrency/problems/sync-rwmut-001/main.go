package main

import (
	"fmt"
	"sync"
	// TODO: добавь нужные импорты
)

// ConfigStore — потокобезопасное хранилище конфигурации.
type ConfigStore struct {
	mu   sync.RWMutex
	data map[string]string
}

func NewConfigStore() *ConfigStore {
	return &ConfigStore{data: make(map[string]string)}
}

// Get возвращает значение по ключу (concurrent-safe).
func (c *ConfigStore) Get(key string) (string, bool) {
	// TODO: реализуй
	return "", false
}

// Set устанавливает значение по ключу (concurrent-safe).
func (c *ConfigStore) Set(key, value string) {
	// TODO: реализуй
}

// GetAll возвращает КОПИЮ всего конфига (concurrent-safe).
func (c *ConfigStore) GetAll() map[string]string {
	// TODO: реализуй
	return nil
}

func main() {
	store := NewConfigStore()

	// Начальные значения
	store.Set("db_host", "localhost")
	store.Set("db_port", "5432")

	// TODO: запусти 10 readers и 2 writers на 500ms
	// Readers: в цикле читают случайный ключ из ["db_host", "db_port", "app_env", "log_level"]
	// Writers: в цикле пишут случайные пары (напр. "app_env"="production", "log_level"="debug")
	// Через 500ms останови всех, выведи итоговый конфиг через GetAll()

	_ = store
	fmt.Println("Final config:")
	// TODO: выведи итоговый конфиг
}
