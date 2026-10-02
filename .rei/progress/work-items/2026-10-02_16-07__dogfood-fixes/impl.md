# Implementación — corregir hallazgos del dogfooding

Work Item: `2026-10-02_16-07__dogfood-fixes` (task)
Plan: `.rei/specs/2026-10-02_16-07__dogfood-fixes/plan.md`

## Resumen

Se corrigieron los dos hallazgos del dogfooding:

1. **Filtro de `rei review-diff`.** El paquete ahora descarta **solo**
   `.rei/progress/**` y conserva el resto: código del proyecto, `AGENTS.md`, la
   spec del Work Item, `.rei/docs/**`, `.rei/agents/**`, `.rei/adapters/**`,
   `.rei/templates/**` y `.rei/config.json`.
2. **Handoff coherente del Leader.** Se documentó en el `## Contrato` de
   `.rei/agents/leader.md` que, ante el `WARN` de `rei validate` por
   incoherencia entre `current.md` (Estado) y `meta.json` (status), el Leader
   debe corregir el campo **Estado** de `current.md` antes de avanzar.

## Archivos modificados

| Archivo | Cambio |
|---------|--------|
| `internal/gitx/gitx.go` | `filterPaths` ahora excluye solo `.rei/progress/**` y perdió el parámetro `id`; se eliminó `filterCode`; rótulos del paquete sin "solo código"; el salto de no rastreados en `--full` omite solo `.rei/progress/**`. |
| `internal/gitx/gitx_test.go` | `TestFilterPaths` cubre la conservación de `.rei/docs`, `.rei/agents`, `.rei/adapters`, `.rei/templates`, `.rei/config.json`, spec de otros Work Items y `AGENTS.md`, y el descarte de `.rei/progress/**`; se eliminó `TestFilterCode` (el filtro ya no existe). |
| `.rei/agents/leader.md` | Delegación: corregir `current.md` ante el `WARN` de `rei validate`. Reglas duras: no avanzar con `current.md` y `meta.json` incoherentes. |
| `.opencode/agents/leader.md` | Regenerado con `rei init opencode` (archivo generado). |

## Detalle de la implementación

- `ReviewDiff` usa la lista ya filtrada (`tracked`) tanto para el resumen como
  para el diff completo, por lo que `.rei/docs/**` y `.rei/agents/**` dejan de
  quedar ocultos. La spec del Work Item sigue apareciendo.
- Los rótulos pasaron de `Resumen (código)` / `Diff completo (solo código; la
  spec se lista arriba)` a `Resumen` / `Diff completo`.
- El bloque de archivos no rastreados en `--full` salta solo
  `.rei/progress/**` (redundante con `filterPaths`, pero explícito).
- No se tocó el comportamiento de base/`last_review_commit`, formato general ni
  exit codes.

## Verificación

- `rei check --quiet` → exit 0.
- `go test ./...` → todo en verde.
- `make build` → compila (exit 0).
- `rei init opencode` regenera `.opencode/agents/leader.md`;
  `rei init opencode --check` → exit 0.
- Comprobación manual de `rei review-diff 2026-10-02_16-07__dogfood-fixes`:
  - el resumen lista `.rei/agents/leader.md` (antes oculto);
  - tras añadir temporalmente una línea a `.rei/docs/project/architecture.md`,
    el resumen la lista; el cambio temporal se revirtió con
    `git checkout -- .rei/docs/project/architecture.md`.

## Estado

Todos los pasos del plan completados. Listo para revisión.
