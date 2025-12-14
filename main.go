// ... imports above
import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os" // <--- ADD THIS IMPORT
    "sync"
    "time"
)

// ... Handlers and data structures ...

// --- 4. Main Function and Router Setup ---

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/requests", CORSHandler(RequestsHandler))

    // --- UPDATED PORT LOGIC ---
    
    // 1. Get the PORT from the environment variable (Render sets this)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080" // Fallback port for local development
    }
    
    // 2. Go's ListenAndServe requires the port to be prefixed with a colon (e.g., :8080)
    listenAddr := ":" + port
    
    fmt.Printf("API server starting on %s\n", listenAddr)
    
    // 3. Start the server using the dynamic port
    if err := http.ListenAndServe(listenAddr, mux); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
