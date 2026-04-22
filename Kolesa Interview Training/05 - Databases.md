---
tags:
  - kolesa
  - postgres
  - databases
topic: Databases
priority: critical
status: ready
---

# Базы данных (фокус: PostgreSQL)

> Более полный материал: [[Go Interview Prep/02 - PostgreSQL]] и [[Go Interview Prep/00 - Cheat Sheet]] (секция PostgreSQL).
> Тут — самое важное, что спросят на 1 этапе.

---

## ACID

- **Atomicity** — транзакция либо вся применилась, либо ничего
- **Consistency** — БД переходит из валидного состояния в валидное (constraints)
- **Isolation** — параллельные транзакции не видят промежуточных состояний друг друга (см. уровни)
- **Durability** — закоммиченная транзакция переживёт крах (WAL)

---

## Уровни изоляции (⭐ часто спрашивают)

| Уровень | Dirty read | Non-repeat | Phantom | Postgres default? |
|---------|:---------:|:---------:|:-------:|:-----------------:|
| Read Uncommitted | ✗ | ✗ | ✗ | (нет в PG) |
| **Read Committed** | ✓ | ✗ | ✗ | ⭐ default |
| Repeatable Read | ✓ | ✓ | ✓ (в PG — без phantoms благодаря MVCC) | |
| Serializable | ✓ | ✓ | ✓ | дорого, но безопасно |

- **Dirty read** — прочитать незакоммиченные изменения другой транзакции
- **Non-repeatable read** — повторить SELECT и получить другой результат (закоммичен UPDATE)
- **Phantom read** — повторить SELECT и увидеть новые строки (закоммичен INSERT)

---

## MVCC (Multi-Version Concurrency Control)

Postgres хранит **несколько версий** одной строки. Каждая строка имеет `xmin` (кем создана), `xmax` (кем удалена/обновлена).

- Читатели **не блокируют** писателей, писатели не блокируют читателей
- **VACUUM** удаляет мёртвые версии (dead tuples)
- **Bloat** — мёртвые версии, которые не почистил autovacuum → таблица пухнет
- `pg_stat_user_tables` — следить за bloat

---

## Индексы

### B-tree (default)
- Equality: `WHERE id = ?`
- Range: `WHERE created_at > ?`
- Сортировка: `ORDER BY` если индекс подходит

### Hash
Только equality. Редко нужен (b-tree почти так же быстр на = и умеет range).

### GIN (⭐ у тебя был trigram в escrow!)
- Полнотекстовый поиск (`tsvector`)
- Массивы
- JSONB
- **Trigram** (`pg_trgm`) — поиск по подстроке (`LIKE '%foo%'`)

### GIST
- Геометрия, полнотекст, диапазоны
- `tstzrange` — исключение пересечений

### BRIN
Для **огромных** таблиц с "естественным порядком" (например timestamps). Очень компактный, но медленный на узких запросах.

### Partial index
```sql
CREATE INDEX ON orders (user_id) WHERE status = 'pending';
```
Меньше размер, быстрее — если запрос всегда с тем же `WHERE`.

### Covering / INCLUDE (index-only scan)
```sql
CREATE INDEX ON orders (user_id) INCLUDE (amount, status);
```
Запрос `SELECT amount FROM orders WHERE user_id=?` не идёт в heap.

---

## EXPLAIN ANALYZE

```sql
EXPLAIN (ANALYZE, BUFFERS, VERBOSE)
SELECT * FROM orders WHERE user_id = 42;
```

**Что читать:**
- `Seq Scan` — полный перебор таблицы (плохо для больших)
- `Index Scan` — по индексу
- `Index Only Scan` — вообще в heap не ходим (хорошо)
- `Bitmap Index Scan + Bitmap Heap Scan` — несколько индексов
- `rows=...` — оценка планировщика; если сильно отличается от `actual rows` — нужен `ANALYZE table`
- `Buffers: shared hit/read` — сколько в кеше (hit), сколько с диска (read)

---

## Блокировки

### Row-level
- `SELECT ... FOR UPDATE` — берёт строку в монопольную блокировку до конца транзакции
- `FOR SHARE` — разделяемая (можно другим читать FOR SHARE, но не UPDATE)
- `FOR UPDATE SKIP LOCKED` — ⭐ для job queues, берёт только свободные
- `FOR UPDATE NOWAIT` — не ждать, если занята — ошибка

### Table-level (обычно не ручные)
- `ACCESS SHARE` — берёт `SELECT`
- `ACCESS EXCLUSIVE` — берут `DROP`, `TRUNCATE`, некоторые `ALTER`

