package doctor

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/k1wi777/my-harness-SDD/internal/check"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

// setupValidHarness crea un harness completo y válido en un directorio temporal,
// sin Work Items y sin repositorio git.
func setupValidHarness(t *testing.T) *paths.Project {
	t.Helper()
	root := t.TempDir()
	p := &paths.Project{Root: root}
	for _, f := range check.RequiredFiles {
		full := filepath.Join(root, f)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for _, d := range []string{".rei/specs", ".rei/progress/work-items"} {
		if err := os.MkdirAll(filepath.Join(root, filepath.FromSlash(d)), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	for _, f := range []string{".rei/progress/current.md", ".rei/progress/history.md"} {
		if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(f)), []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return p
}

func TestRunOK(t *testing.T) {
	p := setupValidHarness(t)
	var buf bytes.Buffer
	if code := Run(p, &buf); code != 0 {
		t.Fatalf("Run = %d, salida:\n%s", code, buf.String())
	}
}

func TestRunMissingFile(t *testing.T) {
	root := t.TempDir()
	p := &paths.Project{Root: root}
	var buf bytes.Buffer
	if code := Run(p, &buf); code == 0 {
		t.Fatalf("esperaba fallo por archivos faltantes, salida:\n%s", buf.String())
	}
}

func TestRunInvalidWorkItem(t *testing.T) {
	p := setupValidHarness(t)
	id := "2026-01-01_00-00__bad"
	if err := os.MkdirAll(p.SpecDir(id), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p.MetaFile(id), []byte("{ inválido"), 0o644); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if code := Run(p, &buf); code == 0 {
		t.Fatalf("esperaba fallo por Work Item inválido, salida:\n%s", buf.String())
	}
}
