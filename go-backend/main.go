package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"
)

func main() {

	schemaData, err := ioutil.ReadFile("schema.graphql")
	if err != nil {
		fmt.Println("Error reading schema:", err)
		return
	}

	// // Parse the schema
	// schema, err := graphql.ParseSchema(string(schemaData), nil)
	// if err != nil {
	// 	fmt.Println("Error parsing schema:", err)
	// 	return
	// }

	initDB()
	defer db.Close()

	parsedSchema := graphql.MustParseSchema(string(schemaData), &Resolver{})

	http.Handle("/graphql", &relay.Handler{Schema: parsedSchema})

	log.Println("Server is running on http://localhost:8080/graphql")
	log.Fatal(http.ListenAndServe(":8080", nil))

	r := mux.NewRouter()
	r.HandleFunc("/api/data", GetData).Methods("GET")

	// CORS configuration
	corsHeaders := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	corsOrigins := handlers.AllowedOrigins([]string{"*"}) // Allow all origins or specify your frontend URL
	corsMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	http.ListenAndServe(":8080", handlers.CORS(corsHeaders, corsOrigins, corsMethods)(r))
}

func GetData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello from backend"})
}
