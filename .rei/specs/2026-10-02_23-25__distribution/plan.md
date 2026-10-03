# Plan — distribución de binarios (goreleaser + GitHub Actions)

## Objetivo

Distribuir `rei` como **binario precompilado** mediante releases de GitHub,
publicadas automáticamente al empujar un tag `v*`. La versión del binario debe
ser **inyectable en tiempo de compilación** (por `-ldflags`), con un valor por
defecto `dev`, y mostrarse en `rei version`.

Puntos acordados:

1. **Versión inyectable.** En `internal/cli`, la versión debe ser una
   **variable de paquete** (no `const`), fijable con
   `-ldflags "-X github.com/k1wi777/my-harness-SDD/internal/cli.version=vX.Y.Z"`.
   Valor por defecto: `dev`. `rei version` la muestra.
2. **`.goreleaser.yaml`.** Compila `linux/darwin/windows` × `amd64/arm64`,
   genera archives (`tar.gz`; `zip` en Windows) y `checksums.txt`, e inyecta la
   versión por `ldflags`.
3. **GitHub Actions.** `.github/workflows/release.yml` se dispara al empujar un
   tag `v*`, ejecuta goreleaser y publica el release.
4. **Documentación de instalación.** Descargar el binario de la release para tu
   SO/arquitectura, ponerlo en `PATH` y verificar con `rei version`; en
   `README.md` y `.rei/docs/usage.md`.

Primer release previsto: **`v0.1.0`**.

Fuera de alcance: instalar goreleaser en este entorno, ejecutar las GitHub
Actions, publicar el release real o firmar artefactos.

## Archivos

### Código del CLI

- `internal/cli/cli.go` — `const version` → `var version` con default `dev`.
- `internal/cli/version_test.go` (nuevo) — test de la versión por defecto y de
  su mutabilidad.

### Build y release

- `Makefile` — variable `VERSION` (default `dev`) y `LDFLAGS`; aplicarla a
  `build`, `install` y `cross`.
- `.goreleaser.yaml` (nuevo) — configuración de goreleaser.
- `.github/workflows/release.yml` (nuevo) — workflow de release por tag.

### Documentación

- `README.md` — sección de **Instalación** desde releases.
- `.rei/docs/usage.md` — actualizar «Primera vez en el repositorio» con la
  instalación desde la release.

### Gestión del Work Item

- `.rei/specs/2026-10-02_23-25__distribution/plan.md` (este archivo).
- `.rei/progress/current.md` y `.rei/specs/2026-10-02_23-25__distribution/meta.json`.

## Cambios

### 1. `internal/cli/cli.go` — versión inyectable

- Reemplazar:

  ```go
  const version = "0.1.0-dev"
  ```

  por:

  ```go
  // version se inyecta en build time con:
  //   -ldflags "-X github.com/k1wi777/my-harness-SDD/internal/cli.version=vX.Y.Z"
  // Por defecto identifica una compilación local sin versión publicada.
  var version = "dev"
  ```

- `rei version` (ya implementado) imprime `rei <version>`; no cambia su
  comportamiento. Import path para `-X`:
  `github.com/k1wi777/my-harness-SDD/internal/cli`.

### 2. Test `internal/cli/version_test.go`

- Añadir un test `TestVersionDefaultYInyectable` que:
  - comprueba que el valor por defecto es `dev`;
  - guarden el valor, lo reasignen (p. ej. a `v0.1.0-test`) y restaurarlo con
    `defer`, demostrando que es una variable mutable (no constante).

  No hace falta capturar stdout: basta con verificar la variable, que es lo que
  consumen tanto `Run` como el linker.

### 3. `Makefile` — inyección de versión

- Añadir cerca de las variables existentes:

  ```make
  VERSION ?= dev
  PKG     := github.com/k1wi777/my-harness-SDD/internal/cli
  LDFLAGS := -X $(PKG).version=$(VERSION)
  ```

- `build`: `$(GO) build -ldflags "$(LDFLAGS)" -o $(BIN_DIR)/$(BINARY) ./cmd/rei`.
- `install`: `$(GO) install -ldflags "$(LDFLAGS)" ./cmd/rei`.
- `cross`: añadir `-ldflags "$(LDFLAGS)"` a cada invocación.
- Mantener `test`, `vet`, `fmt`, `clean` sin cambios. `make build` sin
  `VERSION` produce un binario que reporta `rei dev`.

### 4. `.goreleaser.yaml` (nuevo)

Configuración goreleaser v2:

- `version: 2`, `project_name: rei`.
- `builds`:
  - `id: rei`, `main: ./cmd/rei`, `binary: rei`, `env: [CGO_ENABLED=0]`.
  - `goos: [linux, darwin, windows]`, `goarch: [amd64, arm64]`.
  - `ldflags: [-s -w -X github.com/k1wi777/my-harness-SDD/internal/cli.version={{ .Version }}]`.
