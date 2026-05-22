package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
)

func ListHypervisors(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode([]string{})
}

func AddHypervisor(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "not implemented yet"})
}

func DeleteHypervisor(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = mux.Vars(r)["id"]
    json.NewEncoder(w).Encode(map[string]string{"status": "not implemented yet"})
}

func ListVMs(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = mux.Vars(r)["id"]
    json.NewEncoder(w).Encode([]string{})
}

func UpdateCPU(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = mux.Vars(r)["id"]
    json.NewEncoder(w).Encode(map[string]string{"status": "not implemented yet"})
}

func UpdateMemory(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    _ = mux.Vars(r)["id"]
    json.NewEncoder(w).Encode(map[string]string{"status": "not implemented yet"})
}