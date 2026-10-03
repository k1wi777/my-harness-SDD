package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// buildTarGz construye en memoria un .tar.gz con un único fichero.
func buildTarGz(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	hdr := &tar.Header{Name: name, Mode: 0o755, Size: int64(len(content)), Typeflag: tar.TypeReg}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// buildZip construye en memoria un .zip con un único fichero.
func buildZip(t *testing.T, name string, content []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(name)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := w.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestVerifySHA256(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "asset.tar.gz")
	data := []byte("contenido")
	if err := os.WriteFile(p, data, 0o644); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	if err := verifySHA256(p, []byte(fmt.Sprintf("%x  asset.tar.gz\n", sum)), "asset.tar.gz"); err != nil {
		t.Fatalf("checksum correcto: %v", err)
	}
	bad := sha256.Sum256([]byte("otro"))
	if err := verifySHA256(p, []byte(fmt.Sprintf("%x  asset.tar.gz\n", bad)), "asset.tar.gz"); err == nil {
		t.Fatal("esperaba fallo por checksum incorrecto")
	}
	if err := verifySHA256(p, []byte("abc  otro\n"), "asset.tar.gz"); err == nil {
		t.Fatal("esperaba fallo por entrada ausente")
	}
}

func TestParseChecksums(t *testing.T) {
	sums := parseChecksums([]byte("AAA  rei_0.2.0_linux_amd64.tar.gz\nbbb *otro\n\n"))
	if sums["rei_0.2.0_linux_amd64.tar.gz"] != "aaa" {
		t.Errorf("no parseó el checksum: %v", sums)
	}
	if sums["otro"] != "bbb" {
		t.Errorf("no parseó la entrada con '*': %v", sums)
	}
}

func TestExtractBinaryTarGz(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "rei.tar.gz")
	if err := os.WriteFile(archive, buildTarGz(t, "rei", []byte("NUEVO")), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "out")
	if err := extractBinary(archive, dst, "linux"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != "NUEVO" {
		t.Fatalf("contenido extraído = %q", got)
	}
}

func TestExtractBinaryZip(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "rei.zip")
	if err := os.WriteFile(archive, buildZip(t, "rei.exe", []byte("WIN")), 0o644); err != nil {
		t.Fatal(err)
	}
	dst := filepath.Join(dir, "out.exe")
	if err := extractBinary(archive, dst, "windows"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(dst)
	if string(got) != "WIN" {
		t.Fatalf("contenido extraído = %q", got)
	}
}

func TestExtractBinaryPathTraversalIgnorado(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "rei.tar.gz")
	if err := os.WriteFile(archive, buildTarGz(t, "../../rei", []byte("X")), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := extractBinary(archive, filepath.Join(dir, "out"), "linux"); err == nil {
		t.Fatal("esperaba error: la entrada con '..' debe ignorarse")
	}
}

func TestExtractBinaryFormatoDesconocido(t *testing.T) {
	if err := extractBinary("rei.bin", "out", "linux"); err == nil {
		t.Fatal("esperaba error por formato no soportado")
	}
}
