# Senior Go Developer — Interview Prep
> Вакансия: API Gateway рефакторинг (Go, microservices, observability)
> Формат: Live coding + Теория + System Design + Code Review

---

## 🔗 GitHub Репы

### Go-specific
| Репо | Описание | Ссылка |
|------|----------|--------|
| Devinterview-io/golang-interview-questions | Самая полная база вопросов 2026, runtime/concurrency/modules | https://github.com/Devinterview-io/golang-interview-questions |
| RezaSi/go-interview-practice | 30+ coding challenges, AI simulation, веб-UI на localhost:8080 | https://github.com/RezaSi/go-interview-practice |
| shomali11/go-interview | Технические вопросы решённые на Go | https://github.com/shomali11/go-interview |

### System Design
| Репо | Описание | Ссылка |
|------|----------|--------|
| ashishps1/awesome-system-design-resources | Бесплатные ресурсы для SD, самый звёздный | https://github.com/ashishps1/awesome-system-design-resources |
| donnemartin/system-design-primer | 109k ⭐, Anki flashcards, большая компиляция | https://github.com/donnemartin/system-design-primer |
| Devinterview-io/microservices-interview-questions | Паттерны микросервисов, прямо в тему вакансии | https://github.com/Devinterview-io/microservices-interview-questions |
| madd86/awesome-system-design | Distributed systems ресурсы | https://github.com/madd86/awesome-system-design |

---

## 🟡 Go Concurrency (СЛАБАЯ СТОРОНА — приоритет 1)

### Теория — гарантированные вопросы

**GMP-модель**
- G = goroutine, M = OS thread, P = processor (logical CPU)
- Шедулер распределяет runnable goroutines по доступным M и P
- Кол-во P задаётся `GOMAXPROCS`
- Главная идея: "не общайся через shared memory — общайся через каналы"
- С Go 1.14: preemptive scheduler через compiler-inserted checks (до этого только cooperative)

**Каналы**
- Unbuffered: блокирует sender до тех пор пока receiver не готов
- Buffered: блокирует только когда буфер полный
- `select` — блокируется пока не готов один из case; если несколько готовы — выбирает случайный

**sync примитивы**
- `sync.Mutex` — exclusive lock, только один goroutine
- `sync.RWMutex` — multiple readers / one writer
- `sync.WaitGroup` — ждать завершения N goroutines
- `sync.Once` — выполнить один раз
- `sync/atomic` — atomic операции без mutex, предпочтительнее для счётчиков
- `sync.Pool` — переиспользование объектов, снижает GC pressure

**Context**
- Управляет cancellation signals, deadlines, request-scoped values
- Передаётся через API boundaries и между goroutines
- Необходим для graceful shutdown

### Паттерны конкурентности (знать уметь объяснить + написать)
```
Generator      — функция возвращает канал, пушит последовательность значений
Worker Pool    — N goroutines читают задачи из одного канала
Fan-out        — раздаём работу по нескольким goroutines
Fan-in         — собираем результаты из нескольких каналов в один
Pipeline       — цепочка стадий, каждая читает из предыдущего канала
Rate Limiting  — ticker или token bucket для контроля частоты
```

### Tricky вопросы — ловушки на интервью

**Ловушка #1: deadlock с буферизованным каналом**
```go
// Вопрос: будет ли deadlock?
ch := make(chan int, 4)
go func() {
    for i := range ch { fmt.Println(i) }
}()
for i := 0; i < 5; i++ { ch <- i }
```
Большинство говорят "deadlock на 5-м элементе" — но код работает.
Причина: анонимная горутина обеспечивает concurrent consumption,
буфер не переполняется.

**Ловушка #2: closure в goroutine**
```go
// Классический баг
for i := 0; i < 5; i++ {
    go func() { fmt.Println(i) }() // всегда печатает 5
}
// Правильно:
for i := 0; i < 5; i++ {
    go func(n int) { fmt.Println(n) }(i)
}
```

**Ловушка #3: goroutine leak**
```go
// Leak: никто не читает из канала
ch := make(chan int)
go func() { ch <- 1 }() // висит навсегда
// Решение: всегда закрывать каналы или использовать context
```

### Best practices
- Avoid unbounded goroutines — всегда используй worker pools
- Always close channels when done (иначе goroutine leak)
- Используй context для cancellation
- Детектор гонок: `go run -race ./...` (в 10x больше CPU/RAM — только в dev/test)
- Предпочитай `sync/atomic` для простых счётчиков, mutex для сложной логики

