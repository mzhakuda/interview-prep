# Cowork Prompt — Obsidian Interview Prep Vault
> Использует kepano/obsidian-skills для нативной работы с vault

---

## ШАГ 0 — Установи obsidian-skills ДО запуска промпта

Это обязательно. Скилы учат Cowork правильному синтаксису Obsidian.
Без них callouts, wikilinks и frontmatter могут быть написаны неправильно.

```bash
# В корне твоего Obsidian vault выполни:
git clone https://github.com/kepano/obsidian-skills.git .claude/obsidian-skills

# Или через npx:
npx skills add git@github.com:kepano/obsidian-skills.git
```

После клонирования структура vault должна выглядеть так:
```
your-vault/
├── .claude/
│   └── obsidian-skills/
│       ├── skills/
│       │   ├── obsidian-markdown/SKILL.md  ← callouts, wikilinks, frontmatter
│       │   ├── obsidian-bases/SKILL.md     ← .base файлы для dashboard
│       │   ├── json-canvas/SKILL.md        ← .canvas для mind map
│       │   ├── obsidian-cli/SKILL.md       ← CLI управление vault
│       │   └── defuddle/SKILL.md           ← извлечение контента с веба
│       └── README.md
└── go-senior-interview-prep.md  ← research файл из чата
```

Убедись что оба файла лежат в vault:
- `go-senior-interview-prep.md` (research файл)
- `cowork-prompt.md` (этот файл)

---

## ПРОМПТ (вставлять целиком в Cowork)

---

Перед тем как начать — прочитай SKILL.md файлы из `.claude/obsidian-skills/skills/`:
- `obsidian-markdown/SKILL.md` — обязательно, будешь писать .md файлы
- `obsidian-bases/SKILL.md` — для dashboard в формате .base
- `json-canvas/SKILL.md` — для mind map в формате .canvas
- `defuddle/SKILL.md` — для извлечения контента с веб-страниц

Эти скилы описывают точный синтаксис форматов Obsidian. Следуй им строго.

---

Ты строишь Obsidian vault для подготовки к Senior Golang Developer интервью.

Прочитай `go-senior-interview-prep.md` — там весь research: темы, вопросы, ссылки.

Также учти инсайды с реальных интервью в этом банке:
- Вопросы про опыт и кейсы с работы (STAR формат)
- Go теория — глубоко (GMP, channels, sync примитивы)
- БД — PostgreSQL (MVCC, индексы, locking)
- System Design кейсы
- **Networking**: полный путь от ввода `youtube.com` до загрузки страницы, OSI модель, IPv4 vs IPv6 (вес в байтах и почему именно так)
- **Архитектуры**: Clean Architecture, Hexagonal — плюсы/минусы обеих
- **Live coding**: конкурентность на Go + написать REST сервер/клиент с нуля на `net/http`

---

### ФАЗА 1 — Создай структуру папок

```
Go Interview Prep/
├── 00 - Dashboard.base          ← живой трекер прогресса (Obsidian Bases)
├── 00 - Topic Map.canvas        ← визуальная карта тем (JSON Canvas)
├── 01 - Go Concurrency/
│   ├── 01.1 - GMP Model & Scheduler.md
│   ├── 01.2 - Channels & Select.md
│   ├── 01.3 - Sync Primitives.md
│   ├── 01.4 - Patterns (WorkerPool, FanIn, Pipeline).md
│   └── 01.5 - Tricky Questions & Traps.md
├── 02 - PostgreSQL/
│   ├── 02.1 - MVCC & Transactions.md
│   ├── 02.2 - Indexes Deep Dive.md
│   ├── 02.3 - Locking & Isolation Levels.md
│   └── 02.4 - Performance & Connection Pooling.md
├── 03 - System Design/
│   ├── 03.1 - API Gateway Patterns.md
│   ├── 03.2 - Circuit Breaker & Resilience.md
│   ├── 03.3 - Distributed Patterns.md
│   └── 03.4 - Rate Limiting.md
├── 04 - Networking/
│   ├── 04.1 - URL to Page (Full Journey).md
│   └── 04.2 - OSI Model & IP.md
├── 05 - Architecture/
│   ├── 05.1 - Clean Architecture.md
│   └── 05.2 - Hexagonal Architecture.md
├── 06 - Live Coding/
│   ├── 06.1 - Concurrency Tasks.md
│   └── 06.2 - HTTP Server & REST Client.md
├── 07 - Core Go/
│   ├── 07.1 - GC & Memory Model.md
│   └── 07.2 - Interfaces & Type System.md
└── 08 - Behavioral/
    └── 08.1 - STAR Stories.md
```

