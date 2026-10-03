package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceBinary(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dest := filepath.Join(dir, "rei")
	if err := os.WriteFile(src, []byte("NUEVO"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(dest, []byte("VIEJO"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := replaceBinary(src, dest); err != nil {
		t.Fatalf("replaceBinary: %v", err)
	}
	got, _ := os.ReadFile(dest)
	if string(got) != "NUEVO" {
		t.Fatalf("contenido = %q, want NUEVO", got)
	}
	fi, err := os.Stat(dest)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o755 {
		t.Fatalf("permiso = %v, want 0755", fi.Mode().Perm())
	}
	if _, err := os.Stat(dest + ".bak"); !os.IsNotExist(err) {
		t.Fatalf("el backup debió eliminarse tras el éxito: %v", err)
	}
}

func TestReplaceBinarySinPermiso(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root ignora los permisos de escritura")
	}
	dir := t.TempDir()
	dest := filepath.Join(dir, "rei")
	src := filepath.Join(dir, "src")
	if err := os.WriteFile(dest, []byte("VIEJO"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(src, []byte("NUEVO"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })

	if err := replaceBinary(src, dest); err == nil {
		t.Fatal("esperaba error por falta de permisos")
	}
	if err := os.Chmod(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(dest)
	if string(got) != "VIEJO" {
		t.Fatalf("el binario original no debe cambiar, got %q", got)
	}
}

func TestResolveExecutable(t *testing.T) {
	dir := t.TempDir()
	real := filepath.Join(dir, "rei-real")
	if err := os.WriteFile(real, []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "rei")
	if err := os.Symlink(real, link); err != nil {
		t.Fatal(err)
	}
	got, err := resolveExecutable(func() (string, error) { return link, nil })
	if err != nil {
		t.Fatal(err)
	}
	want, _ := filepath.EvalSymlinks(real)
	if got != want {
		t.Fatalf("resolveExecutable = %q, want %q", got, want)
	}
}
