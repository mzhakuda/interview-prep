---
tags:
  - cheatsheet
  - go
  - postgres
  - backend
topic: Cheat Sheet
priority: critical
status: ready
---

# Backend Interview — Cheat Sheet

---

# GO

---

## GMP — планировщик

- **G** (goroutine) — стек начинается с `2KB`, растёт до `1GB`
- **M** (OS thread) — реальный поток ОС
- **P** (processor) — логический процессор, держит run queue; количество = `GOMAXPROCS` (default = кол-во CPU ядер)
- У каждого P — **local run queue** (макс 256 G); переполнение → global queue
- **Work stealing**: idle P крадёт половину очереди у другого P
- **Preemption**: с Go 1.14 — сигнальная (`SIGURG`), не блокируется на tight loops
- Syscall → M отсоединяется от P, P берёт другой M (или создаёт)
- `runtime.Gosched()` — уступить планировщику (не засыпает)
- `runtime.LockOSThread()` — прибить горутину к конкретному M (для CGo, OpenGL)

---

## Channels

**Поведение nil / closed канала:**

| Операция | nil chan | closed chan |
|----------|----------|------------|
| send | блокирует навсегда | **panic** |
| recv | блокирует навсегда | zero value + `false` |
| close | **panic** | **panic** |

**Buffered vs Unbuffered:**
- Unbuffered: send блокирует до тех пор, пока не появится receiver (синхронный рандеву)
- Buffered `make(chan T, N)`: send блокирует только когда буфер полон
- `len(ch)` — сколько элементов в буфере; `cap(ch)` — размер буфера

**Правила:**
- Close только отправитель, никогда получатель
- Никогда не закрывать канал, в который могут ещё отправить
- `select` при нескольких готовых case выбирает **случайно**
- `default` в select делает его неблокирующим

```go
func producer(ch chan<- int) {}  // только отправка
func consumer(ch <-chan int) {}  // только получение
```

---

## Sync примитивы

**Mutex vs RWMutex:**
- `sync.Mutex` — эксклюзивный lock (Lock/Unlock)
- `sync.RWMutex` — много читателей одновременно ИЛИ один писатель (RLock/Lock)
- Использовать RWMutex когда reads >> writes
- **Не копировать mutex** — передавать только по указателю

**WaitGroup:**
```go
var wg sync.WaitGroup
wg.Add(1)
go func() { defer wg.Done(); ... }()
wg.Wait()
```
- `Add` вызывать ДО запуска горутины (иначе race)

**sync.Once** — выполняет функцию ровно один раз, безопасно при concurrent вызовах

**sync.Pool** — кэш переиспользуемых объектов, снижает GC pressure; объекты могут быть GC-нуты в любой момент

**atomic** — `atomic.AddInt64`, `LoadInt64`, `StoreInt64`, `CompareAndSwapInt64` — счётчики/флаги без mutex

**context.Context:**
- `WithTimeout` / `WithDeadline` / `WithCancel`
- Всегда первый аргумент функции
- `ctx.Done()` — channel, закрывается при отмене
- `ctx.Err()` → `context.Canceled` или `context.DeadlineExceeded`

---

## Goroutine leaks — ловушки

```go
// УТЕЧКА: горутина заблокирована на send, никто не читает
ch := make(chan int)
go func() { ch <- 1 }()

// FIX: буферизованный канал или гарантированный receiver
```

```go
// ЛОВУШКА: closure захватывает переменную цикла
for _, v := range items {
    go func() { process(v) }()  // все горутины видят одно и то же v
}
// FIX — засхадовить:
for _, v := range items {
    v := v
    go func() { process(v) }()
}
```

```go
// PANIC: одновременная запись в map
m := map[string]int{}
go func() { m["a"] = 1 }()
go func() { m["b"] = 2 }()  // fatal error: concurrent map writes
// FIX: sync.RWMutex или sync.Map
```

---

## Stack vs Heap — что где живёт

