package main

import (
    "encoding/json"
    "log"
    "net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "*")
    json.NewEncoder(w).Encode(map[string]string{"message": "Hello World"})
}

func main() {
    http.HandleFunc("/api/hello", helloHandler)
    log.Println("Backend running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}