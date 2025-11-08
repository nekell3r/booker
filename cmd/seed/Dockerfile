FROM golang:1.23-alpine

WORKDIR /app

# Install git for go mod download
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY cmd/seed ./cmd/seed

WORKDIR /app/cmd/seed

CMD ["go", "run", "main.go"]


