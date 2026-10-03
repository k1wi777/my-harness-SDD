# impl.md

## Resumen

Correcciones medias: (1) desbloqueo estándar sin `--force`; (2) tope blando y
estructura de `impl.md`/`review.md` con `WARN` (nunca `FAIL`) en `rei validate`.

## Archivos modificados

- `internal/meta/meta.go`, `internal/meta/meta_test.go`
- `internal/validate/validate.go`, `internal/validate/validate_test.go`
- `.rei/agents/{leader,implementer,reviewer}.md`, `.rei/docs/harness/progress.md`
- Nativos regenerados: `.opencode/agents/*`, `.claude/agents/*`, `CLAUDE.md`

## Cambios

1. `allowedTransitions` admite `blocked -> ready` y `blocked -> in_progress`;
   tests actualizados (rechazado obsoleto -> `blocked -> done`).
2. `validate`: `maxReportLines=40`/`maxReportWords=350`, helper `reportSize` y
   `WARN` si `impl.md`/`review.md` exceden el tope; test `TestReportTooLargeWarns`.
3. Caso H del Leader reescrito con las transiciones estándar de desbloqueo.
4. Contratos de implementer/reviewer y `progress.md` documentan estructura y tope.

## Verificación

- `gofmt -l .` sin salida; `go vet ./...`, `make test`, `make build` OK.
- `rei check --quiet` = 0; `rei validate` = 0.
- `rei init opencode --check` = 0 y `rei init claude --check` = 0.
- Desbloqueo real (temporal): `blocked->in_progress` y `blocked->ready` sin
  `--force`, exit 0; temporal eliminado.
- Tope real: reporte de 45 líneas -> `[WARN]`, sin `FAIL`, exit 0.

## Observaciones

- Nativos regenerados con el CLI; no editados a mano.
- Fuera de alcance: `rei init --update` (no tocado).
