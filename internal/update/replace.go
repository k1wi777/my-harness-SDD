package update

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// resolveExecutable localiza el binario en ejecución y, best-effort, resuelve
// enlaces simbólicos para no reemplazar el enlace sino el archivo real.
func resolveExecutable(exe func() (string, error)) (string, error) {
	path, err := exe()
	if err != nil {
		return "", err
	}
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		path = resolved
	}
	return path, nil
}

// replaceBinary sustituye destPath por el binario srcPath de forma atómica:
// escribe un temporal en el mismo directorio, conserva el permiso de ejecución
// y respalda el original. Si la sustitución falla, restaura el original (R13).
func replaceBinary(srcPath, destPath string) error {
	dir := filepath.Dir(destPath)

	// Pre-comprobar escritura creando el temporal junto al destino (necesario
	// para que os.Rename sea atómico en el mismo filesystem). Si falla, no se
	// toca el binario (R12).
	tmp, err := os.CreateTemp(dir, ".rei-update-*")
	if err != nil {
		return fmt.Errorf("no se puede escribir en %s: %w", dir, err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	src, err := os.Open(srcPath)
	if err != nil {
		tmp.Close()
		return err
	}
	defer src.Close()
	if _, err := io.Copy(tmp, src); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmpName, 0o755); err != nil {
		return err
	}

	backup := destPath + ".bak"
	if err := os.Rename(destPath, backup); err != nil {
		return fmt.Errorf("no se puede reemplazar %s: %w", destPath, err)
	}
	if err := os.Rename(tmpName, destPath); err != nil {
		if rerr := os.Rename(backup, destPath); rerr != nil {
			return fmt.Errorf("fallo al instalar la nueva versión (%v) y al restaurar el original (%v)", err, rerr)
		}
		return fmt.Errorf("fallo al instalar la nueva versión; se restauró el binario original: %w", err)
	}
	// Consumido el temporal en el rename; evita borrarlo de nuevo.
	tmpName = ""
	os.Remove(backup) // best-effort
	return nil
}
