package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Venue DB - use environment variables or defaults
	venueHost := getEnv("VENUE_DB_HOST", "postgres-venue")
	venueUser := getEnv("VENUE_DB_USER", "venue_user")
	venuePass := getEnv("VENUE_DB_PASSWORD", "venue_pass")
	venuePort := getEnv("VENUE_DB_PORT", "5432")
	venueDB := getEnv("VENUE_DB_NAME", "venue")
	
	venueDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		venueUser, venuePass, venueHost, venuePort, venueDB)
	venuePool, err := pgxpool.New(context.Background(), venueDSN)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to venue DB: %v", err))
	}
	defer venuePool.Close()

	// Create sample venue
	venueID := uuid.New().String()
	_, err = venuePool.Exec(context.Background(),
		`INSERT INTO venues (id, name, timezone, address) VALUES ($1, $2, $3, $4)`,
		venueID, "Test Restaurant", "Europe/Moscow", "123 Main St")
	if err != nil {
		panic(err)
	}

	// Create sample room
	roomID := uuid.New().String()
	_, err = venuePool.Exec(context.Background(),
		`INSERT INTO rooms (id, venue_id, name) VALUES ($1, $2, $3)`,
		roomID, venueID, "Main Hall")
	if err != nil {
		panic(err)
	}

	// Create sample tables
	for i := 1; i <= 5; i++ {
		tableID := uuid.New().String()
		_, err = venuePool.Exec(context.Background(),
			`INSERT INTO tables (id, room_id, name, capacity, can_merge, zone) VALUES ($1, $2, $3, $4, $5, $6)`,
			tableID, roomID, fmt.Sprintf("Table %d", i), 4, false, "main")
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Sample data seeded successfully")
}


