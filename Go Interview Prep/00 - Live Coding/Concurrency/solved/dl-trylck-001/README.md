# Problem dl-trylck-001: Deadlock Prevention with TryLock

**Topic:** deadlocks
**Subtopic:** prevention
**Difficulty:** medium
**Time limit:** 12 мин

## Описание

Тот же сценарий: два аккаунта, встречные переводы. Но на этот раз реши проблему **без lock ordering** — используй `TryLock` паттерн.

Идея: захвати первый lock, **попробуй** захватить второй. Если не получилось — отпусти первый и попробуй заново. Это подход "optimistic locking" — никогда не жди второй lock, если не удалось взять сразу.

Go 1.18+ имеет `sync.Mutex.TryLock()` — возвращает `bool` (true если lock захвачен).

## Requirements

- Сигнатура: `func transfer(from, to *Account, amount int)`
- **Не используй** lock ordering (решение через id < id запрещено)
- Используй `TryLock()` для захвата второго мьютекса
- Если `TryLock` вернул false — отпусти первый lock и повтори попытку
- Добавь защиту от livelock (бесконечный retry) — небольшой random backoff
- Оба перевода должны выполниться, баланс сохранён
- Все горутины завершаются

## Expected Output

```
Before: Alice=1000, Bob=1000, Total=2000
Transfer: Alice -> Bob: 300
Transfer: Bob -> Alice: 200
After: Alice=900, Bob=1100, Total=2000
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Паттерн: Lock(from) → TryLock(to) → если false: Unlock(from), sleep(random), retry. Если true: выполни перевод, Unlock обоих.
</details>

<details>
<summary>Hint 2</summary>
Для защиты от livelock: добавь `time.Sleep(time.Duration(rand.Intn(N)) * time.Millisecond)` перед retry. Рандомный backoff разбивает синхронность двух горутин.
</details>

<details>
<summary>Hint 3</summary>
Не используй defer для Unlock — тебе нужен ручной контроль: Unlock(from) при failed TryLock, и Unlock обоих при успехе. Оберни логику в for loop.
</details>
