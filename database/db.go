package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// InitDB initializes the database connection
func InitDB(ctx context.Context, connString string) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}

// RunMigrations runs database migrations
func RunMigrations(ctx context.Context, conn *pgx.Conn) error {
	// Create users table
	_, _ = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name TEXT,
		phnNumber TEXT,
		email TEXT UNIQUE,
		password TEXT
	);
	`)

	// Create doctors table
	_, _ = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS doctors (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		accuracy FLOAT,
		name TEXT,
		phnNumber TEXT,
		speciality TEXT,
		username TEXT UNIQUE,
		email TEXT UNIQUE,
		password TEXT
	);
	`)

	// Create prescriptions table
	_, _ = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS prescriptions (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		userId BIGINT,
		docId BIGINT,
		symptoms TEXT,
		link TEXT
	);
	`)

	// Create items table
	_, _ = conn.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS items (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		presId BIGINT,
		name TEXT,
		type TEXT,
		aiReasons TEXT,
		docReason TEXT
	);
	`)

	// Add missing columns to existing tables
	columnAddStatements := []string{
		"ALTER TABLE doctors ADD COLUMN IF NOT EXISTS accuracy FLOAT",
		"ALTER TABLE doctors ADD COLUMN IF NOT EXISTS speciality TEXT",
		"ALTER TABLE doctors ADD COLUMN IF NOT EXISTS username TEXT",
		"ALTER TABLE doctors ADD COLUMN IF NOT EXISTS email TEXT",
		"ALTER TABLE prescriptions ADD COLUMN IF NOT EXISTS userId BIGINT",
		"ALTER TABLE prescriptions ADD COLUMN IF NOT EXISTS docId BIGINT",
		"ALTER TABLE prescriptions ADD COLUMN IF NOT EXISTS symptoms TEXT",
		"ALTER TABLE prescriptions ADD COLUMN IF NOT EXISTS link TEXT",
		"ALTER TABLE items ADD COLUMN IF NOT EXISTS presId BIGINT",
		"ALTER TABLE items ADD COLUMN IF NOT EXISTS name TEXT",
		"ALTER TABLE items ADD COLUMN IF NOT EXISTS type TEXT",
		"ALTER TABLE items ADD COLUMN IF NOT EXISTS aiReasons TEXT",
		"ALTER TABLE items ADD COLUMN IF NOT EXISTS docReason TEXT",
	}

	for _, stmt := range columnAddStatements {
		_, _ = conn.Exec(ctx, stmt)
	}

	// Create indexes
	indexStatements := []string{
		"CREATE INDEX IF NOT EXISTS idx_prescriptions_userId ON prescriptions(userId)",
		"CREATE INDEX IF NOT EXISTS idx_prescriptions_docId ON prescriptions(docId)",
		"CREATE INDEX IF NOT EXISTS idx_items_presId ON items(presId)",
		"CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)",
		"CREATE INDEX IF NOT EXISTS idx_doctors_email ON doctors(email)",
	}

	for _, stmt := range indexStatements {
		_, _ = conn.Exec(ctx, stmt)
	}

	return nil
}
