---
tags:
  - kolesa
  - testing
  - unit-tests
  - mocks
topic: Testing
priority: high
status: ready
---

# Тесты — Unit, Mocks, Asserts

---

## Пирамида тестов

```
          /\
         /E2E\       мало, дорогие, медленные — только golden path
        /------\
       /Integr. \    средне — реальная БД, реальный Kafka
      /----------\
     /   Unit     \  много, быстрые, дешёвые — ядро
    /--------------\
```

- **Unit:** чистая функция / метод, mock'и зависимостей
- **Integration:** hit реальные зависимости (Postgres через testcontainers)
- **E2E:** весь стек поднят, ходим через HTTP

---

## Unit-тесты в Go — базовое

```go
func TestAdd(t *testing.T) {
    got := Add(2, 3)
    if got != 5 {
        t.Errorf("Add(2,3) = %d, want 5", got)
    }
}
```

### Table-driven (⭐ твой подход в BCC)
```go
func TestParseORA(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    ErrorKind
        wantErr bool
    }{
        {"retryable infra", "ORA-12541", InfrastructureError, false},
        {"business reject", "ORA-20001", NonRetryable, false},
        {"unknown format", "garbage", Unknown, true},
    }
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            got, err := ParseORA(tc.input)
            if (err != nil) != tc.wantErr { t.Fatalf("err: %v", err) }
            if got != tc.want { t.Errorf("got %v, want %v", got, tc.want) }
        })
    }
}
```

### Subtest + `t.Parallel()`
```go
t.Run("case", func(t *testing.T) {
    t.Parallel()  // параллельные subtests
    ...
})
```

---

## Mock vs Stub vs Fake vs Spy

| Тип | Что делает |
|-----|-----------|
| **Stub** | Возвращает заранее подготовленный ответ. Не проверяет вызовы. |
| **Mock** | Проверяет, что и как вызвали (expectations). Тест падает, если вызов не тот. |
| **Fake** | Реальная, но упрощённая реализация (in-memory Repo) |
| **Spy** | Записывает вызовы, проверка постфактум |

**В Go чаще всего:** интерфейсы + `gomock` / `mockery` / руками написанный stub.

### gomock (⭐ у тебя в BCC)
```go
ctrl := gomock.NewController(t)
defer ctrl.Finish()

mockRepo := mocks.NewMockUserRepo(ctrl)
mockRepo.EXPECT().
    GetByID(gomock.Eq("42")).
    Return(User{ID: "42", Name: "Alice"}, nil).
    Times(1)

svc := NewService(mockRepo)
_, err := svc.Get("42")
```

### Почему gomock
- Автогенерация из интерфейса (`mockgen`)
- Строгие expectations (вызов не того аргумента → fail)
- Подходит для больших проектов

### Testify (альтернатива)
- `assert` — продолжает выполнение после fail
- `require` — останавливает на первом fail
- `mock.Mock` — ручной mock, проще чем gomock
```go
assert.Equal(t, expected, actual)
require.NoError(t, err)
```

---

## Best practices

- **Arrange / Act / Assert** — структура теста
- **AAA через комментарии** или пустые строки — читаемость
- **Один assert на логическую проверку**, не 20 в одном тесте
- **Не мокать то, что не владеешь** (не мокай time.Now, стандартную библиотеку) — обёртки через интерфейсы
- **testdata/** — golden-файлы для больших фикстур
- **Test coverage != качество**, но <50% обычно тревожно
- **Не тестируй private функции напрямую** — только через публичный API
- **Детерминированность** — никаких sleep, random без seed, `time.Now()` без инъекции

---

## Integration-тесты — testcontainers

```go
func TestUserRepoIntegration(t *testing.T) {
    ctx := context.Background()
    req := testcontainers.ContainerRequest{
        Image:        "postgres:15-alpine",
        ExposedPorts: []string{"5432/tcp"},
        Env:          map[string]string{"POSTGRES_PASSWORD": "test"},
        WaitingFor:   wait.ForLog("database system is ready"),
    }
    pg, _ := testcontainers.GenericContainer(ctx,
        testcontainers.GenericContainerRequest{ContainerRequest: req, Started: true})
    defer pg.Terminate(ctx)
    // ... connect, migrate, test
}
```

---

## Тестирование concurrency

- **`go test -race`** — обязательно в CI
- Для timing-чувствительных тестов — `golang.org/x/sync/errgroup` и каналы, не sleep
- **eventually-patterns:** `assert.Eventually(t, condFunc, 1*time.Second, 10*time.Millisecond)`

---

## Benchmark-тесты

```go
func BenchmarkParseORA(b *testing.B) {
    for i := 0; i < b.N; i++ {
        _, _ = ParseORA("ORA-12541")
    }
}
```
`go test -bench=. -benchmem` — память + ns/op.

---

## Fuzz-тесты (Go 1.18+)

```go
func FuzzParseORA(f *testing.F) {
    f.Add("ORA-12541")
    f.Fuzz(func(t *testing.T, s string) {
        _, _ = ParseORA(s)  // не должно паниковать
    })
}
```

---

## Типичные вопросы

**1. Чем unit отличается от integration?**
Unit — чистая функция/метод, зависимости замоканы. Integration — поднимает реальные зависимости (БД, брокер).

**2. Что такое coverage и его ловушка?**
% строк, которые выполнились в тестах. Ловушка: 100% coverage не гарантирует правильность — можно покрыть все строки, но не проверить правильные инварианты.

**3. Какие моки использовал?**
gomock с mockgen в BCC Business — 10+ сервисных интерфейсов. Для простых случаев — рукописные stubs. Redis/Postgres — testcontainers в integration-тестах.

**4. Как тестировать private функции?**
Не тестировать напрямую. Если private слишком сложный чтобы покрыть через public API — это сигнал, что надо выделить отдельный пакет/тип.

**5. Как тестировать время?**
Инжектим `Clock interface { Now() time.Time }`. В prod — реальная реализация, в тестах — controlled fake с `clock.Advance(1*time.Hour)`.

**6. Flaky tests — как бороться?**
1. Диагностика: запустить `go test -count=100 -race`
2. Причины: sleep'ы, random без seed, race conditions, shared state между тестами
3. Фикс: eventually-паттерны, уникальные ресурсы на тест, `t.Cleanup()` для тир-дауна

**7. Assert vs Require (testify)?**
`require` останавливает тест — для предусловий (nil-checks). `assert` продолжает — для проверки значений, чтобы видеть все fail'ы за один run.