| | Stack | Heap |
|--|-------|------|
| Аллокация | сдвиг stack pointer (~1ns) | поиск свободного блока (~30-100ns) |
| Освобождение | автоматически при return (frame pop) | GC (периодически) |
| Размер | начинается 2KB, растёт до 1GB per goroutine | ограничен только RAM |
| Locality | данные рядом → cache-friendly | данные разбросаны → cache miss |
| Скорость | **~30-100x быстрее** аллокации | медленнее + GC pressure |

**Stack быстрее в 30-100x** — поэтому escape analysis критичен.

**На stack попадает:**
- Локальные переменные, не выходящие за пределы функции
- Аргументы и return values (small types by value)
- Переменные, размер которых известен compile-time

**На heap попадает (escape на heap):**
- Указатель **возвращается** из функции: `return &obj`
- Передача в **interface**: `fmt.Println(x)` — x escapes (даже `int`!)
- Переменная захвачена **closure** / goroutine
- Размер неизвестен compile-time: `make([]byte, n)` где `n` — переменная
- Объект **слишком большой** для stack frame

```go
// Stack — не escapes
func sum(a, b int) int {
    result := a + b  // result на стеке
    return result    // возвращаем значение, не указатель
}

// Heap — escapes
func newUser() *User {
    u := User{Name: "Alice"}  // u уйдёт на heap
    return &u                  // указатель выходит из функции
}

// Heap — interface escape
func log(v any) { fmt.Println(v) }
x := 42
log(x)  // x escapes to heap (боксируется в interface)

// Heap — closure capture
func counter() func() int {
    n := 0          // n на heap (захвачена closure)
    return func() int { n++; return n }
}
```

```bash
# Посмотреть что escapes
go build -gcflags="-m" ./...
# ./main.go:8:2: moved to heap: u
# ./main.go:14:6: x escapes to heap
```

**struct alignment — экономия памяти:**
```go
// Плохо: 24 байта (padding между полями)
type Bad struct {
    a bool   // 1 + 7 padding
    b int64  // 8
    c bool   // 1 + 7 padding
}

// Хорошо: 16 байт (поля от большего к меньшему)
type Good struct {
    b int64  // 8
    a bool   // 1
    c bool   // 1 + 6 padding
}
```
Правило: **от большего к меньшему**. Инструмент: `go vet -fieldalignment`

---

## GC — сборщик мусора

**Алгоритм: tri-color mark-and-sweep**, конкурентный

**GC roots — откуда начинается обход:**
- **Стеки всех горутин** — локальные переменные, аргументы функций (самый большой источник)
- **Глобальные переменные** — package-level vars
- **CPU регистры** — текущие значения во всех M (OS threads)
- **Finalizer queue** — объекты с `runtime.SetFinalizer`, ожидающие финализации

Всё, что достижимо из roots → **live** (помечается black). Всё, что не достижимо → **white** → удаляется при sweep.

```
GC roots
  │
  ├── goroutine stacks ──→ heap objects ──→ другие heap objects
  │       (G1 stack)            │
  ├── goroutine stacks          └──→ ещё объекты ...
  │       (G2 stack)
  ├── global vars
  └── CPU registers

  Всё недостижимое → white → sweep → освобождено
```

**Фазы:**
1. **Mark Setup** (STW ~10-30μs) — включить write barrier, snapshot всех goroutine стеков
2. **Marking** (concurrent) — обход из roots, раскрашивание grey → black
3. **Mark Termination** (STW ~10-30μs) — дообработать оставшееся grey, выключить write barrier
4. **Sweep** (concurrent) — освобождение white объектов

**Три цвета:**
- **White** — не посещён → кандидат на удаление
- **Grey** — посещён, потомки ещё не проверены → в work queue
- **Black** — посещён, все потомки проверены (live)

**Write barrier** — если горутина создаёт ссылку black → white во время marking, barrier перекрашивает white → grey (не даёт потерять объект)

