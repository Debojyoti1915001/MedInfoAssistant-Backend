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
	createTablesSQL := `
	-- Create users table
	CREATE TABLE IF NOT EXISTS users (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		name TEXT NOT NULL,
		phnNumber TEXT NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);

	-- Create doctors table
	CREATE TABLE IF NOT EXISTS doctors (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		accuracy FLOAT NOT NULL,
		name TEXT NOT NULL,
		phnNumber TEXT NOT NULL,
		speciality TEXT NOT NULL,
		username TEXT UNIQUE NOT NULL,
		email TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);

	-- Create prescriptions table
	CREATE TABLE IF NOT EXISTS prescriptions (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		userId BIGINT NOT NULL REFERENCES users(id),
		docId BIGINT NOT NULL REFERENCES doctors(id),
		symptoms TEXT NOT NULL,
		link TEXT
	);

	-- Create items table
	CREATE TABLE IF NOT EXISTS items (
		id BIGSERIAL PRIMARY KEY,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		presId BIGINT NOT NULL REFERENCES prescriptions(id),
		name TEXT NOT NULL,
		type TEXT NOT NULL CHECK (type IN ('med', 'test')),
		aiReasons TEXT,
		docReason TEXT
	);

	-- Create indexes for better query performance
	CREATE INDEX IF NOT EXISTS idx_prescriptions_userId ON prescriptions(userId);
	CREATE INDEX IF NOT EXISTS idx_prescriptions_docId ON prescriptions(docId);
	CREATE INDEX IF NOT EXISTS idx_items_presId ON items(presId);
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	CREATE INDEX IF NOT EXISTS idx_doctors_email ON doctors(email);
	`
	_, err := conn.Exec(ctx, createTablesSQL)
	return err
}
