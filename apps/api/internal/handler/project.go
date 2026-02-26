package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"dev-portal/api/internal/model"
)

type ProjectHandler struct {
	DB *sql.DB
}

func (h *ProjectHandler) ListProjects(w http.ResponseWriter, r *http.Request) {
	rows, err := h.DB.Query(`
		SELECT p.id, p.name, p.description, p.status,
		       COUNT(k.id) AS key_count,
		       p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN api_keys k ON k.project_id = p.id
		GROUP BY p.id
		ORDER BY p.id
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query projects")
		return
	}
	defer rows.Close()

	projects := []model.Project{}
	for rows.Next() {
		var p model.Project
		var createdAt, updatedAt string
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.KeyCount, &createdAt, &updatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to scan project")
			return
		}
		p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
		p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
		projects = append(projects, p)
	}

	writeJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
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

	res, err := h.DB.Exec(
		"INSERT INTO projects (name, description) VALUES (?, ?)",
		input.Name, input.Description,
	)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			writeError(w, http.StatusBadRequest, "project name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to create project")
		return
	}

	id, _ := res.LastInsertId()
	project, err := h.getProjectByID(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve created project")
		return
	}

	writeJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	project, err := h.getProjectByIDStr(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "project not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to query project")
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var existing model.Project
	var createdAt, updatedAt string
	err := h.DB.QueryRow(`
		SELECT id, name, description, status, created_at, updated_at
		FROM projects WHERE id = ?
	`, id).Scan(&existing.ID, &existing.Name, &existing.Description, &existing.Status, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			writeError(w, http.StatusNotFound, "project not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to query project")
		return
	}

	var input map[string]any
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	setClauses := []string{}
	args := []any{}

	if name, ok := input["name"]; ok {
		nameStr, _ := name.(string)
		nameStr = strings.TrimSpace(nameStr)
		if nameStr == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}
		setClauses = append(setClauses, "name = ?")
		args = append(args, nameStr)
	}
	if desc, ok := input["description"]; ok {
		descStr, _ := desc.(string)
		setClauses = append(setClauses, "description = ?")
		args = append(args, descStr)
	}
	if status, ok := input["status"]; ok {
		statusStr, _ := status.(string)
		if statusStr != "Active" && statusStr != "Archived" {
			writeError(w, http.StatusBadRequest, "invalid status, must be Active or Archived")
			return
		}
		setClauses = append(setClauses, "status = ?")
		args = append(args, statusStr)
	}

	if len(setClauses) == 0 {
		writeError(w, http.StatusBadRequest, "no fields to update")
		return
	}

	setClauses = append(setClauses, "updated_at = datetime('now')")
	args = append(args, id)

	query := "UPDATE projects SET " + strings.Join(setClauses, ", ") + " WHERE id = ?"
	if _, err := h.DB.Exec(query, args...); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			writeError(w, http.StatusBadRequest, "project name already exists")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to update project")
		return
	}

	project, err := h.getProjectByID(existing.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to retrieve updated project")
		return
	}
	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	res, err := h.DB.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}

	affected, _ := res.RowsAffected()
	if affected == 0 {
		writeError(w, http.StatusNotFound, "project not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProjectHandler) getProjectByID(id int64) (*model.Project, error) {
	var p model.Project
	var createdAt, updatedAt string
	err := h.DB.QueryRow(`
		SELECT p.id, p.name, p.description, p.status,
		       COUNT(k.id) AS key_count,
		       p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN api_keys k ON k.project_id = p.id
		WHERE p.id = ?
		GROUP BY p.id
	`, id).Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.KeyCount, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	return &p, nil
}

func (h *ProjectHandler) getProjectByIDStr(idStr string) (*model.Project, error) {
	var p model.Project
	var createdAt, updatedAt string
	err := h.DB.QueryRow(`
		SELECT p.id, p.name, p.description, p.status,
		       COUNT(k.id) AS key_count,
		       p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN api_keys k ON k.project_id = p.id
		WHERE p.id = ?
		GROUP BY p.id
	`, idStr).Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.KeyCount, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	p.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAt)
	p.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAt)
	return &p, nil
}