### Вопросы которые реально задают
1. Объясни GMP модель шедулера
2. Разница buffered vs unbuffered channel
3. Когда использовать channel, когда mutex?
4. Как остановить goroutine? (context cancellation)
5. Что такое goroutine leak и как обнаружить?
6. Реализуй worker pool
7. Что такое data race? Как обнаружить?
8. Объясни sync.Once, где используется?
9. Что такое happens-before в Go memory model?
10. Разница sync.Mutex vs sync.RWMutex, когда что использовать?

---

## 🟡 PostgreSQL / Database (СЛАБАЯ СТОРОНА — приоритет 2)

### MVCC (Multi-Version Concurrency Control)
- Основной механизм конкурентности PostgreSQL
- При UPDATE/DELETE — НЕ перезаписывает строку, создаёт новую версию
- Каждый tuple имеет скрытые колонки `xmin` (создавшая транзакция) и `xmax` (удалившая)
- Читатели не блокируют писателей, писатели не блокируют читателей
- Старые версии ("dead tuples") очищает `autovacuum`

### VACUUM
- `VACUUM` — reclaims место от dead tuples, обновляет visibility map, freezes transaction IDs
- `VACUUM FULL` — перезаписывает всю таблицу, требует exclusive lock, работает медленно — редко
- Обычный `VACUUM` управляется autovacuum автоматически

### Индексы — типы и когда использовать
| Тип | Использование |
|-----|---------------|
| B-tree | Дефолт, большинство запросов (=, <, >, BETWEEN, LIKE 'foo%') |
| Hash | Только equality (=), быстрее B-tree для этого случая |
| GIN | Arrays, full-text search, JSONB |
| GiST | Геометрические данные, custom типы |
| BRIN | Очень большие таблицы с естественным порядком (created_at) |

### Locking
- `SELECT FOR UPDATE` — exclusive lock на строки, блокирует modify и FOR UPDATE у других
- `SELECT FOR SHARE` — shared lock, другие могут читать но не modify
- Deadlock — PostgreSQL автоматически детектит и откатывает одну транзакцию

### Isolation Levels
| Уровень | Dirty Read | Non-Repeatable | Phantom Read |
|---------|-----------|----------------|--------------|
| READ COMMITTED | ✗ | возможен | возможен |
| REPEATABLE READ | ✗ | ✗ | ✗ в PG (строже стандарта) |
| SERIALIZABLE | ✗ | ✗ | ✗ |

### WAL (Write-Ahead Log)
- Изменения записываются в лог ДО применения к данным
- Обеспечивает durability и crash recovery
- Основа для репликации

### Connection Pooling
- PgBouncer — lightweight pooler перед PostgreSQL
- Клиенты подключаются к PgBouncer, он раздаёт соединения из пула
- Режимы: transaction pooling (рекомендован), session pooling
- Снижает overhead и позволяет поддерживать гораздо больше concurrent клиентов

### EXPLAIN / EXPLAIN ANALYZE
```sql
EXPLAIN ANALYZE SELECT * FROM orders WHERE user_id = 123;
-- Показывает: Seq Scan vs Index Scan, actual rows, execution time
-- Seq Scan на большой таблице = нет нужного индекса
```

### Вопросы которые реально задают
1. Объясни MVCC, зачем он нужен?
2. Что такое dead tuples и autovacuum?
3. Какие типы индексов есть в PostgreSQL, когда GIN vs B-tree?
4. Разница REPEATABLE READ и SERIALIZABLE?
5. SELECT FOR UPDATE vs SELECT FOR SHARE?
6. Что такое WAL, зачем он?
7. JSON vs JSONB — разница?
8. Что такое partial index, когда полезен?
9. Как оптимизировать slow query? (EXPLAIN ANALYZE, pg_stat_statements)
10. Что такое connection pooling, зачем PgBouncer?

---

## 🔵 System Design (приоритет 3)

### Под эту конкретную вакансию — API Gateway

**Вопросы с реальных интервью:**
- API Gateway vs Load Balancer — в чём разница?
- Reverse Proxy vs Forward Proxy?
- Как работает Rate Limiter? (token bucket, leaky bucket, sliding window)
- Как работает Circuit Breaker? (Closed → Open → Half-Open)
- Microservices vs Monolith — trade-offs?
- Как обеспечить idempotency в API?
- Retry + exponential backoff — когда и как?
- SSO — как работает?
- Kafka vs RabbitMQ — когда что?

**API Gateway паттерны (прямо из вакансии):**
```
Request Aggregation  — объединение нескольких service calls в один
Response Aggregation — сборка ответов из нескольких сервисов
Caching              — кэш часто запрашиваемых данных
Auth/AuthZ           — JWT или OAuth 2.0 на gateway уровне
Rate Limiting        — защита downstream сервисов
Protocol Translation — REST → gRPC и обратно
```

