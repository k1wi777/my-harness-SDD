# Implementación — `rei status set`

- **Work Item:** `2026-10-02_00-33__rei-status-set`
- **Tipo:** task
- **Estado:** review
- **Agente:** implementer
- **Fecha:** 2026-10-02

## Resumen

Implementado el comando `rei status set <id> <status> [--force]` según
`plan.md`. El comando cambia el `status` de un Work Item en
`.rei/specs/<id>/meta.json` leyendo y reescribiendo la estructura completa
(`meta.Load` + `meta.Save`), de modo que **todos** los campos restantes
(incluidos `base_commit` y `last_review_commit`) se preservan. Valida el id, el
estado destino y la transición estándar de `workflow.md`; `--force` permite
forzar transiciones no estándar. Se añadió la ayuda de `rei help`.

## Archivos

- `internal/meta/meta.go` — `allowedTransitions`, `TransitionAllowed`, `SetStatus`.
- `internal/meta/meta_test.go` — `TestTransitionAllowed`, `TestSetStatusPreservesFields`.
- `internal/cli/cli.go` — dispatch `case "status"`, `cmdStatus`, imports `meta`/`strings`, ayuda.

## Cambios

### `internal/meta/meta.go`

- `allowedTransitions`: tabla de transiciones estándar derivada de
  `workflow.md` (`pending→ready`, `ready→in_progress`,
  `in_progress→review`, `review→{done,changes_requested}`,
  `changes_requested→in_progress`).
- `TransitionAllowed(from, to)`: permite las transiciones de la tabla y
  cualquier `* → blocked`.
- `SetStatus(path, status)`: `Load` + asignación + `Save`; devuelve el estado
  anterior. No edita texto ni usa expresiones regulares.

### `internal/meta/meta_test.go`

- `TestTransitionAllowed`: cubre todas las transiciones permitidas, rechaza
  `pending→in_progress`, `done→in_progress`, `blocked→in_progress`,
  `ready→done`, y comprueba que todo estado puede pasar a `blocked`.
- `TestSetStatusPreservesFields`: guarda un `Meta` con `BaseCommit` y
  `LastReviewCommit` no vacíos, aplica `SetStatus` y verifica el nuevo estado
  y la preservación de `id`, `title`, `description`, `type`, `created_at`,
  `base_commit` y `last_review_commit`.

### `internal/cli/cli.go`

- `Run`: `case "status": return cmdStatus(args[1:])`.
- `cmdStatus`: parsea `set <id> <status> [--force]`; valida estado
  (`meta.IsValidStatus`), existencia del Work Item, no-op, transición
  estándar y `--force`; solo omite la validación de transición al forzar.
- `printHelp`: entrada `status set <id> <status> [--force]` tras `items status`.

## Verificación

Comandos (todo OK):

```text
gofmt -l .          -> sin salida
go build ./...      -> OK
go test ./...       -> OK (todos los paquetes)
make vet            -> OK
make build          -> bin/rei regenerado
rei check --quiet   -> exit 0
```

Verificación manual con Work Item temporal `2026-10-02_00-50__tmp-status`
(creado con `rei new`, con `base_commit` y `last_review_commit` no vacíos,
eliminado después):

```text
1) pending -> ready (permitida)                [0] pending -> ready.
2) ready -> pending --force (no estándar)      [0] ready -> pending (forzado).
3) pending -> in_progress (sin --force)        [1] ERROR: transición no permitida...
4) pending -> in_progress --force              [0] pending -> in_progress (forzado).
5) in_progress -> foo (estado inválido)        [1] ERROR: estado inválido: "foo"
6) in_progress -> in_progress (no-op)          [0] ya está en 'in_progress'.
7) in_progress -> review (permitida)           [0] in_progress -> review.
8) id inexistente                              [1] ERROR: el Work Item ... no existe
9) flag desconocido (--bogus)                  [2] uso: rei status set ...
10) faltan argumentos                          [2] uso: rei status set ...
```

Tras todas las operaciones, `meta.json` del Work Item temporal conservaba
`id`, `title`, `description`, `type`, `created_at`, `base_commit`
(`726a6586...`) y `last_review_commit` (`deadbeef...`); solo cambió `status`.

El Work Item temporal se eliminó por completo:
`.rei/specs/2026-10-02_00-50__tmp-status/` y
`.rei/progress/work-items/2026-10-02_00-50__tmp-status/`.

## Observaciones

- No se creó ningún paquete nuevo; se reutiliza `internal/meta`.
- No se añadieron dependencias.
- El comando no toca `current.md`, ni git, ni ejecuta checks (fuera de alcance).
- `--force` nunca omite la validación de estado válido ni la de existencia del
  Work Item.
