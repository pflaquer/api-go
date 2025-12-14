package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// --- 1. Data Structure ---

// Request represents a single user gig request.
type Request struct {
	ID        int       `json:"id"`
	GigTitle  string    `json:"gig_title"`
	Client    string    `json:"client"`
	Email     string    `json:"email"`
	Details   string    `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

// --- 2. Global State Management ---

// Thread-safe store for all requests. (In-Memory Database)
var (
	requests = []Request{}
	mu       sync.Mutex // Mutex to protect the requests slice from concurrent access
	nextID   = 1
)

// --- 3. Handlers ---

// RequestsHandler handles GET (list all) and POST (create new) requests to /requests.
func RequestsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		listRequests(w, r)
	case "POST":
		createRequest(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// listRequests returns all stored gig requests.
func listRequests(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// Lock the data before reading to ensure thread safety
	mu.Lock()
	defer mu.Unlock()

	if err := json.NewEncoder(w).Encode(requests); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// createRequest handles incoming POST requests to submit a new gig request.
func createRequest(w http.ResponseWriter, r *http.Request) {
	var newRequest Request

	if err := json.NewDecoder(r.Body).Decode(&newRequest); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Basic Validation
	if newRequest.GigTitle == "" || newRequest.Client == "" || newRequest.Email == "" {
		http.Error(w, "Missing required fields (gig_title, client, email)", http.StatusBadRequest)
		return
	}

	// Assign ID and timestamp, and save it thread-safely
	mu.Lock()
	newRequest.ID = nextID
	newRequest.CreatedAt = time.Now()
	requests = append(requests, newRequest)
	nextID++
	mu.Unlock()

	log.Printf("New request created: ID %d, Title: %s", newRequest.ID, newRequest.GigTitle)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newRequest); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// CORSHandler wrapper to add necessary CORS headers.
func CORSHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") 
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// --- 4. Main Function and Router Setup ---

func main() {
	mux := http.NewServeMux()

	// Register the handler with the CORS wrapper
	mux.HandleFunc("/requests", CORSHandler(RequestsHandler))

	port := ":8080"
	fmt.Printf("API server starting on http://localhost%s\n", port)
	
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
