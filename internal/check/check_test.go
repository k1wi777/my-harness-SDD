package check

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

func TestRunOK(t *testing.T) {
	root := t.TempDir()
	p := &paths.Project{Root: root}
	for _, f := range RequiredFiles {
		full := filepath.Join(root, f)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var buf bytes.Buffer
	if code := Run(p, true, &buf); code != 0 {
		t.Fatalf("Run = %d, salida:\n%s", code, buf.String())
	}
}

func TestRunMissingFile(t *testing.T) {
	root := t.TempDir()
	p := &paths.Project{Root: root}
	var buf bytes.Buffer
	if code := Run(p, true, &buf); code == 0 {
		t.Fatalf("esperaba fallo por archivos faltantes")
	}
}

func writeConfig(t *testing.T, p *paths.Project, body string) {
	t.Helper()
	if err := os.MkdirAll(p.ReiDir(), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(p.ReiDir(), "config.json"), []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestRunDeclaredSinChecks(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	var buf bytes.Buffer
	if code := RunDeclared(p, &buf); code != 0 {
		t.Fatalf("RunDeclared = %d, esperaba 0, salida:\n%s", code, buf.String())
	}
	if !strings.Contains(buf.String(), "Sin checks configurados") {
		t.Fatalf("salida inesperada:\n%s", buf.String())
	}
}

func TestRunDeclaredChecksVacios(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	writeConfig(t, p, `{"checks":[]}`)
	var buf bytes.Buffer
	if code := RunDeclared(p, &buf); code != 0 {
		t.Fatalf("RunDeclared = %d, esperaba 0, salida:\n%s", code, buf.String())
	}
	if !strings.Contains(buf.String(), "Sin checks configurados") {
		t.Fatalf("salida inesperada:\n%s", buf.String())
	}
}

func TestRunDeclaredCheckOK(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	writeConfig(t, p, `{"checks":[{"id":"ok","description":"check-ok","command":["sh","-c","exit 0"]}]}`)
	var buf bytes.Buffer
	if code := RunDeclared(p, &buf); code != 0 {
		t.Fatalf("RunDeclared = %d, esperaba 0, salida:\n%s", code, buf.String())
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") || !strings.Contains(out, "check-ok") {
		t.Fatalf("salida inesperada:\n%s", out)
	}
}

func TestRunDeclaredCheckFail(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	writeConfig(t, p, `{"checks":[{"id":"fail","description":"check-fail","command":["sh","-c","exit 1"]}]}`)
	var buf bytes.Buffer
	if code := RunDeclared(p, &buf); code != 1 {
		t.Fatalf("RunDeclared = %d, esperaba 1, salida:\n%s", code, buf.String())
	}
	out := buf.String()
	if !strings.Contains(out, "[FAIL]") || !strings.Contains(out, "check-fail") {
		t.Fatalf("salida inesperada:\n%s", out)
	}
}

func TestRunDeclaredSinEfectos(t *testing.T) {
	p := &paths.Project{Root: t.TempDir()}
	writeConfig(t, p, `{"checks":[{"id":"ok","description":"check-ok","command":["sh","-c","exit 0"]}]}`)
	var buf bytes.Buffer
	if code := RunDeclared(p, &buf); code != 0 {
		t.Fatalf("RunDeclared = %d, esperaba 0, salida:\n%s", code, buf.String())
	}
	if _, err := os.Stat(p.SpecsDir()); !os.IsNotExist(err) {
		t.Errorf("RunDeclared creó .rei/specs/ (err=%v)", err)
	}
	if _, err := os.Stat(p.CurrentFile()); !os.IsNotExist(err) {
		t.Errorf("RunDeclared creó current.md (err=%v)", err)
	}
}
