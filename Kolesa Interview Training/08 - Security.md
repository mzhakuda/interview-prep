---
tags:
  - kolesa
  - security
  - owasp
topic: Security
priority: high
status: ready
---

# Безопасность

> Банковский опыт тут сильное преимущество — упомянуть escrow, RSA digital signing, OTP, compliance.

---

## OWASP Top 10 — must know

### 1. Broken Access Control
Юзер A видит ресурсы юзера B.
- **IDOR** — `/orders/123` без проверки что order принадлежит юзеру
- **Vertical privilege escalation** — user делает admin-actions
- **Защита:** всегда проверять `resource.owner_id == user.id` на каждом endpoint

### 2. Cryptographic Failures
- Пароли в plaintext / MD5 / SHA1
- Секреты в git
- TLS 1.0/1.1, слабые ciphers
- **Защита:** bcrypt/argon2 для паролей, TLS 1.2+ минимум, HSTS

### 3. Injection
- **SQL injection** — `"SELECT * WHERE id=" + userInput` → используй prepared statements (`$1, $2`)
- **Command injection** — `exec.Command("sh", "-c", userInput)` — никогда
- **LDAP / NoSQL injection** — та же история

### 4. Insecure Design
Архитектурный изъян, например отсутствие rate-limiting на login → brute force.

### 5. Security Misconfiguration
- Default passwords
- Открытые admin-panels
- Debug mode в prod
- Listing directories в веб-сервере

### 6. Vulnerable Components
- Устаревшие библиотеки с CVE
- `govulncheck ./...` / `go list -m -u all`
- Dependabot / Renovate

### 7. Authentication Failures
- Brute force (нет rate limit)
- Сессии без expiration
- Password в URL / log

### 8. Software & Data Integrity
- Deserialize untrusted data
- Supply chain attacks (вкоммиченный зловред в зависимости)
- **SLSA**, signed artifacts, SBOM

### 9. Security Logging & Monitoring Failures
- Нет логов auth событий
- Нет алертов на аномалии

### 10. SSRF (Server-Side Request Forgery)
Юзер контролирует URL, куда сервер делает запрос → ходит в internal services (metadata-api облака, etc).
- Allowlist хостов
- Не делать HTTP в localhost / RFC1918

---

## Authentication vs Authorization

- **Authentication** — кто ты (логин/пароль, OTP, certificate)
- **Authorization** — что тебе можно (RBAC, ACL)

### OAuth2 / OIDC / JWT

**OAuth2** — фреймворк авторизации (даёт access_token):
- Client Credentials — server-to-server
- Authorization Code + PKCE — user через браузер (главный для web/mobile)
- Refresh Token — обновление без логина

**OIDC** — OAuth2 + **authentication** (id_token)

**JWT** — формат токена (header.payload.signature):
- Stateless (проверка без похода в auth-сервер)
- Но нельзя "отозвать" — только через короткий TTL + refresh
- Хранить claims, **НЕ секреты** (base64 легко декодируется!)

**Keycloak** — ⭐ твой опыт в Butterfly Effect. Identity provider, поддерживает OAuth2/OIDC/SAML.

---

## Защитные паттерны

### Rate limiting
- **Token bucket** / leaky bucket
- Redis-based: `INCR key / EXPIRE` или RedisCell
- Лимиты по IP / user / API key
- 429 Too Many Requests + `Retry-After`

### CSRF
- Атакующий сайт форсит запрос к твоему (куки автоматически шлются)
- **Защита:** CSRF-token в форме + проверка на backend; SameSite=Strict cookie; CORS

