# Полная настройка через Docker

## Что изменилось

Теперь **всё работает в Docker контейнерах**:
- ✅ Go программы компилируются в контейнерах
- ✅ Proto файлы генерируются в контейнерах
- ✅ Миграции выполняются в контейнерах
- ✅ Seed данные заполняются в контейнерах
- ✅ Тесты запускаются в контейнерах

**Не требуется установка:**
- ❌ Go
- ❌ protoc
- ❌ PostgreSQL клиентов
- ❌ Любых других инструментов

## Быстрый старт

### Вариант 1: Всё в одном (рекомендуется)

```bash
make setup
```

Эта команда:
1. Запустит всю инфраструктуру (Postgres, Redis, Kafka, Prometheus, Grafana, Jaeger)
2. Соберёт и запустит все сервисы (admin-gateway, venue-svc, booking-svc, notify-svc)
3. Выполнит миграции БД
4. Заполнит тестовыми данными

### Вариант 2: Пошагово

```bash
# 1. Запустить инфраструктуру и сервисы
make up

# 2. Подождать пока БД будут готовы (10-15 секунд)
# 3. Выполнить миграции
make migrate

# 4. Заполнить тестовыми данными
make seed
```

## Основные команды

| Команда | Описание |
|---------|----------|
| `make setup` | Полная настройка (инфраструктура + миграции + seed) |
| `make up` | Запустить все сервисы |
| `make down` | Остановить все сервисы |
| `make logs` | Просмотр логов всех сервисов |
| `make migrate` | Выполнить миграции БД |
| `make seed` | Заполнить тестовыми данными |
| `make gen` | Сгенерировать proto файлы (через Docker) |
| `make build` | Пересобрать все Docker образы |
| `make rebuild` | Пересобрать и перезапустить |
| `make test` | Запустить тесты (через Docker) |
| `make status` | Показать статус всех сервисов |

## Доступ к сервисам

После запуска доступны:
- **Admin Gateway**: http://localhost:8080
- **Grafana**: http://localhost:3000 (admin/admin)
- **Jaeger UI**: http://localhost:16686
- **Prometheus**: http://localhost:9090

## Разработка

### Изменение кода

При изменении кода нужно пересобрать образ:

```bash
# Пересобрать конкретный сервис
docker-compose build admin-gateway
docker-compose up -d admin-gateway

# Или пересобрать всё
make rebuild
```

### Просмотр логов

```bash
# Все сервисы
make logs

# Конкретный сервис
docker-compose logs -f admin-gateway
```

### Выполнение команд в контейнере

```bash
# Зайти в контейнер
docker-compose exec admin-gateway sh

# Выполнить команду
docker-compose exec venue-svc sh -c "ls -la"
```

## Очистка

```bash
# Остановить и удалить контейнеры
make down

# Остановить и удалить контейнеры + volumes (удалит все данные!)
docker-compose down -v

# Удалить все образы проекта
docker-compose down --rmi all
```

## Устранение проблем

### Сервисы не запускаются

```bash
# Проверить статус
make status

# Посмотреть логи
make logs

# Пересобрать образы
make rebuild
```

### Миграции не выполняются

```bash
# Убедиться что БД запущены
docker-compose ps

# Выполнить миграции вручную
docker-compose run --rm migrate
```

### Порт занят

Если порт занят, измените его в `docker-compose.yml`:

```yaml
ports:
  - "8081:8080"  # Вместо 8080:8080
```

## Производительность

Для ускорения сборки в Docker можно использовать BuildKit:

```bash
export DOCKER_BUILDKIT=1
export COMPOSE_DOCKER_CLI_BUILD=1
make build
```

Или добавить в `docker-compose.yml`:

```yaml
x-build-args: &build-args
  DOCKER_BUILDKIT: 1
  COMPOSE_DOCKER_CLI_BUILD: 1
```


