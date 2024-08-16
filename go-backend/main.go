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

	// Set up the GraphQL endpoint
	r.Handle("/graphql", &relay.Handler{Schema: schema}).Methods("POST")

	// Set up the REST endpoint
	r.HandleFunc("/api/data", GetData).Methods("GET")

	// CORS configuration
	corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	corsOrigins := handlers.AllowedOrigins([]string{"*"}) // Allow all origins or specify your frontend URL
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	// Start the server with CORS middleware
	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(corsHeaders, corsOrigins, corsMethods)(r)))
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello from backend"})
}
