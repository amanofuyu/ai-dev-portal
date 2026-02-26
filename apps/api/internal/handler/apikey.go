package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"dev-portal/api/internal/model"
)

type ApiKeyHandler struct {
	DB *sql.DB
}

func (h *ApiKeyHandler) ListKeys(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)", projectID).Scan(&exists)
	if err != nil || !exists {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	rows, err := h.DB.Query(`
		SELECT id, key_value, name, is_enabled, last_used_at, project_id, created_at
		FROM api_keys WHERE project_id = ? ORDER BY id
	`, projectID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query keys")
		return
	}
	defer rows.Close()

	keys := []model.ApiKey{}
	for rows.Next() {
		var k model.ApiKey
		var rawKeyValue string
		var createdAt string
		var lastUsedAt sql.NullString
		if err := rows.Scan(&k.ID, &rawKeyValue, &k.Name, &k.IsEnabled, &lastUsedAt, &k.ProjectID, &createdAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan key")
			return
		}
		k.KeyValueMasked = model.MaskKeyValue(rawKeyValue)
		k.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		if lastUsedAt.Valid {
			t, _ := time.Parse("2006-01-02 15:04:05", lastUsedAt.String)
			k.LastUsedAt = &t
		}
		keys = append(keys, k)
	}

	writeJSON(w, http.StatusOK, keys)
}

func (h *ApiKeyHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	var exists bool
	err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = ?)", projectID).Scan(&exists)
	if err != nil || !exists {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	keyValue, err := model.GenerateAPIKey()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate API key")
		return
	}

	res, err := h.DB.Exec(
		"INSERT INTO api_keys (key_value, name, project_id) VALUES (?, ?, ?)",
		keyValue, input.Name, projectID,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create key")
		return
	}

	id, _ := res.LastInsertId()
	var k model.ApiKey
	var createdAt string
	var lastUsedAt sql.NullString
	err = h.DB.QueryRow(`
		SELECT id, key_value, name, is_enabled, last_used_at, project_id, created_at
		FROM api_keys WHERE id = ?
	`, id).Scan(&k.ID, &k.KeyValue, &k.Name, &k.IsEnabled, &lastUsedAt, &k.ProjectID, &createdAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve created key")
		return
	}
	k.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	if lastUsedAt.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", lastUsedAt.String)
		k.LastUsedAt = &t
	}

	writeJSON(w, http.StatusCreated, k)
}

func (h *ApiKeyHandler) UpdateKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var input struct {
		IsEnabled *bool `json:"is_enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if input.IsEnabled == nil {
		writeError(w, http.StatusBadRequest, "is_enabled is required")
		return
	}

	enabled := 0
	if *input.IsEnabled {
		enabled = 1
	}

	res, err := h.DB.Exec("UPDATE api_keys SET is_enabled = ? WHERE id = ?", enabled, id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update key")
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "api key not found")
		return
	}

	var k model.ApiKey
	var rawKeyValue string
	var createdAt string
	var lastUsedAt sql.NullString
	err = h.DB.QueryRow(`
		SELECT id, key_value, name, is_enabled, last_used_at, project_id, created_at
		FROM api_keys WHERE id = ?
	`, id).Scan(&k.ID, &rawKeyValue, &k.Name, &k.IsEnabled, &lastUsedAt, &k.ProjectID, &createdAt)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve updated key")
		return
	}
	k.KeyValueMasked = model.MaskKeyValue(rawKeyValue)
	k.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	if lastUsedAt.Valid {
		t, _ := time.Parse("2006-01-02 15:04:05", lastUsedAt.String)
		k.LastUsedAt = &t
	}

	writeJSON(w, http.StatusOK, k)
}

func (h *ApiKeyHandler) DeleteKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	res, err := h.DB.Exec("DELETE FROM api_keys WHERE id = ?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete key")
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "api key not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ApiKeyHandler) RevealKey(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var revealed model.ApiKeyRevealed
	err := h.DB.QueryRow("SELECT id, key_value FROM api_keys WHERE id = ?", id).Scan(&revealed.ID, &revealed.KeyValue)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "api key not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to query key")
		return
	}

	writeJSON(w, http.StatusOK, revealed)
}
