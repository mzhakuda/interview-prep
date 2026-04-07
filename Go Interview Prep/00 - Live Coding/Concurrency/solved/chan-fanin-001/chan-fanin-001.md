# Problem chan-fanin-001: Merge N каналов

**Topic:** channels
**Subtopic:** fan-in
**Difficulty:** medium
**Time limit:** 12 мин

## Описание

Напиши функцию `merge`, которая принимает произвольное количество каналов `<-chan int` и возвращает один канал `<-chan int`, в который попадают все значения из всех входных каналов.

Когда **все** входные каналы закрыты — выходной канал тоже должен закрыться.

## Requirements

- Сигнатура: `func merge(channels ...<-chan int) <-chan int`
- Горутины не должны утекать — когда все входные каналы закрыты, все внутренние горутины должны завершиться
- Используй `sync.WaitGroup` для отслеживания завершения горутин
- Порядок значений в выходном канале **не важен** (non-deterministic)
- Функция должна корректно работать при 0 входных каналов (вернуть сразу закрытый канал)

## Expected Output

```
Sending: 1 2 3 from ch1
Sending: 10 20 30 from ch2
Sending: 100 200 300 from ch3
Merged (порядок может отличаться): 1 10 100 2 20 200 3 30 300
All channels closed, merged channel closed too.
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Запусти по одной горутине на каждый входной канал. Каждая горутина читает из своего канала и пишет в общий выходной.
</details>

<details>
<summary>Hint 2</summary>
Используй `sync.WaitGroup`: `wg.Add(1)` для каждой горутины, `wg.Done()` когда входной канал закрыт. Отдельная горутина делает `wg.Wait()` и затем `close(out)`.
</details>

<details>
<summary>Hint 3</summary>
Edge case: если `len(channels) == 0`, создай канал, сразу закрой и верни.
</details>
