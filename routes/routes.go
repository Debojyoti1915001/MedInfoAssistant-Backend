package routes

import (
	"net/http"

	"github.com/Debojyoti1915001/MedInfoAssistant-Backend/handlers"
	"github.com/jackc/pgx/v5"
)

func RegisterRoutes(db *pgx.Conn) {
	// Health check route
	http.HandleFunc("/health", handlers.HealthHandler)

	// Authentication check route (no auth required - checks if token is valid)
	http.HandleFunc("/api/auth/check", handlers.AuthCheckHandler(db))

	// User routes
	http.HandleFunc("/api/users", handlers.GetUsersHandler(db))
	http.HandleFunc("/api/users/create", handlers.CreateUserHandler(db))
	http.HandleFunc("/api/users/login", handlers.LoginUserHandler(db))
	http.HandleFunc("/api/users/profile", handlers.AuthMiddleware(handlers.UserProfileHandler(db)))

	// Doctor routes
	http.HandleFunc("/api/doctors", handlers.GetDoctorsHandler(db))
	http.HandleFunc("/api/doctors/create", handlers.CreateDoctorHandler(db))
	http.HandleFunc("/api/doctors/get", handlers.GetDoctorHandler(db))
	http.HandleFunc("/api/doctors/login", handlers.LoginDoctorHandler(db))
	http.HandleFunc("/api/doctors/profile", handlers.AuthMiddleware(handlers.DoctorProfileHandler(db)))
	// Prescription routes
	http.HandleFunc("/api/prescriptions", handlers.GetUserPrescriptionsHandler(db))
	http.HandleFunc("/api/prescriptions/create", handlers.CreatePrescriptionHandler(db))
	http.HandleFunc("/api/prescriptions/get", handlers.GetPrescriptionHandler(db))
	http.HandleFunc("/api/prescriptions/with-items", handlers.GetUserPrescriptionsWithItemsHandler(db))
	http.HandleFunc("/api/prescriptions/seen/update", handlers.UpdatePrescriptionSeenStatusHandler(db))
	http.HandleFunc("/api/doctors/prescriptions-with-items", handlers.GetDoctorPrescriptionsWithItemsHandler(db))

	// Items routes
	http.HandleFunc("/api/items", handlers.GetPrescriptionItemsHandler(db))
	http.HandleFunc("/api/items/create", handlers.CreateItemHandler(db))
	http.HandleFunc("/api/items/get", handlers.GetItemHandler(db))
	http.HandleFunc("/api/items/update", handlers.UpdateItemDocReasonHandler(db))
}
