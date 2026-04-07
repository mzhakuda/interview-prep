package main

import (
	"fmt"
	"strconv"
	"sync"
)

// ⚠️ Этот код зависает. Найди причину и исправь.
func main() {
	var wg sync.WaitGroup
	ch := make(chan string, 5)
	// mu := sync.Mutex{}

	for i := range 5 {
		wg.Add(1)
		go func(ch chan<- string, i int, grp *sync.WaitGroup) {
			defer wg.Done()

			// не понятно зачем тут mutex, у нас нету какой то зависимой от чего то данных.
			// workerNumber мы передаем через goroutine number

			// mu.Lock()
			// defer mu.Unlock()

			msg := fmt.Sprintf("Goroutine %s", strconv.Itoa(i))
			ch <- msg
		}(ch, i, &wg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for {
		select {
		case q, ok := <-ch:
			if !ok {
				return
			}
			fmt.Println(q)
		}
	}
}
