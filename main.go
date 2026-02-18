package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/database"
	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	ctx := context.Background()
	conn, err := database.InitDB(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	if err := database.RunMigrations(ctx, conn); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Database ready")

	routes.RegisterRoutes(conn)

	log.Printf("Server is running on port %s\n", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
