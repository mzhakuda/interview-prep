# Problem eg-cancel-001: Parallel Fetch with Cancellation

**Topic:** errgroup
**Subtopic:** with_context
**Difficulty:** medium
**Time limit:** 12 мин

## Описание

Реализуй параллельную загрузку данных из нескольких "сервисов". Если хотя бы один сервис возвращает ошибку — отмени все остальные запросы и верни эту ошибку.

Функция `fetchService(ctx context.Context, name string)` уже написана — имитирует запрос. Некоторые сервисы "падают".

## Requirements

- Сигнатура: `func fetchAll(ctx context.Context, services []string) ([]string, error)`
- Используй `errgroup.Group` с `errgroup.WithContext`
- Все сервисы запрашиваются параллельно
- При первой ошибке — все остальные горутины должны увидеть отмену через `ctx.Done()`
- Если ошибок нет — верни все результаты (порядок соответствует входному списку)
- Результаты должны быть thread-safe (нельзя писать в shared slice без защиты)

## Expected Output

```
Fetching users...
Fetching orders...
Fetching payments...
Fetching notifications...
users: OK (120ms)
orders: FAILED
payments: cancelled
notifications: cancelled
Error: service "orders" failed
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
errgroup.WithContext(ctx) возвращает группу и derived context. Когда любая горутина в группе возвращает error, этот derived context автоматически отменяется.
</details>

<details>
<summary>Hint 2</summary>
Для сохранения порядка результатов: создай slice нужного размера заранее и пиши по индексу results[i] = value. Каждая горутина пишет в свой индекс — это безопасно без mutex.
</details>

<details>
<summary>Hint 3</summary>
Внутри fetchService проверяй ctx.Err() или используй select с ctx.Done(), чтобы реагировать на отмену.
</details>
