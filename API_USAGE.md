# Как управлять сервисами через API

## Архитектура доступа

**Важно**: `venue-svc` и `booking-svc` - это **внутренние gRPC сервисы**. Они не имеют REST API напрямую.

Все управление идет через **Admin Gateway** (порт 8080), который:
- Предоставляет REST API для клиентов
- Проксирует запросы в gRPC вызовы к внутренним сервисам
- Объединяет функциональность всех сервисов в единый интерфейс

```
Клиент → Admin Gateway (REST) → venue-svc (gRPC)
                              → booking-svc (gRPC)
```

## Доступ к API

**Base URL**: `http://localhost:8080`

**API Info**: `http://localhost:8080/api` - показывает все доступные endpoints

## Аутентификация

Сейчас используется упрощенная аутентификация (для разработки):

```bash
# Получить токен
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin"}'

# Ответ: {"token": "dummy-token"}
```

Для защищенных endpoints используйте заголовок:
```
Authorization: Bearer dummy-token
```

## Управление Venue Service (через Admin Gateway)

### 1. Создать заведение (venue)

```bash
curl -X POST http://localhost:8080/api/v1/venues \
  -H "Authorization: Bearer dummy-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Ресторан 'Уютный'",
    "timezone": "Europe/Moscow",
    "address": "Москва, ул. Примерная, 1"
  }'
```

### 2. Получить список заведений

```bash
curl -X GET "http://localhost:8080/api/v1/venues?limit=10&offset=0" \
  -H "Authorization: Bearer dummy-token"
```

### 3. Создать зал (room) в заведении

```bash
curl -X POST http://localhost:8080/api/v1/venues/{venue_id}/rooms \
  -H "Authorization: Bearer dummy-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Основной зал"
  }'
```

### 4. Создать стол (table) в зале

```bash
curl -X POST http://localhost:8080/api/v1/rooms/{room_id}/tables \
  -H "Authorization: Bearer dummy-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Стол 1",
    "capacity": 4,
    "can_merge": true,
    "zone": "window"
  }'
```

### 5. Проверить доступность столов

```bash
curl -X POST http://localhost:8080/api/v1/availability/check \
  -H "Authorization: Bearer dummy-token" \
  -H "Content-Type: application/json" \
  -d '{
    "venue_id": "venue-1",
    "slot": {
      "date": "2024-01-15",
      "start_time": "19:00",
      "duration_minutes": 120
    },
    "party_size": 4
  }'
```

## Управление Booking Service (через Admin Gateway)

### 1. Создать бронирование

```bash
curl -X POST http://localhost:8080/api/v1/bookings \
  -H "Authorization: Bearer dummy-token" \
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
    "customer_name": "Иван Иванов",
    "customer_phone": "+79001234567",
    "comment": "Столик у окна"
  }'
```

### 2. Получить список бронирований

```bash
curl -X GET "http://localhost:8080/api/v1/bookings?venue_id=venue-1&date=2024-01-15" \
  -H "Authorization: Bearer dummy-token"
```

### 3. Подтвердить бронирование

```bash
curl -X POST http://localhost:8080/api/v1/bookings/{booking_id}/confirm \
  -H "Authorization: Bearer dummy-token"
```

### 4. Отменить бронирование

```bash
curl -X POST http://localhost:8080/api/v1/bookings/{booking_id}/cancel \
  -H "Authorization: Bearer dummy-token" \
  -H "Content-Type: application/json" \
  -d '{
    "reason": "Клиент отменил"
  }'
```

### 5. Отметить гостей как севших

```bash
curl -X POST http://localhost:8080/api/v1/bookings/{booking_id}/seat \
  -H "Authorization: Bearer dummy-token"
```

### 6. Завершить бронирование

```bash
curl -X POST http://localhost:8080/api/v1/bookings/{booking_id}/finish \
  -H "Authorization: Bearer dummy-token"
```

### 7. Отметить как неявку (no-show)

```bash
curl -X POST http://localhost:8080/api/v1/bookings/{booking_id}/no-show \
  -H "Authorization: Bearer dummy-token"
```

## Полный список endpoints

### Venues (заведения)
- `GET /api/v1/venues` - список заведений
- `GET /api/v1/venues/:id` - получить заведение
- `POST /api/v1/venues` - создать заведение
- `PUT /api/v1/venues/:id` - обновить заведение
- `DELETE /api/v1/venues/:id` - удалить заведение

### Rooms (залы)
- `GET /api/v1/venues/:venueId/rooms` - список залов
- `GET /api/v1/rooms/:id` - получить зал
- `POST /api/v1/venues/:venueId/rooms` - создать зал
- `PUT /api/v1/rooms/:id` - обновить зал
- `DELETE /api/v1/rooms/:id` - удалить зал

### Tables (столы)
- `GET /api/v1/rooms/:roomId/tables` - список столов
- `GET /api/v1/tables/:id` - получить стол
- `POST /api/v1/rooms/:roomId/tables` - создать стол
- `PUT /api/v1/tables/:id` - обновить стол
- `DELETE /api/v1/tables/:id` - удалить стол

### Schedule (расписание)
- `GET /api/v1/venues/:venueId/schedule` - получить расписание
- `POST /api/v1/venues/:venueId/schedule` - установить расписание
- `POST /api/v1/venues/:venueId/special-hours` - установить особые часы

### Bookings (бронирования)
- `GET /api/v1/bookings` - список бронирований
- `GET /api/v1/bookings/:id` - получить бронирование
- `POST /api/v1/bookings` - создать бронирование
- `POST /api/v1/bookings/:id/confirm` - подтвердить
- `POST /api/v1/bookings/:id/cancel` - отменить
- `POST /api/v1/bookings/:id/seat` - отметить как севших
- `POST /api/v1/bookings/:id/finish` - завершить
- `POST /api/v1/bookings/:id/no-show` - отметить неявку

### Availability (доступность)
- `POST /api/v1/availability/check` - проверить доступность

## Прямой доступ к gRPC сервисам (для разработки)

Если нужно напрямую обратиться к gRPC сервисам (для отладки):

### Использование grpcurl

```bash
# Установить grpcurl (если нет)
# Windows: choco install grpcurl
# Linux/Mac: go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Получить список методов venue-svc
grpcurl -plaintext localhost:50051 list

# Вызвать метод ListVenues
grpcurl -plaintext -d '{"limit": 10, "offset": 0}' \
  localhost:50051 venue.VenueService/ListVenues

# Вызвать метод ListBookings
grpcurl -plaintext -d '{"venue_id": "venue-1"}' \
  localhost:50052 booking.BookingService/ListBookings
```

## Веб-интерфейс

Откройте в браузере: `http://localhost:8080`

Там должен быть фронтенд для управления (если он реализован в `web/dist`).

## Мониторинг

- **Grafana**: http://localhost:3000 - дашборды с метриками
- **Prometheus**: http://localhost:9090 - метрики напрямую
- **Jaeger**: http://localhost:16686 - трейсинг запросов
- **Kafka UI**: http://localhost:8081 - просмотр событий Kafka

## Примеры использования

См. файл `examples/api-examples.sh` (можно создать) для готовых скриптов.

