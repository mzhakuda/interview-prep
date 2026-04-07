# Problem sync-rwmut-001: Thread-Safe Config Store

**Topic:** sync_primitives
**Subtopic:** rwmutex
**Difficulty:** medium
**Time limit:** 12 мин

## Описание

Реализуй потокобезопасное хранилище конфигурации. Множество горутин постоянно читают конфиг (read-heavy workload), и редко одна горутина обновляет значения.

Используй `sync.RWMutex` чтобы разрешить параллельные чтения, но эксклюзивную запись.

## Requirements

- Структура `ConfigStore` с методами:
  - `Get(key string) (string, bool)` — получить значение по ключу
  - `Set(key, value string)` — установить значение
  - `GetAll() map[string]string` — получить копию всего конфига
- `Get` и `GetAll` могут работать параллельно (read lock)
- `Set` блокирует все операции (write lock)
- `GetAll` возвращает **копию** map, не ссылку на внутренний
- В `main`: запусти 10 readers (читают в цикле), 2 writers (пишут в цикле), работают 500ms, выведи итоговый конфиг

## Expected Output

```
[reader 3] key=db_host -> localhost
[reader 7] key=db_port -> 5432
[writer 1] set app_env = production
[reader 1] key=app_env -> production
...
Final config:
  db_host = localhost
  db_port = 5432
  app_env = production
  ...
```

## Hints (раскрой если застрял)

<details>
<summary>Hint 1</summary>
RWMutex: вызывай RLock/RUnlock для чтения, Lock/Unlock для записи. Используй defer для Unlock.
</details>

<details>
<summary>Hint 2</summary>
GetAll должен скопировать map в новый map внутри read lock. Иначе вызывающий код может читать map без защиты, пока writer пишет — data race.
</details>

<details>
<summary>Hint 3</summary>
Для остановки горутин через 500ms используй context.WithTimeout или time.After + select.
</details>
