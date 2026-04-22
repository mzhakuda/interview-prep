---
tags:
  - kolesa
  - interview
  - core-team
topic: Kolesa Interview
priority: critical
status: ready
date: 2026-04-22
---

# Kolesa Interview — Core Team (1 этап)

**Когда:** 22.04 (ср) 14:30
**Формат:** Google Meet (meet.google.com/rfn-rpdr-hix)
**С кем:** HR Саша + тимлид core команды
**Что:** только прямые тех и кейсовые вопросы (лайфкодинга нет)

---

## Что готовим (по рекомендациям рекрутера)

1. **Один сильный кейс** с метриками — не разделять на команду, всё на "я"
2. **Soft:** вовлечённость, ответственность — [[02 - Behavioral & Ownership]]
3. **Hard:**
   - [[03 - OOP & Design Patterns]]
   - [[04 - Testing]]
   - [[05 - Databases]]
   - [[06 - Linux & Infrastructure]]
   - [[07 - Git CI-CD]]
   - [[08 - Security]]
   - [[09 - Go General]]

---

## Главный кейс (speak-ready)

→ [[01 - Main Case Study]] — PDF streaming, 78% памяти, 5000+ клиентов

Backup кейсы (если попросят другой):
- BPMN → Temporal миграция (reliability)
- Boilerplate для 40+ проектов (leadership/initiative)

Все детали: [[Go Interview Prep/08 - Behavioral/08.1 - STAR Stories]]

---

## Чек-лист за час до созвона

- [ ] Вода, тихая комната, камера
- [ ] Прочитать [[01 - Main Case Study]] один раз вслух
- [ ] Пробежать [[Go Interview Prep/00 - Cheat Sheet]] за 15 минут
- [ ] Открыть [[02 - Behavioral & Ownership]] — перечитать опорные фразы
- [ ] Заготовить 2 вопроса интервьюеру (см. ниже)

---

## Вопросы интервьюеру

1. Какой стек у core-команды? Что считается core — инфраструктурные сервисы или продуктовые?
2. Как выглядит on-call / incident management в Kolesa?
3. Какие главные технические challenges у core-команды сейчас?
4. Сколько этапов дальше и какой фокус у следующего?

---

## Kolesa — контекст (чтобы было понимание)

- Kolesa Group — один из крупнейших classifieds в CA: kolesa.kz (авто), krisha.kz (недвижимость), market.kz
- Стек (публично известно): Go, Python, PHP, PostgreSQL, Kafka, Redis, Kubernetes
- Core-команда скорее всего = платформенная: инфра, shared libs, observability, общие сервисы
- Аналог: платформенные команды в Avito/Ozon