- `archives`:
  - `formats: [tar.gz]` (v2), `name_template` con
    `{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}`.
  - `format_overrides` para `goos: windows` con `formats: [zip]`.
- `checksum`: `name_template: checksums.txt`.
- `changelog`: `sort: asc` (opcional, para un changelog legible).

Nota: goreleaser normaliza el tag `v0.1.0` a `{{ .Version }}` = `0.1.0`, por lo
que el binario del release reportará `rei 0.1.0`.

### 5. `.github/workflows/release.yml` (nuevo)

Workflow disparado por tag:

- `on: push: tags: ['v*']`.
- `permissions: contents: write` (necesario para publicar el release con
  `GITHUB_TOKEN`).
- Job `goreleaser` en `ubuntu-latest`:
  1. `actions/checkout@v4` con `fetch-depth: 0` (goreleaser necesita el
     historial/tags).
  2. `actions/setup-go@v5` con `go-version-file: go.mod` (el módulo declara
     Go 1.27; si el runner no lo soporta, fijar `go-version` explícito).
  3. `goreleaser/goreleaser-action@v6` con
     `distribution: goreleaser`, `version: '~> v2'`,
     `args: release --clean` y `GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}`.

### 6. `README.md` — instalación

- Añadir una sección **Instalación** (antes o dentro de «Primeros pasos») que
  explique:
  1. ir a la página de releases
     (`https://github.com/k1wi777/my-harness-SDD/releases`);
  2. descargar el archive correspondiente a tu SO y arquitectura
     (`rei_<versión>_<os>_<arch>.tar.gz`, o `.zip` en Windows);
  3. descomprimirlo y mover el binario `rei` a un directorio del `PATH`
     (`/usr/local/bin`, `~/.local/bin`, etc.);
  4. verificar con `rei version`.
- Mencionar la alternativa para desarrolladores: `make build` (requiere Go).

### 7. `.rei/docs/usage.md` — instalación

- En «1. Primera vez en el repositorio», actualizar el paso 1 para indicar que
  el binario se descarga de la **release** para tu SO/arquitectura y se coloca
  en el `PATH` (con enlace a la página de releases), manteniendo la alternativa
  `make build` para quien tenga Go.

## Restricciones

- Es una `task`: **no** crear `requirements.md`, `design.md` ni `tasks.md`.
- No cambiar el alcance de `meta.json`.
- La versión debe ser una **variable** de paquete con default `dev`; no una
  `const`. No cambiar el formato de salida de `rei version` (`rei <versión>`).
- No añadir dependencias Go nuevas.
- `dist/` ya está en `.gitignore`; no commitear artefactos de build.
- El workflow solo se dispara con tags `v*`; no crear/modificar otros workflows.
- No ejecutar goreleaser ni las Actions en este entorno; la verificación de
  CI queda fuera de alcance (solo se valida el YAML).
- No romper la suite existente ni `rei check --quiet`.

## Pasos

- [x] 1. Convertir `version` en `internal/cli/cli.go` a `var` con default `dev`
      (comentario con el import path para `-X`).
- [x] 2. Añadir `internal/cli/version_test.go` (default `dev` + mutabilidad).
- [x] 3. Actualizar `Makefile`: `VERSION`, `LDFLAGS` y su uso en `build`,
      `install` y `cross`.
- [x] 4. Crear `.goreleaser.yaml` (builds, archives, checksums y ldflags).
- [x] 5. Crear `.github/workflows/release.yml` (tag `v*` → goreleaser).
- [x] 6. Documentar la instalación desde releases en `README.md`.
- [x] 7. Documentar la instalación desde releases en `.rei/docs/usage.md`.
- [x] 8. Ejecutar `go build` con `-ldflags` y comprobar `rei version`; ejecutar
      `make test`, `make vet` y `rei check --quiet`.
- [x] 9. Validar el YAML de `.goreleaser.yaml` y de `release.yml` (sintaxis y
      estructura); si `goreleaser` no está disponible, validar como YAML y
      revisar las claves contra el esquema v2.

## Verificación

- **Versión por defecto:** `make build` y `bin/rei version` → `rei dev`.
- **Versión inyectada:**
  `go build -ldflags "-X github.com/k1wi777/my-harness-SDD/internal/cli.version=v0.1.0" -o /tmp/rei ./cmd/rei`
  y `/tmp/rei version` → `rei v0.1.0`.
- **YAML:** `.goreleaser.yaml` y `.github/workflows/release.yml` parsean como
  YAML válido; `.goreleaser.yaml` incluye `builds` (6 combinaciones
  SO/arch), `archives`, `checksum` y el `ldflags` con `internal/cli.version`;
  el workflow se dispara con `tags: ['v*']` y usa `goreleaser-action`.
  Si `goreleaser` estuviera instalado: `goreleaser check`.
- **Suite:** `make test` y `make vet` en verde; `rei check --quiet` → 0.
- **Documentación:** `README.md` y `.rei/docs/usage.md` describen la descarga
  del binario de la release y su inclusión en el `PATH`.
