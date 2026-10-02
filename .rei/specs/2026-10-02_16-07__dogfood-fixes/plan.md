# Plan — corregir hallazgos del dogfooding

## Objetivo

Corregir dos hallazgos del dogfooding de subagentes nativos:

1. **Filtro de `rei review-diff`.** Hoy el paquete descarta todo `.rei/` salvo la
   spec del Work Item, por lo que quedan ocultos cambios legítimos en
   `.rei/docs/**`, `.rei/agents/**`, `.rei/adapters/**`, `.rei/templates/**`,
   `.rei/config.json` y `AGENTS.md`. Debe excluir **solo** `.rei/progress/**`
   (bookkeeping de sesión) y conservar el resto: código del proyecto, `AGENTS.md`,
   la spec del Work Item y el resto del harness.
2. **Handoff coherente del Leader.** Al ejecutar `rei validate <id>` tras cada
   subagente puede aparecer un `WARN` de incoherencia entre `current.md` (Estado)
   y `meta.json` (status). El Leader debe corregir `current.md` para alinearlo con
   `meta.json` antes de avanzar, y esto debe quedar documentado en el `## Contrato`
   de `.rei/agents/leader.md`.

## Archivos

- `internal/gitx/gitx.go` — filtros de rutas y armado del paquete de revisión.
- `internal/gitx/gitx_test.go` — pruebas de los filtros.
- `.rei/agents/leader.md` — sección `## Contrato` (Delegación y Reglas duras).

## Cambios

### 1. `internal/gitx/gitx.go`

Criterio único de exclusión: se descarta toda ruta con prefijo `.rei/progress/`;
ninguna otra ruta se descarta.

- `filterPaths` (líneas ~160-176): conservar todas las rutas de entrada salvo
  `.rei/progress/**`. Deja de excluir genéricamente `.rei/` y de tratar la spec
  del Work Item como caso especial. Si el parámetro `id` deja de ser necesario,
  ajustar su firma y las dos llamadas (`tracked` y `untracked`).
- `filterCode` (líneas ~178-187): hoy descarta todo `.rei/`, con lo que las
  docs/agentes del harness no aparecen en el diff. Alinearlo al mismo criterio
  (excluir solo `.rei/progress/**`) o eliminarlo y reutilizar la lista ya
  filtrada. Actualizar los rótulos `Resumen (código)` y
  `Diff completo (solo código)` para que no prometan "solo código".
- Bloque `full` del diff de archivos no rastreados (línea ~136): el salto que hoy
  usa `strings.HasPrefix(u, ".rei/")` debe omitir solo `.rei/progress/**`, de modo
  que los archivos nuevos bajo `.rei/docs`, `.rei/agents`, etc. aparezcan con su
  diff.
- Actualizar los comentarios de ambas funciones para reflejar el nuevo criterio.

### 2. `.rei/agents/leader.md` (`## Contrato`)

- En **Delegación**, junto al punto que ordena ejecutar
  `rei validate <work-item-id>` tras cada subagente, añadir: si `rei validate`
  emite el `WARN` de incoherencia entre `current.md` (Estado) y `meta.json`
  (status), el Leader corrige `current.md` para alinearlo con `meta.json` antes de
  avanzar (campo **Estado** de `.rei/progress/current.md` con el status de
  `meta.json`).
- Reflejar la misma obligación en **Reglas duras** (no avanzar mientras
  `current.md` y `meta.json` estén incoherentes).

### 3. `internal/gitx/gitx_test.go`

- Actualizar `TestFilterPaths` al nuevo criterio: se conserva `.rei/docs/**`,
  `.rei/specs/**` (incluida la spec de otros Work Items) y `AGENTS.md`; solo se
  descarta `.rei/progress/**`.
- Ajustar `TestFilterCode` según el diseño elegido (o eliminarlo si el filtro
  desaparece), cubriendo explícitamente que `.rei/docs`, `.rei/agents`,
  `.rei/config.json` y `AGENTS.md` se conservan, y que `.rei/progress/**` se
  descarta.

## Restricciones

- Fuera de alcance: slimming de Contratos y otros runtimes. No editar a mano
  adaptadores generados (p. ej. `.opencode/agents/leader.md`); el cambio se
  realiza solo en la fuente `.rei/agents/leader.md`.
- No cambiar el comportamiento de `rei review-diff` más allá de la exclusión de
  `.rei/progress/**` (base/last_review_commit, formato y exit codes intactos).
- No alterar el alcance de `meta.json`.
- `rei check --quiet` debe seguir pasando y `go test ./...` debe quedar en verde.

## Pasos

- [x] 1. Ajustar `filterPaths` para excluir únicamente `.rei/progress/**`.
- [x] 2. Alinear `filterCode` al mismo criterio (o eliminarlo) y actualizar rótulos.
- [x] 3. Corregir el salto de no rastreados del modo `--full` a solo `.rei/progress/**`.
- [x] 4. Actualizar `gitx_test.go` con el nuevo comportamiento de los filtros.
- [x] 5. Documentar en `.rei/agents/leader.md` (`## Contrato`) la corrección de
  `current.md` cuando `rei validate` reporte incoherencia con `meta.json`.
- [x] 6. Ejecutar `rei check --quiet` y `go test ./...`.