**Настройка:**
- **`GOGC=100`** (default): GC при удвоении heap от live heap; `GOGC=50` — чаще, меньше памяти; `GOGC=off` — выключить
- **`GOMEMLIMIT`** (Go 1.19+): soft лимит; GC агрессивнее при приближении; рекомендация для контейнеров: `container_limit × 0.9`
- STW паузы — **sub-millisecond** в современном Go

```bash
# Профилирование
go tool pprof http://localhost:6060/debug/pprof/heap      # heap allocations
go tool pprof http://localhost:6060/debug/pprof/goroutine # goroutine leaks
go test -bench=. -benchmem -memprofile=mem.prof           # benchmark + memory
go tool pprof mem.prof
```

---

## Базовые типы данных

**Числа:**
```
int8 / int16 / int32 / int64
uint8 / uint16 / uint32 / uint64
int / uint          — размер зависит от платформы (64bit → 64bit)
float32 / float64
complex64 / complex128
byte  = uint8       — алиас, используется для raw bytes
rune  = int32       — алиас, используется для Unicode code point
uintptr             — целое для хранения указателя (unsafe арифметика)
```

**Булево и строки:**
```
bool                — true / false
string              — immutable, UTF-8 bytes, {ptr, len}
```

**Нулевые значения (zero values):**
```
int, float → 0
bool       → false
string     → ""  (пустая строка, не nil)
pointer    → nil
slice      → nil  (len=0, cap=0, ptr=nil)
map        → nil  (обращение к nil map → panic при записи)
channel    → nil
interface  → nil
```

---

## char / rune / byte

Go не имеет типа `char`. Вместо него:

| Тип | Алиас | Размер | Для чего |
|-----|-------|--------|---------|
| `byte` | `uint8` | 1 байт | raw bytes, ASCII символы |
| `rune` | `int32` | 4 байта | Unicode code point (любой символ) |

```go
var b byte = 'A'       // 65
var r rune = 'Я'       // 1071 (Unicode)
var r2 rune = '🔥'     // 128293

s := "Привет"
len(s)                 // 12 — байт (UTF-8: каждый кирилл. = 2 байта)
len([]rune(s))         // 6  — символов (рун)

// Итерация по байтам:
for i := 0; i < len(s); i++ { _ = s[i] }  // byte

// Итерация по рунам:
for i, r := range s { _, _ = i, r }        // rune, i — байтовый offset
```

---

## String — immutable

- Внутри: `{ptr *byte, len int}` — 16 байт на 64-bit
- **Immutable**: `s[0] = 'a'` → compile error
- Substring `s[1:5]` — **не копирует**, делит backing array → большая строка держится в памяти
- Конкатенация `s + t` — создаёт новую строку (аллокация)
- Много конкатенаций → `strings.Builder` (не аллоцирует каждый раз)

```go
// Конвертации
[]byte("hello")         // string → []byte (копия)
string([]byte{104,105}) // []byte → string (копия)
string(rune(1071))      // rune → string: "Я"
[]rune("Привет")        // string → []rune (копия, для посимвольной работы)

// Сравнение — по значению
"abc" == "abc"   // true
"abc" < "abd"    // true (лексикографически)
```

---

## Array — value type, фиксированный размер

- **Размер — часть типа**: `[3]int` и `[4]int` — разные типы
- **Mutable** — элементы можно менять
- **Value type**: присваивание и передача в функцию = **полная копия**
- Редко используется напрямую — чаще слайсы

```go
a := [3]int{1, 2, 3}
b := a              // b — независимая копия
b[0] = 99
fmt.Println(a[0])   // 1 (a не изменился)

// Размер известен compile-time:
var arr [1024]byte  // на стеке (если не escapes)
```

---

## Slice — внутренности

**SliceHeader — 24 байта:**
```
┌──────────────────────────────────────┐
│ ptr  *T  (8 байт) → backing array   │
│ len  int (8 байт) → кол-во элементов│
│ cap  int (8 байт) → размер массива  │
└──────────────────────────────────────┘
```

