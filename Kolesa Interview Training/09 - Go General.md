---
tags:
  - kolesa
  - go
  - fundamentals
topic: Go General
priority: critical
status: ready
---

# Go — общие вопросы

> Полный материал: [[Go Interview Prep/00 - Cheat Sheet]], [[Go Interview Prep/01 - Go Concurrency]], [[Go Interview Prep/07 - Core Go]].
> Тут — сконцентрированные ответы на то, что гарантированно спросят.

---

## Базовые типы

| Тип | Размер / поведение |
|-----|---|
| `int` / `uint` | 32 или 64 бита (зависит от платформы) |
| `int32`, `int64` | фикс размер |
| `float32`, `float64` | IEEE 754 |
| `string` | immutable, pointer + length |
| `[]byte` | slice, mutable |
| `rune` | alias для `int32`, 1 unicode code point |
| `byte` | alias для `uint8` |

### String vs []byte
- `string` immutable, копирование при конверсии в `[]byte` (есть оптимизации компилятора)
- Индекс `s[i]` → `byte` (не rune!) — для unicode итерируем через `for i, r := range s`

---

## Slice — внутренности

```go
type slice struct {
    ptr    *T   // указатель на underlying array
    len    int  // текущая длина
    cap    int  // capacity (место до конца массива)
}
```

- `append` — если `len == cap`: аллоцирует новый массив (обычно x2 до 1024 элементов, потом +25%)
- `s[i:j]` — **шарит underlying array** с оригинальным slice → мутация видна обоим
- **Pitfall:** `append` может вернуть новый slice, всегда `s = append(s, x)`
- **Утечка памяти:** маленький slice от большого массива держит весь массив живым → `copy` в новый slice

---

## Map

- Hash-table с bucket'ами
- `m["k"]` → `value, ok` — idiom для проверки наличия
- **Не потокобезопасная** — для concurrent использовать `sync.Map` или mutex
- **Итерация в случайном порядке** — Go умышленно рандомизирует чтобы люди не полагались
- Значение `nil`-map — можно читать (вернёт zero), но **запись паникует**

---

## Goroutines & Scheduler (GMP)

- **G** (goroutine) — стек начинается с 2KB, растёт до 1GB
- **M** (OS thread)
- **P** (logical processor) — держит run-queue, количество = `GOMAXPROCS`
- **Work stealing** — idle P крадёт задачи у другого P
- **Preemption** — с Go 1.14 сигнальная (`SIGURG`), не блокируется на tight loop

### goroutine vs thread
- Стек goroutine маленький (2KB) → можно запускать миллионы
- Поток ОС = ~1-2MB стек → несколько тысяч максимум
- Переключение goroutine — cheap (в user space), поток — дорого (kernel)

---

## Channels

| Операция | nil chan | closed chan |
|----------|----------|-------------|
| send | блок навсегда | **panic** |
| recv | блок навсегда | zero value + `ok=false` |
| close | **panic** | **panic** |

- Unbuffered: рандеву (send блокирует до recv)
- Buffered `make(chan T, N)`: send блок только когда буфер полон
- `select` с несколькими готовыми — рандомно
- `select { default: }` — неблокирующий

**Правила:**
- Close **только отправитель**
- Никогда не close если могут ещё отправить

---

## Sync

- **`sync.Mutex`** — эксклюзив. Не копировать (передавать по указателю).
- **`sync.RWMutex`** — много читателей ИЛИ один писатель. Использовать когда reads >> writes.
- **`sync.WaitGroup`** — `Add(1)` ДО запуска goroutine, `Done()` в `defer`, `Wait()` снаружи.
- **`sync.Once`** — ровно один раз, safe для concurrent `Do`.
- **`sync.Pool`** — кеш переиспользуемых объектов, снижает GC pressure. Объекты могут быть собраны GC.
- **`sync/atomic`** — атомарные операции на примитивах (int32, pointer). Быстрее mutex для счётчиков.

---

## Context

- **Cancellation** — propagate cancel сигнала по дереву вызовов
- `context.Background()` — корневой
- `context.WithCancel(parent)` — отменяемый вручную
- `context.WithTimeout(parent, d)` — таймаут
- `context.WithDeadline(parent, t)` — дедлайн (абсолютное время)
- `context.WithValue(parent, key, val)` — передача данных (⚠️ не злоупотреблять!)

**Правила:**
- Первый аргумент функции: `ctx context.Context`
- Не хранить в структурах (кроме специальных случаев)
- **Проверять `<-ctx.Done()`** в длинных операциях
- `ctx.Err()` — причина отмены (`Canceled` / `DeadlineExceeded`)

---

## Error handling

