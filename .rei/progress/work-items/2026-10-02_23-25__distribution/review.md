# Revisión — distribución de binarios (goreleaser + GitHub Actions)

> Work Item: `2026-10-02_23-25__distribution` (task).
> Reviewer · resultado: **aprobado**.

## Alcance revisado

Revisión contra `plan.md` (9 pasos), `impl.md` y el paquete de
`rei review-diff 2026-10-02_23-25__distribution` (base
`9e7758c79e18efdb02d6e6fb207180e793a3066b`). Tipo `task`: sin
`requirements.md`/`design.md`/`tasks.md`, conforme a la planificación.

## Comprobaciones

### Objetivo y planificación

- [x] **§1 Versión inyectable.** `internal/cli/cli.go`: `const version =
      "0.1.0-dev"` → `var version = "dev"` con comentario del import path para
      `-X`. `Run` conserva el formato `rei <versión>`.
- [x] **§2 Test.** `internal/cli/version_test.go` (nuevo):
      `TestVersionDefaultYInyectable` comprueba el default `dev` y la
      mutabilidad (reasigna y restaura con `defer`).
- [x] **§3 Makefile.** `VERSION ?= dev`, `PKG`, `LDFLAGS := -X
      $(PKG).version=$(VERSION)`; aplicado a `build`, `install` y las cinco
      invocaciones de `cross`. `test`/`vet`/`fmt`/`clean` intactos.
- [x] **§4 `.goreleaser.yaml`.** v2, `project_name: rei`; build `main:
      ./cmd/rei`, `CGO_ENABLED=0`, `goos: [linux, darwin, windows]` ×
      `goarch: [amd64, arm64]` (6 combinaciones); `ldflags` con `-s -w -X
      github.com/k1wi777/my-harness-SDD/internal/cli.version={{ .Version }}`;
      archives `formats: [tar.gz]` con override `zip` para Windows y
      `name_template` `{{ .ProjectName }}_{{ .Version }}_{{ .Os }}_{{ .Arch }}`;
      `checksum.name_template: checksums.txt`; `changelog.sort: asc`.
- [x] **§5 `.github/workflows/release.yml`.** `on.push.tags: ['v*']`,
      `permissions.contents: write`; job en `ubuntu-latest` con
      `actions/checkout@v4` (`fetch-depth: 0`), `actions/setup-go@v5`
      (`go-version-file: go.mod`) y `goreleaser/goreleaser-action@v6`
      (`version: '~> v2'`, `args: release --clean`, `GITHUB_TOKEN`).
- [x] **§6/§7 Documentación.** `README.md` estrena sección **Instalación**
      (URL de releases, patrón `rei_<versión>_<os>_<arch>.tar.gz` / `.zip` en
      Windows, mover a `PATH`, `rei version`, alternativa `make build`).
      `.rei/docs/usage.md` actualiza el paso 1 de «Primera vez en el
      repositorio» con la descarga desde la release.

### Restricciones

- [x] No se crearon `requirements.md`/`design.md`/`tasks.md` (es una `task`).
- [x] `version` es una **variable** de paquete (no `const`) con default `dev`;
      el formato de salida `rei <versión>` no cambió.
- [x] No se añadieron dependencias Go.
- [x] `dist/` sigue en `.gitignore`; no se han introducido artefactos de build
      al repositorio. El workflow solo se dispara con tags `v*`.
- [x] No se ejecutaron goreleaser ni las GitHub Actions reales (fuera de
      alcance); la verificación de CI se limita a validar el YAML, como indica
      el plan.

### Verificación

`verification.md` no define aún checkpoints `V*` (documento pendiente de
personalización y `.rei/config.json` sin checks); se ejecutaron las
verificaciones declaradas en `plan.md`:

- [x] `make build` → `bin/rei version` = `rei dev`.
- [x] `go build -ldflags "-X github.com/k1wi777/my-harness-SDD/internal/cli.version=v0.1.0" -o /tmp/rei-review ./cmd/rei` → `/tmp/rei-review version` = `rei v0.1.0`.
- [x] `.goreleaser.yaml` y `.github/workflows/release.yml` parsean como YAML
      válido (`npx --yes js-yaml`, ambos OK) y contienen las claves del plan.
- [x] `gofmt -l .` → sin salida; `go vet ./...` → exit 0; `make test` → todos
      los paquetes `ok`.
- [x] `rei check --quiet` → exit 0.
- [x] `rei validate 2026-10-02_23-25__distribution` → `Resultado: OK`.

Evidencia en `.rei/progress/work-items/2026-10-02_23-25__distribution/impl.md`.

## Desviaciones

Ninguna. La implementación corresponde exactamente a los 9 pasos del plan, sin
ampliar alcance. Nota: `package-lock.json` figura como archivo sin rastrear en
`review-diff`, pero está fechado el 10-sep-2026, es previo a este Work Item
(ya señalado como preexistente y ajeno en revisiones anteriores) y no forma
parte del diff de este trabajo.

## Conclusión

El Work Item cumple el objetivo y las restricciones, y todas las verificaciones
declaradas pasan. **Aprobado.**
