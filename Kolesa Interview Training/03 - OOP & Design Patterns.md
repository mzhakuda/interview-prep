---
tags:
  - kolesa
  - oop
  - design-patterns
topic: OOP & Patterns
priority: high
status: ready
---

# ООП & Паттерны проектирования

> В Go нет классического ООП, но собеседования это всё равно проверяют. Важно уметь перевести концепции ООП в идиоматичный Go.

---

## SOLID в Go (обязательный минимум)

### S — Single Responsibility
Одна структура / пакет = одна причина меняться.
```go
// плохо — делает всё
type UserService struct{}
func (s *UserService) Create(u User) {}
func (s *UserService) SendEmail(u User) {}
func (s *UserService) LogToFile(u User) {}

// хорошо
type UserRepo interface { Create(u User) error }
type Mailer interface { Send(to, body string) error }
```

### O — Open/Closed
Открыто для расширения, закрыто для модификации.
- В Go достигается через **интерфейсы**: новое поведение = новая реализация интерфейса.

### L — Liskov Substitution
Реализация интерфейса должна работать везде, где ожидается интерфейс. Если `io.Reader.Read` в одной реализации иногда возвращает 0 байт без `io.EOF` — это нарушение LSP.

### I — Interface Segregation
**Маленькие интерфейсы в Go — канон.** `io.Reader`, `io.Writer` — 1 метод.
> *"The bigger the interface, the weaker the abstraction" — Rob Pike*

### D — Dependency Inversion
Зависеть от интерфейсов, не от реализаций. **Интерфейс объявляется на стороне потребителя** (это Go-идиома).
```go
// в пакете service/user — потребитель
type userRepo interface { GetByID(id string) (User, error) }
type Service struct { repo userRepo }

// в пакете storage/postgres — реализация (не знает про интерфейс)
type Repo struct { db *sql.DB }
func (r *Repo) GetByID(id string) (User, error) { ... }
```

---

## Паттерны — с примерами из твоего опыта

### 1. Circuit Breaker ⭐ (ты сам реализовывал)
**Зачем:** защита от каскадных отказов при интеграции с внешними системами.
**Состояния:** Closed → Open → Half-Open.
**Где применял:** 10+ интеграций в BCC Business (CBS, CEA, GovBus) — per-client конфиг, мониторинг состояний.

### 2. Decorator (generic handler в Go)
**Где применял:** Generic handler decorator в BCC — композируемые middleware для auth/validation/logging с санитизацией полей.
```go
func WithAuth[T any](next Handler[T]) Handler[T] { ... }
func WithLogging[T any](next Handler[T]) Handler[T] { ... }
```

### 3. Repository
Абстракция над хранилищем. В тестах меняется на mock.

### 4. Strategy
Разные алгоритмы за одним интерфейсом. Пример — ORA-код парсер: разные стратегии для классификации ошибок.

### 5. Factory
Создание объектов с разной конфигурацией.
```go
func NewHTTPClient(opts ...Option) *Client { ... }  // functional options
```

### 6. Observer / Pub-Sub
RabbitMQ/Kafka. Ты делал notification-сервис через Temporal — тоже pub-sub по сути.

### 7. Saga / Orchestration vs Choreography ⭐ (Temporal)
**Orchestration:** центральный координатор (Temporal workflow) — у тебя depoisits, document exchange.
**Choreography:** события в очереди, сервисы подписаны — RabbitMQ/Kafka.

### 8. Observer: Singleton в Go
`sync.Once` — ленивая инициализация one-time. Глобальных синглтонов избегать.

### 9. Builder / Functional Options
**Идиомна для Go** вместо классического Builder:
```go
type Option func(*Server)
func WithTimeout(d time.Duration) Option { return func(s *Server) { s.timeout = d } }
srv := NewServer(WithTimeout(5*time.Second), WithAuth(...))
```

### 10. Adapter
Твой `io.Reader` поверх gRPC stream — классический Adapter.

---

## Композиция vs Наследование

Go **не поддерживает наследование** — только композицию (embedding).
```go
type Base struct { ID string }
type User struct {
    Base              // embedding — все методы Base доступны на User
    Email string
}
```
**Плюсы:** гибче, нет diamond problem, явные зависимости.
**Минусы:** больше boilerplate при большом числе embedded-структур.

---

## Интерфейсы в Go

- **Duck typing**: реализация неявная, не надо `implements`
- **Empty interface** `interface{}` / `any` — убегать от этого, использовать generics
- **Type assertion:** `v, ok := x.(SomeType)` — всегда с `ok`, panic иначе
- **Type switch:**
```go
switch v := x.(type) {
case int:    ...
case string: ...
default:     ...
}
```
- **nil-interface trap:** интерфейс != nil если содержит typed nil
```go
var p *MyStruct = nil
var i interface{} = p
fmt.Println(i == nil) // false! ловушка
```

---

## Типичные вопросы

**1. Разница абстрактного класса и интерфейса?**
В Go нет абстрактных классов. Аналог "абстрактного класса" = интерфейс + embedded реализация по умолчанию.

**2. Почему предпочитаем композицию наследованию?**
- Нет diamond problem
- Легче заменить поведение (inject другую зависимость)
- Явные связи читаются проще

**3. DI в Go без фреймворка?**
Конструктор + интерфейсы + functional options. Фреймворки типа `wire` / `fx` — только если большая кодобаза. На старте — ручной DI.

**4. Когда использовать generics?**
- Когда нужно одно и то же поведение для разных типов без потери типобезопасности
- Коллекции, middleware, общие утилиты
- **Не использовать** для "метапрограммирования" — Go это не любит

**5. Как применить паттерн Strategy в Go?**
Через интерфейс + конструктор, который принимает нужную реализацию. Никаких `switch case` в бизнес-логике.

---

## Go-анти-паттерны (показывают опыт)

- **God package** — один пакет с кучей несвязанного (util, common)
- **Stringly-typed API** — всё string, нет типов
- **Игнор ошибок** через `_ = err`
- **Панические "обработчики"** вместо возврата error
- **Over-engineering с интерфейсами** — интерфейс под один тип, "чтоб было"
