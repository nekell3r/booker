# Все команды работают через Docker - не требуется установка Go, protoc и других инструментов локально

.PHONY: up down logs migrate seed gen clean test build rebuild

# Generate protobuf files using Docker
gen:
	@echo "Generating protobuf files using Docker..."
	@docker run --rm -v ${PWD}:/workspace -w /workspace golang:1.21-alpine sh -c "\
		apk add --no-cache protobuf protobuf-dev && \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
		mkdir -p pkg/proto/common pkg/proto/venue pkg/proto/booking && \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto/common/*.proto && \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto/venue/*.proto && \
		protoc --go_out=. --go_opt=paths=source_relative \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative \
			proto/booking/*.proto"

# Start infrastructure and services
up:
	@echo "Starting infrastructure and services..."
	@docker-compose --profile infra-min --profile apps up -d
	@echo "Waiting for services to be ready..."
	@sleep 5
	@echo "Services started! Admin Gateway: http://localhost:8080"

# Start full infrastructure
up-full:
	@echo "Starting full infrastructure..."
	@docker-compose --profile infra-full --profile apps up -d
	@sleep 5
	@echo "Services started!"

# Stop all services
down:
	@echo "Stopping all services..."
	@docker-compose --profile infra-min --profile apps --profile tools down -v


# Run migrations (requires infrastructure to be running)
migrate:
	@echo "Running migrations..."
	@docker-compose --profile infra-min run --rm migrate

# Seed sample data (requires infrastructure and migrations)
seed:
	@echo "Seeding sample data..."
	@docker-compose --profile infra-min run --rm seed

# Run tests using Docker
test:
	@echo "Running tests..."
	@docker run --rm -v ${PWD}:/workspace -w /workspace golang:1.21-alpine \
		sh -c "go mod download && go test ./..."

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@docker run --rm -v ${PWD}:/workspace -w /workspace alpine \
		sh -c "rm -rf pkg/proto/*.pb.go pkg/proto/*_grpc.pb.go"

# Build all services (rebuild images)
build:
	@echo "Building all services..."
	@docker-compose build

# Rebuild and restart services
rebuild:
	@echo "Rebuilding and restarting services..."
	@docker-compose build --no-cache
	@docker-compose --profile infra-min --profile apps up -d

# Full setup: start infrastructure, run migrations, seed data
setup:
	@echo "Starting full setup..."
	@echo "Step 1: Starting infrastructure and services..."
	@docker-compose --profile infra-min --profile apps up -d
	@echo "Step 2: Waiting for databases to be ready..."
	@sleep 15
	@echo "Step 3: Running migrations..."
	@docker-compose --profile infra-min run --rm migrate || (echo "Migration failed, retrying..." && sleep 5 && docker-compose --profile infra-min run --rm migrate)
	@echo "Step 4: Seeding sample data..."
	@docker-compose --profile infra-min run --rm seed
	@echo ""
	@echo "✅ Setup complete!"
	@echo "📊 Admin Gateway: http://localhost:8080"
	@echo "📈 Grafana: http://localhost:3000 (admin/admin)"
	@echo "🔍 Jaeger: http://localhost:16686"
	@echo "📉 Prometheus: http://localhost:9090"
	@echo "📨 Kafka UI: http://localhost:8081"

# Show service status
status:
	@docker-compose ps

# View logs
logs:
	@docker-compose --profile infra-min --profile apps logs -f admin-gateway venue-svc booking-svc

logs-all:
	@docker-compose --profile infra-min --profile apps logs -f

logs-gateway:
	@docker-compose --profile infra-min --profile apps logs -f admin-gateway

logs-venue:
	@docker-compose --profile infra-min --profile apps logs -f venue-svc

logs-booking:
	@docker-compose --profile infra-min --profile apps logs -f booking-svc

logs-notify:
	@docker-compose --profile infra-min --profile apps logs -f notify-svc

# Restart a specific service
restart:
	@docker-compose restart $(SERVICE)

# Execute command in a service container
exec:
	@docker-compose exec $(SERVICE) $(CMD)
