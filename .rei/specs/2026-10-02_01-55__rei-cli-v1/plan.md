# Plan — rei cli v1

## Objetivo

Completar la CLI v1 de REI Harness (sin `rei init` ni adaptadores de runtime) con
tres capacidades:

1. **`rei item show <id>`** — ficha de un Work Item en un solo bloque: metadatos
   (`id`, `title`, `type`, `status`, `created_at`, `base_commit`,
   `last_review_commit`), documentos presentes de planificación/impl/review,
   si es la sesión activa (estado y agente) y los hallazgos de
   `validate.WorkItem`. **Solo lectura**.
2. **`rei test`** — ejecuta **únicamente** los checks declarados en
   `.rei/config.json`, **sin efectos secundarios** (no crea estructura, no
   inicializa git, no valida integridad). Salida visible por check; `exit 0` si
   todos pasan, `1` si alguno falla; si no hay checks, informa y `exit 0`.
   No debe confundirse con los tests de REI (`make test`, desarrollo).
3. **Ayuda por comando** — `rei <cmd> --help` (o `-h`) muestra el uso y detalle
   de cada comando existente; `rei help` general se actualiza con los comandos
   nuevos (`test`, `item`).

Fuera de alcance: `rei init`, adaptadores de runtime, cambios de comportamiento
en comandos existentes, nuevas dependencias externas.

## Archivos

- `internal/show/show.go` — **nuevo**. Paquete de presentación de la ficha
  (`item show`). Aísla la lógica del CLI y permite tests unitarios.
- `internal/show/show_test.go` — **nuevo**. Tests del paquete.
- `internal/check/check.go` — **modificado**. Exportar el runner de un check y
  añadir `RunDeclared` (checks de `config.json` sin side effects).
- `internal/check/check_test.go` — **modificado**. Tests de `RunDeclared`.
- `internal/cli/cli.go` — **modificado**. Dispatch de `item` y `test`; se retira
  `printHelp` (se mueve a `help.go`).
- `internal/cli/help.go` — **nuevo**. Tabla única de comandos + ayuda general y
  por comando (evita duplicar texto).
- `internal/cli/help_test.go` — **nuevo**. Tests de la tabla/ayuda.
- `.rei/specs/2026-10-02_01-55__rei-cli-v1/plan.md` — este archivo.

No modificar otros archivos.

## Cambios

### 1. `internal/show/show.go` — `rei item show <id>`

Función pública:

```go
// WorkItem imprime la ficha de un Work Item y devuelve el código de salida.
// Es de solo lectura: no escribe archivos ni ejecuta comandos.
func WorkItem(p *paths.Project, id string, out io.Writer) int
```

Comportamiento:

1. Si no existe `p.MetaFile(id)`: `ERROR: el Work Item "<id>" no existe` en
   `stderr` y `return 1`.
2. `meta.Load`; si falla (JSON inválido): `ERROR: <err>` en `stderr` y `return 1`.
3. Imprime el bloque:
   - **Metadatos**: una línea por campo; `base_commit`/`last_review_commit`
     vacíos se muestran como `(vacío)`.
   - **Documentos**: según `m.Type`:
     - planning: `plan.md` (task) o `requirements.md`, `design.md`, `tasks.md`
       (feature);
     - `impl.md` y `review.md` en `p.WorkItemDir(id)`;
     - `spec.md` en `p.WorkItemDir(id)` (registro de bloqueo del Spec Author).
     Cada uno con `sí`/`no` (existencia con `os.Stat`).
   - **Sesión**: `state.ReadSession(p)`; si `s.Active() && s.WorkItem == id`
     → `activa (estado: <estado>, agente: <agente>)`; si no → `no es la sesión
     activa`; error de lectura → `no se pudo leer current.md`.
   - **Validación**: `validate.WorkItem(p, id)`; una línea por issue con
     `[<Level>] <Message>` (reutiliza `validate.LevelFail`/`LevelWarn`).
4. `return 0` si no hay `validate.LevelFail`; `return 1` si hay al menos uno.

