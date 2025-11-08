package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Venue DB
	venueDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("VENUE_DB_USER"), os.Getenv("VENUE_DB_PASSWORD"),
		os.Getenv("VENUE_DB_HOST"), os.Getenv("VENUE_DB_PORT"), os.Getenv("VENUE_DB_NAME"))

	venueDB, err := sql.Open("pgx", venueDSN)
	if err != nil {
		panic(err)
	}
	defer venueDB.Close()

	// Read and execute venue migrations
	venueSQL, err := os.ReadFile("../../migrations/001_venue_schema.sql")
	if err != nil {
		panic(err)
	}

	if _, err := venueDB.Exec(string(venueSQL)); err != nil {
		panic(err)
	}

	fmt.Println("Venue migrations applied")

	// Booking DB
	bookingDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("BOOKING_DB_USER"), os.Getenv("BOOKING_DB_PASSWORD"),
		os.Getenv("BOOKING_DB_HOST"), os.Getenv("BOOKING_DB_PORT"), os.Getenv("BOOKING_DB_NAME"))

	bookingDB, err := sql.Open("pgx", bookingDSN)
	if err != nil {
		panic(err)
	}
	defer bookingDB.Close()

	// Read and execute booking migrations
	bookingSQL, err := os.ReadFile("../../migrations/002_booking_schema.sql")
	if err != nil {
		panic(err)
	}

	if _, err := bookingDB.Exec(string(bookingSQL)); err != nil {
		panic(err)
	}

	fmt.Println("Booking migrations applied")
}