- Ошибки — **значения**, возвращаются обычно последним аргументом
- `errors.Is(err, target)` — проверка на sentinel error
- `errors.As(err, &target)` — unwrap в конкретный тип
- `fmt.Errorf("context: %w", err)` — оборачивание

### Идиомы
```go
if err != nil {
    return fmt.Errorf("create user: %w", err)
}
```

- **panic** — только для "не должно случиться" (nil pointer, invalid state). Не для бизнес-ошибок.
- **recover** — только на границе (например HTTP handler не должен крашить процесс)

---

## Goroutine leaks

Классический случай:
```go
ch := make(chan int)
go func() { ch <- compute() }()
// ... используем ctx.Done, но не читаем из ch → goroutine блокирует send навсегда
```

**Как избежать:**
- Buffered channel на 1 если не факт что прочитают
- `select { case ch <- v: case <-ctx.Done(): return }`
- Всегда думать "кто закроет / прочитает"

---

## GC (garbage collector)

- **Concurrent mark-sweep**, tri-color
- Цель: `GOGC=100` (default) — запускается когда heap вырос x2 от baseline
- **Stop-the-world** — очень короткие паузы (sub-ms)
- **Write barriers** во время mark-фазы

### Снижение GC pressure
- `sync.Pool`
- Переиспользование slices (`s = s[:0]`)
- Избегать лишних allocations (предаллокация, указатели vs value-копии)
- `pprof heap` для диагностики

---

## Stack vs Heap

- Go-компилятор делает **escape analysis** — если значение не escape'ит за пределы функции → stack
- Escape: возврат указателя, передача в goroutine, в interface{}, в channel
- `go build -gcflags="-m"` — видим что escape'ит

---

## Generics (1.18+)

```go
func Max[T constraints.Ordered](a, b T) T {
    if a > b { return a }
    return b
}
```
- Type parameters в `[]`
- Constraints = интерфейсы (могут иметь `~int | string`)
- Осторожно: generics **не** полиморфизм в classical OOP — не замена интерфейсам

---

## Interfaces

- Implicit implementation (duck typing)
- **Empty interface `interface{}` / `any`** — можно принять что угодно, но теряем типы
- Type assertion: `v, ok := x.(T)`
- Type switch: `switch v := x.(type) { case int: ... }`

### Nil-interface trap
```go
var p *MyErr = nil
var err error = p
err != nil  // TRUE! (typed nil)
```
Решение: `if p == nil { return nil }` на уровне конкретного типа.

---

## Топ-10 вопросов которые спросят

**1. Что такое goroutine, чем отличается от потока?**
Lightweight user-level поток, 2KB стек, multiplex на OS threads через scheduler. Переключение cheap, можно запустить миллионы.

**2. Разница buffered / unbuffered channel?**
Unbuffered — send блок до recv (синхронный handshake). Buffered — блок только когда буфер полон.

**3. Как передать cancellation между goroutines?**
`context.Context` — `WithCancel` / `WithTimeout`, в длинных функциях `<-ctx.Done()`.

**4. Как избежать race condition?**
- `go test -race`
- sync.Mutex / RWMutex
- каналы (share memory by communicating)
- atomic для простых счётчиков

**5. Defer — как работает?**
LIFO (последний deferred выполняется первым). Вычисляется аргументы при defer, выполняется при return. Осторожно в tight loop (deferred освобождается только при возврате из функции).

**6. slice vs array?**
Array — фиксированный размер, value-тип (копируется). Slice — дескриптор (ptr + len + cap), reference-like.

**7. Error handling — почему не исключения?**
Go-философия: ошибки — часть API, должны быть explicit. Panics — для catastrophic. Легче reason about, но boilerplate.

**8. Как работает `select`?**
Выбирает готовый `case` из нескольких каналов. Несколько готовых — рандом. `default` делает неблокирующим. Пустой `select{}` — блок навсегда.

**9. Как делать HTTP-сервер с graceful shutdown?**
```go
srv := &http.Server{Addr: ":8080"}
go srv.ListenAndServe()
ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM)
defer stop()
<-ctx.Done()
shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
srv.Shutdown(shutdownCtx)
```

**10. Методы с pointer receiver vs value receiver?**
- **Pointer receiver** — если метод мутирует, или если структура большая, или для консистентности (если хоть один метод pointer — все должны быть)
- **Value receiver** — маленькие immutable структуры (time.Time, UUID)

---

## Проектная специфика (упомянуть если спросят)

- **gRPC bidirectional streaming** — io.Pipe, zero-copy, chunks
- **Temporal workflows** — orchestration, retry policies, ContextPropagator
- **Generic handler decorator** — composable middleware через generics
- **Circuit breaker** — fallback при деградации внешних сервисов
- **Colvir ORA-коды** — классификация retryable vs non-retryable errors
