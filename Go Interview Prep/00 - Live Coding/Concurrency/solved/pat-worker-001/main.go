package main

import (
	"fmt"
	"sync"
	"time"
	// TODO: добавь нужные импорты
)

// process имитирует обработку URL (не меняй эту функцию)
func process(url string, workerID int) string {
	fmt.Printf("Worker %d processing: %s\n", workerID, url)
	time.Sleep(50 * time.Millisecond) // имитация работы
	return "processed: " + url
}

// workerPool запускает numWorkers воркеров для обработки urls.
// Возвращает все результаты. Порядок не важен.
// Все горутины должны завершиться, каналы — корректно закрыться.
func workerPool(urls []string, numWorkers int) []string {
	// TODO: реализуй
	// 1. Создай каналы jobs и results
	// 2. Запусти producer — горутину, которая отправит все urls в jobs
	// 3. Запусти numWorkers воркеров
	// 4. Обеспечь закрытие results после завершения всех воркеров
	// 5. Собери и верни все результаты

	jobs := make(chan string)
	results := make(chan string)

	go func() {
		for _, url := range urls {
			jobs <- url
		}
		close(jobs)
	}()

	wg := sync.WaitGroup{}
	wg.Add(numWorkers)

	for i := range numWorkers {
		go func(workerIndex int) {
			defer wg.Done()
			for job := range jobs {
				results <- process(job, workerIndex)
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	doneWorks := make([]string, 0)
	for doneWork := range results {
		doneWorks = append(doneWorks, doneWork)
	}
	return doneWorks
}

func main() {
	urls := []string{
		"https://example.com/a",
		"https://example.com/b",
		"https://example.com/c",
		"https://example.com/d",
		"https://example.com/e",
	}

	results := workerPool(urls, 3)

	fmt.Println("\nResults:")
	for _, r := range results {
		fmt.Println(" ", r)
	}
	fmt.Println("All done.")
}
