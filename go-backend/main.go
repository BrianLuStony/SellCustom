package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"go-backend/db"
	"go-backend/resolvers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
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

	// Load the GraphQL schema from the file
	schemaPath := "schema/schema.graphql"
	schemaString := LoadSchema(schemaPath)

	// Parse the schema
	schema := graphql.MustParseSchema(schemaString, &resolvers.Resolver{})

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

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello from backend"})
}