**Backing array в памяти:**
```
s := []int{10, 20, 30, 40, 50}
         cap=5, len=5

backing array:  [ 10 | 20 | 30 | 40 | 50 ]
                  ^ptr

b := s[1:3]     // b.ptr → &array[1], len=2, cap=4
                  [ 10 | 20 | 30 | 40 | 50 ]
                         ^ptr(b)
                         |←len=2→|←cap=4→|

b[0] = 99       // меняет array[1] — видно в s тоже!
```

**Как растёт cap при append:**
```
cap < 256  → удваивается (×2)
cap ≥ 256  → растёт на ~25% + сглаживание
```

```go
// append — самое важное поведение:
a := []int{1, 2, 3}   // len=3, cap=3

b := append(a, 4)      // cap превышен → новый backing array
b[0] = 99
fmt.Println(a[0])      // 1 — a и b теперь НЕЗАВИСИМЫ

c := a[:2]             // len=2, cap=3, ТОТЖЕ backing array
c = append(c, 99)      // cap не превышен → пишет в a[2]!
fmt.Println(a)         // [1 2 99] — a изменился!

// Защита через full slice expression:
c = a[:2:2]            // len=2, cap=2 — cap обрезан
c = append(c, 99)      // теперь новый backing array, a не меняется
```

**nil vs empty — разница важна в API:**
```go
var s []int          // nil,   len=0, cap=0, s == nil → true
s = []int{}          // empty, len=0, cap=0, s == nil → false
s = make([]int, 0)   // empty, len=0, cap=0, s == nil → false

// Оба безопасны для range и append
// НО: json.Marshal(nil slice) → "null"
//     json.Marshal(empty slice) → "[]"
```

---

## Map — внутренности

**hmap struct** (runtime/map.go):
```
┌────────────────────────────────────────┐
│ count     int        — кол-во элементов│
│ flags     uint8      — concurrent lock │
│ B         uint8      — log2(кол-во buckets): 2^B buckets
│ noverflow uint16     — overflow buckets│
│ hash0     uint32     — seed для хэша   │
│ buckets   *bmap      — массив buckets  │
│ oldbuckets *bmap     — старые buckets (при росте)
│ nevacuate  uintptr   — прогресс эвакуации
└────────────────────────────────────────┘
```

**Bucket (bmap) — 8 пар key/value:**
```
┌──────────────────────────────────────────┐
│ tophash  [8]uint8   — верхние 8 бит хэша│
│ keys     [8]K       — 8 ключей           │
│ values   [8]V       — 8 значений         │
│ overflow *bmap      — указатель на след. bucket (при коллизии)
└──────────────────────────────────────────┘
```

**Как работает lookup `m["key"]`:**
```
1. hash := hashFunc("key", m.hash0)
2. bucket_index = hash & (2^B - 1)    // младшие B бит
3. tophash = hash >> 56                // верхние 8 бит
4. Идём по bucket, сравниваем tophash → если совпал, сравниваем key
5. Если bucket заполнен → идём по overflow цепочке
```

**Load factor и рост:**
- Load factor ≈ **6.5** (среднее кол-во элементов на bucket)
- При превышении → **evacuation**: создаётся новый массив buckets ×2, данные переносятся постепенно (incremental, не STW)
- Во время эвакуации `oldbuckets` и `buckets` оба живут

**Почему порядок итерации случаен:**
- Намеренно рандомизирован с Go 1.0 (чтобы не полагались на порядок)
- Итерация начинается со случайного bucket + случайного offset внутри bucket

**Почему не thread-safe:**
- При конкурентной записи нет лока — `fatal error: concurrent map read and map write`
- Решение: `sync.RWMutex` или `sync.Map`

```go
m := map[string]int{"a": 1}
var m2 map[string]int       // nil map

// nil map:
_ = m2["x"]    // ok → 0 (zero value)
m2["x"] = 1    // PANIC

// Безопасное чтение:
v, ok := m["key"]           // ok=false если нет

// Удаление:
delete(m, "key")             // ок если ключа нет

// Ключи только comparable:
// ✓ bool, int, float, string, pointer, array, struct (без slice/map полей)
// ✗ slice, map, func
```

