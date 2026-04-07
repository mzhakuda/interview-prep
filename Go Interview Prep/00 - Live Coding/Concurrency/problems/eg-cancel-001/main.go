package main

import (
	"context"
	"fmt"
	"time"
	// TODO: добавь нужные импорты
)

// fetchService имитирует запрос к сервису. Не меняй эту функцию.
func fetchService(ctx context.Context, name string) (string, error) {
	// Имитация: orders всегда падает через 100ms, остальные отвечают за 80-200ms
	durations := map[string]time.Duration{
		"users":         80 * time.Millisecond,
		"orders":        100 * time.Millisecond,
		"payments":      150 * time.Millisecond,
		"notifications": 200 * time.Millisecond,
	}

	d := durations[name]
	if d == 0 {
		d = 100 * time.Millisecond
	}

	fmt.Printf("Fetching %s...\n", name)

	select {
	case <-time.After(d):
		if name == "orders" {
			return "", fmt.Errorf("service %q failed", name)
		}
		return fmt.Sprintf("%s: OK (%dms)", name, d.Milliseconds()), nil
	case <-ctx.Done():
		fmt.Printf("%s: cancelled\n", name)
		return "", ctx.Err()
	}
}

// fetchAll запрашивает все сервисы параллельно.
// При первой ошибке — отменяет остальные и возвращает ошибку.
// При успехе — возвращает результаты в порядке входного списка.
func fetchAll(ctx context.Context, services []string) ([]string, error) {
	// TODO: реализуй с errgroup.WithContext
	return nil, nil
}

func main() {
	ctx := context.Background()
	services := []string{"users", "orders", "payments", "notifications"}

	results, err := fetchAll(ctx, services)
	if err != nil {
		fmt.Printf("\nError: %v\n", err)
		return
	}

	fmt.Println("\nAll results:")
	for _, r := range results {
		fmt.Println(" ", r)
	}
}
