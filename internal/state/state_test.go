package state

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

func TestReadSessionActive(t *testing.T) {
	root := t.TempDir()
	progress := filepath.Join(root, ".rei", "progress")
	if err := os.MkdirAll(progress, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Sesión actual\n\n" +
		"- **Work Item:** 2026-10-01_10-00__demo\n" +
		"- **Tipo:** task\n" +
		"- **Estado:** in_progress\n" +
		"- **Inicio:** x\n" +
		"- **Agente activo:** implementer\n"
	if err := os.WriteFile(filepath.Join(progress, "current.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	p := &paths.Project{Root: root}
	s, err := ReadSession(p)
	if err != nil {
		t.Fatal(err)
	}
	if !s.Active() {
		t.Fatal("la sesión debería estar activa")
	}
	if s.WorkItem != "2026-10-01_10-00__demo" || s.State != "in_progress" || s.Agent != "implementer" {
		t.Fatalf("sesión inesperada: %+v", s)
	}
}

func TestReadSessionIdle(t *testing.T) {
	root := t.TempDir()
	progress := filepath.Join(root, ".rei", "progress")
	if err := os.MkdirAll(progress, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "- **Work Item:** _ninguno_\n- **Estado:** _—_ \n"
	if err := os.WriteFile(filepath.Join(progress, "current.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	s, _ := ReadSession(&paths.Project{Root: root})
	if s.Active() {
		t.Fatal("no debería haber sesión activa")
	}
}
