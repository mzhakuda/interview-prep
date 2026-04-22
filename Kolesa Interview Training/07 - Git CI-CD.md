---
tags:
  - kolesa
  - git
  - cicd
topic: Git & CI/CD
priority: high
status: ready
---

# Git & CI/CD

---

## Git — базовое

### Команды, без которых никак
```bash
git status
git log --oneline --graph --all
git diff HEAD~1 HEAD
git stash / git stash pop
git restore <file>                  # отменить изменения в файле
git restore --staged <file>         # убрать из staged
```

### Ветки
```bash
git switch -c feature/x             # создать и перейти
git branch -d feature/x             # удалить ветку
git branch -D feature/x             # force-удалить
```

### Merge vs Rebase

**Merge:**
- Сохраняет историю как есть, создаёт merge-commit
- Безопасен (не переписывает историю)
- Но история "кудрявая"

**Rebase:**
- Переписывает коммиты поверх другой ветки — линейная история
- Не делать rebase на **публичных** ветках (main) — ломает всем другим
- `git rebase -i` — интерактивный: squash / fixup / reword / drop

**Правило:** feature-ветки rebase → main / main merge в feature. Или squash-merge в PR.

### Работа с конфликтами
```bash
git merge feature  # конфликт
# редактируем файлы, маркеры <<<<<<<
git add <file>
git merge --continue
# или git merge --abort
```

### Полезное
- `git cherry-pick <sha>` — взять один коммит из другой ветки
- `git revert <sha>` — создать обратный коммит (не переписывает историю)
- `git reset --soft HEAD~1` — отменить последний коммит, оставив изменения в staged
- `git reset --hard HEAD~1` — ⚠️ удалить последний коммит + изменения
- `git reflog` — история всех HEAD (даже после reset/rebase можно восстановить)
- `git bisect` — бинарный поиск коммита с багом

### `.gitignore`
- Не класть в репо: secrets (`.env`), сгенерированные файлы, `node_modules`, `vendor` (если не принято)
- `git rm --cached file` — убрать уже закоммиченный файл + не трекать

### Commits — best practices
- **Conventional commits:** `feat:`, `fix:`, `refactor:`, `test:`, `chore:`, `docs:`
- Атомарные коммиты — один коммит = одна логическая единица
- Сообщение: что и **зачем**, не "обновил код"
- Squash перед merge в main для чистоты

---

## Git Flow стратегии

### Trunk-based (⭐ современный подход)
- Одна main-ветка
- Короткоживущие feature-branches (<1-2 дня)
- Часто мерджим в main, feature flags для незаконченного
- Деплой main = prod

### GitFlow (классика, тяжелее)
- `main`, `develop`, `feature/*`, `release/*`, `hotfix/*`
- Нужен при редких релизах

### GitHub Flow
- `main` + короткие feature-branches → PR → merge
- Деплой после merge

---

## CI/CD — концепции

### CI (Continuous Integration)
Каждый push — автозапуск:
1. **Lint** (`golangci-lint run`)
2. **Build**
3. **Test** (`go test ./... -race -cover`)
4. **Security scan** (gosec, Trivy)
5. **Build Docker image**

### CD (Continuous Delivery / Deployment)
- **Delivery** — артефакты готовы, деплой по кнопке
- **Deployment** — автодеплой в prod после merge

### Стратегии деплоя

| Стратегия | Суть | Риск |
|-----------|------|------|
| **Rolling** | по одному поду меняем | быстрый rollback, но mixed versions |
| **Blue-Green** | параллельный prod, переключение роутинга | нужно x2 инфры |
| **Canary** | 5% → 25% → 100%, мониторим метрики | нужны хорошие метрики |
| **Feature flags** | код задеплоен, фича выключена | независимо от деплоя |

### Rollback
- Helm: `helm rollback <release> <revision>`
- kubectl: `kubectl rollout undo deployment/X`
- Git-driven (ArgoCD): revert-коммит и он разворачивает

---

## GitLab CI (твой опыт в BCC)

```yaml
stages: [lint, test, build, deploy]

variables:
  GO_VERSION: "1.22"

lint:
  stage: lint
  image: golangci/golangci-lint:latest
  script: golangci-lint run

test:
  stage: test
  image: golang:$GO_VERSION
  services: [postgres:15]
  script:
    - go test ./... -race -coverprofile=cover.out
    - go tool cover -func=cover.out

build:
  stage: build
  image: docker:latest
  services: [docker:dind]
  script:
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA

deploy-prod:
  stage: deploy
  only: [main]
  when: manual  # ручное подтверждение для prod
  script: helm upgrade app ./chart --set image.tag=$CI_COMMIT_SHA
```

---

## Code Review

- **Atomic PRs** — одна задача = один PR, <400 строк
- **Description:** что, зачем, как тестировалось
- **Feature flags** чтобы мерджить рано
- Ревью нужно проходить быстро: тоо long PR → тёмная бездна

---

## Типичные вопросы

**1. Rebase vs Merge — что ты используешь?**
В feature-branches делаю rebase on main, чтобы подтянуть последнее. В main — только merge или squash-merge. Правило: никогда не rebase-ить публичные ветки.

**2. Git-flow vs Trunk-based?**
Зависит от релизной политики. Если деплоим часто (каждый день/несколько раз в день) — trunk-based с feature flags. Если редко и есть фиксированные релизы — GitFlow. В BCC у нас был trunk-like: main → prod, feature branches <2 дней, feature-flags для выключения.

**3. Что такое semantic versioning?**
`MAJOR.MINOR.PATCH` — breaking.feature.fix. 2.3.1 → 2.3.2 багфикс, 2.4.0 новая фича backward-compat, 3.0.0 breaking.

**4. Как обеспечить что тесты зелёные в main?**
- Protected branches — merge только через PR
- Required checks — CI должен пройти
- Required review — как минимум 1 approve

**5. Обнаружил регрессию, как найти коммит?**
`git bisect start → git bisect bad HEAD → git bisect good <known-good-sha>` — бинпоиск.

**6. Долго стартующий CI — что делать?**
- Docker layer caching
- Кешировать `go mod cache` между запусками
- Параллелить этапы (lint и test одновременно)
- Матрицы — разбить тесты по пакетам
- Для lint+test не пересобирать — использовать `golang` base image

**7. Как управлять secrets в CI?**
- GitLab variables (masked + protected для prod)
- Vault-integration через jwt
- Никогда не echo-ать их в логи
- Rotation политика

**8. Отличие CI от CD?**
CI (Continuous Integration) — интеграция кода: build, test на каждый push. CD (Continuous Delivery) — артефакты готовы к деплою в любой момент. Continuous Deployment — автодеплой без ручного шага.

**9. Что такое blue-green деплой?**
Два одинаковых prod-окружения (blue, green). Деплой в неактивное → прогоняем smoke-тесты → переключаем роутинг трафика. Быстрый rollback = переключаем обратно.

**10. Monorepo vs Polyrepo?**
- Mono: атомарные изменения через сервисы, проще shared libs. Но тяжёлые билды, need smart CI.
- Poly: независимость, маленькие builds. Но версионирование shared libs.
- В BCC скорее polyrepo по сервисам + boilerplate шарился.
