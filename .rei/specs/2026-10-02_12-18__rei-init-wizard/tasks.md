# Tasks — rei init wizard

> Work Item: `2026-10-02_12-18__rei-init-wizard` (feature).
>
> Orden de ejecución. Cada tarea indica los requisitos que implementa.
> El Implementer marca `[x]` **inmediatamente** al terminar cada tarea.

---

## Fase 1 — Plantillas y scaffold reutilizable

- [x] **T1 — Añadir los marcadores `<!-- REI:PENDIENTE -->` a los documentos de personalización:** un marcador en `AGENTS.md` §2, otro en `AGENTS.md` §3, y al menos uno en `.rei/docs/project/architecture.md`, `.rei/docs/project/conventions.md` y `.rei/docs/project/verification.md`. Colocarlos junto al contenido de ejemplo, sin romper el Markdown ni el render. (R18)

- [x] **T2 — Extraer `check.EnsureStructure`** en `internal/check/check.go`: concentra los pasos 2–4 de `check.Run` (crear `.rei/specs/`, `.rei/progress/work-items/`, `current.md` e `history.md` desde plantilla, e inicializar git si falta). `check.Run` pasa a llamarla sin cambiar su salida ni su código de salida. Devolver `1` si no puede crear los archivos base; el fallo de git es aviso. (R2, R3, R4, R5, R6, R7)

## Fase 2 — Núcleo del CLI (`initwizard`)

- [x] **T3 — Crear `internal/initwizard/initwizard.go` con la tabla canónica `Docs` y `Pending`:** `marker = "<!-- REI:PENDIENTE -->"`; `Pending(p) ([]Doc, error)` lee cada documento y lo marca pendiente si contiene el marcador (archivo inexistente ⇒ pendiente). Incluir los pasos asociados y las etiquetas legibles. (R8, R19)

- [x] **T4 — Implementar `Init(p, out) int`:** llama a `check.EnsureStructure`, detecta pendientes, imprime cabecera, resultado del scaffold, estado por documento y el plan de pasos ordenado (tabla del diseño). Idempotente: no sobrescribe archivos existentes. Devuelve `0` en el camino correcto y `1` si falla la creación base. (R1, R2, R3, R4, R5, R6, R7, R9, R10, R11, R13)

- [x] **T5 — Implementar `Status(p, out) int`:** reporte de solo lectura de los documentos pendientes/completos; `0` si no queda ninguno, `1` si queda al menos uno. No escribe archivos. (R14, R15, R16)

- [x] **T6 — Tests de `internal/initwizard`** (`t.TempDir()` + `paths.Project{Root: ...}`):
  - detección: documento con marcador ⇒ pendiente; sin marcador ⇒ completo; `AGENTS.md` con un marcador ⇒ pendiente;
  - `Init` crea `.rei/specs/`, `.rei/progress/work-items/`, `current.md`, `history.md` y es idempotente (segunda ejecución no falla ni sobrescribe contenido modificado);
  - `Status` devuelve `1` con pendientes y `0` sin ellos, sin crear/alterar archivos. (R1–R19)

## Fase 3 — Integración con el CLI

- [x] **T7 — Añadir el dispatch de `init` en `internal/cli/cli.go`:** `cmdInit(rest)` con el modo por defecto (`initwizard.Init`) y el subcomando `status` (`initwizard.Status`). Argumentos no reconocidos ⇒ uso por `stderr` y código `2`. (R1, R12, R14, R17)

- [x] **T8 — Añadir la entrada `init` a la tabla de `internal/cli/help.go`:** `usage: "init [status]"`, descripción breve y detalle (propósito, subcomandos, códigos de salida `0`/`1`/`2` y relación con el rol `initializer`). (R36, R37)

- [x] **T9 — Tests de CLI** en `internal/cli/help_test.go` (y/o un test del dispatch): la ayuda general incluye `init`; `commandHelpText("init")` menciona `status` y el rol `initializer`; `printCommandHelp("init")` devuelve `0`; un comando `init` con argumento inválido devuelve `2`. (R12, R17, R36, R37)

## Fase 4 — Diagnóstico y rol `initializer`

- [x] **T10 — Extender `internal/doctor/doctor.go`:** sección de personalización que reporta cada documento pendiente como `[WARN]` y, si no hay ninguno, `[OK] Documentación del proyecto personalizada.`; no incrementa `failures`. Añadir los casos a `internal/doctor/doctor_test.go`. (R33, R34, R35)

- [x] **T11 — Crear `.rei/agents/initializer.md`** con frontmatter (`name: initializer`), precondición (hay marcadores pendientes), protocolo por paso (leer destino, inferir del repo, entrevista, sugerencias marcadas, interpretar, escribir en el formato de cada documento, eliminar marcadores, confirmar), reanudabilidad (primer pendiente vía `rei init status`), regla de no inventar, alcance de escritura limitado, actualización de `.rei/config.json` en el paso de verificación y cierre con `rei init status` + `rei doctor`. Prohibición explícita de implementar código. (R21, R22, R23, R24, R25, R26, R27, R28, R29, R30, R31, R32)

- [x] **T12 — Declarar `.rei/agents/initializer.md` en `check.RequiredFiles`** (`internal/check/check.go`) y actualizar los tests correspondientes si comprueban la lista. (R41)

- [x] **T13 — Integrar en `AGENTS.md`:** añadir la instrucción al Leader de delegar en el rol `initializer` cuando `rei init status` reporte pendientes, antes de continuar con otro trabajo. (R39)

- [x] **T14 — Verificar la invariante de `rei check`:** comprobar que `rei check --quiet` devuelve `0` con marcadores pendientes presentes; ajustar si algún paso los tratara como fallo. No añadir el chequeo de marcadores a `check.Run`. (R40)

## Fase 5 — Validación final

- [x] **T15 — Ejecutar la batería de validación** y corregir lo que falle:
  `gofmt -l internal` (sin salida), `make vet`, `make test`, `make build`,
  `rei check --quiet` (`exit 0`).
  Verificación manual:
  - `rei init` crea la estructura, informa git y muestra pendientes + plan;
  - ejecutarlo de nuevo no rompe nada;
  - `rei init status` devuelve `0`/`1` según haya marcadores;
  - `rei doctor` reporta los pendientes como `[WARN]` sin fallar;
  - `rei help` incluye `init` y `rei help init` documenta `status` y el rol. (R1–R41)

- [x] **T16 — Marcar como completadas (`[x]`) todas las tareas** de este archivo a medida que se terminan, dejando `tasks.md` al día. (trazabilidad)
