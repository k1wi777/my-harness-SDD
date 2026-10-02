package check

import (
	"bytes"
	"os"
	"path/filepath"
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
