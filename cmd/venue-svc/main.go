package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"booker/cmd/venue-svc/config"
	"booker/cmd/venue-svc/repository"
	"booker/cmd/venue-svc/service"
	"booker/pkg/kafka"
	"booker/pkg/metrics"
	"booker/pkg/redis"
	"booker/pkg/tracing"
	bookingpb "booker/pkg/proto/booking"
	venuepb "booker/pkg/proto/venue"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	// Logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if cfg.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	// Tracing
	shutdown, err := tracing.InitTracer("venue-svc", cfg.JaegerEndpoint)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize tracer")
	}
	defer shutdown()

	// PostgreSQL
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDB)
	
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pool.Close()

	// Redis
	redisClient := redis.NewClient(cfg.RedisAddr, cfg.RedisPassword)

	// Kafka Producer with retry logic
	kafkaBrokers := []string{cfg.KafkaBrokers}
	var producer *kafka.Producer
	maxRetries := 20
	retryDelay := 3 * time.Second
	log.Info().Strs("brokers", kafkaBrokers).Msg("Attempting to connect to Kafka...")
	for i := 0; i < maxRetries; i++ {
		var err error
		producer, err = kafka.NewProducer(kafkaBrokers)
		if err == nil {
			log.Info().Msg("Kafka producer connected successfully")
			break
		}
		if i < maxRetries-1 {
			log.Warn().Err(err).Int("attempt", i+1).Int("max_retries", maxRetries).Dur("retry_delay", retryDelay).Msg("Failed to create Kafka producer, retrying...")
			time.Sleep(retryDelay)
		} else {
			log.Fatal().Err(err).Int("total_attempts", maxRetries).Msg("Failed to create Kafka producer after all retries")
		}
	}
	if producer == nil {
		log.Fatal().Msg("Kafka producer is nil after retry loop")
	}
	defer producer.Close()

	// Repository
	repo := repository.New(pool, redisClient)

	// Booking service client
	bookingConn, err := grpc.Dial(
		cfg.BookingSvcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to booking service")
	}
	defer bookingConn.Close()
	bookingClient := bookingpb.NewBookingServiceClient(bookingConn)

	// Service
	svc := service.New(repo, producer, bookingClient, cfg)

	// Start metrics server
	startMetricsServer(cfg.MetricsPort)

	// gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(metrics.UnaryServerMetricsInterceptor("venue-svc")),
	)
	venuepb.RegisterVenueServiceServer(s, svc)

	// Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		log.Info().Int("port", cfg.Port).Msg("Venue service started")
		if err := s.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("Server failed")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("Shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		log.Info().Msg("Server stopped gracefully")
	case <-shutdownCtx.Done():
		log.Warn().Msg("Shutdown timeout, forcing stop")
		s.Stop()
	}
}