**Circuit Breaker состояния:**
```
CLOSED      → нормальная работа, запросы проходят
OPEN        → сервис упал, запросы сразу отклоняются (fail fast)
HALF-OPEN   → тестовые запросы, если успешны → CLOSED
```

**Distributed Systems паттерны (из JD "будет плюсом"):**
- **Retry** — повтор с backoff при transient failures
- **Backoff** — exponential delay между retry (+ jitter)
- **Idempotency** — одинаковый результат при повторном запросе (idempotency key)
- **Saga** — distributed transaction через sequence of local transactions
- **Outbox Pattern** — надёжная публикация событий в message broker

### Ресурсы для SD
- https://github.com/ashishps1/awesome-system-design-resources
- https://github.com/donnemartin/system-design-primer
- ByteByteGo (Alex Xu) — лучшая книга по SD interviews

---

## 🟢 Core Go теория (приоритет 4)

### Вопросы которые точно спросят
1. Как работает GC в Go? (tri-color mark-and-sweep, stop-the-world)
2. Интерфейсы в Go — как реализованы внутри? (iface, eface)
3. Разница между pointer receiver и value receiver?
4. Что такое escape analysis?
5. Как работает `defer`? (LIFO, выполняется после return)
6. `new()` vs `make()` — разница?
7. Как работает `init()`?
8. Пустой интерфейс `interface{}` vs `any`?
9. Generics в Go — когда использовать?
10. Как профилировать Go приложение? (pprof, trace)

### Slices — частые ловушки
```go
// Ловушка: append может создать новый underlying array
a := make([]int, 3, 5)
b := a[1:3]
b[0] = 99 // меняет и a тоже! (shared underlying array)

// Безопасное копирование:
b := make([]int, len(a))
copy(b, a)
```

### Error handling
```go
// Go 1.13+ — wrapping errors
if err := doSomething(); err != nil {
    return fmt.Errorf("context: %w", err)
}

// Unwrapping
if errors.Is(err, ErrNotFound) { ... }
var myErr *MyError
if errors.As(err, &myErr) { ... }
```

---

## 🟢 Microservices / Transport (приоритет 5)

### gRPC vs REST
| | REST | gRPC |
|--|------|------|
| Protocol | HTTP/1.1 | HTTP/2 |
| Format | JSON | Protobuf (binary) |
| Performance | Ниже | Выше |
| Streaming | Ограничено | Bidirectional |
| Contract | OpenAPI | .proto файлы |

### Observability (прямо из JD)
- **Tracing** — OpenTelemetry/Jaeger, trace_id через все сервисы
- **Metrics** — Prometheus (pull), counter/gauge/histogram
- **Logging** — structured logging (JSON), correlation с trace_id
- **Grafana** — дашборды по Prometheus метрикам

### Testing patterns
- **Unit tests** — мокирование интерфейсов (`testify/mock`)
- **Integration tests** — реальные зависимости (testcontainers-go)
- **Contract tests** — Pact, проверка API контракта между сервисами

---

## 📋 Итоговый план подготовки

| Приоритет | Тема | Ресурс | Статус |
|-----------|------|--------|--------|
| 🔴 1 | Go concurrency: GMP, channels, patterns | RezaSi/go-interview-practice | ☐ |
| 🔴 2 | PostgreSQL: MVCC, индексы, locking | secondtalent.com/interview-guide/postgresql | ☐ |
| 🟡 3 | System Design: API Gateway, circuit breaker | ashishps1/awesome-system-design-resources | ☐ |
| 🟡 4 | Go теория: GC, interfaces, scheduler | Devinterview-io/golang-interview-questions | ☐ |
| 🟢 5 | Microservices паттерны + observability | Devinterview-io/microservices-interview-questions | ☐ |

---

## 🔗 Все ссылки одним списком

### GitHub
- https://github.com/RezaSi/go-interview-practice
- https://github.com/Devinterview-io/golang-interview-questions
- https://github.com/shomali11/go-interview
- https://github.com/ashishps1/awesome-system-design-resources
- https://github.com/donnemartin/system-design-primer
- https://github.com/Devinterview-io/microservices-interview-questions
- https://github.com/madd86/awesome-system-design

### Статьи / гайды
- https://www.secondtalent.com/interview-guide/golang/ (23 Advanced Go questions)
- https://www.secondtalent.com/interview-guide/postgresql/ (30 Advanced PG questions)
- https://dev.to/crusty0gphr/tricky-golang-interview-questions-part-4-concurrent-consumption-34oe
- https://medium.com/@abhigyandwivedi/golang-concurrency-interview-questions-and-answers-80c688904471
- https://dsysd-dev.medium.com/20-advanced-questions-asked-for-a-senior-developer-position-interview-1a65203e5d5e
