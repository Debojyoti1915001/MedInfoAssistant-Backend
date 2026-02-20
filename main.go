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

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight (OPTIONS request)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {

	log.Println("Starting application...")

	// Load .env (only used locally)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Get PORT (Render provides this automatically)
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

	// Wrap mux with CORS middleware
	handler := enableCORS(http.DefaultServeMux)

	// Start server
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