Reutiliza: `meta`, `state`, `validate`, `paths`. No reimplementa validación ni
lectura de sesión.

Formato propuesto (ejemplo real de este Work Item):

```text
Work Item: 2026-10-02_01-55__rei-cli-v1
  title:              rei cli v1
  type:               task
  status:             ready
  created_at:         2026-10-02T01:55:00-05:00
  base_commit:        (vacío)
  last_review_commit: (vacío)

Documentos:
  plan.md             sí
  requirements.md     no
  design.md           no
  tasks.md            no
  impl.md             no
  review.md           no
  spec.md             no

Sesión: activa (estado: ready, agente: _—_)

Validación:
  (sin hallazgos)
Resultado: OK
```

Se listan siempre los siete documentos (existan o no) para que la ficha sea
comparativa entre Work Items.

### 2. `internal/check/check.go` — `rei test`

**Refactor sin cambiar `check.Run`:**

- Exportar el runner actual `runCheck` como
  `Exec(dir string, command []string, out io.Writer) bool` (misma semántica:
  comando vacío → `false`; `Dir = dir`; captura stdout+stderr; en fallo imprime
  la salida indentada). `check.Run` (paso 5) pasa a llamar `Exec`.
- Añadir:

```go
// RunDeclared ejecuta solo los checks de .rei/config.json, sin efectos
// secundarios (no crea estructura, no inicializa git, no valida integridad).
// NO es la suite de tests de REI (esa es `make test`).
func RunDeclared(p *paths.Project, out io.Writer) int
```

Lógica de `RunDeclared` (solo `config.Load` + `Exec`):

1. `config.Load(p.ReiDir()/config.json)`:
   - error → `[FAIL]  .rei/config.json inválido: <err>` y `return 1`;
   - `len(cfg.Checks) == 0` → `[INFO]  Sin checks configurados
     (.rei/config.json). Nada que ejecutar.` y `return 0`.
2. Por cada check (`desc = Description`, o `ID` si vacío):
   `Exec(p.Root, ch.Command, out)`; si pasa `[OK]  <desc>`, si falla
   `[FAIL]  <desc>` y `exit = 1` (la salida del comando la emite `Exec`).
3. `return exit`.

Prohibido en `RunDeclared`: `os.MkdirAll`, `template.*`, `gitx.*`,
`check.RequiredFiles`, `validate.*`, y llamar a `check.Run`.

### 3. `internal/cli/help.go` — ayuda por comando

Tabla única que alimenta la ayuda general y la de cada comando:

```go
type commandInfo struct {
    name   string   // "check"
    usage  string   // "check [--quiet]"
    short  string   // descripción de una línea
    detail []string // líneas adicionales (subcomandos, notas, códigos)
}

var commands = []commandInfo{ /* check, doctor, new, session, items, status,
                                 commit, validate, review-diff, test, item,
                                 version, help */ }
```

- `printHelp()` → cabecera + bloque `Comandos:` construido desde la tabla
  (`usage` + `short`), más la línea de uso `rei <comando> --help`.
- `commandHelpText(name string) (string, bool)` → texto de un comando a partir
  de la tabla (testeable sin escribir a stdout).
- `printCommandHelp(name string) int` → imprime `rei <name> — <short>`, `Uso:`
  y `detail`; si `name` no está en la tabla, imprime `comando desconocido: <name>`
  en `stderr` + ayuda general y `return 2`; si existe, `return 0`.

Detalle por comando (resumen):

