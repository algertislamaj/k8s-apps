package handlers

import (
    "encoding/json"
    "net/http"
)

func Hello(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(map[string]string{"message": "Hello World"})
}