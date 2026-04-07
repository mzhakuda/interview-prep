package main

import (
	"context"
	"fmt"
	"sync"
	// TODO: добавь нужные импорты
)

// merge принимает context и произвольное количество каналов.
// Возвращает один канал со всеми значениями.
// При отмене контекста — все горутины завершаются, out закрывается.
func merge(ctx context.Context, channels ...<-chan int) <-chan int {
	// TODO: реализуй
	// 1. Создай выходной канал
	// 2. Edge case: 0 каналов
	// 3. Для каждого канала — горутина с select на ctx.Done()
	// 4. Закрой out когда все горутины завершились
	out := make(chan int)
	if len(channels) == 0 {
		close(out)
		return out
	}

	wg := sync.WaitGroup{}
	for _, channel := range channels {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return

				case val, ok := <-channel:
					if !ok {
						return
					}

					select {
					case <-ctx.Done():
						return
					case out <- val:
						continue
					}
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func main() {
	// === Тест 1: нормальное завершение ===
	fmt.Println("=== Normal ===")
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for _, v := range []int{1, 2, 3} {
			ch1 <- v
		}
		close(ch1)
	}()
	go func() {
		for _, v := range []int{10, 20, 30} {
			ch2 <- v
		}
		close(ch2)
	}()

	ctx := context.Background()
	for val := range merge(ctx, ch1, ch2) {
		fmt.Println(val)
	}
	fmt.Println("All done.")

	// === Тест 2: cancel после 3 значений ===
	fmt.Println("\n=== Cancel after 3 ===")
	ch3 := make(chan int)
	ch4 := make(chan int)

	go func() {
		for _, v := range []int{1, 2, 3, 4, 5} {
			ch3 <- v
		}
		close(ch3)
	}()
	go func() {
		for _, v := range []int{10, 20, 30, 40, 50} {
			ch4 <- v
		}
		close(ch4)
	}()

	ctx2, cancel := context.WithCancel(context.Background())
	count := 0
	for val := range merge(ctx2, ch3, ch4) {
		fmt.Println(val)
		count++
		if count == 3 {
			cancel()
		}
	}
	fmt.Println("Cancelled after 3 values.")
}
