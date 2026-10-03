package initwizard

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	reiskel "github.com/k1wi777/my-harness-SDD"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

// ResolveRootFrom determina el proyecto destino partiendo de dir: el ancestro
// más cercano que contenga .rei/ o, si no existe, dir mismo.
func ResolveRootFrom(dir string) (*paths.Project, error) {
	p, err := paths.FindFrom(dir)
	if err == nil {
		return p, nil
	}
	if errors.Is(err, paths.ErrNotFound) {
		return &paths.Project{Root: dir}, nil
	}
	return nil, err
}

// ResolveRoot aplica ResolveRootFrom sobre el directorio de trabajo actual.
func ResolveRoot() (*paths.Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return ResolveRootFrom(cwd)
}

// InstallSkeleton copia el esqueleto embebido en p.Root sin sobrescribir los
// archivos que ya existan. Devuelve 0 si lo completa y 1 si no pudo crear algún
// archivo o directorio.
//
// La omisión se decide por existencia (no comparando contenido): es lo que
// garantiza la idempotencia y que `rei init` nunca destruya personalizaciones
// del usuario (AGENTS.md, .rei/config.json, ...).
func InstallSkeleton(p *paths.Project, out io.Writer) int {
	exit := 0
	created, skipped := 0, 0
	baseline := map[string]string{}

	walkErr := fs.WalkDir(reiskel.Skeleton, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == "." {
			return nil
		}
		dest := filepath.Join(p.Root, filepath.FromSlash(path))

		if d.IsDir() {
			if err := os.MkdirAll(dest, 0o755); err != nil {
				fmt.Fprintf(out, "[FAIL]  No se pudo crear %s: %v\n", path, err)
				exit = 1
			}
			return nil
		}

		embedded, err := reiskel.Skeleton.ReadFile(path)
		if err != nil {
			fmt.Fprintf(out, "[FAIL]  No se pudo leer %s del esqueleto: %v\n", path, err)
			exit = 1
			return nil
		}

		if current, statErr := os.ReadFile(dest); statErr == nil {
			if bytes.Equal(current, embedded) {
				// Sin cambios respecto al embebido: queda como baseline.
				baseline[path] = hashBytes(embedded)
			}
			if path == "AGENTS.md" {
				fmt.Fprintf(out, "[WARN]  %s ya existe; se conserva sin cambios.\n", path)
			}
			fmt.Fprintf(out, "[SKIP]  %s (ya existe)\n", path)
			skipped++
			return nil
		} else if !errors.Is(statErr, fs.ErrNotExist) {
			fmt.Fprintf(out, "[FAIL]  No se pudo comprobar %s: %v\n", path, statErr)
			exit = 1
			return nil
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			fmt.Fprintf(out, "[FAIL]  No se pudo crear el directorio de %s: %v\n", path, err)
			exit = 1
			return nil
		}
		if err := os.WriteFile(dest, embedded, 0o644); err != nil {
			fmt.Fprintf(out, "[FAIL]  No se pudo crear %s: %v\n", path, err)
			exit = 1
			return nil
		}
		// Recién creado a partir del embebido: idéntico por construcción.
		baseline[path] = hashBytes(embedded)
		fmt.Fprintf(out, "[OK]    %s\n", path)
		created++
		return nil
	})
	if walkErr != nil {
		fmt.Fprintf(out, "[FAIL]  No se pudo recorrer el esqueleto: %v\n", walkErr)
		exit = 1
	}

	if err := saveManifest(p, baseline); err != nil {
		fmt.Fprintf(out, "[WARN]  No se pudo guardar el manifiesto %s: %v\n", ManifestPath, err)
	}

	fmt.Fprintf(out, "%d archivo(s) creado(s), %d omitido(s).\n", created, skipped)
	return exit
}
