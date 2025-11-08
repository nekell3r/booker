package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"booker/pkg/redis"
)

type Repository struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func New(db *pgxpool.Pool, redis *redis.Client) *Repository {
	return &Repository{
		db:    db,
		redis: redis,
	}
}

// Venue operations
func (r *Repository) CreateVenue(ctx context.Context, name, timezone, address string) (string, error) {
	id := uuid.New().String()
	_, err := r.db.Exec(ctx,
		`INSERT INTO venues (id, name, timezone, address, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, NOW(), NOW())`,
		id, name, timezone, address)
	return id, err
}

func (r *Repository) GetVenue(ctx context.Context, id string) (*Venue, error) {
	var v Venue
	err := r.db.QueryRow(ctx,
		`SELECT id, name, timezone, address, created_at, updated_at
		 FROM venues WHERE id = $1`, id).
		Scan(&v.ID, &v.Name, &v.Timezone, &v.Address, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *Repository) ListVenues(ctx context.Context, limit, offset int32) ([]*Venue, int32, error) {
	var total int32
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM venues`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, name, timezone, address, created_at, updated_at
		 FROM venues ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var venues []*Venue
	for rows.Next() {
		var v Venue
		if err := rows.Scan(&v.ID, &v.Name, &v.Timezone, &v.Address, &v.CreatedAt, &v.UpdatedAt); err != nil {
			return nil, 0, err
		}
		venues = append(venues, &v)
	}

	return venues, total, nil
}

func (r *Repository) UpdateVenue(ctx context.Context, id, name, address string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE venues SET name = $1, address = $2, updated_at = NOW() WHERE id = $3`,
		name, address, id)
	return err
}

func (r *Repository) DeleteVenue(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM venues WHERE id = $1`, id)
	return err
}

// Room operations
func (r *Repository) CreateRoom(ctx context.Context, venueID, name string) (string, error) {
	id := uuid.New().String()
	_, err := r.db.Exec(ctx,
		`INSERT INTO rooms (id, venue_id, name, created_at, updated_at)
		 VALUES ($1, $2, $3, NOW(), NOW())`,
		id, venueID, name)
	return id, err
}

func (r *Repository) GetRoom(ctx context.Context, id string) (*Room, error) {
	var room Room
	err := r.db.QueryRow(ctx,
		`SELECT id, venue_id, name, created_at, updated_at
		 FROM rooms WHERE id = $1`, id).
		Scan(&room.ID, &room.VenueID, &room.Name, &room.CreatedAt, &room.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *Repository) ListRooms(ctx context.Context, venueID string, limit, offset int32) ([]*Room, int32, error) {
	var total int32
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM rooms WHERE venue_id = $1`, venueID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, venue_id, name, created_at, updated_at
		 FROM rooms WHERE venue_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		venueID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var rooms []*Room
	for rows.Next() {
		var room Room
		if err := rows.Scan(&room.ID, &room.VenueID, &room.Name, &room.CreatedAt, &room.UpdatedAt); err != nil {
			return nil, 0, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, total, nil
}

func (r *Repository) UpdateRoom(ctx context.Context, id, name string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE rooms SET name = $1, updated_at = NOW() WHERE id = $2`,
		name, id)
	return err
}

func (r *Repository) DeleteRoom(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM rooms WHERE id = $1`, id)
	return err
}

// Table operations
func (r *Repository) CreateTable(ctx context.Context, roomID, name string, capacity int32, canMerge bool, zone string) (string, error) {
	id := uuid.New().String()
	_, err := r.db.Exec(ctx,
		`INSERT INTO tables (id, room_id, name, capacity, can_merge, zone, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		id, roomID, name, capacity, canMerge, zone)
	
	// Invalidate cache
	r.redis.Del(ctx, fmt.Sprintf("layout:%s", roomID))
	
	return id, err
}

func (r *Repository) GetTable(ctx context.Context, id string) (*Table, error) {
	var t Table
	err := r.db.QueryRow(ctx,
		`SELECT id, room_id, name, capacity, can_merge, zone, created_at, updated_at
		 FROM tables WHERE id = $1`, id).
		Scan(&t.ID, &t.RoomID, &t.Name, &t.Capacity, &t.CanMerge, &t.Zone, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repository) ListTables(ctx context.Context, roomID, venueID string, limit, offset int32) ([]*Table, int32, error) {
	var query string
	var args []interface{}
	
	if roomID != "" {
		query = `SELECT COUNT(*) FROM tables WHERE room_id = $1`
		args = []interface{}{roomID}
	} else if venueID != "" {
		query = `SELECT COUNT(*) FROM tables t JOIN rooms r ON t.room_id = r.id WHERE r.venue_id = $1`
		args = []interface{}{venueID}
	} else {
		query = `SELECT COUNT(*) FROM tables`
		args = []interface{}{}
	}
	
	var total int32
	err := r.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	if roomID != "" {
		query = `SELECT id, room_id, name, capacity, can_merge, zone, created_at, updated_at
		 FROM tables WHERE room_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{roomID, limit, offset}
	} else if venueID != "" {
		query = `SELECT t.id, t.room_id, t.name, t.capacity, t.can_merge, t.zone, t.created_at, t.updated_at
		 FROM tables t JOIN rooms r ON t.room_id = r.id WHERE r.venue_id = $1 ORDER BY t.created_at DESC LIMIT $2 OFFSET $3`
		args = []interface{}{venueID, limit, offset}
	} else {
		query = `SELECT id, room_id, name, capacity, can_merge, zone, created_at, updated_at
		 FROM tables ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		args = []interface{}{limit, offset}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var tables []*Table
	for rows.Next() {
		var t Table
		if err := rows.Scan(&t.ID, &t.RoomID, &t.Name, &t.Capacity, &t.CanMerge, &t.Zone, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, 0, err
		}
		tables = append(tables, &t)
	}

	return tables, total, nil
}

func (r *Repository) UpdateTable(ctx context.Context, id, name string, capacity int32, zone string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE tables SET name = $1, capacity = $2, zone = $3, updated_at = NOW() WHERE id = $4`,
		name, capacity, zone, id)
	
	// Invalidate cache
	var roomID string
	r.db.QueryRow(ctx, `SELECT room_id FROM tables WHERE id = $1`, id).Scan(&roomID)
	if roomID != "" {
		r.redis.Del(ctx, fmt.Sprintf("layout:%s", roomID))
	}
	
	return err
}

func (r *Repository) DeleteTable(ctx context.Context, id string) error {
	var roomID string
	r.db.QueryRow(ctx, `SELECT room_id FROM tables WHERE id = $1`, id).Scan(&roomID)
	
	_, err := r.db.Exec(ctx, `DELETE FROM tables WHERE id = $1`, id)
	
	// Invalidate cache
	if roomID != "" {
		r.redis.Del(ctx, fmt.Sprintf("layout:%s", roomID))
	}
	
	return err
}

// Models
type Venue struct {
	ID        string
	Name      string
	Timezone  string
	Address   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Room struct {
	ID        string
	VenueID   string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Table struct {
	ID        string
	RoomID    string
	Name      string
	Capacity  int32
	CanMerge  bool
	Zone      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

