# Проблемы при старте сервисов

## Почему сервис падает?

### Проблема: Race Condition при старте

Когда вы запускаете `make restart-all`, происходит следующее:

1. **Все контейнеры перезапускаются одновременно**:
   ```
   venue-svc    → стартует сразу
   booking-svc  → стартует сразу  
   redpanda     → стартует сразу (но медленнее!)
   ```

2. **Сервисы пытаются подключиться к Kafka при старте**:
   ```go
   // В main() функции venue-svc/main.go
   producer, err := kafka.NewProducer(kafkaBrokers)
   if err != nil {
       log.Fatal().Err(err).Msg("Failed to create Kafka producer")  // ← ПАДАЕТ ЗДЕСЬ!
   }
   ```

3. **Kafka (Redpanda) еще не готова**:
   - Redpanda запускается медленнее (инициализация кластера)
   - Healthcheck проверяет готовность: `rpk cluster health`
   - Это может занять 10-50 секунд (5 retries × 10s interval)
   - Но сервисы стартуют быстрее и пытаются подключиться сразу

4. **Результат**: 
   ```
   venue-svc → подключается к redpanda:9092 → connection refused → Fatal → контейнер падает
   ```

## Почему нужны ретраи?

### Проблема с `depends_on`

В `docker-compose.yml` есть зависимость:
```yaml
venue-svc:
  depends_on:
    redpanda:
      condition: service_healthy  # ← Ждет healthcheck
```

**НО** это не гарантирует, что:
1. Kafka полностью готова принимать подключения
2. Healthcheck может пройти, но порт 9092 еще не открыт
3. При перезапуске (`restart`) зависимости не всегда соблюдаются строго

### Решение: Retry в коде

Ретраи дают время Kafka запуститься:

```go
maxRetries := 20
retryDelay := 3 * time.Second
// Всего до 60 секунд ожидания

for i := 0; i < maxRetries; i++ {
    producer, err := kafka.NewProducer(kafkaBrokers)
    if err == nil {
        break  // Успех!
    }
    time.Sleep(retryDelay)  // Ждем и пробуем снова
}
```

**Преимущества**:
- ✅ Сервис не падает сразу
- ✅ Дает время Kafka запуститься
- ✅ Работает даже если healthcheck прошел рано
- ✅ Защита от временных сетевых проблем

## Альтернативные решения

### 1. Улучшить healthcheck Redpanda

```yaml
healthcheck:
  test: ["CMD-SHELL", "rpk cluster health | grep -q 'Healthy' && nc -z localhost 9092 || exit 1"]
  # Проверяет не только health, но и доступность порта
```

### 2. Lazy Connection (ленивое подключение)

Вместо подключения при старте, подключаться при первом использовании:

```go
type Service struct {
    producer *kafka.Producer
    mu sync.Mutex
}

func (s *Service) getProducer() (*kafka.Producer, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if s.producer == nil {
        var err error
        s.producer, err = kafka.NewProducer(brokers)
        if err != nil {
            return nil, err
        }
    }
    return s.producer, nil
}
```

**Плюсы**: Сервис стартует даже если Kafka недоступна
**Минусы**: Сложнее, нужно обрабатывать ошибки при публикации

### 3. Circuit Breaker Pattern

Использовать паттерн Circuit Breaker для автоматических ретраев:

```go
breaker := circuit.NewBreaker(circuit.Config{
    MaxFailures: 5,
    Timeout: 30 * time.Second,
})

producer, err := breaker.Execute(func() (interface{}, error) {
    return kafka.NewProducer(brokers)
})
```

## Текущее решение (рекомендуемое)

**Комбинация подходов**:

1. ✅ `depends_on: service_healthy` - базовая защита
2. ✅ Retry в коде - дополнительная защита
3. ✅ Логирование попыток - для отладки

Это дает:
- Надежный старт при первом запуске
- Устойчивость к перезапускам
- Понятные логи при проблемах

## Время запуска

При идеальных условиях:
- Redpanda healthcheck: ~10-20 секунд
- Retry логика: до 60 секунд (20 попыток × 3 сек)
- **Итого**: сервисы могут стартовать до 60 секунд, но обычно ~10-30 секунд

Это нормально для микросервисной архитектуры!

