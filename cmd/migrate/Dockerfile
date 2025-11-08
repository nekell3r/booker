FROM golang:1.23-alpine

WORKDIR /app

# Install postgres client and git for migrations
RUN apk add --no-cache postgresql-client git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY cmd/migrate ./cmd/migrate
COPY migrations ./migrations

WORKDIR /app/cmd/migrate

CMD ["go", "run", "main.go"]


