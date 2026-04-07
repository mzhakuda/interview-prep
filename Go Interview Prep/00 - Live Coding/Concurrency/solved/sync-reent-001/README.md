# Problem sync-reent-001: Fix the Reentrant Lock Deadlock

**Topic:** sync_primitives
**Subtopic:** mutex
**Difficulty:** medium
**Time limit:** 12 мин

## Описание

Дан thread-safe кэш с методами `Get`, `Set`, `Delete`, `Keys`, и `Cleanup`. Код содержит **deadlock из-за повторного захвата мьютекса** (reentrant locking). Всё компилируется, но `Cleanup` зависает.

Твоя задача:
1. Найди **все** места где происходит повторный захват lock
2. Рефактори код так, чтобы deadlock был невозможен
3. Все методы должны остаться thread-safe для внешних вызовов

## Requirements

- Исправь deadlock, не убирая thread-safety
- Все публичные методы (`Get`, `Set`, `Delete`, `Keys`, `Cleanup`) должны быть concurrent-safe
- `Cleanup` должен удалить все expired записи
- Не используй `sync.RWMutex` (оставь обычный `Mutex`) — фокус на правильной структуре lock'ов
- Напиши комментарий объясняющий твой подход

## Expected Output

```
Set: a=1, b=2, c=3
Keys before cleanup: [a b c]
Cleanup: deleted 2 expired keys
Keys after cleanup: [b]
Get b: 2, found: true
Get a: , found: false
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Паттерн: раздели каждый метод на публичный (берёт lock) и приватный (работает без lock). Публичный вызывает приватный под lock'ом. Внутренние методы вызывают только приватные версии.
</details>

<details>
<summary>Hint 2</summary>
Например: `Get` берёт lock и вызывает `getLocked`. `Cleanup` берёт lock и вызывает `deleteLocked` для каждого expired ключа. Никто не берёт lock дважды.
</details>