---

## Mutability — сводная таблица

| Тип | Mutable | Reference/Value | Zero value | Nil возможен |
|-----|---------|-----------------|------------|-------------|
| `int`, `float`, `bool` | да | value | `0` / `false` | нет |
| `string` | **нет** (immutable) | value (ptr+len) | `""` | нет |
| `array [N]T` | да (элементы) | **value** (копируется!) | `[N]zero` | нет |
| `slice []T` | да | reference (ptr+len+cap) | `nil` | да |
| `map[K]V` | да | reference | `nil` | да |
| `chan T` | — | reference | `nil` | да |
| `*T` (pointer) | да | reference | `nil` | да |
| `struct` | да | value (копируется!) | все поля zero | нет |
| `interface` | — | reference (itab+ptr) | `nil` | да |

---

## Interfaces & Type System

**Внутренняя структура:**
- `iface` (непустой interface): `{itab *itab, data unsafe.Pointer}` — itab содержит тип + таблицу методов
- `eface` (пустой interface `any`): `{_type *_type, data unsafe.Pointer}`

**Nil interface ловушка:**
```go
var p *MyStruct = nil
var i MyInterface = p
fmt.Println(i == nil)  // FALSE — interface содержит тип, data == nil
```
Interface равен nil только если оба поля (тип И данные) nil.

**Value vs Pointer receiver:**
- Pointer receiver: может мутировать, удовлетворяет interface только для указателя
- Value receiver: не мутирует, удовлетворяет interface для обоих
- Правило: если хотя бы один метод — pointer receiver → все должны быть pointer

---

## Slices

**Внутренняя структура:** `{ptr *T, len int, cap int}`

```go
s := make([]int, 3, 5)  // len=3, cap=5
s = append(s, 1)         // len=4, cap=5 — без аллокации
s = append(s, 1, 2)      // len=6, cap≥10 — новая аллокация, cap удваивается
```

**Ловушки:**
- Два слайса от одного массива **делят backing array** — мутации видны в обоих
- `append` может вернуть новый слайс (старый не изменится если cap превышен)
- `copy(dst, src)` — копирует min(len(dst), len(src)) без shared backing
- Nil слайс: `len=0, cap=0, ptr=nil` — безопасно append и range

---

---

# POSTGRESQL

---

## MVCC

**MVCC** (Multi-Version Concurrency Control) — основной механизм конкурентности PostgreSQL.
PostgreSQL хранит несколько версий одной строки одновременно.  
Каждая транзакция видит свой consistent snapshot данных.

- **MVCC** = читатели не блокируют писателей, писатели не блокируют читателей
- Каждая строка имеет `xmin` (создана транзакцией) и `xmax` (удалена/обновлена транзакцией)
- **Dead tuples**: старые версии строк после UPDATE/DELETE — не удаляются физически до VACUUM
- **VACUUM**: удаляет dead tuples, обновляет visibility map, не блокирует таблицу
- **VACUUM FULL**: перезаписывает таблицу — освобождает место на диске — требует эксклюзивную блокировку
- **HOT** (Heap Only Tuple): UPDATE создаёт новую версию строки на той же странице → индекс не обновляется (быстрее)
- **Transaction ID wraparound**: TxID 32-битный (лимит 4B) → VACUUM FREEZE помечает строки как замороженные

---

## ACID

- **Atomicity** — транзакция либо целиком выполнена, либо отменена; ROLLBACK отменяет всё
- **Consistency** — БД переходит из одного валидного состояния в другое; constraints не нарушаются
- **Isolation** — конкурентные транзакции не видят промежуточные состояния друг друга
- **Durability** — закоммиченная транзакция переживает сбой (WAL → fsync)

**WAL** (Write-Ahead Log): изменения пишутся в WAL до записи в data pages → crash recovery воспроизводит WAL

---

## Уровни изоляции

