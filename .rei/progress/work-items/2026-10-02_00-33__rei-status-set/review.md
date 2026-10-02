# Revisión — `rei status set`

- **Work Item:** `2026-10-02_00-33__rei-status-set`
- **Tipo:** `task`
- **Estado final:** `done`
- **Base revisada:** `726a658661373a8caaec8134187061213739cb46`
- **Fecha:** 2026-10-02

## Alcance revisado

Paquete de `rei review-diff 2026-10-02_00-33__rei-status-set` (exit 0) y su
ampliación `--full`:

- `internal/meta/meta.go` — modificado (+38: `allowedTransitions`, `TransitionAllowed`, `SetStatus`).
- `internal/meta/meta_test.go` — modificado (+73: `TestTransitionAllowed`, `TestSetStatusPreservesFields`).
- `internal/cli/cli.go` — modificado (+71: imports, `case "status"`, `cmdStatus`, ayuda).
- `.rei/specs/2026-10-02_00-33__rei-status-set/{meta.json,plan.md}` — spec del Work Item.
- `.rei/progress/current.md` — actualizado por el workflow (no por el comando).
- `package-lock.json` — sin rastrear, preexistente (sep 10) y ajeno a este Work Item.

## Objetivo cumplido

El comando `rei status set <id> <status> [--force]` está implementado según
`plan.md`:

- **Preservación total de campos:** `meta.SetStatus` hace `Load` + asignación de
  `Status` + `Save` sobre el struct completo; no hay edición de texto ni
  expresiones regulares. `base_commit` y `last_review_commit` se conservan.
- **Validación de id:** si no existe `.rei/specs/<id>/meta.json` →
  `ERROR: el Work Item "<id>" no existe`, salida `1`.
- **Validación de estado:** `meta.IsValidStatus`; estado inválido → salida `1`.
- **Validación de transición:** `meta.TransitionAllowed` con la tabla estándar de
  `workflow.md` (`pending→ready`, `ready→in_progress`, `in_progress→review`,
  `review→{done,changes_requested}`, `changes_requested→in_progress`) y
  `*→blocked`. Transición no estándar sin `--force` → salida `1`.
- **`--force`:** solo omite la validación de transición (marca `(forzado)`); no
  omite la validación de estado ni la de existencia.
- **No-op idempotente:** mismo estado → `ya está en '<status>'`, salida `0`.
- **Uso/códigos:** subcomando distinto de `set`, argumentos insuficientes o flag
  desconocido → uso en `stderr` + salida `2`.
- **Ayuda:** `printHelp` incluye `status set <id> <status> [--force]`.

## Restricciones respetadas

- **No edita `current.md`:** el comando solo abre/reescribe
  `.rei/specs/<id>/meta.json`; verificado por hash de `current.md` antes/después
  de un `status set` (sin cambios, ver más abajo).
- **No toca git:** no hay llamadas a `gitx` ni a `git` en `cmdStatus`/`SetStatus`;
  `HEAD` permaneció en `726a6586...` tras todas las ejecuciones.
- **No ejecuta checks:** `cmdStatus` no llama a `check.Run`.
- **No crea paquetes ni dependencias:** reutiliza `internal/meta` (stdlib +
  paquetes internos).
- **No toca otros archivos/Work Items:** solo el `meta.json` del id indicado.
- **Estilo:** mensajes `ERROR:` a `stderr`, `uso:` para errores de uso, igual que
  el resto del CLI.

## Verificaciones (checkpoints y comandos solicitados)

`.rei/config.json` no declara checkpoints (`"checks": []`), por lo que el
checkpoint rápido efectivo es `rei check`.

| Verificación | Comando | Resultado |
|--------------|---------|-----------|
| Checkpoint rápido | `rei check --quiet` | exit `0` ✅ |
| Tests | `make test` | exit `0`; `internal/meta` ok ✅ |
| Vet | `make vet` | exit `0` ✅ |
| Build implícito | `go build ./...` (vía tests/vet) | OK ✅ |

### Verificación independiente (Work Item temporal)

Creado `2026-10-02_99-99__review-tmp` (`rei new`, con `base_commit` y
`last_review_commit` fijados a `726a6586...`), se ejecutaron los escenarios y se
eliminó después (spec y carpeta de progreso):

```text
pending -> ready (permitida)              [0] pending -> ready.
ready -> pending (sin --force)            [1] ERROR: transición no permitida...
ready -> pending --force (no estándar)    [0] ready -> pending (forzado).
pending -> in_progress (sin --force)      [1] ERROR: transición no permitida...
pending -> in_progress --force            [0] pending -> in_progress (forzado).
in_progress -> bogus (estado inválido)    [1] ERROR: estado inválido: "bogus"
in_progress -> in_progress (no-op)        [0] ya está en 'in_progress'.
in_progress -> review (permitida)         [0] in_progress -> review.
id inexistente                            [1] ERROR: el Work Item "no-existe" no existe
flag desconocido (--bogus)                [2] uso: rei status set ...
faltan argumentos                         [2] uso: rei status set ...
```

Tras todas las operaciones, el `meta.json` temporal conservaba `id`, `title`,
`description`, `type`, `created_at`, `base_commit`
(`726a658661373a8caaec8134187061213739cb46`) y `last_review_commit`
(`726a658661373a8caaec8134187061213739cb46`); solo cambió `status`. El Work Item
temporal se eliminó por completo.

Comprobación de efectos colaterales: `sha256` de `.rei/progress/current.md` y
`HEAD` sin variación tras ejecutar `rei status set`.

## Observaciones

- La tabla de transiciones coincide exactamente con `workflow.md` (incluida la
  rama `review → changes_requested` y `* → blocked`).
- `cmdStatus` ignora argumentos posicionales sobrantes (`set <id> <status> extra`),
  comportamiento no especificado en `plan.md`; no afecta a los requisitos ni a
  los códigos de salida documentados, por lo que no es bloqueante.
- `package-lock.json` figura como cambio sin rastrear ajeno al Work Item; no se
  le atribuye impacto.

## Acciones requeridas

Ninguna.

## Resultado

Aprobado. El objetivo se cumple, las restricciones se respetan, la verificación
independiente confirma el comportamiento esperado y todas las verificaciones
(`make test`, `make vet`, `rei check --quiet`) pasan.
