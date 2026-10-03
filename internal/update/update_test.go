package update

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// releaseFixture simula la API de GitHub y la descarga de assets con httptest.
type releaseFixture struct {
	server *httptest.Server
	tag    string
	files  map[string][]byte
}

func newReleaseFixture(t *testing.T, tag string, files map[string][]byte) *releaseFixture {
	t.Helper()
	f := &releaseFixture{tag: tag, files: files}
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/"+repoPath+"/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		type ga struct {
			Name string `json:"name"`
			URL  string `json:"browser_download_url"`
		}
		assets := make([]ga, 0, len(f.files))
		for name := range f.files {
			assets = append(assets, ga{Name: name, URL: f.server.URL + "/dl/" + name})
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"tag_name": f.tag, "assets": assets})
	})
	mux.HandleFunc("/dl/", func(w http.ResponseWriter, r *http.Request) {
		data, ok := f.files[strings.TrimPrefix(r.URL.Path, "/dl/")]
		if !ok {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write(data)
	})
	f.server = httptest.NewServer(mux)
	t.Cleanup(f.server.Close)
	return f
}

// validReleaseFiles construye un asset válido + checksums.txt para un tag/SO.
func validReleaseFiles(t *testing.T, tag, goos, goarch string) map[string][]byte {
	t.Helper()
	name := assetName(strings.TrimPrefix(tag, "v"), goos, goarch)
	var archive []byte
	if goos == "windows" {
		archive = buildZip(t, "rei.exe", []byte("NUEVO"))
	} else {
		archive = buildTarGz(t, "rei", []byte("NUEVO"))
	}
	sum := sha256.Sum256(archive)
	return map[string][]byte{
		name:            archive,
		"checksums.txt": []byte(fmt.Sprintf("%x  %s\n", sum, name)),
	}
}

func linuxOpts(fx *releaseFixture, out *bytes.Buffer, current string) Options {
	return Options{
		CurrentVersion: current,
		Stdout:         out,
		APIBase:        fx.server.URL,
		Client:         fx.server.Client(),
		GOOS:           "linux",
		GOARCH:         "amd64",
	}
}

func TestRunUpToDate(t *testing.T) {
	fx := newReleaseFixture(t, "v0.1.0", nil)
	var buf bytes.Buffer
	opts := linuxOpts(fx, &buf, "v0.1.0")
	if code := Run(opts); code != 0 {
		t.Fatalf("Run up-to-date = %d, want 0", code)
	}
	if !strings.Contains(buf.String(), "última versión") {
		t.Fatalf("salida inesperada: %q", buf.String())
	}
}

func TestRunCheckNoEscribe(t *testing.T) {
	fx := newReleaseFixture(t, "v0.2.0", validReleaseFiles(t, "v0.2.0", "linux", "amd64"))
	dest := filepath.Join(t.TempDir(), "rei")
	if err := os.WriteFile(dest, []byte("VIEJO"), 0o755); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	opts := linuxOpts(fx, &buf, "v0.1.0")
	opts.Check = true
	opts.Executable = func() (string, error) { return dest, nil }
	if code := Run(opts); code != 0 {
		t.Fatalf("Run --check = %d, want 0", code)
	}
	if !strings.Contains(buf.String(), "Última release disponible") {
		t.Fatalf("salida inesperada: %q", buf.String())
	}
	if got, _ := os.ReadFile(dest); string(got) != "VIEJO" {
		t.Fatalf("--check no debe escribir el binario, got %q", got)
	}
}

func TestRunDevInstala(t *testing.T) {
	fx := newReleaseFixture(t, "v0.2.0", validReleaseFiles(t, "v0.2.0", "linux", "amd64"))
	dest := filepath.Join(t.TempDir(), "rei")
	if err := os.WriteFile(dest, []byte("VIEJO"), 0o755); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	opts := linuxOpts(fx, &buf, "dev")
	opts.Executable = func() (string, error) { return dest, nil }
	if code := Run(opts); code != 0 {
		t.Fatalf("Run dev = %d, want 0; salida: %q", code, buf.String())
	}
	if !strings.Contains(buf.String(), "AVISO") {
		t.Fatalf("dev debe avisar: %q", buf.String())
	}
	if got, _ := os.ReadFile(dest); string(got) != "NUEVO" {
		t.Fatalf("binario no actualizado, got %q", got)
	}
}

