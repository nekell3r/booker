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
	"booker/pkg/redis"
	"booker/pkg/tracing"
	venuepb "booker/pkg/proto/venue"
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

	// Kafka Producer
	kafkaBrokers := []string{cfg.KafkaBrokers}
	producer, err := kafka.NewProducer(kafkaBrokers)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create Kafka producer")
	}
	defer producer.Close()

	// Repository
	repo := repository.New(pool, redisClient)

	// Service
	svc := service.New(repo, producer, cfg)

	// Start metrics server
	startMetricsServer(cfg.MetricsPort)

	// gRPC Server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to listen")
	}

	s := grpc.NewServer()
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