| Comando | `usage` | `detail` clave |
|---------|---------|----------------|
| `check` | `check [--quiet]` | `--quiet, -q` omite los `[OK]`. |
| `doctor` | `doctor` | Solo lectura. |
| `new` | `new <id> <feature\|task> [title]` | Crea `meta.json` y carpeta de progreso. |
| `session` | `session [status]` | `start <id> <type>`, `archive`, `reset`. |
| `items` | `items status` | Lista IDs y estados. |
| `status` | `status set <id> <status> [--force]` | `--force` permite transiciones no estándar. |
| `commit` | `commit set <id> <base_commit\|last_review_commit>` | Registra HEAD de git. |
| `validate` | `validate [<id>]` | Sin `<id>`, usa la sesión activa. |
| `review-diff` | `review-diff <id> [--full]` | Paquete de revisión por diff. |
| `test` | `test` | Ejecuta `.rei/config.json`; **NO** es `make test`; sin side effects; `0`/`1`; sin checks → informa y `0`. |
| `item` | `item show <id>` | Ficha: metadatos, documentos, sesión y validación; solo lectura. |
| `version` | `version` | Muestra la versión. |
| `help` | `help [<comando>]` | Ayuda general o de un comando. |

### 4. `internal/cli/cli.go` — dispatch y flags de ayuda

- Al inicio de `Run`, tras `len(args) == 0`: obtener `cmd := args[0]`,
  `rest := args[1:]`; si algún elemento de `rest` es `--help`/`-h`, llamar a
  `printCommandHelp(cmd)` y devolver su código.
- `case "help"`: si `rest` no está vacío → `printCommandHelp(rest[0])`; si no →
  `printHelp()` y `return 0`.
- `case "test": return cmdTest(rest)`.
- `case "item": return cmdItem(rest)`.

Implementaciones:

```go
func cmdTest(args []string) int {
    if len(args) > 0 {
        fmt.Fprintln(os.Stderr, "uso: rei test")
        return 2
    }
    p, code := project()
    if p == nil { return code }
    return check.RunDeclared(p, os.Stdout)
}

func cmdItem(args []string) int {
    if len(args) != 2 || args[0] != "show" {
        fmt.Fprintln(os.Stderr, "uso: rei item show <id>")
        return 2
    }
    p, code := project()
    if p == nil { return code }
    return show.WorkItem(p, args[1], os.Stdout)
}
```

- Eliminar `printHelp` de `cli.go` (vive en `help.go`).
- `--help`/`-h` **no** debe romper usos existentes como `rei check --quiet`
  (solo actúan cuando aparecen como argumento tras el comando).

### 5. Tests

`internal/show/show_test.go` (estilo `t.TempDir()` / `paths.Project{Root: ...}`):

- Work Item `task` con `plan.md`, `impl.md` y sesión activa coincidente: la
  salida contiene `id`, `status`, `plan.md sí` y `Sesión: activa`; `return 0`.
- Work Item `feature` sin `requirements.md`/`design.md`/`tasks.md`: aparecen como
  `no` y hay `[FAIL]` → `return 1`.
- `id` inexistente → `return 1` y error por `stderr`.
- `meta.json` inválido → `return 1`.

`internal/check/check_test.go` (añadir):

- `RunDeclared` sin `.rei/config.json` (o `checks: []`) → `return 0` e imprime
  `Sin checks configurados`.
- Un check que pasa (`[]string{"sh", "-c", "exit 0"}`) → `return 0` y `[OK]`.
- Un check que falla (`[]string{"sh", "-c", "exit 1"}`) → `return 1` y `[FAIL]`.
- Verificar que `RunDeclared` **no** crea `.rei/specs/` ni `current.md`
  (comprobar con `os.Stat` tras la llamada). Los comandos de prueba asumen
  entorno POSIX (dev; el cross-build a Windows no ejecuta tests).

`internal/cli/help_test.go`:

- `commandHelpText` devuelve texto para `check`, `doctor`, `new`, `session`,
  `items`, `status`, `commit`, `validate`, `review-diff`, `test`, `item`.
- `printCommandHelp("no-existe")` → `2`.
- La ayuda de `test` menciona `make test` (para evitar la confusión).
- La ayuda general (`commandHelpText` de la tabla) incluye `test` e `item`.

## Códigos de salida

