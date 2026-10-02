package paths

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindFrom(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".rei"), 0o755); err != nil {
		t.Fatal(err)
	}
	nested := filepath.Join(root, "a", "b")
	if err := os.MkdirAll(nested, 0o755); err != nil {
		t.Fatal(err)
	}
	p, err := FindFrom(nested)
	if err != nil {
		t.Fatalf("FindFrom: %v", err)
	}
	if p.Root != root {
		t.Fatalf("root = %q, want %q", p.Root, root)
	}
}

func TestFindFromNotFound(t *testing.T) {
	dir := t.TempDir()
	if _, err := FindFrom(dir); err == nil {
		t.Fatal("esperaba error cuando no hay .rei/")
	}
}
