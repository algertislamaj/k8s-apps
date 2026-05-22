package models

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    connStr := fmt.Sprintf(
        "host=%s port=5432 user=%s password=%s dbname=%s sslmode=disable",
        getEnv("POSTGRES_HOST", "postgres"),
        getEnv("POSTGRES_USER", "vmmanager"),
        getEnv("POSTGRES_PASSWORD", "changeme123"),
        getEnv("POSTGRES_DB", "vmmanager"),
    )

    var err error
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }

    if err = DB.Ping(); err != nil {
        log.Fatal("Database unreachable:", err)
    }

    createTables()
    log.Println("Database connected")
}

func createTables() {
    query := `
    CREATE TABLE IF NOT EXISTS hypervisors (
        id          SERIAL PRIMARY KEY,
        name        VARCHAR(100) NOT NULL,
        type        VARCHAR(20) NOT NULL,   -- 'esxi' or 'hyperv'
        host        VARCHAR(255) NOT NULL,
        username    VARCHAR(100) NOT NULL,
        created_at  TIMESTAMP DEFAULT NOW()
    );

    CREATE TABLE IF NOT EXISTS hypervisor_secrets (
        hypervisor_id   INT REFERENCES hypervisors(id) ON DELETE CASCADE,
        secret_name     VARCHAR(255) NOT NULL,   -- k8s secret name
        PRIMARY KEY (hypervisor_id)
    );
    `
    if _, err := DB.Exec(query); err != nil {
        log.Fatal("Failed to create tables:", err)
    }
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}