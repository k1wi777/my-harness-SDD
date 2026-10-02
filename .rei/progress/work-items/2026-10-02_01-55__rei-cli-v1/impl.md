# Implementación — `rei cli v1`

- **Work Item:** `2026-10-02_01-55__rei-cli-v1`
- **Tipo:** task
- **Estado:** review
- **Agente:** implementer
- **Fecha:** 2026-10-02
- **Plan:** `plan.md` (11 pasos, todos `[x]`)

## Resumen

Implementados los tres comandos de la CLI v1 según `plan.md`:

1. `rei item show <id>` — ficha de solo lectura de un Work Item (metadatos,
   documentos, sesión activa y hallazgos de `validate`). Nuevo paquete
   `internal/show`.
2. `rei test` — ejecuta **solo** los checks de `.rei/config.json`, sin efectos
   secundarios, con salida visible y exit `0`/`1`; sin checks informa y sale `0`.
   Se reutiliza el runner de `internal/check` exportado como `Exec` y se añade
   `RunDeclared`.
3. Ayuda por comando `rei <cmd> --help` y `rei help <comando>`, con la ayuda
   general consolidada en una tabla única (`internal/cli/help.go`) que ahora
   incluye `test` e `item`.

No se modificó la planificación; solo se marcaron los pasos completados. Sin
dependencias externas nuevas (solo stdlib y paquetes internos).

## Archivos

- `internal/show/show.go` — **nuevo**. `WorkItem(p, id, out) int` + helpers.
- `internal/show/show_test.go` — **nuevo**. 4 casos (task activa, feature sin
  planificación, id inexistente, `meta.json` inválido).
- `internal/check/check.go` — **modificado**. `runCheck` → `Exec`; nuevo
  `RunDeclared`; `Run` sigue llamando a `Exec`.
- `internal/check/check_test.go` — **modificado**. Tests de `RunDeclared`
  (sin checks, `checks: []`, check OK, check FAIL, sin side effects).
- `internal/cli/help.go` — **nuevo**. `commandInfo`, `commands`, `helpText`,
  `printHelp`, `commandHelpText`, `printCommandHelp`.
- `internal/cli/help_test.go` — **nuevo**. Casos de tabla/ayuda.
- `internal/cli/cli.go` — **modificado**. Dispatch de `--help`/`-h`, `help`
  con subcomando, `test`, `item`; `cmdTest`/`cmdItem`; import `internal/show`;
  se retiró `printHelp` (ahora en `help.go`).
- `.rei/specs/2026-10-02_01-55__rei-cli-v1/plan.md` — pasos marcados `[x]`.

## Cambios

### `internal/show/show.go`

- `WorkItem(p, id, out) int`: comprueba existencia de `meta.json` (si no,
  `ERROR: el Work Item "..." no existe` por stderr y `1`), carga con
  `meta.Load` (JSON inválido → `ERROR` y `1`) e imprime:
  - **Metadatos**: `title`, `type`, `status`, `created_at`, con
    `base_commit`/`last_review_commit` vacíos mostrados como `(vacío)`.
  - **Documentos**: los siete en orden fijo (`plan.md`, `requirements.md`,
    `design.md`, `tasks.md`, `impl.md`, `review.md`, `spec.md`) con `sí`/`no`.
    Los de planificación se resuelven en `.rei/specs/<id>/` (`docPath`) y el
    resto en `.rei/progress/work-items/<id>/`.
  - **Sesión**: `state.ReadSession`; `activa (estado: ..., agente: ...)` si
    coincide el Work Item, `no es la sesión activa` si no, o
    `no se pudo leer current.md` ante error.
  - **Validación**: `validate.WorkItem`, una línea `[LEVEL] mensaje` por hallazgo
    o `(sin hallazgos)`; `Resultado: OK`/`FALLO`.
- Código de salida `1` si hay al menos un `validate.LevelFail`; `0` en caso
  contrario. Solo lectura: sin `os.WriteFile`, `os.MkdirAll` ni ejecución de
  comandos.

### `internal/check/check.go`

- `Exec(dir string, command []string, out io.Writer) bool`: runner exportado
  (antes `runCheck`) con la misma semántica (comando vacío → `false`;
  `Dir = dir`; captura stdout+stderr; en fallo imprime la salida indentada).
  `Run` (paso 5) pasa a llamar `Exec`.
