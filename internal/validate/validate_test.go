package validate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/k1wi777/my-harness-SDD/internal/meta"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

const testID = "2026-10-01_10-00__demo"

func writeMeta(t *testing.T, p *paths.Project, m *meta.Meta) {
	t.Helper()
	if err := os.MkdirAll(p.SpecDir(m.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(p.WorkItemDir(m.ID), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := m.Save(p.MetaFile(m.ID)); err != nil {
		t.Fatal(err)
	}
}

func TestFeatureReadyMissingDocs(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	writeMeta(t, p, &meta.Meta{
		ID: testID, Title: "t", Description: "d",
		Type: "feature", Status: meta.StatusReady, CreatedAt: "2026-10-01T10:00:00-05:00",
	})
	issues, err := WorkItem(p, testID)
	if err != nil {
		t.Fatal(err)
	}
	if countLevel(issues, LevelFail) == 0 {
		t.Fatal("esperaba FAIL por documentos de planificación faltantes")
	}
}

func TestTaskDoneOK(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	writeMeta(t, p, &meta.Meta{
		ID: testID, Title: "t", Description: "d",
		Type: "task", Status: meta.StatusDone, CreatedAt: "2026-10-01T10:00:00-05:00",
	})
	mustWrite(t, filepath.Join(p.SpecDir(testID), "plan.md"), "- [x] 1\n")
	mustWrite(t, filepath.Join(p.WorkItemDir(testID), "impl.md"), "impl\n")
	mustWrite(t, filepath.Join(p.WorkItemDir(testID), "review.md"), "review\n")

	issues, err := WorkItem(p, testID)
	if err != nil {
		t.Fatal(err)
	}
	if n := countLevel(issues, LevelFail); n != 0 {
		t.Fatalf("esperaba 0 FAIL, got %d: %+v", n, issues)
	}
}

func TestEmptyDescriptionFails(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	writeMeta(t, p, &meta.Meta{
		ID: testID, Title: "t", Type: "task", Status: meta.StatusPending,
		CreatedAt: "2026-10-01T10:00:00-05:00",
	})
	issues, _ := WorkItem(p, testID)
	if countLevel(issues, LevelFail) == 0 {
		t.Fatal("esperaba FAIL por description vacía")
	}
}

func TestReportTooLargeWarns(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	writeMeta(t, p, &meta.Meta{
		ID: testID, Title: "t", Description: "d",
		Type: "task", Status: meta.StatusReview, CreatedAt: "2026-10-01T10:00:00-05:00",
	})
	mustWrite(t, filepath.Join(p.SpecDir(testID), "plan.md"), "- [x] 1\n")
	big := strings.Repeat("palabra ", maxReportWords+1)
	mustWrite(t, filepath.Join(p.WorkItemDir(testID), "impl.md"), big)

	issues, err := WorkItem(p, testID)
	if err != nil {
		t.Fatal(err)
	}
	if n := countLevel(issues, LevelFail); n != 0 {
		t.Fatalf("esperaba 0 FAIL, got %d: %+v", n, issues)
	}
	if countLevel(issues, LevelWarn) == 0 {
		t.Fatalf("esperaba WARN por impl.md demasiado grande: %+v", issues)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func countLevel(issues []Issue, level Level) int {
	n := 0
	for _, is := range issues {
		if is.Level == level {
			n++
		}
	}
	return n
}
