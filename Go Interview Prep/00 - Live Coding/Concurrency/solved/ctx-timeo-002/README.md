# Problem ctx-timeo-002: Slow Function Wrapper

**Topic:** context
**Subtopic:** timeout
**Difficulty:** medium
**Time limit:** 10 мин
**Source:** [Habr/MTS](https://habr.com/ru/companies/ru_mts/articles/909158/)

## Описание

У тебя есть "медленная" функция `slowFunc()` которая возвращает `int64`. Ты **не можешь** изменить её сигнатуру — она не принимает context.

Напиши обёртку `ctxFunc(ctx context.Context) (int64, error)`, которая:
- Запускает `slowFunc()` в фоне
- Если результат пришёл до отмены контекста — возвращает его
- Если контекст отменён раньше — возвращает `0, ctx.Err()`
- **Не допускает goroutine leak** в обоих сценариях

## Requirements

- Сигнатура: `func ctxFunc(ctx context.Context) (int64, error)`
- `slowFunc()` вызывается ровно один раз
- При отмене контекста — горутина с `slowFunc()` не должна "висеть" навечно (think about channel buffering)
- В `main`: продемонстрируй оба сценария — успех и таймаут

## Expected Output

```
=== Scenario 1: enough time ===
Result: 42, err: <nil>

=== Scenario 2: timeout ===
Result: 0, err: context deadline exceeded
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Используй select с двумя cases: получение результата из канала и ctx.Done().
</details>

<details>
<summary>Hint 2</summary>
Если канал unbuffered и ctx отменяется раньше — горутина с slowFunc() завершится, попытается отправить в канал, и заблокируется навечно (goroutine leak). Подумай о буферизации.
</details>