func TestRunWindowsAbreNavegador(t *testing.T) {
	fx := newReleaseFixture(t, "v0.2.0", nil)
	var opened string
	var buf bytes.Buffer
	opts := linuxOpts(fx, &buf, "v0.1.0")
	opts.GOOS = "windows"
	opts.OpenBrowser = func(u string) error { opened = u; return nil }
	if code := Run(opts); code != 0 {
		t.Fatalf("Run windows = %d, want 0", code)
	}
	if opened != releasesPage {
		t.Fatalf("OpenBrowser(%q), want %q", opened, releasesPage)
	}
	if !strings.Contains(buf.String(), releasesPage) {
		t.Fatalf("debe imprimir la URL de releases: %q", buf.String())
	}
}

func TestRunCheckIntegridadFallida(t *testing.T) {
	files := validReleaseFiles(t, "v0.2.0", "linux", "amd64")
	name := assetName("0.2.0", "linux", "amd64")
	files["checksums.txt"] = []byte("deadbeef  " + name + "\n")
	fx := newReleaseFixture(t, "v0.2.0", files)

	dest := filepath.Join(t.TempDir(), "rei")
	if err := os.WriteFile(dest, []byte("VIEJO"), 0o755); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	opts := linuxOpts(fx, &buf, "v0.1.0")
	opts.Executable = func() (string, error) { return dest, nil }
	if code := Run(opts); code != 1 {
		t.Fatalf("Run con checksum inválido = %d, want 1", code)
	}
	if !strings.Contains(buf.String(), "SHA-256") {
		t.Fatalf("debe informar del fallo de integridad: %q", buf.String())
	}
	if got, _ := os.ReadFile(dest); string(got) != "VIEJO" {
		t.Fatalf("el binario no debe modificarse, got %q", got)
	}
}

func TestRunSinRed(t *testing.T) {
	var buf bytes.Buffer
	opts := Options{
		CurrentVersion: "v0.1.0",
		Stdout:         &buf,
		APIBase:        "http://127.0.0.1:1",
		GOOS:           "linux",
		GOARCH:         "amd64",
	}
	if code := Run(opts); code != 1 {
		t.Fatalf("Run sin red = %d, want 1", code)
	}
}

func TestRunAssetAusente(t *testing.T) {
	fx := newReleaseFixture(t, "v0.2.0", map[string][]byte{"checksums.txt": []byte("x  y\n")})
	var buf bytes.Buffer
	opts := linuxOpts(fx, &buf, "v0.1.0")
	if code := Run(opts); code != 1 {
		t.Fatalf("Run sin asset = %d, want 1", code)
	}
	if !strings.Contains(buf.String(), "AVISO") {
		t.Fatalf("debe avisar de la ausencia del asset: %q", buf.String())
	}
}

func TestRunCheckNetwork(t *testing.T) {
	fx := newReleaseFixture(t, "v0.2.0", nil)
	var buf bytes.Buffer
	opts := linuxOpts(fx, &buf, "v0.1.0")
	opts.Check = true
	if code := Run(opts); code != 0 {
		t.Fatalf("--check = %d, want 0", code)
	}
}

// TestRunNoReleases comprueba que un 404 de la API (sin releases publicadas) se
// trata como un estado normal: mensaje claro y salida 0.
func TestRunNoReleases(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	t.Cleanup(srv.Close)

	var buf bytes.Buffer
	opts := Options{
		CurrentVersion: "v0.1.0",
		Stdout:         &buf,
		APIBase:        srv.URL,
		Client:         srv.Client(),
		GOOS:           "linux",
		GOARCH:         "amd64",
	}
	if code := Run(opts); code != 0 {
		t.Fatalf("Run sin releases = %d, want 0; salida: %q", code, buf.String())
	}
	if !strings.Contains(buf.String(), "No hay releases publicadas") {
		t.Fatalf("debe informar de la ausencia de releases: %q", buf.String())
	}
}
