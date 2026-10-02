# Tasks — rei init installer

> Work Item: `2026-10-02_13-21__rei-init-installer` (feature).
>
> Orden de ejecución. Cada tarea indica los requisitos que implementa.
> El Implementer marca `[x]` **inmediatamente** al terminar cada tarea.

---

## Fase 1 — Empaquetado del esqueleto

- [x] **T1 — Crear `embed.go` en la raíz del repositorio (`package rei`)** con la directiva `//go:embed all:.rei/docs all:.rei/agents all:.rei/templates all:.rei/config.json AGENTS.md` y la variable exportada `var Skeleton embed.FS`. Asegurar que el patrón incluye `AGENTS.md`, `.rei/config.json`, `.rei/docs/**`, `.rei/agents/**` y `.rei/templates/**`, y excluye `.rei/specs/` y `.rei/progress/`. (R1, R2, R4)

- [x] **T2 — Añadir `embed_test.go` (raíz, `package rei`)**: comprobar que `Skeleton` contiene `AGENTS.md`, `.rei/config.json`, al menos un archivo bajo `.rei/docs/`, `.rei/agents/` y `.rei/templates/`, y que `.rei/specs` y `.rei/progress` no existen en el `embed.FS`. (R1, R2)

## Fase 2 — Instalador en `initwizard`

- [x] **T3 — Implementar la resolución del destino** en `internal/initwizard/installer.go`: `ResolveRootFrom(dir) (*paths.Project, error)` reutiliza `paths.FindFrom` y, ante `paths.ErrNotFound`, devuelve `&paths.Project{Root: dir}`; `ResolveRoot()` aplica `ResolveRootFrom` sobre el directorio actual. (R5)

- [x] **T4 — Implementar `InstallSkeleton(p, out) int`** en `internal/initwizard/installer.go`: recorrer `reiskel.Skeleton` con `fs.WalkDir`, crear directorios con `MkdirAll`, escribir archivos con `0644` solo si no existen, omitir los existentes, emitir aviso específico para `AGENTS.md` existente y reportar creados/omitidos con un resumen. Devolver `1` si algún archivo o directorio no puede crearse; `0` en caso contrario. (R3, R6, R7, R8, R9, R10)

- [x] **T5 — Integrar el despliegue en `Init`** (`internal/initwizard/initwizard.go`): llamar a `InstallSkeleton` antes de `check.EnsureStructure`; conservar la detección de pendientes y el plan de pasos; devolver `1` si falla el despliegue, `EnsureStructure` o `Pending`. (R11, R12, R13, R14, R15, R16, R17, R18, R19, R22)

- [x] **T6 — Tests de `internal/initwizard`** (`installer_test.go`, con `t.TempDir()`):
  - `InstallSkeleton` crea todos los archivos del esqueleto y `Init` crea la estructura de estado;
  - una segunda ejecución no falla ni sobrescribe un archivo modificado manualmente (`AGENTS.md`, `.rei/config.json`);
  - `AGENTS.md` preexistente emite el aviso y no se toca;
  - el resumen refleja los conteos de creados/omitidos;
  - `ResolveRootFrom` con `.rei/` ancestro devuelve ese root y sin `.rei/` devuelve el directorio dado. (R3, R5, R6, R7, R8, R9, R10, R11, R12)

## Fase 3 — CLI y ayuda

- [x] **T7 — Actualizar `cmdInit`** (`internal/cli/cli.go`): el modo instalador obtiene el destino con `initwizard.ResolveRoot()` (no exige `.rei/` previo); `rei init status` sigue usando `project()`; se mantiene la validación de argumentos y el código `2` para uso incorrecto. (R5, R20, R21, R23)

- [x] **T8 — Actualizar la ayuda de `init`** (`internal/cli/help.go` y `internal/cli/help_test.go`): documentar el despliegue del esqueleto, la creación de estructura y git, y los códigos de salida. (R26)

## Fase 4 — Integridad y validación

- [x] **T9 — Añadir un test de integridad**: instalar en un directorio temporal vacío, verificar que existen todos los archivos de `check.RequiredFiles` más `.rei/config.json`, que `rei check` (o `check.Run`) no reporta faltantes, y que `rei init status` reporta la personalización como pendiente. (R24, R25)

- [x] **T10 — Ejecutar la batería de validación** y corregir lo que falle:
  `gofmt -l internal` (sin salida), `make vet`, `make test`, `make build`,
  `rei check --quiet` (`exit 0`).
  Verificación manual:
  - en un directorio vacío, `rei init` despliega el esqueleto, crea `.rei/specs/`, `.rei/progress/work-items/`, `current.md`, `history.md` e inicializa git;
  - ejecutar `rei init` de nuevo reporta archivos omitidos y no modifica nada;
  - `AGENTS.md` preexistente avisa y se conserva;
  - `rei init status` devuelve `1` (marcadores) y `rei help init` menciona el despliegue. (R1–R27)

- [x] **T11 — Marcar como completadas (`[x]`) todas las tareas** de este archivo a medida que se terminan, dejando `tasks.md` al día. (trazabilidad)
