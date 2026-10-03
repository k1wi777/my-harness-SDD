# Implementación — `rei update [--check]`

> Work Item: `2026-10-03_00-34__rei-update-binary`
> Tipo: feature · Estado: review

## Resumen

Se implementó `rei update [--check]`: consulta `releases/latest` de GitHub,
compara semver con la versión en ejecución y, en Unix/macOS, descarga el asset
del SO/arch, verifica SHA-256 contra `checksums.txt`, extrae y reemplaza el
binario con backup. Windows guía al navegador; `--check` solo informa; `dev`
avisa y ofrece instalar; sin red/permisos degrada con código 1.

## Archivos

- `internal/update/semver.go` — parseo/comparación semver y detección de `dev`.
- `internal/update/github.go` — `releases/latest`, tipos y selección de
  asset/`checksums.txt`.
- `internal/update/archive.go` — descarga, parseo de checksums, SHA-256 y
  extracción `.tar.gz`/`.zip` (con guarda anti path-traversal).
- `internal/update/replace.go` — resolución del ejecutable y sustitución
  atómica (temp + chmod 0755 + rename con backup/restauración).
- `internal/update/browser.go` — `openBrowser` best-effort por SO.
- `internal/update/update.go` — `Options` con costuras inyectables (incl.
  `APIBase`) y `Run` con las ramas y códigos 0/1.
- `internal/cli/cli.go` — `case "update"` y `cmdUpdate` (otro arg → 2).
- `internal/cli/help.go` — entrada `update [--check]` documentada.
- `internal/update/*_test.go`, `internal/cli/update_test.go` — tests con
  `httptest` y `t.TempDir()`.

## Verificación

- `gofmt -l .` sin salida; `go vet ./...` OK; `make test` OK (todos los
  paquetes); `make build` OK; `rei check --quiet` exit 0.
- `rei help` y `rei update --help` incluyen `update`/`--check`.
- `rei update --check` real: `ERROR ... 404 Not Found` (el repo no publica
  release), exit 1; degrada con aviso sin descargar ni escribir. `rei update
  bogus` → uso + exit 2.
- Tests cubren semver, selección asset/checksums, sha256 ok/fallida, tar.gz/zip,
  path traversal, reemplazo con backup y permisos, y ramas dev/--check/Windows/
  up-to-date/sin red/asset ausente.

## Observaciones

- Sin dependencias nuevas (solo librería estándar).
- La costura `Options.APIBase` se añadió para poder usar `httptest` (requisito
  de testabilidad del diseño).
