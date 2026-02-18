package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/database"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/routes"
	"github.com/joho/godotenv"
)

func main() {

	log.Println("Starting application...")

	// Load .env (only used locally)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get PORT from Render
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // fallback for local
	}

	// Get DATABASE_URL
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// Initialize DB
	ctx := context.Background()
	conn, err := database.InitDB(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	// Run migrations
	if err := database.RunMigrations(ctx, conn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database ready")

	// Register routes BEFORE starting server
	routes.RegisterRoutes(conn)

	log.Printf("Server is running on port %s\n", port)

	// Start server (ONLY ONCE)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
