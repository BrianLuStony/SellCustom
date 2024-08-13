package main

import (
    "encoding/json"
    "net/http"
)

func main() {
    http.HandleFunc("/api/data", func(w http.ResponseWriter, r *http.Request) {
        data := map[string]string{"message": "Hello from Go!"}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(data)
    })

    http.ListenAndServe(":8080", nil)
}

