package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"go-backend/db"
	"go-backend/resolvers"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
	"github.com/skip2/go-qrcode"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type QRCode struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expiresAt"`
	UserID    string    `json:"userId"`
}

type UploadedImage struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	URL    string `json:"url"`
}

var (
	users          = make(map[string]User)
	qrCodes        = make(map[string]QRCode)
	uploadedImages = make(map[string][]UploadedImage)
	store          = sessions.NewCookieStore([]byte("secret-key"))
	mu             sync.Mutex
)

func LoadSchema(schemaPath string) string {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}
	return string(schemaBytes)
}

func main() {
	// Initialize the database connection
	db.InitDB()
	db.CheckTables()

	// Load the GraphQL schema from the file
	schemaPath := "schema/schema.graphql"
	schemaString := LoadSchema(schemaPath)

	// Parse the schema
	schema := graphql.MustParseSchema(schemaString, resolvers.NewResolver())

	// Create a new mux router
	r := mux.NewRouter()

	// CORS configuration
	corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	corsOrigins := handlers.AllowedOrigins([]string{"http://localhost:5173", "https://yourdomain.com"}) // Update with your frontend URLs
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// Set up the GraphQL endpoint with CORS
	graphqlHandler := handlers.CORS(corsHeaders, corsOrigins, corsMethods)(&relay.Handler{Schema: schema})
	r.Handle("/graphql", graphqlHandler).Methods("POST", "OPTIONS")

	// Add a specific handler for OPTIONS requests to the GraphQL endpoint
	r.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	}).Methods("OPTIONS")

	r.HandleFunc("/register", registerHandler).Methods("POST")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/generate-qr", generateQRHandler).Methods("GET")
	r.HandleFunc("/upload/{id}", uploadHandler).Methods("POST")
	r.HandleFunc("/user/images", getUserImagesHandler).Methods("GET")

	// Set up the REST endpoint
	r.HandleFunc("/api/data", GetData).Methods("GET")

	// Apply CORS middleware to the entire router
	corsRouter := handlers.CORS(corsHeaders, corsOrigins, corsMethods)(r)

	// Determine port for HTTP service
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	// Start the server with the CORS-enabled router
	log.Printf("Server is running on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, corsRouter))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.ID = uuid.New().String()
	users[user.ID] = user

	session, _ := store.Get(r, "session")
	session.Values["userID"] = user.ID
	session.Save(r, w)

	json.NewEncoder(w).Encode(user)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// In a real application, you would validate the credentials here
	// For this example, we'll just find a user with the given username
	var user User
	for _, u := range users {
		if u.Username == credentials.Username {
			user = u
			break
		}
	}

	if user.ID == "" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	session, _ := store.Get(r, "session")
	session.Values["userID"] = user.ID
	session.Save(r, w)

	json.NewEncoder(w).Encode(user)
}

func generateQRHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := uuid.New().String()
	expiresAt := time.Now().Add(5 * time.Minute)
	uploadURL := fmt.Sprintf("http://localhost:3000/upload/%s", id)

	qrCode := QRCode{
		ID:        id,
		URL:       uploadURL,
		ExpiresAt: expiresAt,
		UserID:    userID,
	}

	mu.Lock()
	qrCodes[id] = qrCode
	mu.Unlock()

	qr, err := qrcode.Encode(uploadURL, qrcode.Medium, 256)
	if err != nil {
		http.Error(w, "Failed to generate QR code", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(qr)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	qrCode, exists := qrCodes[id]
	mu.Unlock()

	if !exists || time.Now().After(qrCode.ExpiresAt) {
		http.Error(w, "Invalid or expired QR code", http.StatusBadRequest)
		return
	}

	// Handle file upload here
	// For simplicity, we'll just create a dummy UploadedImage
	imageID := uuid.New().String()
	image := UploadedImage{
		ID:     imageID,
		UserID: qrCode.UserID,
		URL:    fmt.Sprintf("http://example.com/images/%s", imageID),
	}

	mu.Lock()
	uploadedImages[qrCode.UserID] = append(uploadedImages[qrCode.UserID], image)
	delete(qrCodes, id)
	mu.Unlock()

	response := map[string]string{"message": "File uploaded successfully", "imageUrl": image.URL}
	json.NewEncoder(w).Encode(response)
}

func getUserImagesHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	userID, ok := session.Values["userID"].(string)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	images := uploadedImages[userID]
	mu.Unlock()

	json.NewEncoder(w).Encode(images)
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello from backend"})
}
