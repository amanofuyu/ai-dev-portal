package main

import (
	"log"
	"net/http"
	"os"

	"dev-portal/api/internal/db"
	"dev-portal/api/internal/handler"
	"dev-portal/api/internal/middleware"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7080"
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "dev.db"
	}

	database, err := db.Init(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	if err := db.Seed(database); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}

	ph := &handler.ProjectHandler{DB: database}
	kh := &handler.ApiKeyHandler{DB: database}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /projects", ph.ListProjects)
	mux.HandleFunc("POST /projects", ph.CreateProject)
	mux.HandleFunc("GET /projects/{id}", ph.GetProject)
	mux.HandleFunc("PATCH /projects/{id}", ph.UpdateProject)
	mux.HandleFunc("DELETE /projects/{id}", ph.DeleteProject)

	mux.HandleFunc("GET /projects/{projectId}/keys", kh.ListKeys)
	mux.HandleFunc("POST /projects/{projectId}/keys", kh.CreateKey)
	mux.HandleFunc("PATCH /keys/{id}", kh.UpdateKey)
	mux.HandleFunc("DELETE /keys/{id}", kh.DeleteKey)
	mux.HandleFunc("GET /keys/{id}/reveal", kh.RevealKey)

	wrapped := middleware.CORS(mux)

	log.Printf("API server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, wrapped); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