| Comando | Situación | Código |
|---------|-----------|--------|
| `item show` | ficha mostrada sin `FAIL` de validación | `0` |
| `item show` | Work Item inexistente, `meta.json` ilegible/inválido, o ≥1 `FAIL` | `1` |
| `item show` | uso incorrecto (falta `show`/`<id>` o sobran args) | `2` |
| `test` | todos los checks pasan, o no hay checks | `0` |
| `test` | algún check falla, o `.rei/config.json` inválido | `1` |
| `test` | argumentos extra | `2` |
| `rei <cmd> --help` | comando conocido | `0` |
| `rei <cmd> --help` | comando desconocido | `2` |

## Restricciones

- `item show` es **solo lectura**: prohibido `os.WriteFile`, `os.MkdirAll`,
  git mutante o ejecutar checks.
- `test` **no** puede crear estructura, inicializar git, validar integridad ni
  llamar a `check.Run`; solo `config.Load` + `Exec`.
- `test` **no** son los tests de REI: dejar explícito en comentarios y ayuda que
  estos son `make test` (dev, requiere Go).
- No cambiar la salida ni el comportamiento de los comandos existentes
  (`check`, `doctor`, `new`, `session`, `items`, `status`, `commit`, `validate`,
  `review-diff`), salvo la ayuda general (consolidada en tabla).
- Reutilizar antes que duplicar: `meta`, `state`, `validate`, `config`, `check`.
- Sin dependencias externas nuevas (solo stdlib y paquetes internos).
- Mensajes y ayuda en español, coherentes con el resto del CLI.
- `--help`/`-h` se interpretan como petición de ayuda cuando aparecen tras el
  comando.

## Validaciones

Ejecutar desde la raíz con
`PATH="$HOME/.local/bin:/tmp/opencode/go/bin:$PATH"`:

- `gofmt -l internal` → sin salida.
- `make vet` → OK.
- `make test` → OK (sin regresiones; nuevos tests pasan).
- `make build` → OK.
- `rei check --quiet` → `exit 0`.
- Manual:
  - `rei item show 2026-10-02_01-55__rei-cli-v1` muestra metadatos, documentos,
    sesión y validación; `rei item show no-existe` → `exit 1`.
  - `rei test` con `.rei/config.json` vacío → informa y `exit 0`; añadir
    temporalmente un check que falle y comprobar `exit 1` (revertir).
  - `rei <cmd> --help` para los 11 comandos y `rei help test`.
  - `rei help` lista `test` e `item`.

## Pasos

- [x] 1. Crear `internal/show/show.go` con `WorkItem(p, id, out) int` y helpers de existencia/impresión.
- [x] 2. Implementar la ficha: metadatos, documentos (planning según `type`, `impl.md`, `review.md`, `spec.md`), sesión y validación; aplicar códigos `0`/`1`.
- [x] 3. Crear `internal/show/show_test.go` con los cuatro casos descritos.
- [x] 4. En `internal/check/check.go`, exportar `Exec` (renombrar `runCheck`) y actualizar el paso 5 de `Run`.
- [x] 5. Añadir `RunDeclared(p, out) int` a `internal/check/check.go` usando solo `config.Load` + `Exec`.
- [x] 6. Añadir a `internal/check/check_test.go` los tests de `RunDeclared` (sin checks, check OK, check FAIL, sin side effects).
- [x] 7. Crear `internal/cli/help.go` con `commandInfo`, `commands`, `printHelp`, `commandHelpText` y `printCommandHelp` (incluyendo `test` e `item`).
- [x] 8. Crear `internal/cli/help_test.go` con los casos de tabla y ayuda.
- [x] 9. En `internal/cli/cli.go`: dispatch de `--help`/`-h` por comando, `case "help"` con subcomando, `case "test"`, `case "item"`, `cmdTest` y `cmdItem`; importar `internal/show`; eliminar `printHelp`.
- [x] 10. Ejecutar las validaciones (`gofmt`, `make vet`, `make test`, `make build`, `rei check --quiet`) y la verificación manual; corregir lo que falle.
- [x] 11. Marcar cada paso como completado (`[x]`) en este archivo a medida que se termina.
