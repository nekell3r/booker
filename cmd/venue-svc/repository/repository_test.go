package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVenueModel tests the Venue model structure
func TestVenueModel(t *testing.T) {
	venue := &Venue{
		ID:        uuid.New().String(),
		Name:      "Test Venue",
		Timezone:  "UTC",
		Address:   "123 Main St",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.NotEmpty(t, venue.ID)
	assert.Equal(t, "Test Venue", venue.Name)
	assert.Equal(t, "UTC", venue.Timezone)
	assert.Equal(t, "123 Main St", venue.Address)
}

// TestRoomModel tests the Room model structure
func TestRoomModel(t *testing.T) {
	room := &Room{
		ID:        uuid.New().String(),
		VenueID:   "venue-1",
		Name:      "Main Room",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.NotEmpty(t, room.ID)
	assert.Equal(t, "venue-1", room.VenueID)
	assert.Equal(t, "Main Room", room.Name)
}

// TestTableModel tests the Table model structure
func TestTableModel(t *testing.T) {
	table := &Table{
		ID:        uuid.New().String(),
		RoomID:    "room-1",
		Name:      "Table 1",
		Capacity:  4,
		CanMerge:  true,
		Zone:      "window",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	assert.NotEmpty(t, table.ID)
	assert.Equal(t, "room-1", table.RoomID)
	assert.Equal(t, "Table 1", table.Name)
	assert.Equal(t, int32(4), table.Capacity)
	assert.True(t, table.CanMerge)
	assert.Equal(t, "window", table.Zone)
}

// TestRepository_New tests repository creation
func TestRepository_New(t *testing.T) {
	// This test would require a real database connection
	// For now, we test the structure
	// In a real test, we'd create a test database connection
	// repo := New(testDB, redisClient)
	// assert.NotNil(t, repo)
}

// Integration test helpers (would be in integration test file)
func setupTestDB(t *testing.T) (*Repository, func()) {
	// Setup test database using testcontainers or in-memory database
	// Return repository and cleanup function
	t.Skip("Integration test - requires database")
	return nil, func() {}
}

func TestRepository_CreateVenue_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()
	id, err := repo.CreateVenue(ctx, "Test Venue", "UTC", "123 Main St")
	require.NoError(t, err)
	assert.NotEmpty(t, id)

	// Verify venue was created
	venue, err := repo.GetVenue(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, "Test Venue", venue.Name)
	assert.Equal(t, "UTC", venue.Timezone)
	assert.Equal(t, "123 Main St", venue.Address)
}

func TestRepository_CreateRoom_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create venue first
	venueID, err := repo.CreateVenue(ctx, "Test Venue", "UTC", "123 Main St")
	require.NoError(t, err)

	// Create room
	roomID, err := repo.CreateRoom(ctx, venueID, "Main Room")
	require.NoError(t, err)
	assert.NotEmpty(t, roomID)

	// Verify room was created
	room, err := repo.GetRoom(ctx, roomID)
	require.NoError(t, err)
	assert.Equal(t, "Main Room", room.Name)
	assert.Equal(t, venueID, room.VenueID)
}

func TestRepository_CreateTable_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create venue and room first
	venueID, err := repo.CreateVenue(ctx, "Test Venue", "UTC", "123 Main St")
	require.NoError(t, err)

	roomID, err := repo.CreateRoom(ctx, venueID, "Main Room")
	require.NoError(t, err)

	// Create table
	tableID, err := repo.CreateTable(ctx, roomID, "Table 1", 4, true, "window")
	require.NoError(t, err)
	assert.NotEmpty(t, tableID)

	// Verify table was created
	table, err := repo.GetTable(ctx, tableID)
	require.NoError(t, err)
	assert.Equal(t, "Table 1", table.Name)
	assert.Equal(t, int32(4), table.Capacity)
	assert.True(t, table.CanMerge)
	assert.Equal(t, "window", table.Zone)
	assert.Equal(t, roomID, table.RoomID)
}

func TestRepository_ListTables_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	repo, cleanup := setupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	// Create venue, room, and tables
	venueID, err := repo.CreateVenue(ctx, "Test Venue", "UTC", "123 Main St")
	require.NoError(t, err)

	roomID, err := repo.CreateRoom(ctx, venueID, "Main Room")
	require.NoError(t, err)

	_, err = repo.CreateTable(ctx, roomID, "Table 1", 4, true, "window")
	require.NoError(t, err)

	_, err = repo.CreateTable(ctx, roomID, "Table 2", 2, false, "corner")
	require.NoError(t, err)

	// List tables by room
	tables, total, err := repo.ListTables(ctx, roomID, "", 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int32(2), total)
	assert.Equal(t, 2, len(tables))

	// List tables by venue
	tables, total, err = repo.ListTables(ctx, "", venueID, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int32(2), total)
	assert.Equal(t, 2, len(tables))
}

