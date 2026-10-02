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

// allowedTransitions describe las transiciones estándar de workflow.md.
// Cualquier estado puede pasar a blocked.
var allowedTransitions = map[string][]string{
	StatusPending:          {StatusReady},
	StatusReady:            {StatusInProgress},
	StatusInProgress:       {StatusReview},
	StatusReview:           {StatusDone, StatusChangesRequested},
	StatusChangesRequested: {StatusInProgress},
}

// TransitionAllowed indica si from -> to es una transición estándar.
func TransitionAllowed(from, to string) bool {
	if to == StatusBlocked {
		return true
	}
	for _, s := range allowedTransitions[from] {
		if s == to {
			return true
		}
	}
	return false
}

// SetStatus actualiza el status de un meta.json preservando el resto de campos.
// Devuelve el estado anterior.
func SetStatus(path, status string) (string, error) {
	m, err := Load(path)
	if err != nil {
		return "", err
	}
	previous := m.Status
	m.Status = status
	if err := m.Save(path); err != nil {
		return "", err
	}
	return previous, nil
}
