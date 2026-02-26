package db

import (
	"database/sql"
	"fmt"
	"log"

	"dev-portal/api/internal/model"
)

type seedKey struct {
	Name      string
	IsEnabled bool
}

var seedData = []struct {
	Name        string
	Description string
	Status      string
	Keys        []seedKey
}{
	{
		Name:        "Cloud Platform",
		Description: "Cloud infrastructure management",
		Status:      "Active",
		Keys: []seedKey{
			{Name: "Production Key", IsEnabled: true},
			{Name: "Staging Key", IsEnabled: true},
			{Name: "Legacy Key", IsEnabled: false},
			{Name: "Testing Key", IsEnabled: true},
		},
	},
	{
		Name:        "Mobile App",
		Description: "Mobile application backend",
		Status:      "Active",
		Keys: []seedKey{
			{Name: "iOS Key", IsEnabled: true},
			{Name: "Android Key", IsEnabled: true},
			{Name: "Deprecated Key", IsEnabled: false},
		},
	},
}

func Seed(db *sql.DB) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count); err != nil {
		return fmt.Errorf("check projects table: %w", err)
	}
	if count > 0 {
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, p := range seedData {
		res, err := tx.Exec(
			"INSERT INTO projects (name, description, status) VALUES (?, ?, ?)",
			p.Name, p.Description, p.Status,
		)
		if err != nil {
			return fmt.Errorf("insert project %q: %w", p.Name, err)
		}
		projectID, _ := res.LastInsertId()

		for _, k := range p.Keys {
			keyValue, err := model.GenerateAPIKey()
			if err != nil {
				return fmt.Errorf("generate key: %w", err)
			}
			enabled := 0
			if k.IsEnabled {
				enabled = 1
			}
			if _, err := tx.Exec(
				"INSERT INTO api_keys (key_value, name, is_enabled, project_id) VALUES (?, ?, ?, ?)",
				keyValue, k.Name, enabled, projectID,
			); err != nil {
				return fmt.Errorf("insert key %q: %w", k.Name, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	log.Println("Seed data inserted: 2 projects, 7 API keys")
	return nil
}