| Уровень | Dirty Read | Non-Repeatable Read | Phantom Read |
|---------|-----------|---------------------|-------------|
| READ UNCOMMITTED | возможен | возможен | возможен |
| **READ COMMITTED** *(дефолт PG)* | ✗ | возможен | возможен |
| REPEATABLE READ | ✗ | ✗ | ✗ в PG* |
| SERIALIZABLE | ✗ | ✗ | ✗ |

*PG REPEATABLE READ использует снапшоты → phantoms тоже предотвращаются

- **Dirty Read**: читаешь незакоммиченные данные другой транзакции
- **Non-Repeatable Read**: один и тот же SELECT возвращает разные данные внутри одной транзакции
- **Phantom Read**: строки появляются/исчезают из-за INSERT/DELETE другой транзакции

---

## Индексы

| Тип | Когда использовать | Особенности |
|-----|--------------------|-------------|
| **B-tree** | `=`, `<`, `>`, `BETWEEN`, `LIKE 'foo%'`, `ORDER BY` | дефолт; почти всегда правильный выбор |
| **GIN** | arrays, JSONB, full-text search (`@@`), `@>` | медленно строится, быстро ищет; для "содержит" |
| **GiST** | геометрия, range types, nearest-neighbor | lossy (нужен recheck); хорош для PostGIS |
| **BRIN** | огромные таблицы с естественным физическим порядком (timestamps, ID) | крошечный; грубый; быстрые вставки |

**Partial index:**
```sql
-- Индексирует только активных пользователей — меньше, быстрее
CREATE INDEX idx_active ON users(email) WHERE active = true;
```

**Covering index** (Index Only Scan):
```sql
-- Запрос отвечается только из индекса, без обращения к heap
CREATE INDEX idx_cover ON orders(user_id) INCLUDE (status, created_at);
```

**Составной индекс**: порядок важен — `(a, b, c)` используется для запросов по `a`, `a+b`, `a+b+c` — но не по `b` отдельно

---

## Блокировки

```sql
SELECT ... FOR UPDATE           -- эксклюзивная блокировка строки
SELECT ... FOR UPDATE NOWAIT    -- немедленно упасть если заблокировано
SELECT ... FOR UPDATE SKIP LOCKED  -- пропустить заблокированные строки (паттерн очереди)
SELECT ... FOR SHARE            -- разделяемая читающая блокировка
```

**Deadlock:**
- PG обнаруживает автоматически (~1с) → убивает одну транзакцию
- Предотвращение: всегда захватывать блокировки в одном порядке
- `lock_timeout = '2s'` — упасть если не удаётся получить блокировку за N времени
- DDL (ALTER TABLE) → AccessExclusive — блокирует всё, включая SELECT

---

## EXPLAIN ANALYZE

```sql
EXPLAIN ANALYZE SELECT ...;
```

| Узел | Что значит |
|------|-----------|
| **Seq Scan** | полный перебор таблицы — ок для маленьких таблиц или > 10-20% строк |
| **Index Scan** | использует индекс + fetch каждой строки из heap |
| **Index Only Scan** | все данные в индексе (covering index) — самый быстрый |
| **Bitmap Index Scan** | собирает указатели из индекса, потом batch fetch из heap |
| **Nested Loop** | для каждой строки outer — probe inner (хорош для маленьких наборов) |
| **Hash Join** | строит hash table из inner, probe outer (хорош для больших) |

**Чтение cost:** `cost=start..total rows=N width=байты`
- `actual time=start..end rows=N loops=K` — умножить на loops для реального времени
- Расхождение rows estimate vs actual → устаревшая статистика → `ANALYZE`

---

## Connection Pooling (PgBouncer)

| Режим | Соединение с PG | Prepared stmts | Когда использовать |
|-------|----------------|----------------|-------------------|
| **Session** | 1 клиент = 1 PG соединение на всю сессию | ✓ | мало клиентов |
| **Transaction** | соединение возвращается после COMMIT/ROLLBACK | ✗ | микросервисы, много коротких транзакций |
| **Statement** | соединение после каждого statement | ✗ | только autocommit запросы |

