package template

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

func TestRenderFile(t *testing.T) {
	dir := t.TempDir()
	tpl := filepath.Join(dir, "t.md")
	if err := os.WriteFile(tpl, []byte("Hola {{NAME}} y {{OTRO}}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(dir, "out.md")
	if err := RenderFile(tpl, target, map[string]string{"NAME": "mundo"}); err != nil {
		t.Fatal(err)
	}
	data, _ := os.ReadFile(target)
	if got := string(data); !strings.Contains(got, "Hola mundo y {{OTRO}}") {
		t.Fatalf("got %q", got)
	}
}

func TestStartAndArchive(t *testing.T) {
	root := t.TempDir()
	p := &paths.Project{Root: root}
	tpl := p.TemplatesDir()
	if err := os.MkdirAll(tpl, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(p.ProgressDir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tpl, "current.md"), []byte("- **Work Item:** {{WORK_ITEM}}\n- **Estado:** {{STATE}}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(tpl, "history.md"), []byte("# Bitácora\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := StartSession(p, "2026-10-01_10-00__x", "task"); err != nil {
		t.Fatal(err)
	}
	cur, _ := os.ReadFile(p.CurrentFile())
	if !strings.Contains(string(cur), "2026-10-01_10-00__x") {
		t.Fatalf("StartSession: %s", cur)
	}

	id, err := ArchiveSession(p)
	if err != nil {
		t.Fatal(err)
	}
	if id != "2026-10-01_10-00__x" {
		t.Fatalf("ArchiveSession id = %q", id)
	}
	cur, _ = os.ReadFile(p.CurrentFile())
	if strings.Contains(string(cur), "2026-10-01_10-00__x") {
		t.Fatal("current.md no se restableció")
	}
	hist, _ := os.ReadFile(p.HistoryFile())
	if !strings.Contains(string(hist), "2026-10-01_10-00__x") {
		t.Fatal("history.md no contiene la entrada")
	}
}