### Deadlocks
Postgres детектирует deadlock и **кильнёт одну из транзакций** с ошибкой. На клиенте — retry.

### Как избежать
- Все транзакции берут локи в **одинаковом порядке**
- Короткие транзакции
- **Не делать HTTP-вызовы внутри транзакции** ← твой опыт из BCC (фикс через Temporal)

---

## Транзакции + приложение

- **Savepoints** — вложенные "под-транзакции" внутри большой транзакции
- **Advisory locks** — `pg_advisory_lock(id)` для app-level lock'ов (например, leader election)

---

## Connection pool — PgBouncer

- Postgres жрёт ~10MB RAM на коннекшн → нельзя держать тысячи
- **PgBouncer** перед Postgres:
  - **session mode** — как был (1 клиент = 1 backend), безопасно
  - **transaction mode** — коннекшн возвращается в пул после COMMIT (экономично, но **ломает prepared statements и SET session**)
  - **statement mode** — самое экономичное, но очень много ограничений

---

## Схема и миграции

- **`goose` / `golang-migrate`** — версионирование миграций
- Миграции — **idempotent + forward-only**
- Не менять уже применённую миграцию — пиши новую
- **Zero-downtime миграции:**
  - ADD COLUMN NULLABLE → backfill → NOT NULL
  - Rename: копия колонки → dual-write → переключение → drop старой
  - `CREATE INDEX CONCURRENTLY` (не блокирует таблицу)

---

## Репликация

- **Streaming replication** — master → hot standby (async/sync)
- **Logical replication** — на уровне рядов (pgoutput), для частичной репликации
- **Read replicas** — routing по типу запроса; осторожно с replication lag

---

## Типичные вопросы

**1. Разница INNER / LEFT / RIGHT / FULL JOIN?**
- INNER — только совпадающие пары
- LEFT — все из левой + совпадения; не совпало — NULL
- RIGHT — зеркально
- FULL — все из обеих + NULL на непарные

**2. GROUP BY vs DISTINCT?**
`DISTINCT` — просто убирает дубликаты. `GROUP BY` позволяет агрегаты. Часто взаимозаменяемы на одной колонке, но `GROUP BY` более мощный.

**3. WHERE vs HAVING?**
`WHERE` — фильтр ДО агрегации. `HAVING` — фильтр ПОСЛЕ агрегации (по результатам `GROUP BY`).

**4. Как оптимизировал медленный запрос?**
Пример из Gekata: ES-схема + GIN индексы → +10-20% скорости. Общий план: `EXPLAIN ANALYZE` → смотрю seq scan / высокое rows → добавляю индекс / переписываю условие / partial index. Иногда — переписать JOIN на CTE или наоборот.

**5. Когда NoSQL, когда SQL?**
- **SQL** — транзакции, сложные JOIN, консистентность, relational данные
- **NoSQL (Mongo/DynamoDB)** — schema-less, документы, очень высокая запись, шардирование из коробки
- **Redis** — кеш, session, rate limit, leaderboard

**6. CAP-теорема?**
Из {Consistency, Availability, Partition tolerance} в распределённой системе при сетевом разбиении выбираем 2 из 3. Фактически выбор между C и A при партиции.

**7. Что сделать если запрос медленный и индекс не помогает?**
- Проверить `EXPLAIN ANALYZE` — может planner идёт в seq scan из-за плохой статистики → `ANALYZE`
- Проверить selectivity индекса
- Partial index
- Денормализация / материализованные views
- Партиционирование таблицы
- Кеш в Redis перед БД

**8. Idempotency в API (⭐ банк)**
- `Idempotency-Key` header, сохраняем в БД `(key, response)` с TTL
- При повторе — отдаём закешированный ответ, не выполняем операцию снова
- Уникальный индекс на `(idempotency_key)` — атомарный insert через `ON CONFLICT`

**9. Pagination — offset vs keyset?**
- **Offset/Limit:** `LIMIT 20 OFFSET 10000` — медленно на больших страницах
- **Keyset (seek):** `WHERE (created_at, id) > (?, ?) ORDER BY created_at, id LIMIT 20` — стабильно быстрое

---

## Специфичное к твоему опыту (упомянуть)

- **Full-text search на казахском** — кастомная нормализационная функция + 12 GIN trigram индексов для поиска без диакритики
- **Circuit breaker** перед интеграциями с CBS / CEA / GovBus
- **Temporal вместо BPMN** — HTTP вызовы вне DB-транзакций (классический fix deadlock)
