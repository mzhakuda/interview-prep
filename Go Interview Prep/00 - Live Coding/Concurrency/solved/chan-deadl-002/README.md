# Problem chan-deadl-002: Find the Deadlock

**Topic:** deadlocks
**Subtopic:** detection
**Difficulty:** medium
**Time limit:** 8 мин
**Source:** [Habr/MTS](https://habr.com/ru/companies/ru_mts/articles/909158/)

## Описание

Дан код: 5 горутин отправляют сообщения в канал, main читает из канала через `select`. Код **зависает** — найди почему и исправь.

## Requirements

- Определи причину зависания (напиши комментарий)
- Исправь код так, чтобы все 5 сообщений были выведены
- `wg.Wait()` должен дождаться завершения всех горутин
- Программа должна корректно завершиться
- Mutex в этой задаче **не нужен** — убери его если считаешь лишним, или объясни зачем он

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
Посмотри на for + select: когда этот цикл завершится? Что будет после него?
</details>

<details>
<summary>Hint 2</summary>
Канал никогда не закрывается. select без default блокируется навечно. Кто и когда должен закрыть канал?
</details>
