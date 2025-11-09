FROM golang:1.23-alpine AS protoc-builder

# Install protoc, git and plugins
RUN apk add --no-cache protobuf protobuf-dev git
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /app
COPY proto/ ./proto/
COPY go.mod go.sum ./

# Generate proto files
RUN mkdir -p pkg/proto/common pkg/proto/venue pkg/proto/booking
RUN protoc --proto_path=proto --go_out=pkg/proto --go_opt=paths=source_relative \
    --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative \
    proto/common/*.proto
RUN protoc --proto_path=proto --go_out=pkg/proto --go_opt=paths=source_relative \
    --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative \
    proto/venue/*.proto proto/common/*.proto
RUN protoc --proto_path=proto --go_out=pkg/proto --go_opt=paths=source_relative \
    --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative \
    proto/booking/*.proto proto/common/*.proto

FROM golang:1.23-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Build cache invalidation (increment to force rebuild)
ARG CACHE_BUST=1
RUN echo "Cache bust: $CACHE_BUST"

# Copy source code first
COPY cmd/ ./cmd/

# Copy pkg subdirectories (create pkg structure)
COPY pkg/kafka/ ./pkg/kafka/
COPY pkg/tracing/ ./pkg/tracing/
COPY pkg/redis/ ./pkg/redis/
COPY pkg/metrics/ ./pkg/metrics/

# Copy proto generated files from previous stage (after other pkg subdirs)
COPY --from=protoc-builder /app/pkg/proto ./pkg/proto

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/admin-gateway ./cmd/admin-gateway

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/bin/admin-gateway .
COPY web/dist ./web/dist

EXPOSE 8080

CMD ["./admin-gateway"]
