## Задача 2: Bike Rental System

> **Тип:** IoT + Geospatial + Real-time Tracking
> **Сложность:** Medium-Hard
> **Статус:** TODO

---

## Условие

Город запускает сервис аренды велосипедов по подписке.

**Масштаб:** 10M users, 20M bikes, 200K stations

**Функциональные требования:**
1. Просмотр ближайших станций
2. Отстыковать байк и начать поездку
3. Байк шлет GPS координаты каждые 20 сек
4. При возврате — summary (дистанция, маршрут, время)

**Типы байков:**
- Электрический — цена по дистанции
- Механический — бесплатно

**Два таска:**
- Таск 1: System Architecture — компоненты, как данные текут
- Таск 2: Domain Class Diagram — ООП модель (не схема БД!)

---

## Новые концепции (по сравнению с задачей 1)

```
IoT / GPS:     байк шлет данные каждые 20 сек → не REST API
               Протокол: MQTT (легкий IoT протокол)
               байк → MQTT broker → Kafka

Geospatial:    "ближайшие станции" → не обычный SELECT
               Нужны гео-индексы: PostGIS или Redis Geo

Time-series:   GPS координаты → много мелких записей по времени
               InfluxDB / TimescaleDB

Domain Model:  Class diagram ≠ схема БД
               Наследование: ElectricBike extends Bike
               Интерфейсы, методы, связи между классами
```

---

## Инсайты из реального интервью

```
Q: "Как байк общается с Kafka?"
A: MQTT протокол → MQTT broker → Kafka
   (не SSH, не REST — MQTT специально для IoT)

Q: "Как решить cascade failure?"
A: Event-driven + timeouts + circuit breaker

Ловушка: Domain class diagram ≠ схема БД
  В class diagram: наследование, интерфейсы, методы
  В БД: таблицы, foreign keys, индексы
```

---

_Разбор будет добавлен в процессе работы над задачей_