---

### ФАЗА 2 — Собери вопросы с веба через defuddle

Используй defuddle для извлечения чистого markdown с каждой страницы.
Defuddle убирает навигацию, рекламу — остаётся только контент.

**Go вопросы:**
```
defuddle https://www.secondtalent.com/interview-guide/golang/
defuddle https://dsysd-dev.medium.com/20-advanced-questions-asked-for-a-senior-developer-position-interview-1a65203e5d5e
defuddle https://medium.com/@abhigyandwivedi/golang-concurrency-interview-questions-and-answers-80c688904471
defuddle https://dev.to/crusty0gphr/tricky-golang-interview-questions-part-4-concurrent-consumption-34oe
```

**PostgreSQL:**
```
defuddle https://www.secondtalent.com/interview-guide/postgresql/
```

**System Design:**
```
defuddle https://dev.to/kumarkalyan/17-best-free-github-repositories-to-crack-system-design-interviews-h5p
```

Для каждого URL: извлеки контент → отбери вопросы уровня Senior → распредели по нужным файлам.

Критерии отбора:
- Уровень Senior (не "что такое горутина")
- Специфичны для тем из JD: API Gateway, observability, рефакторинг legacy
- Tricky / неочевидные — где легко ошибиться на интервью
- Практические — можно написать код или нарисовать схему
- Минимум 10, максимум 20 вопросов на файл

---

### ФАЗА 3 — Формат каждого .md файла

Читай `obsidian-markdown/SKILL.md` для точного синтаксиса callouts и frontmatter.
Используй этот шаблон для каждого файла:

```markdown
---
tags:
  - go
  - concurrency
  - senior
topic: GMP Model
priority: critical
status: not-started
source: secondtalent.com
---

# GMP Model & Scheduler

Связанные темы: [[01.2 - Channels & Select]] | [[01.3 - Sync Primitives]]

---

> [!question]- 1. Объясни GMP модель Go шедулера. Что происходит когда горутина блокируется на syscall?
>
> **G** (Goroutine) — лёгкий поток, управляется Go runtime, начальный стек ~2KB
> **M** (Machine) — реальный OS поток
> **P** (Processor) — логический процессор, содержит очередь runnable goroutines. Кол-во = GOMAXPROCS
>
> Когда G блокируется на syscall:
> 1. M отсоединяется от P
> 2. P подхватывает другой свободный M (или создаёт новый)
> 3. Заблокированная G остаётся на M без P
> 4. После разблокировки G ищет свободный P, иначе уходит в глобальную очередь

> [!hint]- Подсказка если завис
> Подумай: что произойдёт с другими горутинами пока одна ждёт I/O?
> Ключ: P не простаивает, он переходит к другому M

> [!example]- Код
> ```go
> runtime.GOMAXPROCS(4)               // 4 P = 4 параллельных потока
> fmt.Println(runtime.NumGoroutine()) // кол-во живых горутин
> ```

> [!warning]- Частая ошибка на интервью
> Путают "параллелизм" и "конкурентность". Go гарантирует конкурентность,
> параллелизм — только если GOMAXPROCS > 1 и есть несколько ядер.

---

> [!question]- 2. Следующий вопрос...
```

**Правила frontmatter** (строго):
- `status`: `not-started` | `in-progress` | `done`
- `priority`: `critical` | `important` | `normal`
- `topic`: короткое название темы
- `tags`: список тегов

