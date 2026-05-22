package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"vm-manager/models"

	"github.com/gorilla/mux"
)

type HypervisorInput struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func ListHypervisors(w http.ResponseWriter, r *http.Request) {
	rows, err := models.DB.Query(`
        SELECT id, name, type, host, username, created_at FROM hypervisors
    `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var hypervisors []map[string]interface{}
	for rows.Next() {
		var id int
		var name, hvType, host, username, createdAt string
		rows.Scan(&id, &name, &hvType, &host, &username, &createdAt)
		hypervisors = append(hypervisors, map[string]interface{}{
			"id":         id,
			"name":       name,
			"type":       hvType,
			"host":       host,
			"username":   username,
			"created_at": createdAt,
		})
	}

	if hypervisors == nil {
		hypervisors = []map[string]interface{}{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hypervisors)
}

func AddHypervisor(w http.ResponseWriter, r *http.Request) {
	var input HypervisorInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	// Store hypervisor in DB
	var id int
	err := models.DB.QueryRow(`
        INSERT INTO hypervisors (name, type, host, username)
        VALUES ($1, $2, $3, $4) RETURNING id
    `, input.Name, input.Type, input.Host, input.Username).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store password in k8s secret
	secretName := fmt.Sprintf("hypervisor-%d", id)
	if err := createK8sSecret(secretName, input.Password); err != nil {
		// rollback DB insert
		models.DB.Exec(`DELETE FROM hypervisors WHERE id = $1`, id)
		http.Error(w, "failed to store credentials", http.StatusInternalServerError)
		return
	}

	models.DB.Exec(`
        INSERT INTO hypervisor_secrets (hypervisor_id, secret_name)
        VALUES ($1, $2)
    `, id, secretName)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id, "status": "created"})
}

func DeleteHypervisor(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var secretName string
	err := models.DB.QueryRow(`
        SELECT secret_name FROM hypervisor_secrets WHERE hypervisor_id = $1
    `, id).Scan(&secretName)
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete k8s secret
	if secretName != "" {
		deleteK8sSecret(secretName)
	}

	models.DB.Exec(`DELETE FROM hypervisors WHERE id = $1`, id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
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
