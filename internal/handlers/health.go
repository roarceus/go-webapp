package handlers

import (
	"log"
	"net/http"
	"webapp-hello-world/internal/database"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Check for query parameters in the URL
	if len(r.URL.Query()) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: Parameters are not allowed"))
		return
	}

	// Check for payload (body content)
	if r.ContentLength > 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request: Body is not allowed"))
		return
	}

	// Set headers
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Attempt to connect to DB
	db := database.GetDB()

	// Check if DB is nil (connection failed)
	if db == nil {
		log.Printf("Health check failed: Database connection unavailable")
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
		return
	}

	// Try to insert health check record
	_, err := db.Exec("INSERT INTO webapp.health_checks (checked_at) VALUES (CURRENT_TIMESTAMP)")
	if err != nil {
		// Log the error for debugging purposes
		log.Printf("Health check failed: %v", err)
		// Return 503 Service Unavailable if there is an error
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Service Unavailable"))
		return
	}

	// Return 200 OK if everything works fine
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
