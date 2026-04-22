---
tags:
  - kolesa
  - linux
  - infra
  - docker
  - kubernetes
topic: Linux & Infra
priority: high
status: ready
---

# Linux & Инфраструктура

---

## Linux — базовые команды, которые ждут

### Процессы / система
- `ps aux | grep X` — список процессов
- `top` / `htop` — реалтайм мониторинг CPU/RAM
- `kill -9 PID` — SIGKILL, `kill -15` — SIGTERM (graceful)
- `strace -p PID` — трейс syscall'ов процесса
- `lsof -p PID` — какие файлы/сокеты открыты
- `ulimit -a` — лимиты процесса (file descriptors, memory)
- `nproc` — сколько CPU
- `free -h` — память
- `df -h` — диск по разделам
- `du -sh /path` — размер директории

### Сеть
- `netstat -tulpn` / `ss -tulpn` — открытые порты
- `curl -v url` — HTTP-запрос с хедерами
- `dig`, `nslookup` — DNS
- `tcpdump -i eth0 port 80` — снифер
- `traceroute` — маршрут пакета
- `iptables -L` — firewall rules

### Файлы / текст
- `grep -r "foo" .` — рекурсивный поиск
- `find . -name "*.go" -mtime -1` — файлы .go модифицированные за последний день
- `awk '{print $2}'` — 2-й столбец
- `sed 's/old/new/g'` — замена
- `tail -f /var/log/app.log` — real-time tail
- `less +F file` — like tail -f но можно скроллить
- `xargs` — из stdin в аргументы команды

### Сигналы в Go (⭐ graceful shutdown)
- `SIGTERM` (15) — дефолт от `docker stop`, `kubectl delete pod`
- `SIGINT` (2) — Ctrl+C
- `SIGKILL` (9) — нельзя перехватить, процесс мгновенно умирает
- `SIGHUP` (1) — обычно "перечитать конфиг"
- `SIGURG` — ⚠️ использует go-runtime для preemption (не перехватывать!)

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
defer stop()
<-ctx.Done()
// graceful shutdown
```

---

## Docker

### Образы
- **Multi-stage build** — must have для Go:
```dockerfile
FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app ./cmd/server

FROM gcr.io/distroless/static:nonroot
COPY --from=build /app /app
USER nonroot
ENTRYPOINT ["/app"]
```
- **CGO_ENABLED=0** — static binary, можно на distroless/scratch
- **-ldflags="-s -w"** — убрать debug info (~20-30% меньше)
- **distroless** — нет shell, минимум CVE

### Что размер образа уменьшает
- Multi-stage
- Alpine / distroless / scratch base
- `.dockerignore` — не класть `.git`, `node_modules`
- Не устанавливать build-tools в финальный образ

### Layers
- Каждая `RUN/COPY/ADD` — новый layer
- Cache инвалидация — если слой изменился, все ниже пересобираются
- Поэтому `COPY go.* + go mod download` ДО `COPY .` — кеш сохранится при изменении кода

### Команды
```bash
docker ps -a               # все контейнеры
docker logs -f <id>        # logs
docker exec -it <id> sh    # shell внутрь
docker inspect <id>        # JSON с метаданными
docker stats               # использование CPU/RAM
docker system prune -a     # очистка всего
```

---

## Kubernetes

### Объекты — must know
- **Pod** — минимальная единица (1+ контейнеров)
- **Deployment** — управляет replica'ми Pod'ов, rolling update
- **StatefulSet** — для stateful нагрузки (БД), stable identity
- **Service** — стабильный endpoint для набора Pod'ов (ClusterIP / NodePort / LoadBalancer)
- **Ingress** — HTTP-роутинг снаружи в Service
- **ConfigMap** — конфиг (env, файлы)
- **Secret** — секреты (base64, не зашифровано в etcd по дефолту!)
- **Job** / **CronJob** — разовые / периодические задачи
- **HPA** — автоскейл по CPU/memory/custom метрикам
- **Namespace** — изоляция ресурсов

### Probes (⭐ важно)
- **livenessProbe** — если падает, под перезапускается
- **readinessProbe** — если падает, под выкидывается из Service (не получает трафик)
- **startupProbe** — для медленно стартующих приложений, отсрочивает liveness

```yaml
readinessProbe:
  httpGet: {path: /healthz, port: 8080}
  initialDelaySeconds: 5
  periodSeconds: 10
```

### Resources / QoS
```yaml
resources:
  requests: {cpu: 100m, memory: 128Mi}  # гарантия
  limits:   {cpu: 500m, memory: 256Mi}  # максимум
