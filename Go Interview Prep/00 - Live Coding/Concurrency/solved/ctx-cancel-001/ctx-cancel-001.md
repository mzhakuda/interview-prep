# Problem ctx-cancel-001: Fan-in с cancellation

**Topic:** context
**Subtopic:** cancel
**Difficulty:** medium
**Time limit:** 15 мин

## Описание

Доработай функцию `merge` из предыдущей задачи: добавь поддержку `context.Context` для отмены.

Если контекст отменяется — все внутренние горутины должны завершиться, выходной канал закрыться. Никаких goroutine leaks.

## Requirements

- Сигнатура: `func merge(ctx context.Context, channels ...<-chan int) <-chan int`
- При отмене контекста все горутины завершаются
- Выходной канал закрывается в любом случае (нормальное завершение или cancel)
- Горутины не должны утекать — ни при cancel, ни при нормальном завершении
- Edge case: 0 каналов — вернуть закрытый канал
- Значения, отправленные до cancel, должны быть доступны для чтения

## Expected Output

```
=== Normal (all channels close) ===
1
10
100
...
All done.

=== Cancel after 3 values ===
1
10
100
Cancelled after 3 values.
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Внутри каждой горутины замени `for val := range channel` на цикл с `select` — два case: чтение из канала и `<-ctx.Done()`.
</details>

<details>
<summary>Hint 2</summary>
При отправке в `out` тоже нужен `select` с `ctx.Done()`. Иначе горутина застрянет на `out <- val` если reader уже ушёл.
</details>

<details>
<summary>Hint 3</summary>
Паттерн внутри горутины:
```go
for {
    select {
    case val, ok := <-ch:
        if !ok { return }
        select {
        case out <- val:
        case <-ctx.Done():
            return
        }
    case <-ctx.Done():
        return
    }
}
```
Два select — один для чтения, вложенный для записи. Оба слушают cancel.
</details>