- `RunDeclared(p, out) int`: usa solo `config.Load` + `Exec`. Config inválida →
  `[FAIL]  .rei/config.json inválido: ...` y `1`; sin checks →
  `[INFO]  Sin checks configurados (.rei/config.json). Nada que ejecutar.` y
  `0`; por check, `[OK]`/`[FAIL]` con `Description` (o `ID`) y `exit = 1` si
  alguno falla. Prohibido: `os.MkdirAll`, `template.*`, `gitx.*`,
  `check.RequiredFiles`, `validate.*` y `check.Run`.

### `internal/cli/help.go`

- Tabla única `commands` (`commandInfo{name, usage, short, detail}`) con 13
  comandos (`check`, `doctor`, `new`, `session`, `items`, `status`, `commit`,
  `validate`, `review-diff`, `test`, `item`, `version`, `help`).
- `helpText()` construye la ayuda general (uso + `usage` + `short`) y
  `printHelp()` la imprime; la ayuda general lista `test` e `item`.
- `commandHelpText(name)` devuelve el texto de un comando (testeable, sin
  escribir a stdout): cabecera `rei <name> — <short>`, `Uso:` y `detail`.
- `printCommandHelp(name)` imprime la ayuda; si el comando no existe, escribe
  `comando desconocido: <name>` por stderr, muestra la ayuda general y devuelve
  `2`.
- La ayuda de `test` deja explícito que **no** es la suite de REI (`make test`).

### `internal/cli/cli.go`

- `Run`: tras `len(args) == 0`, calcula `cmd := args[0]`, `rest := args[1:]`; si
  algún elemento de `rest` es `--help`/`-h`, llama a `printCommandHelp(cmd)`.
- `case "help"`: con `rest` no vacío → `printCommandHelp(rest[0])`; sin él →
  `printHelp()` y `0`.
- `case "test": return cmdTest(rest)` y `case "item": return cmdItem(rest)`.
- `cmdTest`: argumentos extra → `uso: rei test` y `2`; si no,
  `check.RunDeclared(p, os.Stdout)`.
- `cmdItem`: exige `show <id>` (si no, `uso: rei item show <id>` y `2`); si no,
  `show.WorkItem(p, args[1], os.Stdout)`.
- Eliminado `printHelp` de `cli.go` (vive en `help.go`); import `internal/show`.

## Verificación

Validaciones ejecutadas desde la raíz (todas OK):

```text
gofmt -l internal    -> sin salida
make vet             -> OK
make test            -> OK (todos los paquetes, sin regresiones)
make build           -> OK (bin/rei regenerado)
rei check --quiet    -> exit 0
```

Verificación manual:

**(a) `item show`**

```text
rei item show 2026-10-02_01-55__rei-cli-v1   -> exit 0
  Work Item: 2026-10-02_01-55__rei-cli-v1
    title: rei cli v1 / type: task / status: in_progress
    base_commit: e4680010c952f682b6830408be4396b4b5864665
    last_review_commit: (vacío)
  Documentos: plan.md sí; requirements/design/tasks/impl/review/spec no
  Sesión: activa (estado: in_progress, agente: implementer)
  Validación: (sin hallazgos) / Resultado: OK

rei item show no-existe                        -> exit 1
  ERROR: el Work Item "no-existe" no existe
```

**(b) `rei test`** (config revertido a `{ "checks": [] }` tras la prueba)

```text
rei test (vacío)          -> exit 0  [INFO]  Sin checks configurados (...). Nada que ejecutar.
rei test (check temporal) -> exit 1  [FAIL]  check temporal que falla
rei test (restaurado)     -> exit 0  [INFO]  Sin checks configurados (...). Nada que ejecutar.
```

**(c) Ayuda por comando**

```text
rei <cmd> --help para check, doctor, new, session, items, status, commit,
  validate, review-diff, test, item, version, help  -> exit 0
rei no-existe --help                                -> exit 2
rei help test / rei help item                       -> exit 0
rei help no-existe                                  -> exit 2
```

**(d) `rei help`** lista la ayuda general con `test` e `item` (`exit 0`).

## Observaciones

- No se cambió el comportamiento de los comandos existentes; solo se consolidó
  la ayuda general en la tabla de `help.go` y se añadió el dispatch de ayuda.
- `--help`/`-h` solo se interpretan como ayuda cuando aparecen tras el comando,
  por lo que `rei check --quiet` sigue funcionando igual.
- `item show` y `test` son de solo lectura/sin efectos secundarios; se verificó
  en tests que `RunDeclared` no crea `.rei/specs/` ni `current.md`.
- Los tests de `RunDeclared` asumen entorno POSIX (`sh -c`); el cross-build a
  Windows no ejecuta tests, coherente con lo previsto en el plan.
- `rei test` no sustituye a `make test`; la ayuda lo deja explícito.
