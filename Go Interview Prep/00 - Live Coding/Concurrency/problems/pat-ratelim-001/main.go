package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Effector func(context.Context) (string, error)

// Throttle оборачивает функцию rate limiter'ом на основе token bucket.
// ⚠️ Код содержит минимум 2 concurrency бага. Найди и исправь.
func Throttle(e Effector, max uint, refill uint, d time.Duration) Effector {
	var tokens = max
	var once sync.Once

	return func(ctx context.Context) (string, error) {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		once.Do(func() {
			ticker := time.NewTicker(d)

			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						tokens = min(tokens+refill, max)
					}
				}
			}()
		})

		if tokens <= 0 {
			return "", fmt.Errorf("too many calls")
		}

		tokens--
		return e(ctx)
	}
}

func main() {
	// Простая функция-эффектор
	myFunc := func(ctx context.Context) (string, error) {
		return "ok", nil
	}

	throttled := Throttle(myFunc, 5, 2, 100*time.Millisecond)

	// Запусти 20 горутин, делающих запросы
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ctx := context.Background()
			result, err := throttled(ctx)
			if err != nil {
				fmt.Printf("call %d: %v\n", id, err)
			} else {
				fmt.Printf("call %d: %s\n", id, result)
			}
		}(i)
	}

	wg.Wait()
	fmt.Println("\n--- waiting for refill ---")
	time.Sleep(150 * time.Millisecond)

	// Ещё 3 вызова после refill
	for i := 20; i < 23; i++ {
		ctx := context.Background()
		result, err := throttled(ctx)
		if err != nil {
			fmt.Printf("call %d: %v\n", i, err)
		} else {
			fmt.Printf("call %d: %s\n", i, result)
		}
	}
}
