# Booker - Система бронирования столов

Микросервисная система бронирования столов на Go с использованием gRPC, Kafka, Redis и Postgres.

## Архитектура

- **admin-gateway**: REST API + фронтенд для администраторов
- **venue-svc**: Управление заведениями, залами, столами и расписанием
- **booking-svc**: Ядро бронирований с машиной состояний и Redis holds
- **notify-svc**: Уведомления о событиях бронирований

## Быстрый старт

### Требования

- **Docker** и **Docker Compose**
- Больше ничего не нужно! Все работает в контейнерах.

### Запуск (всё в одном)

```bash
make setup
```

Эта команда запустит всю инфраструктуру, сервисы, выполнит миграции и заполнит тестовыми данными.

### Пошаговый запуск

1. Запустить инфраструктуру и сервисы:
```bash
make up
```

2. Выполнить миграции:
```bash
make migrate
```

3. Заполнить тестовыми данными:
```bash
make seed
```

> **Примечание**: Proto файлы генерируются автоматически при сборке Docker образов. Если нужно сгенерировать локально: `make gen` (тоже через Docker).

### Просмотр логов

```bash
# Все микросервисы (следить в реальном времени)
make logs

# Конкретный сервис:
make logs-gateway    # admin-gateway
make logs-venue     # venue-svc
make logs-booking   # booking-svc
make logs-notify    # notify-svc

# Через docker-compose:
docker-compose --profile apps logs -f admin-gateway
docker-compose --profile apps logs --tail=50 admin-gateway
docker-compose --profile apps logs --timestamps -f venue-svc
```

### Доступ к сервисам

- Admin Gateway: http://localhost:8080
- Grafana: http://localhost:3000 (admin/admin)
- Jaeger: http://localhost:16686
- Prometheus: http://localhost:9090
- Kafka UI: http://localhost:8081

## Структура проекта

```
.
├── cmd/
│   ├── admin-gateway/    # REST API + фронтенд
│   ├── venue-svc/         # Сервис заведений
│   ├── booking-svc/       # Сервис бронирований
│   └── notify-svc/         # Сервис уведомлений
├── pkg/
│   ├── proto/             # Сгенерированные proto файлы
│   ├── kafka/             # Kafka producer/consumer
│   ├── redis/             # Redis клиент
│   └── tracing/           # OpenTelemetry трейсинг
├── proto/                 # Proto определения
├── migrations/            # SQL миграции
└── deploy/                # Конфигурации для инфраструктуры
```

## API

### Создание бронирования

```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "venue_id": "venue-1",
    "table": {
      "venue_id": "venue-1",
      "room_id": "room-1",
      "table_id": "table-1"
    },
    "slot": {
      "date": "2024-01-15",
      "start_time": "19:00",
      "duration_minutes": 120
    },
    "party_size": 4,
    "customer_name": "John Doe",
    "customer_phone": "+1234567890"
  }'
```

## Разработка

### Генерация proto файлов

```bash
make gen
```

### Запуск тестов

```bash
make test
```

### Просмотр логов

```bash
make logs
```

### Остановка

```bash
make down
```

## Лицензия

MIT