**Правила wikilinks** (из obsidian-markdown SKILL):
- Каждый файл должен иметь секцию "Связанные темы" с wikilinks на смежные файлы
- Формат: `[[имя файла без расширения]]`

---

### ФАЗА 4 — Специальные файлы

#### 04.1 - URL to Page (Full Journey).md
Это реальный вопрос с интервью — сделай максимально детально:

```markdown
> [!question]- Опиши полный путь от ввода youtube.com до отображения страницы. На каком уровне OSI работает каждый шаг?
>
> 1. **Browser cache** — есть ли IP в кэше браузера?
> 2. **DNS Resolution** (Layer 7 → Layer 3)
>    - OS cache → /etc/hosts → Recursive resolver → Root NS → .com TLD → Authoritative NS
>    - Результат: IP адрес youtube.com
> 3. **TCP Handshake** (Layer 4 — Transport)
>    - SYN → SYN-ACK → ACK
> 4. **TLS Handshake** (Layer 6 — Presentation)
>    - ClientHello → ServerHello → Certificate → Key Exchange → Finished
> 5. **HTTP/2 Request** (Layer 7 — Application)
>    - GET / HTTP/2, заголовки Host/Accept-Encoding/Cookie
> 6. **Load Balancer** → роутинг на backend
> 7. **API Gateway** → auth, rate limiting, routing к микросервисам
> 8. **Server response** → HTML, статус 200
> 9. **Browser rendering**
>    - HTML parse → DOM → CSS parse → CSSOM → Render Tree → Layout → Paint

> [!question]- IPv4 vs IPv6: сколько байт весит каждый, почему и как ты это понял?
>
> **IPv4**: 32 бита = **4 байта**
> - Формат: 4 октета по 8 бит: `192.168.1.1`
> - 2³² ≈ 4.3 млрд адресов — исчерпано в 2011 году
>
> **IPv6**: 128 бит = **16 байт** (в 4 раза больше)
> - Формат: 8 групп по 16 бит: `2001:0db8::8a2e:0370:7334`
> - 2¹²⁸ ≈ 340 ундециллионов адресов
>
> **Мнемоника**: IPv4 = 4 числа × 8 бит = 32 бит = 4 байта.
> IPv6 = 4× больше групп × 2× шире = 128 бит = 16 байт.
```

#### 06.1 - Concurrency Tasks.md
Задачи для live coding — без решения сначала, решение скрыто:

```markdown
> [!question]- ЗАДАЧА: Реализуй Worker Pool. N воркеров, канал задач, канал результатов, graceful shutdown через context.

> [!warning]- Типичные ошибки
> - WaitGroup.Add() после go func() — race condition
> - Не обработать ctx.Done() внутри воркера — утечка горутин
> - close(results) до wg.Wait() — panic: send on closed channel

> [!example]- Решение
> ```go
> func WorkerPool(ctx context.Context, numWorkers int, jobs <-chan int) <-chan int {
>     results := make(chan int, numWorkers)
>     var wg sync.WaitGroup
>
>     for i := 0; i < numWorkers; i++ {
>         wg.Add(1)
>         go func() {
>             defer wg.Done()
>             for {
>                 select {
>                 case job, ok := <-jobs:
>                     if !ok {
>                         return
>                     }
>                     results <- job * job
>                 case <-ctx.Done():
>                     return
>                 }
>             }
>         }()
>     }
>
>     go func() {
>         wg.Wait()
>         close(results)
>     }()
>
>     return results
> }
> ```

---

> [!question]- ЗАДАЧА: Напиши HTTP сервер на net/http с logging middleware, GET /health и POST /echo

> [!example]- Решение
> ```go
> func loggingMiddleware(next http.Handler) http.Handler {
>     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
>         start := time.Now()
>         next.ServeHTTP(w, r)
>         log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
>     })
> }
>
> func main() {
>     mux := http.NewServeMux()
>     mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
>         json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
>     })
>     mux.HandleFunc("POST /echo", func(w http.ResponseWriter, r *http.Request) {
>         body, _ := io.ReadAll(r.Body)
>         w.Header().Set("Content-Type", "application/json")
>         w.Write(body)
>     })
>     http.ListenAndServe(":8080", loggingMiddleware(mux))
> }
> ```
```

