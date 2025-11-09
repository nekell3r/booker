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

# Pull all required images (with manual retry on failure)
pull-images:
	@echo "Pulling Docker images (this may take a while)..."
	@docker-compose --profile infra-min --profile apps pull || (echo "⚠️  Pull failed. This is often due to network issues." && echo "💡 Try running this command again, or check your network connection." && exit 1)
	@echo "✅ Images pulled successfully!"

# Full setup: start infrastructure, run migrations, seed data
setup:
	@echo "Starting full setup..."
	@echo "Step 1: Pulling Docker images (if needed)..."
	@docker-compose --profile infra-min --profile apps pull || echo "⚠️  Some images failed to pull, continuing anyway..."
	@echo "Step 2: Starting infrastructure and services..."
	@docker-compose --profile infra-min --profile apps up -d || (echo "❌ Failed to start services. This might be due to:" && echo "   - Network issues pulling images" && echo "   - Port conflicts" && echo "   - Docker daemon not running" && echo "" && echo "💡 Try running 'make pull-images' first, then 'make up'" && exit 1)
	@echo "Step 3: Waiting for databases to be ready..."
	@sleep 15
	@echo "Step 4: Running migrations..."
	@docker-compose --profile infra-min run --rm migrate || (echo "Migration failed, retrying..." && sleep 5 && docker-compose --profile infra-min run --rm migrate)
	@echo "Step 5: Seeding sample data..."
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

# Restart all application services
restart:
	@echo "Restarting all application services..."
	@docker-compose --profile infra-min --profile apps restart

# Restart a specific service
restart-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please specify SERVICE variable. Example: make restart-service SERVICE=admin-gateway"; \
		exit 1; \
	fi
	@echo "Restarting $(SERVICE)..."
	@docker-compose --profile infra-min --profile apps restart $(SERVICE)

# Restart all services (infrastructure + apps)
restart-all:
	@echo "Restarting all services (infrastructure + apps)..."
	@docker-compose --profile infra-min --profile apps restart

# Rebuild and restart a specific service
rebuild-service:
	@if [ -z "$(SERVICE)" ]; then \
		echo "❌ Please specify SERVICE variable. Example: make rebuild-service SERVICE=admin-gateway"; \
		exit 1; \
	fi
	@echo "Rebuilding and restarting $(SERVICE)..."
	@docker-compose --profile infra-min --profile apps build $(SERVICE)
	@docker-compose --profile infra-min --profile apps up -d $(SERVICE)

# Execute command in a service container
exec:
	@if [ -z "$(SERVICE)" ] || [ -z "$(CMD)" ]; then \
		echo "❌ Please specify SERVICE and CMD variables. Example: make exec SERVICE=admin-gateway CMD='sh'"; \
		exit 1; \
	fi
	@docker-compose exec $(SERVICE) $(CMD)
