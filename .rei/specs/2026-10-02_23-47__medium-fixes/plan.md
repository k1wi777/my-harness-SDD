# plan.md

## Objetivo

Dos correcciones medias, sin solaparse:

1. **Desbloqueo sin `--force`**: `internal/meta` hoy solo admite `X -> blocked`; salir
   de `blocked` exige `--force`. Añadir transiciones estándar `blocked -> in_progress`
   y `blocked -> ready`, y actualizar el **Caso H** del `## Contrato` de
   `.rei/agents/leader.md` para que el Leader guíe al usuario (retomar la
   implementación, volver a planificación/revisión de la planificación, o cancelar).
2. **Topes/estructura de reportes**: definir una estructura concisa y un **tope
   blando** (~40 líneas / ~350 palabras) para `impl.md` y `review.md`, documentarlo
   en los Contratos de `implementer`/`reviewer` y en
   `.rei/docs/harness/progress.md`, y emitir un **`WARN` (sin fallar)** en
   `validate` si se excede.

## Archivos

Código y tests:

- `internal/meta/meta.go`
- `internal/meta/meta_test.go`
- `internal/validate/validate.go`
- `internal/validate/validate_test.go`

Contratos y documentación:

- `.rei/agents/leader.md` (Caso H)
- `.rei/agents/implementer.md` (estructura/tope de `impl.md`)
- `.rei/agents/reviewer.md` (estructura/tope de `review.md`)
- `.rei/docs/harness/progress.md` (secciones `impl.md` y `review.md`)

Generados nativos (no se editan a mano; se regeneran con el CLI):

- `.opencode/agents/*.md`
- `.claude/agents/*.md` (y/o `CLAUDE.md` si el adaptador lo regenera)

Fuera de alcance: `rei init --update`.

## Cambios

### 1. Transiciones desde `blocked` (`internal/meta/meta.go`)

- En `allowedTransitions`, añadir la entrada:
  `StatusBlocked: {StatusReady, StatusInProgress}`.
- Ajustar el comentario de `allowedTransitions` para reflejar que `blocked` puede
  reanudarse a `ready` o `in_progress` (y que `* -> blocked` sigue permitido por
  `TransitionAllowed`).
- No cambiar la firma de `TransitionAllowed` ni `SetStatus`.

### 2. Tests de transiciones (`internal/meta/meta_test.go`)

- Mover `{StatusBlocked, StatusInProgress}` de la lista `rejected` a `allowed`.
- Añadir `{StatusBlocked, StatusReady}` a `allowed`.
- Sustituir el caso rechazado eliminado por una transición no estándar que siga
  siéndolo, p. ej. `{StatusBlocked, StatusDone}`.

### 3. Aviso de tamaño de reportes (`internal/validate/validate.go`)

- Definir constantes del tope blando: `maxReportLines = 40` y
  `maxReportWords = 350`.
- Añadir un helper que, dado el contenido, devuelva `(lines, words)`
  (`lines` con `strings.Split` sobre el contenido sin el salto final; `words`
  con `len(strings.Fields(...))`).
- En `WorkItem`, tras las comprobaciones de presencia de `impl.md`/`review.md`,
  para cada archivo que exista en `p.WorkItemDir(id)`:
  - leerlo; si la lectura falla, no convertir el fallo en FAIL (no es el
    propósito), solo omitir el aviso;
  - si `lines > maxReportLines || words > maxReportWords`, añadir un `warn` del
    estilo `"%s excede el tope blando (~%d líneas / ~%d palabras): %d líneas, %d palabras"`.
- El aviso es **`WARN`**: nunca incrementa FAIL ni cambia el código de salida por
  sí mismo.

### 4. Test del aviso de tamaño (`internal/validate/validate_test.go`)

- Añadir un test (p. ej. `TestReportTooLargeWarns`) que cree un Work Item `task`
  en estado `review` con `plan.md`, un `impl.md` que exceda el tope (muchas
  líneas o palabras) y sin `review.md` (no requerido en `review`).
- Comprobar que aparece al menos un `WARN` de tamaño y que el número de `FAIL`
  es 0.

### 5. Caso H del Leader (`.rei/agents/leader.md`)

