package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
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

	r.HandleFunc("/register", registerHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/login", loginHandler).Methods("POST", "OPTIONS")
	r.HandleFunc("/generate-qr", generateQRHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/upload/{id}", uploadHandlerGET).Methods("GET", "OPTIONS")
	r.HandleFunc("/upload/{id}", uploadHandlerPOST).Methods("POST", "OPTIONS")
	r.HandleFunc("/user/images", getUserImagesHandler).Methods("GET", "OPTIONS")

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
	session.Options.SameSite = http.SameSiteNoneMode
	session.Options.Secure = true
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
	session.Options.SameSite = http.SameSiteNoneMode
	session.Options.Secure = true
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
	uploadURL := fmt.Sprintf("https://%s/upload/%s", r.Host, id) // Use the actual host

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
func uploadHandlerGET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	_, exists := qrCodes[id]
	mu.Unlock()

	if !exists {
		http.Error(w, "Invalid or expired QR code", http.StatusBadRequest)
		return
	}

	// Serve a simple HTML page for file upload
	html := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>Upload Image</title>
	</head>
	<body>
		<h1>Upload Image</h1>
		<form action="/upload/{{.ID}}" method="post" enctype="multipart/form-data">
			<input type="file" name="file" accept="image/*">
			<button type="submit">Upload</button>
		</form>
	</body>
	</html>
	`

	tmpl, err := template.New("upload").Parse(html)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, struct{ ID string }{ID: id})
}

func uploadHandlerPOST(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	mu.Lock()
	qrCode, exists := qrCodes[id]
	mu.Unlock()

	if !exists || time.Now().After(qrCode.ExpiresAt) {
		http.Error(w, "Invalid or expired QR code", http.StatusBadRequest)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	uploadURL := "https://upload.uploadcare.com/base/"
	apiKey := "3cb90b7ae33a12d4d266"

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	err = writer.WriteField("UPLOADCARE_PUB_KEY", apiKey)
	if err != nil {
		http.Error(w, "Error creating form", http.StatusInternalServerError)
		return
	}
	filePart, err := writer.CreateFormFile("file", fileHeader.Filename)
	if err != nil {
		http.Error(w, "Error creating form file part", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(filePart, file)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	writer.Close()
	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		http.Error(w, "Error creating upload request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Error uploading to Uploadcare", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Error from Uploadcare", resp.StatusCode)
		return
	}

	var uploadResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&uploadResponse); err != nil {
		http.Error(w, "Error decoding upload response", http.StatusInternalServerError)
		return
	}

	imageURL := fmt.Sprintf("https://ucarecdn.com/%s/", uploadResponse["file"])

	// Create the UploadedImage entry
	imageID := uuid.New().String()
	image := UploadedImage{
		ID:     imageID,
		UserID: qrCode.UserID,
		URL:    imageURL,
	}

	// // Here you would typically save the file to a storage service
	// // For this example, we'll just create a dummy UploadedImage
	// imageID := uuid.New().String()
	// image := UploadedImage{
	// 	ID:     imageID,
	// 	UserID: qrCode.UserID,
	// 	URL:    fmt.Sprintf("http://example.com/images/%s", imageID),
	// }

	mu.Lock()
	uploadedImages[qrCode.UserID] = append(uploadedImages[qrCode.UserID], image)
	delete(qrCodes, id)
	mu.Unlock()

	response := map[string]string{"message": "File uploaded successfully", "imageUrl": image.URL}
	w.Header().Set("Content-Type", "application/json")
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
