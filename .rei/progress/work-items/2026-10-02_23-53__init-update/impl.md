# Implementación — `rei init --update [--force]`

## Resumen

Implementados los 8 pasos de `plan.md`. `rei init --update [--force]` refresca el
esqueleto embebido en un proyecto ya inicializado combinando reglas marker-aware
(docs de personalización) y un manifiesto de hashes
(`.rei/install-manifest.json`) para el resto del harness. `--force` sobrescribe
todo; `.rei/specs/**` y `.rei/progress/**` nunca se tocan.

## Archivos

- `manifest.go` (nuevo): `ManifestPath`, `hashBytes`, `loadManifest` (nil si
  falta/inválido) y `saveManifest` (claves ordenadas, indentado, 0644).
- `installer.go`: `InstallSkeleton` registra el baseline (archivos idénticos al
  embebido) y guarda el manifiesto; fallo de guardado = `[WARN]`, sin cambiar exit.
- `update.go` (nuevo): `Update(p, force, out)` con la tabla de decisión, reporte
  `[NEW]/[UPD]/[OK]/[SKIP]`, reescritura del manifiesto y resumen.
- `cli.go`: `cmdInit` detecta `--update`, valida args (`--update [--force]`,
  order-independent), resuelve con `project()` y regenera nativos solo si
  `.opencode/`/`.claude/` existen.
- `help.go`: uso y detalle de `--update`/`--force`/manifiesto.
- Tests: `manifest_test.go` y `update_test.go` (nuevos); `help_test.go` y
  `cli_test.go` ajustados.

## Verificación

- V1 `go test ./...` → 0. V2 `gofmt -l .` limpio y `go vet ./...` → 0.
- V3 `make build` → `bin/rei` compilado.
- V4 en `/tmp/opencode/v4-update`: personalizado + `leader.md` + `config.json`
  conservados `[SKIP]`; docs con marcador y harness con baseline intacto `[UPD]`;
  `--force` sobrescribe (0 omitidos); `specs/`+`progress/` sin cambios; 2ª
  `--update` idempotente (0/0/25/0); fuera de proyecto REI → exit 1; nativos
  regenerados.
- V5 `rei check --quiet`, `rei init opencode --check` y `claude --check` → 0.
## Observaciones

- Proyectos previos sin manifiesto: la primera `--update` conserva con aviso el
  harness que difiera; `--force` para el primer refresco.
- Sin dependencias externas nuevas ni cambios en `embed.go`/planificación.
