package model

import (
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"
)

type Project struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	KeyCount    int       `json:"key_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ApiKey struct {
	ID             int64      `json:"id"`
	KeyValue       string     `json:"key_value,omitempty"`
	KeyValueMasked string     `json:"key_value_masked,omitempty"`
	Name           string     `json:"name"`
	IsEnabled      bool       `json:"is_enabled"`
	LastUsedAt     *time.Time `json:"last_used_at"`
	ProjectID      int64      `json:"project_id"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ApiKeyRevealed struct {
	ID       int64  `json:"id"`
	KeyValue string `json:"key_value"`
}

func GenerateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "sk-" + hex.EncodeToString(bytes), nil
}

func MaskKeyValue(key string) string {
	prefix := "sk-"
	if !strings.HasPrefix(key, prefix) || len(key) < len(prefix)+8 {
		return "sk-******"
	}
	body := key[len(prefix):]
	return prefix + body[:4] + "******" + body[len(body)-4:]
}
