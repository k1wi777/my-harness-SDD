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

func TestTransitionAllowed(t *testing.T) {
	allowed := []struct{ from, to string }{
		{StatusPending, StatusReady},
		{StatusReady, StatusInProgress},
		{StatusInProgress, StatusReview},
		{StatusReview, StatusDone},
		{StatusReview, StatusChangesRequested},
		{StatusChangesRequested, StatusInProgress},
	}
	for _, tc := range allowed {
		if !TransitionAllowed(tc.from, tc.to) {
			t.Errorf("%s -> %s debería estar permitida", tc.from, tc.to)
		}
	}

	rejected := []struct{ from, to string }{
		{StatusPending, StatusInProgress},
		{StatusDone, StatusInProgress},
		{StatusBlocked, StatusInProgress},
		{StatusReady, StatusDone},
	}
	for _, tc := range rejected {
		if TransitionAllowed(tc.from, tc.to) {
			t.Errorf("%s -> %s NO debería estar permitida", tc.from, tc.to)
		}
	}

	for _, from := range ValidStatuses {
		if !TransitionAllowed(from, StatusBlocked) {
			t.Errorf("%s -> blocked debería estar permitida", from)
		}
	}
}

func TestSetStatusPreservesFields(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "meta.json")
	m := &Meta{
		ID:               "2026-10-02_00-33__demo",
		Title:            "Demo",
		Description:      "desc",
		Type:             "task",
		Status:           StatusPending,
		CreatedAt:        "2026-10-02T00:33:00-05:00",
		BaseCommit:       "726a658661373a8caaec8134187061213739cb46",
		LastReviewCommit: "deadbeef",
	}
	if err := m.Save(path); err != nil {
		t.Fatal(err)
	}

	prev, err := SetStatus(path, StatusReady)
	if err != nil {
		t.Fatal(err)
	}
	if prev != StatusPending {
		t.Fatalf("estado anterior = %q, want %q", prev, StatusPending)
	}

	got, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != StatusReady {
		t.Fatalf("status = %q, want %q", got.Status, StatusReady)
	}
	if got.ID != m.ID || got.Title != m.Title || got.Description != m.Description ||
		got.Type != m.Type || got.CreatedAt != m.CreatedAt ||
		got.BaseCommit != m.BaseCommit || got.LastReviewCommit != m.LastReviewCommit {
		t.Fatalf("campos no preservados: got %+v, want %+v", got, m)
	}
}