Reescribir el **Caso H — `status == blocked`** para que, tras informar el motivo,
guíe al usuario con las transiciones ahora estándar (sin `--force`):

- **Retomar la implementación** → `rei status set <id> in_progress`, luego relanzar
  `implementer`.
- **Volver a la planificación / revisar la planificación** →
  `rei status set <id> ready`, luego relanzar `spec_author` si procede.
- **Cancelar** → no existe estado de cancelación; dejar el Work Item en `blocked`
  y acordarlo con el usuario (no forzar transiciones a `done`).
- Mantener la regla de NO avanzar sin decisión explícita del usuario.

### 6. Estructura y tope en los Contratos de implementer/reviewer

- `.rei/agents/implementer.md`, en **Cierre**: indicar que `impl.md` sigue una
  estructura concisa (resumen, archivos modificados, cambios, verificación,
  observaciones) y un tope blando de ~40 líneas / ~350 palabras; `rei validate`
  avisa (`WARN`) si se excede.
- `.rei/agents/reviewer.md`, en **Cierre**: indicar que `review.md` sigue una
  estructura concisa (resultado, verificaciones, observaciones, acciones
  requeridas si aplica) y el mismo tope blando; `rei validate` avisa (`WARN`) si
  se excede.

### 7. Documentación del proyecto (`.rei/docs/harness/progress.md`)

- En las secciones `impl.md` y `review.md`, fijar la estructura mínima concisa y
  el tope blando (~40 líneas / ~350 palabras), indicando que superarlo produce un
  `WARN` en `rei validate`, nunca un fallo.

### 8. Regenerar nativos

- Ejecutar `rei init opencode` y `rei init claude` para regenerar
  `.opencode/agents/*` y `.claude/agents/*` desde los Contratos editados.
- Verificar después con `rei init opencode --check` y `rei init claude --check`
  (ambos deben devolver 0).

## Restricciones

- **NO** implementar ni tocar `rei init --update`.
- **NO** modificar el alcance descrito en `meta.json`.
- **NO** cambiar la firma pública de `TransitionAllowed`/`SetStatus` ni el
  formato de `meta.json`.
- El aviso de tamaño es **`WARN`, nunca `FAIL`**: no debe romper `rei check` ni
  `rei validate` en Work Items válidos.
- Los archivos `.opencode/agents/*` y `.claude/agents/*` son generados: **no
  editarlos a mano**; regenerarlos con el CLI.
- Los tests deben pasar con `go test ./...`; el repositorio no tiene checks
  declarados en `.rei/config.json`.

## Pasos

- [x] 1. Añadir `blocked -> ready` y `blocked -> in_progress` en
  `internal/meta/meta.go` y actualizar el comentario de `allowedTransitions`.
- [x] 2. Actualizar `internal/meta/meta_test.go` (mover/añadir casos permitidos y
  reemplazar el rechazado obsoleto).
- [x] 3. Implementar el tope blando y el `WARN` de tamaño en
  `internal/validate/validate.go` (constantes + helper + aviso para `impl.md` y
  `review.md`).
- [x] 4. Añadir el test del aviso de tamaño en
  `internal/validate/validate_test.go`.
- [x] 5. Actualizar el **Caso H** en `.rei/agents/leader.md`.
- [x] 6. Documentar estructura y tope en `.rei/agents/implementer.md` y
  `.rei/agents/reviewer.md`.
- [x] 7. Documentar estructura y tope en `.rei/docs/harness/progress.md`.
- [x] 8. Regenerar nativos con `rei init opencode` y `rei init claude`.
- [x] 9. Verificar: `go test ./...`, `rei check --quiet`, `rei validate
  <work-item-id>`, `rei init opencode --check` y `rei init claude --check`.

## Verificación

- `go test ./...` (o `make test`) pasa.
- `rei check --quiet` devuelve 0.
- `rei init opencode --check` devuelve 0 y `rei init claude --check` devuelve 0.
- Un `rei status set <id> in_progress` desde `blocked` funciona sin `--force`
  (cubierto por el test de `internal/meta`).
- `rei validate <work-item-id>` no emite `FAIL` y sí emite `WARN` cuando un
  `impl.md`/`review.md` supera el tope.
