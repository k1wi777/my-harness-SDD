# Implementación — distribución de binarios (goreleaser + GitHub Actions)

- **Work Item:** `2026-10-02_23-25__distribution`
- **Tipo:** task
- **Implementer:** implementer
- **Estado:** review
- **Base commit:** `9e7758c79e18efdb02d6e6fb207180e793a3066b`

## Resumen

Se implementó la distribución de `rei` como binario precompilado siguiendo
exactamente los 9 pasos de `plan.md`:

1. La versión del CLI dejó de ser una constante y pasó a ser una **variable de
   paquete** con valor por defecto `dev`, inyectable en build time con
   `-ldflags -X .../internal/cli.version=...`.
2. Se añadió un test que verifica el valor por defecto y la mutabilidad.
3. El `Makefile` inyecta la versión (`VERSION`, `LDFLAGS`) en `build`, `install`
   y `cross`.
4. Se añadió `.goreleaser.yaml` (goreleaser v2) para linux/darwin/windows ×
   amd64/arm64 con archives, checksums y `ldflags`.
5. Se añadió `.github/workflows/release.yml` que se dispara con tags `v*` y
   ejecuta goreleaser publicando el release.
6. Se documentó la instalación desde releases en `README.md`.
7. Se actualizó la instalación en `.rei/docs/usage.md`.

## Archivos modificados / creados

| Archivo | Cambio |
|---------|--------|
| `internal/cli/cli.go` | `const version = "0.1.0-dev"` → `var version = "dev"` con comentario del import path para `-X`. |
| `internal/cli/version_test.go` | **Nuevo.** `TestVersionDefaultYInyectable`: default `dev` + mutabilidad. |
| `Makefile` | Nuevas variables `VERSION ?= dev`, `PKG`, `LDFLAGS`; aplicadas a `build`, `install` y `cross`. |
| `.goreleaser.yaml` | **Nuevo.** Config goreleaser v2. |
| `.github/workflows/release.yml` | **Nuevo.** Workflow de release por tag. |
| `README.md` | Nueva sección **Instalación** desde releases. |
| `.rei/docs/usage.md` | Paso 1 de «Primera vez en el repositorio» actualizado a descarga desde release. |
| `.rei/specs/2026-10-02_23-25__distribution/plan.md` | Pasos marcados `[x]`. |

## Cambios realizados

### 1. Versión inyectable — `internal/cli/cli.go`

```go
// version se inyecta en build time con:
//
//	-ldflags "-X github.com/k1wi777/my-harness-SDD/internal/cli.version=vX.Y.Z"
//
// Por defecto identifica una compilación local sin versión publicada.
var version = "dev"
```

`Run` no cambió: sigue imprimiendo `rei <version>`.

### 2. Test — `internal/cli/version_test.go`

Comprueba que el default es `dev`, guarda el valor original, lo reasigna a
`v0.1.0-test` y lo restaura con `defer`. Demuestra que la variable es mutable y
por tanto sobrescribible por el linker.

### 3. Makefile

```make
VERSION ?= dev
PKG     := github.com/k1wi777/my-harness-SDD/internal/cli
LDFLAGS := -X $(PKG).version=$(VERSION)
```

`build` e `install` y las cinco invocaciones de `cross` usan
`-ldflags "$(LDFLAGS)"`.

### 4. `.goreleaser.yaml`

- `version: 2`, `project_name: rei`.
- Build `id: rei`, `main: ./cmd/rei`, `CGO_ENABLED=0`, `goos`
  linux/darwin/windows, `goarch` amd64/arm64 (6 combinaciones).
- `ldflags: -s -w -X github.com/k1wi777/my-harness-SDD/internal/cli.version={{ .Version }}`.
- Archives `tar.gz` con `name_template`
  `{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}` y override `zip`
  para Windows.
- `checksum.name_template: checksums.txt`; `changelog.sort: asc`.

### 5. `.github/workflows/release.yml`

- `on: push: tags: ['v*']`, `permissions: contents: write`.
- Job `goreleaser` en `ubuntu-latest`: `actions/checkout@v4` (`fetch-depth: 0`),
  `actions/setup-go@v5` (`go-version-file: go.mod`) y
  `goreleaser/goreleaser-action@v6` (`~> v2`, `release --clean`,
  `GITHUB_TOKEN`).

### 6. Documentación

- `README.md`: nueva sección **Instalación** con la URL de releases, el patrón
  de archive `rei_<versión>_<os>_<arch>.tar.gz` (`.zip` en Windows), mover el
  binario al `PATH` y verificar con `rei version`; alternativa `make build`.
- `.rei/docs/usage.md`: el paso 1 apunta a la página de releases y mantiene la
  alternativa `make build`.

## Proceso de verificación

Comandos ejecutados (todos con exit 0):

| Verificación | Comando | Resultado |
|--------------|---------|-----------|
| Build por defecto | `make build` | `bin/rei version` → `rei dev` |
| Versión inyectada | `go build -ldflags "-X github.com/k1wi777/my-harness-SDD/internal/cli.version=v0.1.0" -o /tmp/rei-dist ./cmd/rei` | `/tmp/rei-dist version` → `rei v0.1.0` |
| Formato | `gofmt -l .` | sin salida |
| Vet | `go vet ./...` | exit 0 |
| Suite | `make test` | todos los paquetes `ok` |
| Check del harness | `rei check --quiet` | exit 0 |
| YAML `.goreleaser.yaml` | `js-yaml` (PyYAML no disponible) | YAML válido; goreleaser v2, 6 combinaciones SO/arch, archives, checksum, ldflags con `internal/cli.version` |
| YAML `release.yml` | `js-yaml` | YAML válido; `on.push.tags: ['v*']`, `permissions.contents: write`, `goreleaser-action@v6` |

## Observaciones

- `goreleaser` no está instalado en el entorno (fuera de alcance), por lo que la
  validación se hizo como YAML y revisión manual de las claves contra el esquema
  v2.
- `PyYAML` no está disponible en el sistema; la validación YAML se realizó con
  `js-yaml` (instalado en `/tmp/opencode/yamlcheck`), de forma temporal y sin
  afectar al repositorio.
- goreleaser normaliza el tag `v0.1.0` a `{{ .Version }}` = `0.1.0`, por lo que
  el binario del release reportará `rei 0.1.0`.
- No se añadieron dependencias Go ni se modificó el alcance de `meta.json`.
- No se ejecutaron goreleaser ni las GitHub Actions reales.