### CORS
- Браузерная защита от cross-origin запросов
- **Preflight (OPTIONS)** для не-simple методов
- `Access-Control-Allow-Origin: https://trusted.com` (не `*` для auth endpoint'ов!)

### XSS (Cross-Site Scripting)
- Вставка злого `<script>` в HTML
- **Защита:** encoding при рендере, `Content-Security-Policy` header, не `innerHTML` c user-данными

### TLS / mTLS
- **TLS** — сервер доказывает свою идентичность клиенту (сертификат)
- **mTLS** — обе стороны. Часто для service-to-service в k8s-cluster (Istio/Linkerd)

---

## Банковская специфика (⭐ твой опыт BCC)

- **RSA digital signing** — подписание документов приватным ключом клиента, верификация публичным
- **OTP** — one-time password (TOTP по RFC 6238 / HOTP)
- **Двухфакторная аутентификация**
- **Audit trail** — неизменяемый лог кто/когда/что сделал (compliance, ЦБ РФ / КЗ)
- **Idempotency** для платежей — `Idempotency-Key` header, дедуп на стороне сервера
- **Escrow** — деньги блокируются до выполнения условий, безопасные расчёты между сторонами
- **PCI DSS** — стандарт обработки карточных данных (если сталкивался — упомянуть tokenization)

---

## Secrets management

- **Vault (HashiCorp)** — dynamic secrets, rotation, audit
- **Cloud-native:** AWS Secrets Manager, GCP Secret Manager
- **k8s Secret** — НЕ достаточно (base64 в etcd), нужен encryption-at-rest или external secret

---

## Безопасность Go-кода

- `crypto/rand` для случайных байт (не `math/rand`!)
- `bcrypt` для паролей (`golang.org/x/crypto/bcrypt`)
- **Constant-time comparison** `hmac.Equal` / `subtle.ConstantTimeCompare` — против timing-атак
- `html/template` вместо `text/template` для HTML-рендера (auto-escape)
- Context deadlines — защита от slow-loris
- `net/http.Server` с таймаутами (`ReadTimeout`, `WriteTimeout`, `IdleTimeout`)

### Пример — правильная проверка токена
```go
// НЕПРАВИЛЬНО (timing attack)
if token == expected { ... }

// ПРАВИЛЬНО
if subtle.ConstantTimeCompare([]byte(token), []byte(expected)) == 1 { ... }
```

---

## Типичные вопросы

**1. Чем отличается аутентификация от авторизации?**
Authentication — proof identity. Authorization — policy check "что тебе можно".

**2. Stateful vs stateless session?**
- **Stateful:** session_id → lookup в Redis/БД. Можно инвалидировать. Нужен shared storage.
- **Stateless (JWT):** токен сам содержит claims. Не надо lookup. Но отозвать нельзя (только через короткий TTL + refresh + блеклист коротких).

**3. Где JWT хранить на фронте?**
- **Access token** — в memory (через closure / React state). **НЕ в localStorage** (XSS).
- **Refresh token** — `httpOnly; Secure; SameSite=Strict` cookie.

**4. Как защитить API?**
1. HTTPS (TLS 1.2+)
2. Auth (OAuth2/JWT/API key)
3. Authorization (RBAC на каждом endpoint)
4. Rate limiting
5. Input validation (strict schemas, OpenAPI)
6. Output encoding (XSS)
7. Secrets в Vault, не в коде
8. Audit logs

**5. SQL injection — как победить?**
Prepared statements (`$1, $2, $3`). НЕ конкатенация строк. ORM тоже по сути prepared statements. + input validation.

**6. CORS зачем?**
Браузерная защита: скрипт с `evil.com` не может делать авторизованные запросы на `bank.com`. Сервер контролирует, какие origins разрешены.

**7. Что такое SSRF и как защищаться?**
Юзер подаёт URL, сервер идёт по нему и возможно достаёт internal resources (cloud metadata-endpoint `169.254.169.254`, internal DBs). Защита: allowlist hosts, блок RFC1918, проверка резолва DNS.

**8. Как сделать password reset безопасно?**
1. User вводит email → всегда отвечаем "если email есть, мы отправим" (нет user enumeration)
2. Генерим **cryptographic-secure random** токен (32+ байт)
3. Храним в БД **хеш токена** (не сам токен), TTL 15 минут, one-time-use
4. Ссылка с токеном по HTTPS
5. После сброса инвалидируем все сессии

**9. Что делать если секрет закоммичен в git?**
1. **Rotate** сразу — старый считать скомпрометированным
2. Удалить из истории: `git filter-repo` (не решает утечку, но убирает из свежих клонов)
3. Форс-пуш главной ветки (осторожно!)
4. Уведомить команду

**10. HTTPS — что гарантирует?**
- Confidentiality (шифрование)
- Integrity (MAC)
- Authentication server'а (certificate, но не клиента — это mTLS)
- НЕ защищает от XSS, SQLi, broken auth.
