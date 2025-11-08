-- Venue service migrations

-- Venues
CREATE TABLE IF NOT EXISTS venues (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    timezone VARCHAR(50) NOT NULL,
    address TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Rooms
CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(36) PRIMARY KEY,
    venue_id VARCHAR(36) NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rooms_venue_id ON rooms(venue_id);

-- Tables
CREATE TABLE IF NOT EXISTS tables (
    id VARCHAR(36) PRIMARY KEY,
    room_id VARCHAR(36) NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    capacity INTEGER NOT NULL,
    can_merge BOOLEAN NOT NULL DEFAULT FALSE,
    zone VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tables_room_id ON tables(room_id);

-- Opening hours
CREATE TABLE IF NOT EXISTS opening_hours (
    id VARCHAR(36) PRIMARY KEY,
    venue_id VARCHAR(36) NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    weekday INTEGER NOT NULL CHECK (weekday >= 0 AND weekday <= 6),
    open_time TIME NOT NULL,
    close_time TIME NOT NULL,
    UNIQUE(venue_id, weekday)
);

CREATE INDEX IF NOT EXISTS idx_opening_hours_venue_id ON opening_hours(venue_id);

-- Special hours
CREATE TABLE IF NOT EXISTS special_hours (
    id VARCHAR(36) PRIMARY KEY,
    venue_id VARCHAR(36) NOT NULL REFERENCES venues(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    open_time TIME,
    close_time TIME,
    is_closed BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(venue_id, date)
);

CREATE INDEX IF NOT EXISTS idx_special_hours_venue_id ON special_hours(venue_id);
CREATE INDEX IF NOT EXISTS idx_special_hours_date ON special_hours(date);



