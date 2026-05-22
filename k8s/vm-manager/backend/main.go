package main

import (
    "log"
    "net/http"
    "vm-manager/handlers"

    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/api/hello", handlers.Hello).Methods("GET")

    // Hypervisors
    r.HandleFunc("/api/hypervisors", handlers.ListHypervisors).Methods("GET")
    r.HandleFunc("/api/hypervisors", handlers.AddHypervisor).Methods("POST")
    r.HandleFunc("/api/hypervisors/{id}", handlers.DeleteHypervisor).Methods("DELETE")

    // VMs
    r.HandleFunc("/api/hypervisors/{id}/vms", handlers.ListVMs).Methods("GET")
    r.HandleFunc("/api/vms/{id}/cpu", handlers.UpdateCPU).Methods("PUT")
    r.HandleFunc("/api/vms/{id}/memory", handlers.UpdateMemory).Methods("PUT")

    log.Println("Backend running on :8080")
    log.Fatal(http.ListenAndServe(":8080", r))
}