package meta

import (
	"encoding/json"
	"os"
)

// Estados válidos de un Work Item (ver workflow.md).
const (
	StatusPending          = "pending"
	StatusReady            = "ready"
	StatusInProgress       = "in_progress"
	StatusReview           = "review"
	StatusDone             = "done"
	StatusBlocked          = "blocked"
	StatusChangesRequested = "changes_requested"
)

// Meta representa .rei/specs/<id>/meta.json.
type Meta struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Description      string `json:"description"`
	Type             string `json:"type"`
	Status           string `json:"status"`
	CreatedAt        string `json:"created_at"`
	BaseCommit       string `json:"base_commit"`
	LastReviewCommit string `json:"last_review_commit"`
}

// Load lee y parsea un meta.json.
func Load(path string) (*Meta, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var m Meta
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

// Save escribe el meta.json preservando el orden de campos.
func (m *Meta) Save(path string) error {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// ValidStatuses contiene todos los estados admitidos.
var ValidStatuses = []string{
	StatusPending, StatusReady, StatusInProgress,
	StatusReview, StatusDone, StatusBlocked, StatusChangesRequested,
}

// ValidTypes contiene los tipos admitidos.
var ValidTypes = []string{"feature", "task"}

func IsValidStatus(s string) bool {
	for _, v := range ValidStatuses {
		if v == s {
			return true
		}
	}
	return false
}

func IsValidType(s string) bool {
	for _, v := range ValidTypes {
		if v == s {
			return true
		}
	}
	return false
}
