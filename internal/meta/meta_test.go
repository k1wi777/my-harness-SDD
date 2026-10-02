package meta

import (
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "meta.json")
	m := &Meta{
		ID:          "2026-10-01_10-00__demo",
		Title:       "Demo",
		Description: "desc",
		Type:        "task",
		Status:      StatusPending,
		CreatedAt:   "2026-10-01T10:00:00-05:00",
	}
	if err := m.Save(path); err != nil {
		t.Fatal(err)
	}
	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if *got != *m {
		t.Fatalf("got %+v, want %+v", got, m)
	}
}

func TestIsValid(t *testing.T) {
	if !IsValidStatus(StatusInProgress) {
		t.Error("in_progress debería ser válido")
	}
	if IsValidStatus("nope") {
		t.Error("'nope' no debería ser válido")
	}
	if !IsValidType("feature") || IsValidType("x") {
		t.Error("validación de type incorrecta")
	}
}
