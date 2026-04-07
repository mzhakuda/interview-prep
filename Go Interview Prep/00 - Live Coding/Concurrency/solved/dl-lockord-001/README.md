# Problem dl-lockord-001: Find and Fix the Deadlock

**Topic:** deadlocks
**Subtopic:** prevention
**Difficulty:** medium
**Time limit:** 10 мин

## Описание

Дан код банковской системы: две горутины одновременно делают перевод между счетами. Код содержит **deadlock**. Твоя задача:

1. Объясни **почему** возникает deadlock (в комментарии)
2. Исправь код, чтобы deadlock был невозможен
3. Переводы должны остаться атомарными (баланс не теряется и не создаётся)

## Requirements

- Найди и исправь deadlock в функции `transfer`
- Оба перевода должны выполниться корректно
- Итоговая сумма балансов должна остаться неизменной (conservation)
- Программа должна завершиться за < 1 секунды
- Напиши комментарий объясняющий причину deadlock и твой фикс

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
Deadlock возникает из-за нарушения lock ordering: горутина 1 берёт lock A потом lock B, горутина 2 берёт lock B потом lock A. Каждая держит один lock и ждёт другой.
</details>

<details>
<summary>Hint 2</summary>
Классическое решение: всегда захватывай мьютексы в одном и том же порядке. Нужен детерминированный критерий — например, по ID аккаунта (меньший ID лочится первым).
</details>
