# plan.md — pase final de pulido

## Objetivo

Cerrar cuatro puntos de pulido del CLI `rei` antes de dar por terminado el
proyecto:

1. `internal/validate`: evaluar el tope de tamaño de `impl.md`/`review.md`
   **solo en Work Items activos** (`in_progress`, `review`,
   `changes_requested`). Así `rei doctor` deja de emitir `WARN` por el
   histórico en `done`/`pending`/`ready`/`blocked`.
2. `internal/update`: mensaje claro y comprensible cuando la API responde
   **404 / no hay releases publicadas**, en lugar de
   `respuesta inesperada de GitHub: 404 Not Found`.
3. Documentar `rei update [--check]` en `README.md` y `.rei/docs/usage.md`.
4. Subir el tope blando de reportes a **~80 líneas / ~600 palabras** y reforzar
   la instrucción en los Contratos de `.rei/agents/implementer.md` y
   `.rei/agents/reviewer.md` (y `.rei/docs/harness/progress.md`); regenerar los
   adaptadores nativos y verificar `--check == 0`.

## Archivos

| Archivo | Cambio |
|---------|--------|
| `internal/validate/validate.go` | Tope solo para estados activos + subir constantes. |
| `internal/validate/validate_test.go` | Test: `done` con reporte grande NO emite `WARN`. |
| `internal/update/github.go` | Detectar 404 y devolver un error centinela claro. |
| `internal/update/update.go` | Manejar el centinela con mensaje comprensible y código 0. |
| `internal/update/update_test.go` | Test: API 404 → mensaje claro, código 0. |
| `README.md` | Documentar `rei update [--check]`. |
| `.rei/docs/usage.md` | Documentar `rei update [--check]`. |
| `.rei/agents/implementer.md` | Tope `~80`/`~600` en el Contrato. |
| `.rei/agents/reviewer.md` | Tope `~80`/`~600` en el Contrato. |
| `.rei/docs/harness/progress.md` | Tope `~80`/`~600` (dos referencias). |
| `.opencode/agents/*`, `.claude/agents/*`, `CLAUDE.md` | Regenerados vía `rei init`. |

## Cambios

### 1. `internal/validate/validate.go`

- Subir las constantes del tope blando:
  - `maxReportLines = 80`
  - `maxReportWords = 600`
- Añadir un helper `isActiveStatus(status string) bool` que devuelva `true`
  solo para `meta.StatusInProgress`, `meta.StatusReview` y
  `meta.StatusChangesRequested`.
- Envolver el bucle que mide `impl.md`/`review.md` con
  `if isActiveStatus(m.Status) { ... }`. El resto de la validación no cambia.
- El `WARN` sigue siendo blando (nunca `FAIL`) y usa las constantes, por lo que
  el texto del aviso se actualiza solo.

### 2. `internal/update/github.go`

- Declarar un error centinela, p. ej.
  `var errNoReleases = errors.New("no hay releases publicadas")`.
- En `fetchLatestRelease`, antes del `return` genérico de estado no-200:
  - si `resp.StatusCode == http.StatusNotFound`, devolver
    `nil, errNoReleases`.
  - cualquier otro estado conserva el error genérico actual.
- Importar `errors`.

### 3. `internal/update/update.go`

- Tras `fetchLatestRelease`, si `errors.Is(err, errNoReleases)`:
  - imprimir un mensaje comprensible, p. ej.
    `No hay releases publicadas todavía. Nada que actualizar.`
  - devolver **0**.
- **Justificación del código 0:** la ausencia de releases no es un fallo del
  comando ni del entorno; es un estado normal («nada que actualizar»),
  equivalente a «ya estás en la última versión». Devolver 1 reservaría el error
  para fallos reales (red, HTTP 5xx, JSON inválido), que mantienen su código 1.
- El resto de la lógica (dev, up-to-date, `--check`, Windows, assets) no cambia.

### 4. Documentación de `rei update`

- `README.md`: nueva sección bajo **Instalación** que explique `rei update`
  (descarga e instala la última release, verifica checksums, reemplazo atómico)
  y `rei update --check` (solo informa si hay una versión nueva). Indicar que
  `--check` devuelve 0 cuando no hay nada que actualizar.
- `.rei/docs/usage.md`: añadir una subsección (en la Sección 1, junto a la
  instalación) con el mismo contenido orientado a personas.

### 5. Tope blando en contratos y documentación

- `.rei/docs/harness/progress.md`: actualizar las dos menciones
  `~40 líneas / ~350 palabras` → `~80 líneas / ~600 palabras`.
- `.rei/agents/implementer.md`: en **Cierre** (paso 2) actualizar el tope a
  `~80 líneas / ~600 palabras`, manteniendo la instrucción de concreción.
- `.rei/agents/reviewer.md`: en **Cierre** (paso 1) actualizar el tope a
  `~80 líneas / ~600 palabras`.
- Regenerar adaptadores: `rei init opencode` y `rei init claude`. Los archivos
  `.opencode/agents/*`, `.claude/agents/*` y `CLAUDE.md` se generan desde los
  Contratos; **no editarlos a mano**.

## Restricciones

- No tocar `rei cost` ni la selección de modelo por rol (fuera de alcance).
- No modificar el alcance de `meta.json`.
- El `WARN` de tamaño debe seguir siendo blando: nunca `FAIL`.
- No editar a mano los archivos generados (`.opencode/`, `.claude/`,
  `CLAUDE.md`); usar `rei init`.
- No modificar reportes históricos de `.rei/progress/work-items/` ya cerrados.
- No introducir dependencias externas nuevas.

## Pasos

- [x] 1. `internal/validate/validate.go`: subir `maxReportLines` a 80 y
  `maxReportWords` a 600; añadir `isActiveStatus` y aplicar el tope solo a
  `in_progress`/`review`/`changes_requested` (punto 1).
- [x] 2. `internal/validate/validate_test.go`: añadir un test que con `status =
  done` y un `impl.md` grande **no** emita `WARN`; verificar que el test de
  `review` sigue emitiéndolo (punto 1).
- [x] 3. `internal/update/github.go`: definir `errNoReleases` y devolverlo en
  HTTP 404 (punto 2).
- [x] 4. `internal/update/update.go`: manejar `errNoReleases` con el mensaje
  «No hay releases publicadas todavía. Nada que actualizar.» y salida 0
  (punto 2).
- [x] 5. `internal/update/update_test.go`: añadir test con servidor httptest que
  responda 404 en `/releases/latest`; comprobar mensaje claro y código 0
  (punto 2).
- [x] 6. `README.md`: documentar `rei update [--check]` (punto 3).
- [x] 7. `.rei/docs/usage.md`: documentar `rei update [--check]` (punto 3).
- [x] 8. `.rei/docs/harness/progress.md`: actualizar el tope a ~80 líneas /
  ~600 palabras en las dos menciones (punto 4).
- [x] 9. `.rei/agents/implementer.md` y `.rei/agents/reviewer.md`: actualizar el
  tope e instrucción de concreción en sus Contratos (punto 4).
- [x] 10. Regenerar nativos con `rei init opencode` y `rei init claude`
  (punto 4).
- [x] 11. Verificar: `go test ./internal/validate/... ./internal/update/...`;
  `rei doctor` sin `WARN` de tamaño del histórico; `rei update --check` con
  mensaje claro; `rei init opencode --check` = 0 y `rei init claude --check` = 0;
  `rei check --quiet` = 0.
