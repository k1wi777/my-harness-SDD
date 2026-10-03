# impl.md — pase final de pulido

## Resumen

Pase final de pulido (4 puntos) implementado según `plan.md`, los 11 pasos
completados. Type `task`.

1. El tope blando de tamaño de `impl.md`/`review.md` se evalúa **solo** en Work
   Items activos (`in_progress`, `review`, `changes_requested`); `done`,
   `pending`, `ready` y `blocked` dejan de medirse → desaparece el ruido de
   `rei doctor`.
2. `rei update` trata la ausencia de releases (HTTP 404) como estado normal:
   mensaje claro y salida 0.
3. `rei update [--check]` documentado en `README.md` y `.rei/docs/usage.md`.
4. Tope blando subido a ~80 líneas / ~600 palabras, reforzado en los Contratos
   de `implementer`/`reviewer` y en `progress.md`; adaptadores nativos
   regenerados.

## Archivos modificados

- `internal/validate/validate.go` — constantes 80/600, helper `isActiveStatus`,
  tope envuelto en `if isActiveStatus(m.Status)`.
- `internal/validate/validate_test.go` — `TestReportTooLargeDoneNoWarn`.
- `internal/update/github.go` — centinela `errNoReleases`, 404 → centinela.
- `internal/update/update.go` — maneja `errNoReleases` (mensaje + return 0).
- `internal/update/update_test.go` — `TestRunNoReleases` (httptest 404).
- `README.md`, `.rei/docs/usage.md` — documentación de `rei update [--check]`.
- `.rei/docs/harness/progress.md` — tope actualizado (2 menciones).
- `.rei/agents/implementer.md`, `.rei/agents/reviewer.md` — tope + concreción.
- `.opencode/agents/{implementer,reviewer}.md`, `.claude/agents/{implementer,reviewer}.md`
  — regenerados vía `rei init` (no editados a mano).

## Cambios

- `validate.go`: `maxReportLines = 80`, `maxReportWords = 600`; nuevo
  `isActiveStatus` que devuelve true solo para los tres estados activos; el
  bucle de medición queda condicionado a ese helper. El `WARN` sigue siendo
  blando (nunca `FAIL`).
- `github.go`: `var errNoReleases = errors.New("no hay releases publicadas")`;
  en `fetchLatestRelease`, `StatusNotFound` devuelve `nil, errNoReleases` antes
  del `return` genérico.
- `update.go`: si `errors.Is(err, errNoReleases)` imprime
  «No hay releases publicadas todavía. Nada que actualizar.» y devuelve 0. El
  resto de la lógica (dev, up-to-date, `--check`, Windows, assets) intacta.
- Docs y Contratos: contenido del punto 3 y tope ~80/~600 en el punto 4.

## Verificación

- `gofmt -l .` → sin salida (exit 0).
- `go vet ./...` → exit 0.
- `make test` → exit 0 (incluye `internal/validate` e `internal/update`).
- `make build` → exit 0.
- `rei doctor` → sin `WARN` de tamaño del histórico; solo quedan los
  `WARN` de personalización pendiente (esperados).
- `rei update --check` → «No hay releases publicadas todavía. Nada que
  actualizar.», exit 0.
- `rei init opencode --check` = 0 y `rei init claude --check` = 0.
- `rei check --quiet` → exit 0.

## Observaciones

- No se tocó `rei cost` ni la selección de modelo por rol (fuera de alcance).
- No se modificaron reportes históricos de Work Items cerrados.
- No se editaron a mano los archivos generados; se regeneraron con `rei init`.
- Sin dependencias externas nuevas (`errors` es de la stdlib).
