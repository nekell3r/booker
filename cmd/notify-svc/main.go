package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

	// Kafka consumer
	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	if brokers[0] == "" {
		brokers = []string{"localhost:9092"}
	}

	handler := &BookingEventHandler{}
	consumer, err := kafka.NewConsumer(brokers, "notify-svc-group", handler)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create consumer")
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


