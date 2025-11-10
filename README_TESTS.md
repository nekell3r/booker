# Тесты для проекта Booker

Этот документ описывает структуру тестов в проекте и способы их запуска.

## Запуск тестов через Docker

Все тесты запускаются через Docker, так как проект полностью контейнеризован.

### Основные команды

```bash
# Запустить все тесты
make test

# Запустить тесты с покрытием
make test-coverage

# Запустить тесты с подробным выводом
make test-verbose

# Запустить только юнит-тесты (без интеграционных)
make test-unit

# Запустить тесты для конкретного пакета
make test-package PACKAGE=./cmd/booking-svc/service
make test-package PACKAGE=./cmd/venue-svc/repository

# Запустить тесты и сгенерировать HTML отчет о покрытии
make test-coverage-html
```

### Примеры использования

```bash
# Все тесты
make test

# Только быстрые юнит-тесты
make test-unit

# Тесты конкретного сервиса
make test-package PACKAGE=./cmd/booking-svc/service
make test-package PACKAGE=./cmd/venue-svc/service
make test-package PACKAGE=./cmd/admin-gateway/handlers

# Тесты с покрытием
make test-coverage

# HTML отчет о покрытии (создаст файл coverage.html)
make test-coverage-html
```

## Структура тестов

### Юнит-тесты

Юнит-тесты находятся рядом с тестируемым кодом и имеют суффикс `_test.go`:

- ✅ `cmd/booking-svc/repository/repository_test.go` - тесты для booking repository
- ✅ `cmd/venue-svc/repository/repository_test.go` - тесты для venue repository
- ✅ `pkg/kafka/producer_test.go` - тесты структур Kafka событий
- ✅ `pkg/redis/client_test.go` - тесты форматов Redis ключей

### Интеграционные тесты

Интеграционные тесты находятся в директории `integration/`:

- ✅ `integration/api_test.go` - базовые интеграционные тесты для API endpoints

**Примечание:** Текущие интеграционные тесты проверяют базовую структуру API. Полноценные интеграционные тесты с реальными сервисами требуют testcontainers. Интеграционные тесты пропускаются при использовании флага `-short` (команда `make test-unit`).

### Тестовые утилиты

Моки и тестовые утилиты находятся в `internal/testutil/`:

- `internal/testutil/mocks.go` - моки для репозиториев, клиентов и сервисов

## Запуск тестов напрямую через Docker

Если нужно запустить тесты напрямую через `docker run`:

```bash
# Все тесты
docker run --rm -v ${PWD}:/workspace -w /workspace golang:1.23-alpine \
  sh -c "apk add --no-cache git && go mod download && go test ./..."

# Только юнит-тесты
docker run --rm -v ${PWD}:/workspace -w /workspace golang:1.23-alpine \
  sh -c "apk add --no-cache git && go mod download && go test -short ./..."

# Конкретный пакет
docker run --rm -v ${PWD}:/workspace -w /workspace golang:1.23-alpine \
  sh -c "apk add --no-cache git && go mod download && go test -v ./cmd/booking-svc/service"
```

## Использование моков

Пример использования моков в тестах:

```go
func TestService_CreateBooking(t *testing.T) {
    repo := &testutil.MockBookingRepository{}
    repo.CreateBookingFunc = func(ctx context.Context, booking *repository.Booking) error {
        // Ваша логика мока
        return nil
    }
    
    // Использование мока в тесте
    err := repo.CreateBooking(ctx, booking)
    assert.NoError(t, err)
}
```

## Зависимости для тестов

Проект использует следующие библиотеки для тестирования:

- `github.com/stretchr/testify/assert` - assertions
- `github.com/stretchr/testify/require` - require assertions (останавливают тест при ошибке)
- `github.com/stretchr/testify/mock` - моки (для будущего использования)

## Примечания

1. **Proto файлы** должны быть сгенерированы перед запуском тестов. Они генерируются автоматически при сборке Docker-образов, или можно запустить `make gen`.

2. **Интеграционные тесты** требуют запущенной инфраструктуры (PostgreSQL, Redis, Kafka). Используйте `make test-unit` для пропуска интеграционных тестов.

3. **Моки** в `internal/testutil/mocks.go` упрощены. Для production рекомендуется использовать более продвинутые инструменты генерации моков (например, `mockgen`).

4. **Покрытие кода**: HTML отчет о покрытии создается в файле `coverage.html` в корне проекта после выполнения `make test-coverage-html`.

## Текущее состояние

**Реализовано:**
- ✅ Repository тесты (booking, venue)
- ✅ Component тесты (kafka, redis)  
- ✅ Базовые integration тесты (API health check)
- ✅ Тестовые моки для repositories и gRPC клиентов
- ✅ Команды для запуска тестов через Makefile
- ✅ Поддержка Windows (docker-compose вместо прямых docker run)

**TODO (для будущего):**
- [ ] Service layer тесты (требуют полную реализацию всех gRPC методов в моках)
- [ ] Handlers тесты (требуют полную реализацию всех gRPC методов в моках)
- [ ] Интеграционные тесты с testcontainers (PostgreSQL, Redis, Kafka)
- [ ] Тесты для middleware
- [ ] Тесты для конфигурации
- [ ] Увеличение покрытия кода
- [ ] Benchmarks (тесты производительности)
