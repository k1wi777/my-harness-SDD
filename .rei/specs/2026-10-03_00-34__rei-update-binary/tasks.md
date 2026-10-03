# Tasks — `rei update` (auto-actualización del binario)

> Work Item: `2026-10-03_00-34__rei-update-binary`
> Tipo: feature
>
> El Implementer marca cada tarea `[x]` inmediatamente al terminarla.

## Tareas

- [x] T1 — Crear `internal/update/semver.go`: normalizar (`v` opcional) y comparar versiones `major.minor.patch`; detectar `dev`/no parseable. (R3, R10)
- [x] T2 — Crear `internal/update/github.go`: consultar `releases/latest` (`User-Agent`, JSON mínimo), seleccionar el asset `<os>_<arch>` y localizar `checksums.txt`. (R2, R5)
- [x] T3 — Crear `internal/update/archive.go`: descarga con timeout, parseo de `checksums.txt` y verificación SHA-256; abortar sin escribir si falla. (R6, R11)
- [x] T4 — Añadir la extracción del binario a `internal/update/archive.go` para `.tar.gz` y `.zip`, con protección contra path traversal y limitada a `rei`/`rei.exe`. (R6)
- [x] T5 — Crear `internal/update/replace.go`: comprobar escritura, escribir temporal en el directorio destino, `chmod 0755` y `rename` con backup; restaurar el original si falla. (R7, R12, R13)
- [x] T6 — Crear `internal/update/browser.go`: `openBrowser` best-effort por SO (`xdg-open`/`open`/`cmd /c start`). (R8)
- [x] T7 — Crear `internal/update/update.go`: `Options` con costuras inyectables y `Run`, orquestando ramas `dev`, `--check`, sin actualización, Windows y Unix/macOS con códigos 0/1. (R1, R4, R8, R9, R10, R12)
- [x] T8 — Cablear `cmdUpdate` en `internal/cli/cli.go`: `case "update"`, validar `--check` (otro argumento → 2), pasar `version` a `update.Run`. (R1, R15)
- [x] T9 — Añadir la entrada `update` a la tabla `commands` de `internal/cli/help.go` (uso `update [--check]` y detalle de degradación). (R14)
- [x] T10 — Tests de `internal/update/*_test.go` con `httptest` y `t.TempDir()`: comparación semver, selección de asset/checksums, verificación ok y fallida, extracción tar.gz/zip, reemplazo con backup, ramas `dev`/`--check`/Windows/up-to-date y códigos de salida. (R2–R13)
- [x] T11 — Tests de `internal/cli`: despacho de `update`, ayuda del comando y código 2 ante argumentos inválidos. (R1, R14, R15)
- [x] T12 — Verificación final: `gofmt`, `go vet ./...`, `make test` y comprobación manual de `rei update --check`. (R1, R14)
