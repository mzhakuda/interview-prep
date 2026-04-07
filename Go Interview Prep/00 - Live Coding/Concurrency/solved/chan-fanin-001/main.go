package main

import (
	"fmt"
	"sync"
	// TODO: добавь нужные импорты
)

// merge принимает произвольное количество каналов и возвращает один,
// в который попадают все значения из всех входных каналов.
// Когда все входные каналы закрыты — выходной тоже закрывается.
// Горутины не должны утекать.
func merge(channels ...<-chan int) <-chan int {
	// TODO: реализуй
	// 1. Создай выходной канал
	// 2. Обработай edge case: 0 каналов
	// 3. Для каждого входного канала запусти горутину
	// 4. Обеспечь закрытие выходного канала когда все входные закрыты

	out := make(chan int)
	if len(channels) == 0 {
		close(out)
		return out
	}

	wg := sync.WaitGroup{}
	for _, channel := range channels {
		wg.Go(func() {
			for val := range channel {
				out <- val
			}
		})
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	// Создаём 3 канала и наполняем их
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

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

	go func() {
		for _, v := range []int{100, 200, 300} {
			ch3 <- v
		}
		close(ch3)
	}()

	// Мержим и читаем все значения
	for val := range merge(ch1, ch2, ch3) {
		fmt.Println(val)
	}

	fmt.Println("All channels closed, merged channel closed too.")
}
