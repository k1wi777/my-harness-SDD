package update

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// maxChecksumsBytes acota la descarga del archivo de checksums.
const maxChecksumsBytes = 1 << 20

// httpGet realiza un GET con User-Agent y exige un 200.
func httpGet(client *http.Client, url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("respuesta inesperada al descargar %s: %s", url, resp.Status)
	}
	return resp, nil
}

// downloadToFile descarga url en dst.
func downloadToFile(client *http.Client, url, dst string) error {
	resp, err := httpGet(client, url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// downloadToMemory descarga url en memoria (para archivos pequeños).
func downloadToMemory(client *http.Client, url string) ([]byte, error) {
	resp, err := httpGet(client, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(io.LimitReader(resp.Body, maxChecksumsBytes))
}

// parseChecksums parsea el formato goreleaser "<sha256>  <nombre>".
func parseChecksums(data []byte) map[string]string {
	sums := map[string]string{}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(strings.TrimSpace(line))
		if len(fields) < 2 {
			continue
		}
		name := strings.TrimPrefix(fields[len(fields)-1], "*")
		sums[name] = strings.ToLower(fields[0])
	}
	return sums
}

// verifySHA256 verifica el hash del archivo dst contra la entrada de checksums
// correspondiente a name. Devuelve error si falta la entrada o no coincide.
func verifySHA256(dst string, checksums []byte, name string) error {
	want, ok := parseChecksums(checksums)[name]
	if !ok {
		return fmt.Errorf("checksums.txt no contiene una entrada para %s", name)
	}
	f, err := os.Open(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	got := hex.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(got), []byte(want)) {
		return fmt.Errorf("la verificación SHA-256 de %s falló", name)
	}
	return nil
}

// extractBinary extrae el binario rei/rei.exe desde un .tar.gz o .zip.
func extractBinary(archivePath, dst, goos string) error {
	want := "rei"
	if goos == "windows" {
		want = "rei.exe"
	}
	switch {
	case strings.HasSuffix(archivePath, ".zip"):
		return extractFromZip(archivePath, dst, want)
	case strings.HasSuffix(archivePath, ".tar.gz"):
		return extractFromTarGz(archivePath, dst, want)
	default:
		return fmt.Errorf("formato de archivo no soportado: %s", filepath.Base(archivePath))
	}
}

// safeMember rechaza rutas absolutas o con ".." (defensa ante path traversal).
func safeMember(name string) bool {
	if filepath.IsAbs(name) {
		return false
	}
	for _, part := range strings.Split(filepath.ToSlash(name), "/") {
		if part == ".." {
			return false
		}
	}
	return true
}

func writeEntry(dst string, r io.Reader) error {
	f, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

func extractFromTarGz(archivePath, dst, want string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if hdr.Typeflag != tar.TypeReg || !safeMember(hdr.Name) {
			continue
		}
		if filepath.Base(hdr.Name) != want {
			continue
		}
		return writeEntry(dst, tr)
	}
	return fmt.Errorf("%s no contiene el binario %s", filepath.Base(archivePath), want)
}

func extractFromZip(archivePath, dst, want string) error {
	zr, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zr.Close()
	for _, zf := range zr.File {
		if zf.FileInfo().IsDir() || !safeMember(zf.Name) {
			continue
		}
		if filepath.Base(zf.Name) != want {
			continue
		}
		rc, err := zf.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		return writeEntry(dst, rc)
	}
	return fmt.Errorf("%s no contiene el binario %s", filepath.Base(archivePath), want)
}