- Transaction mode — самый распространённый в backend сервисах
- Prepared statements + transaction mode → нужен `DEALLOCATE` или отключить prepared stmts

---

---

# gRPC vs REST

| | REST | gRPC |
|--|------|------|
| Протокол | HTTP/1.1 или HTTP/2 | HTTP/2 только |
| Формат | JSON (текст) | Protobuf (бинарный) |
| Схема | Опционально (OpenAPI) | Обязательно (`.proto`) |
| Кодогенерация | Опционально | Встроенная |
| Стриминг | Ограниченный (SSE, WS) | 4 типа нативно |
| Браузер | Нативно | Нужен прокси (grpc-web) |
| Производительность | — | ~5-10x меньше payload |

**4 типа gRPC стриминга:**
1. **Unary** — один запрос, один ответ
2. **Server streaming** — один запрос, поток ответов
3. **Client streaming** — поток запросов, один ответ
4. **Bidirectional** — обе стороны стримят одновременно

**Когда gRPC:** внутренние микросервисы, высокий throughput, стриминг, строгие контракты
**Когда REST:** публичные API, браузерные клиенты, простой CRUD

---

# Redis

**Структуры данных:**

| Тип | Команды | Когда использовать |
|-----|---------|-------------------|
| String | GET/SET/INCR/EXPIRE | кэш, счётчики, сессии |
| List | LPUSH/RPUSH/LPOP/LRANGE | очереди, последние N элементов |
| Set | SADD/SMEMBERS/SINTER | уникальные теги, связи |
| Sorted Set | ZADD/ZRANGE/ZRANGEBYSCORE | лидерборды, rate limiting |
| Hash | HSET/HGET/HGETALL | объекты, профили пользователей |
| Stream | XADD/XREAD/XACK | event log, message queue |

**Политики вытеснения** (когда достигнут `maxmemory`):

| Политика | Поведение |
|----------|----------|
| `noeviction` | возвращать ошибку на запись |
| `allkeys-lru` | вытеснять наименее недавно использованные (любой ключ) |
| `volatile-lru` | LRU среди ключей с TTL |
| `allkeys-lfu` | вытеснять наименее часто используемые |
| `volatile-ttl` | вытеснять ключ с наименьшим TTL |

**Персистентность:**
- **RDB** (снапшот): point-in-time дамп; быстрый рестарт; можно потерять последние N минут
- **AOF** (Append-Only File): логирует каждую запись; медленнее; почти нет потерь данных
- **RDB + AOF** = рекомендуется для production

---

# RabbitMQ

**Типы exchange:**

| Тип | Маршрутизация | Когда использовать |
|-----|---------------|--------------------|
| **Direct** | точное совпадение routing key | task queues, point-to-point |
| **Topic** | паттерн (`*.logs`, `#`) | маршрутизация событий по категории |
| **Fanout** | broadcast всем привязанным очередям | нотификации, pub/sub |
| **Headers** | совпадение по заголовкам сообщения | сложная маршрутизация без routing key |

**Подтверждения:**
- `auto-ack`: сообщение удаляется при доставке (риск: потеря при падении consumer)
- `manual-ack`: consumer отправляет `ack` после обработки (безопасно)
- `nack + requeue=true`: вернуть в очередь
- `nack + requeue=false`: отправить в Dead Letter Exchange (DLX)

**Dead Letter Exchange (DLX):** очередь с `x-dead-letter-exchange` — сообщения уходят туда при nack+no-requeue, истечении TTL, переполнении очереди

**Prefetch count (QOS):** ограничивает сколько неподтверждённых сообщений может держать consumer

---

# Docker / Kubernetes

**Docker:**
- `COPY` vs `ADD`: предпочитать COPY — ADD автораспаковывает архивы
- Multi-stage build: `FROM golang:1.22 AS builder` → `FROM alpine` (маленький образ)
- Layer caching: COPY исходники ПОСЛЕ установки зависимостей
- Не запускать от root: `USER nobody` в production

**Kubernetes:**

