# Problem pat-worker-001: Worker Pool

**Topic:** patterns
**Subtopic:** worker_pool
**Difficulty:** medium
**Time limit:** 15 мин

## Описание

Реализуй worker pool: N воркеров обрабатывают задачи из общего канала и отправляют результаты в выходной канал.

У тебя есть список URL-ов (строки). Каждый "воркер" — горутина, которая берёт URL из входного канала, вызывает функцию `process(url)` (уже написана — имитирует обработку), и отправляет результат в выходной канал.

## Requirements

- Сигнатура: `func workerPool(urls []string, numWorkers int) []string`
- Запусти ровно `numWorkers` горутин-воркеров
- Воркеры читают задачи из общего канала `jobs`
- Результаты собираются в канал `results`
- Функция возвращает все результаты (порядок не важен)
- Все горутины должны завершиться — никаких leaks
- Каналы должны быть корректно закрыты

## Expected Output

```
Worker 1 processing: https://example.com/a
Worker 2 processing: https://example.com/b
Worker 3 processing: https://example.com/c
Worker 1 processing: https://example.com/d
Worker 2 processing: https://example.com/e
Results (порядок может отличаться):
  processed: https://example.com/a
  processed: https://example.com/b
  processed: https://example.com/c
  processed: https://example.com/d
  processed: https://example.com/e
All done.
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Тебе нужны два канала: `jobs chan string` (входной) и `results chan string` (выходной). Один producer наполняет jobs, N воркеров читают из jobs и пишут в results.
</details>

<details>
<summary>Hint 2</summary>
Кто закрывает jobs? Producer — после того как отправил все URL. Кто закрывает results? Нужно дождаться завершения всех воркеров (WaitGroup), и только потом close(results).
</details>

<details>
<summary>Hint 3</summary>
Сбор результатов: main не может range по results пока results не закрыт. Но close(results) должен произойти после wg.Wait(). Значит close(results) должен быть в отдельной горутине, а сбор — в main.
</details>