```
- **Requests** учитывает scheduler при размещении
- **Limits** — hard cap; превышение memory → OOMKilled
- **QoS классы:** Guaranteed (requests==limits) > Burstable > BestEffort

### Rolling update
- `maxSurge` — сколько сверх replicas во время апдейта
- `maxUnavailable` — сколько может быть недоступно
- Стратегии: `RollingUpdate` (default) / `Recreate`

### Команды
```bash
kubectl get pods -n <ns>
kubectl describe pod <name>
kubectl logs <pod> -c <container> --previous  # логи предыдущего падения
kubectl exec -it <pod> -- sh
kubectl top pod                               # CPU/RAM
kubectl rollout status deployment/<name>
kubectl rollout undo deployment/<name>        # откат
kubectl port-forward pod/<name> 8080:8080
```

### OOMKilled — откуда?
- Память процесса > `limits.memory` → kubelet убивает
- Причины: memory leak, не тот heap-профиль, слишком низкий limit
- Диагностика: `kubectl describe pod` → `Last State: Terminated, Reason: OOMKilled`, потом `pprof`/`profile heap`

---

## Observability (⭐ часто спрашивают)

### 3 столпа
1. **Metrics** — Prometheus (counter / gauge / histogram / summary)
2. **Logs** — structured (JSON), ELK / Loki
3. **Traces** — OpenTelemetry / Jaeger, distributed tracing

### Метрики — что трекать (Golden Signals)
- **Latency** (histogram `_duration_seconds_bucket`)
- **Traffic** (RPS, counter)
- **Errors** (error rate, counter с label `code`)
- **Saturation** (CPU/RAM/connections)

### Trace context propagation (⭐ ты делал!)
В BCC ты написал custom Temporal `ContextPropagator` для `X-Request-ID` — это ровно про distributed tracing. Можно упомянуть.

---

## CI/CD (вкратце — полнее в [[07 - Git CI-CD]])

- **Build** — Docker image в registry
- **Test** — unit + lint (golangci-lint) + race detector
- **Security scan** — Trivy / gosec
- **Deploy** — Helm / ArgoCD / plain kubectl

---

## Типичные вопросы

**1. Разница Docker и VM?**
VM — полная ОС, свой kernel. Docker — контейнер, общий kernel с хостом, изоляция через namespaces + cgroups. Контейнеры **быстрее стартуют и легче по ресурсам**.

**2. Что такое cgroups и namespaces?**
- **namespaces** — изолируют ресурсы (PID, net, mount, user, uts, ipc). Контейнер не видит процессы хоста.
- **cgroups** — ограничивают ресурсы (CPU, RAM, IO).

**3. Граceful shutdown в Go-сервисе на Kubernetes?**
1. Kubernetes шлёт SIGTERM
2. Из readinessProbe начинает выдавать failure → Service перестаёт слать трафик
3. Приложение ловит SIGTERM → отменяет `context` → `http.Server.Shutdown(ctx)` → дожидается активных запросов
4. По `terminationGracePeriodSeconds` (default 30s) — SIGKILL, если не закрылся
Твой graceful shutdown manager в BCC — 5 минут, дожидается HTTP + Temporal + БД.

**4. Readiness vs Liveness?**
- **Liveness fail → pod restart**
- **Readiness fail → pod выкидывается из Service**, но НЕ рестартится
- Оба не должны проверять внешние зависимости (БД) иначе все поды умрут при deg БД

**5. Почему в лимитах указать CPU важно?**
Без CPU limit один под может задавить ноду. Но есть нюанс: CPU throttling через CFS может давать латенси хвосты → для latency-sensitive сервисов иногда ставят только requests.

**6. Как дебажить pod который crash-loop'ит?**
```bash
kubectl describe pod <name>               # смотрим events
kubectl logs <name> --previous            # логи убитого пода
kubectl get events --sort-by=.lastTimestamp
```
Частые причины: OOMKilled, плохой readiness, миссинг config/secret, image pull error.

**7. Как секреты правильно хранить?**
- `Secret` в k8s — base64, **не** шифруется (нужен etcd encryption-at-rest)
- Лучше — **external secret manager** (Vault, AWS SM, GCP SM) + CSI-driver или ExternalSecret
- Никогда в git!

**8. Чем отличаются Deployment и StatefulSet?**
- **Deployment** — pod'ы взаимозаменяемые, random имена, shared storage не нужен
- **StatefulSet** — стабильные имена (`pg-0`, `pg-1`), стабильный storage, порядок создания/удаления
