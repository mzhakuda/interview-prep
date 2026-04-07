# Problem pat-ratelim-001: Find Bugs in Rate Limiter

**Topic:** patterns
**Subtopic:** rate_limiter
**Difficulty:** hard
**Time limit:** 15 мин
**Source:** [Habr/MTS](https://habr.com/ru/companies/ru_mts/articles/909158/)

## Описание

Дан rate limiter на основе token bucket. Код содержит **минимум 2 concurrency бага**. Твоя задача:

1. Найди **все** баги (напиши комментарии к каждому)
2. Исправь код
3. Продемонстрируй работу в `main`

## Requirements

- Найди и исправь все concurrency баги
- Rate limiter должен:
  - Иметь `max` токенов изначально
  - Пополняться на `refill` токенов каждые `d`
  - Не превышать `max` при пополнении
  - Отклонять вызовы когда токенов нет
- Напиши комментарии объясняющие каждый найденный баг
- В `main`: запусти 20 горутин, делающих запросы через throttled функцию

## Expected Behavior

```
call 1: ok (tokens left ~)
call 2: ok
...
call N: too many calls
(после refill)
call M: ok (tokens replenished)
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Переменная `tokens` читается и пишется из нескольких горутин без защиты. Это data race. Подумай чем защитить: mutex или atomic.
</details>

<details>
<summary>Hint 2</summary>
Подумай про жизненный цикл горутины с ticker. Что произойдёт если ctx — это child context из конкретного запроса? Когда этот запрос завершится, что случится с refill горутиной?
</details>

<details>
<summary>Hint 3</summary>
once.Do запускает refill горутину с контекстом первого вызова. Если этот контекст отменится — refill остановится навсегда, даже для последующих вызовов с живым контекстом.
</details>