#### 08.1 - STAR Stories.md

```markdown
---
tags: [behavioral, star]
topic: Behavioral
priority: important
status: not-started
---

# Behavioral — STAR Stories

Интервьюер: "расскажи кейс с работы и как решал"
Подготовь минимум 3 истории заранее по STAR фреймворку.

Связанные темы: [[00 - Dashboard]]

---

> [!question]- Расскажи о самом сложном техническом решении которое ты принимал

> [!hint]- STAR Framework
> - **S**ituation — контекст, что за система, масштаб, команда
> - **T**ask — что конкретно нужно было решить
> - **A**ction — что ты лично сделал (технически, детально)
> - **R**esult — измеримый результат (latency, uptime, delivery time...)

**Шаблон для заполнения:**

| # | Тема | Situation | Task | Action | Result |
|---|------|-----------|------|--------|--------|
| 1 | Рефакторинг / легаси | | | | |
| 2 | Production инцидент | | | | |
| 3 | Технический выбор / trade-off | | | | |
| 4 | Оптимизация производительности | | | | |
| 5 | Несогласие с командой | | | | |

**Подсказка по материалу из BCC:**
- Temporal workflows — что именно строил, какие edge cases?
- gRPC сервисы — какие trade-offs принимал vs REST?
- Kubernetes — что деплоил, что ломалось в проде?
- RabbitMQ — какие паттерны использовал, были ли проблемы с durability?
- Elasticsearch — как строил индексы, как оптимизировал запросы?
```

---

### ФАЗА 5 — Dashboard (00 - Dashboard.base)

Читай `obsidian-bases/SKILL.md` для точного синтаксиса `.base` файла.

Создай `00 - Dashboard.base` с несколькими view:
- **Все темы** — table view, все файлы vault, столбцы: topic, priority, status
- **Не начато** — filter: status = "not-started", sort: priority desc
- **В процессе** — filter: status = "in-progress"
- **Готово** — filter: status = "done"

Используй точный синтаксис из SKILL.md — не угадывай формат.

---

### ФАЗА 6 — Mind Map (00 - Topic Map.canvas)

Читай `json-canvas/SKILL.md` для точного синтаксиса.

Создай `00 - Topic Map.canvas`:
- Центральный text node: "Senior Go Interview"
- 8 file nodes вокруг — по одному первому файлу из каждой папки
- Цвет edges по приоритету: красный = critical (01, 02, 06), жёлтый = important, зелёный = normal
- Edges между связанными темами: 01→06 (Concurrency→Live Coding), 03→04 (SD→Networking)

---

### ФИНАЛЬНЫЕ ИНСТРУКЦИИ

1. **Читай SKILL.md файлы первыми** — до создания любых файлов
2. **Используй defuddle для всех URL** — не пытайся парсить страницы вручную
3. **Wikilinks обязательны** — каждый файл должен ссылаться на смежные темы через `[[...]]`
4. **Frontmatter на каждом файле** — status, priority, topic, tags
5. **Не дублируй** — вопрос из research файла → обогати кодом/hint/warning, не копируй
6. **Senior уровень** — пропускай базовые вопросы
7. **После каждой папки** — сообщи сколько вопросов добавил и откуда
8. **В конце** — сообщи общую статистику: файлов создано, вопросов всего

**Порядок выполнения:**
1. Прочитать все SKILL.md файлы
2. Прочитать `go-senior-interview-prep.md`
3. Создать структуру папок
4. Запустить defuddle по всем URL, собрать вопросы
5. Создать все .md файлы
6. Создать `00 - Dashboard.base`
7. Создать `00 - Topic Map.canvas`

Перед началом каждой фазы кратко скажи что собираешься делать.
