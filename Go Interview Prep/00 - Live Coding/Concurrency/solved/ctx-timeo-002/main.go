package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// slowFunc — медленная функция, сигнатуру менять нельзя.
func slowFunc() int64 {
	time.Sleep(time.Duration(200+rand.Intn(300)) * time.Millisecond)
	return 42
}

// ctxFunc — обёртка над slowFunc с поддержкой context.
// Если контекст отменён раньше — вернуть 0, ctx.Err().
// Не допускать goroutine leak.
func ctxFunc(ctx context.Context) (int64, error) {
	result := make(chan int64, 1)
	go func() {
		result <- slowFunc()
	}()

	for {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case slowVal := <-result:
			return slowVal, nil
		}
	}

}

func main() {
	// Scenario 1: достаточно времени (1 секунда)
	fmt.Println("=== Scenario 1: enough time ===")
	ctx1, cancel1 := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel1()
	result, err := ctxFunc(ctx1)
	fmt.Printf("Result: %d, err: %v\n\n", result, err)

	// Scenario 2: таймаут (50ms — slowFunc не успеет)
	fmt.Println("=== Scenario 2: timeout ===")
	ctx2, cancel2 := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel2()
	result, err = ctxFunc(ctx2)
	fmt.Printf("Result: %d, err: %v\n", result, err)
}
