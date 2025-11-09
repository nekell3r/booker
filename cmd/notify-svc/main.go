package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"booker/pkg/kafka"
	"booker/pkg/tracing"
)

func main() {
	// Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Tracing
	jaegerEndpoint := os.Getenv("JAEGER_ENDPOINT")
	shutdown, err := tracing.InitTracer("notify-svc", jaegerEndpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize tracer")
	}
	defer shutdown()

	// Kafka consumer with retry logic
	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	if brokers[0] == "" {
		brokers = []string{"localhost:9092"}
	}

	handler := &BookingEventHandler{}
	var consumer *kafka.Consumer
	maxRetries := 20
	retryDelay := 3 * time.Second
	log.Info().Strs("brokers", brokers).Msg("Attempting to connect to Kafka...")
	for i := 0; i < maxRetries; i++ {
		var err error
		consumer, err = kafka.NewConsumer(brokers, "notify-svc-group", handler)
		if err == nil {
			log.Info().Msg("Kafka consumer connected successfully")
			break
		}
		if i < maxRetries-1 {
			log.Warn().Err(err).Int("attempt", i+1).Int("max_retries", maxRetries).Dur("retry_delay", retryDelay).Msg("Failed to create Kafka consumer, retrying...")
			time.Sleep(retryDelay)
		} else {
			log.Fatal().Err(err).Int("total_attempts", maxRetries).Msg("Failed to create Kafka consumer after all retries")
		}
	}
	if consumer == nil {
		log.Fatal().Msg("Kafka consumer is nil after retry loop")
	}
	defer consumer.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	topics := []string{"booking.confirmed", "booking.cancelled", "booking.no_show"}

	go func() {
		for {
			if err := consumer.Consume(ctx, topics); err != nil {
				log.Error().Err(err).Msg("Consumer error")
			}
		}
	}()

	log.Info().Strs("topics", topics).Msg("Notify service started")

	<-ctx.Done()
	log.Info().Msg("Shutting down...")
}

type BookingEventHandler struct{}

func (h *BookingEventHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *BookingEventHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *BookingEventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Info().
			Str("topic", message.Topic).
			Str("key", string(message.Key)).
			Msg("Received booking event")

		// TODO: Send notification (email/SMS/Telegram)
		// For MVP, just log

		session.MarkMessage(message, "")
	}
	return nil
}


