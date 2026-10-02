package validate

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/k1wi777/my-harness-SDD/internal/meta"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
	"github.com/k1wi777/my-harness-SDD/internal/state"
)

// Level clasifica un hallazgo.
type Level string

const (
	LevelFail Level = "FAIL"
	LevelWarn Level = "WARN"
)

// Issue es un hallazgo de la validación.
type Issue struct {
	Level   Level
	Message string
}

// WorkItem comprueba la consistencia interna de un Work Item.
func WorkItem(p *paths.Project, id string) ([]Issue, error) {
	var issues []Issue
	fail := func(format string, a ...any) {
		issues = append(issues, Issue{LevelFail, fmt.Sprintf(format, a...)})
	}
	warn := func(format string, a ...any) {
		issues = append(issues, Issue{LevelWarn, fmt.Sprintf(format, a...)})
	}

	metaPath := p.MetaFile(id)
	if !fileExists(metaPath) {
		fail("no existe %s", metaPath)
		return issues, nil
	}

	m, err := meta.Load(metaPath)
	if err != nil {
		fail("meta.json inválido: %v", err)
		return issues, nil
	}

	for _, kv := range []struct{ name, val string }{
		{"id", m.ID},
		{"title", m.Title},
		{"description", m.Description},
		{"type", m.Type},
		{"status", m.Status},
		{"created_at", m.CreatedAt},
	} {
		if strings.TrimSpace(kv.val) == "" {
			fail("campo vacío en meta.json: %s", kv.name)
		}
	}

	if m.ID != "" && m.ID != id {
		fail("id de meta.json ('%s') no coincide con la carpeta ('%s')", m.ID, id)
	}
	if m.Type != "" && !meta.IsValidType(m.Type) {
		fail("type inválido: '%s'", m.Type)
	}
	if m.Status != "" && !meta.IsValidStatus(m.Status) {
		fail("status inválido: '%s'", m.Status)
	}
	if m.CreatedAt != "" {
		if prefix := createdAtPrefix(id); prefix != "" && !strings.HasPrefix(m.CreatedAt, prefix) {
			warn("created_at ('%s') no coincide con el prefijo del id (%s...)", m.CreatedAt, prefix)
		}
	}

	if fi, err := os.Stat(p.WorkItemDir(id)); err != nil || !fi.IsDir() {
		warn("no existe %s", p.WorkItemDir(id))
	}

	if needsPlanning(m.Status) {
		switch m.Type {
		case "feature":
			for _, f := range []string{"requirements.md", "design.md", "tasks.md"} {
				if !fileExists(filepath.Join(p.SpecDir(id), f)) {
					fail("falta %s (feature en estado %s)", f, m.Status)
				}
			}
		case "task":
			if !fileExists(filepath.Join(p.SpecDir(id), "plan.md")) {
				fail("falta plan.md (task en estado %s)", m.Status)
			}
		}
	}

	switch m.Status {
	case meta.StatusReview, meta.StatusDone, meta.StatusChangesRequested:
		if !fileExists(filepath.Join(p.WorkItemDir(id), "impl.md")) {
			fail("falta impl.md (estado %s)", m.Status)
		}
	}
	if m.Status == meta.StatusDone && !fileExists(filepath.Join(p.WorkItemDir(id), "review.md")) {
		fail("falta review.md (estado done)")
	}

	if s, err := state.ReadSession(p); err == nil && s.WorkItem == id {
		if s.State != "" && s.State != m.Status {
			warn("current.md dice Estado='%s' pero meta.json dice status='%s'", s.State, m.Status)
		}
	}

	return issues, nil
}

func needsPlanning(status string) bool {
	switch status {
	case meta.StatusReady, meta.StatusInProgress, meta.StatusReview,
		meta.StatusChangesRequested, meta.StatusDone:
		return true
	}
	return false
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}

// createdAtPrefix deriva "YYYY-MM-DDTHH:mm" del id "YYYY-MM-DD_HH-mm__slug".
func createdAtPrefix(id string) string {
	prefix := id
	if i := strings.Index(id, "__"); i >= 0 {
		prefix = id[:i]
	}
	i := strings.Index(prefix, "_")
	if i < 0 {
		return ""
	}
	datePart := prefix[:i]
	timePart := prefix[i+1:]
	if datePart == "" || timePart == "" {
		return ""
	}
	return datePart + "T" + strings.Replace(timePart, "-", ":", 1)
}