| Концепция | Что это |
|-----------|---------|
| Pod | минимальная единица; 1+ контейнеров, общая сеть/storage |
| Deployment | управляет ReplicaSet; rolling updates; rollback |
| Service | стабильный DNS + load balancing для группы pod |
| ConfigMap | несекретный конфиг; env vars или volume |
| Secret | base64-закодированный секретный конфиг (не зашифрован по умолчанию!) |
| HPA | Horizontal Pod Autoscaler; масштабирует по CPU/memory/custom метрикам |

**Probes:**

| Probe | Назначение | При неудаче |
|-------|-----------|-------------|
| **Liveness** | Жив ли под? (deadlock, зависание) | Перезапустить контейнер |
| **Readiness** | Готов ли принимать трафик? | Убрать из Service endpoints |
| **Startup** | Закончил ли стартовать? (медленный старт) | Блокирует liveness проверки |

**Resources:**
- `requests` — гарантированные ресурсы (используется для scheduling)
- `limits` — максимум (CPU throttled; memory → OOM kill)

---

# HTTP & Сети

**Версии HTTP:**

| | HTTP/1.1 | HTTP/2 | HTTP/3 |
|--|---------|--------|--------|
| Транспорт | TCP | TCP | QUIC (UDP) |
| Мультиплексинг | Нет (HOL blocking) | Да (streams) | Да (независимые streams) |
| Заголовки | Текст, повторяются | HPACK сжатие | QPACK сжатие |
| TLS | Опционально | Де-факто обязательно | Встроен |
| Применение | legacy | API, браузеры | мобильные, нестабильные сети |

**TCP vs UDP:**
- **TCP**: с установкой соединения, упорядоченный, надёжный → HTTP, DB, SSH
- **UDP**: без соединения, без гарантий, быстрый → DNS, видеостриминг, игры, QUIC

**TLS 1.3 Handshake:**
1. Client → ClientHello (поддерживаемые cipher suites, key share)
2. Server → ServerHello + Certificate + key share
3. Client проверяет сертификат → вычисляет session keys
4. Оба отправляют Finished → **1-RTT** (vs 2-RTT в TLS 1.2)

**HTTP статус коды:**
- `200` OK, `201` Created, `204` No Content
- `301` Permanent redirect, `304` Not Modified (кэш актуален)
- `400` Bad Request, `401` Unauthorized, `403` Forbidden, `404` Not Found
- `409` Conflict, `429` Too Many Requests
- `500` Internal Server Error, `502` Bad Gateway, `503` Service Unavailable, `504` Gateway Timeout

---

# CAP теорема

- **C**onsistency — каждый read видит последний write (или ошибку)
- **A**vailability — каждый запрос получает ответ (без ошибки, может быть устаревший)
- **P**artition Tolerance — система работает при сетевом разрыве

**Правило**: при partition нужно выбрать C или A. P — не опциональна в distributed системах.

| Выбор | Поведение | Примеры |
|-------|----------|---------|
| **CP** | возвращает ошибку вместо устаревших данных | etcd, Zookeeper, Consul, HBase |
| **AP** | возвращает потенциально устаревшие данные | Cassandra, DynamoDB, CouchDB |
| **CA** | невозможна в distributed системе | single-node Postgres (не distributed) |

**Модели консистентности (от слабой к сильной):**
- **Eventual** — все узлы сойдутся в итоге (без гарантий по времени) — Cassandra, DNS
- **Read-your-writes** — ты всегда видишь свои собственные записи
- **Monotonic read** — если прочитал значение X, не прочитаешь более старое
- **Strong / Linearizable** — операции атомарны, глобально упорядочены — etcd, Zookeeper
- **Serializable** — транзакции выглядят как последовательные — PostgreSQL SERIALIZABLE

**Реальные примеры:**
- Банковские переводы → Strong / Serializable
- Лента соцсетей → Eventual достаточно
- Auth сессии → минимум Read-your-writes
- Distributed конфиг (etcd) → CP (Linearizable)
