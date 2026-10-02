# Implementación — rei init wizard

> Work Item: `2026-10-02_12-18__rei-init-wizard` (feature).
> Implementado por el rol `implementer` siguiendo `requirements.md`, `design.md`
> y `tasks.md` (T1–T16, todas completadas).

---

## Resumen

Se implementó el wizard `rei init` en dos niveles:

1. **Nivel determinista (CLI).** Nuevo paquete `internal/initwizard` y comando
   `rei init [status]` que crea la estructura base (reutilizando
   `check.EnsureStructure`), detecta los documentos de personalización
   pendientes (marcador `<!-- REI:PENDIENTE -->`) y muestra el plan de pasos.
   No usa IA, red ni dependencias externas salvo `git` (opcional).
2. **Nivel IA (rol `initializer`).** Nuevo `.rei/agents/initializer.md` con la
   entrevista guiada, reanudabilidad por marcadores, sugerencias explícitas,
   alcance de escritura limitado y cierre con `rei init status` + `rei doctor`.

El estado del wizard vive en los marcadores de los propios documentos (sin
archivo de estado), tal y como fija el diseño (D2).

---

## Archivos

### Nuevos

| Archivo | Propósito |
|---------|-----------|
| `internal/initwizard/initwizard.go` | `Docs`, `Pending`, `Init`, `Status` y tabla canónica del plan de pasos. |
| `internal/initwizard/initwizard_test.go` | Tests de detección, scaffold idempotente y `Status`. |
| `.rei/agents/initializer.md` | Rol de IA para la personalización guiada. |

### Modificados

| Archivo | Cambio |
|---------|--------|
| `internal/check/check.go` | Extraída `EnsureStructure` (pasos 2–4 de `Run`) y añadido `.rei/agents/initializer.md` a `RequiredFiles`. |
| `internal/check/check_test.go` | Tests de `EnsureStructure` (OK e idempotencia). |
| `internal/cli/cli.go` | Import de `initwizard` y dispatch `init` / `init status` (`cmdInit`). |
| `internal/cli/help.go` | Entrada `init [status]` en la tabla de ayuda. |
| `internal/cli/help_test.go` | Tests de ayuda y de argumentos inválidos de `init`. |
| `internal/doctor/doctor.go` | Sección de personalización (WARN, sin incrementar fallos). |
| `internal/doctor/doctor_test.go` | Tests de pendientes (WARN) y de documentación personalizada. |
| `AGENTS.md` | Marcadores en §2 y §3; instrucción al Leader de delegar en `initializer`. |
| `.rei/docs/project/architecture.md` | Marcador `<!-- REI:PENDIENTE -->`. |
| `.rei/docs/project/conventions.md` | Marcador `<!-- REI:PENDIENTE -->`. |
| `.rei/docs/project/verification.md` | Marcador `<!-- REI:PENDIENTE -->`. |

---

## Cambios por tarea

- **T1 — Marcadores (R18).** `AGENTS.md` con dos marcadores (uno en §2 y otro en
  §3) y un marcador en cada uno de `architecture.md`, `conventions.md` y
  `verification.md`, colocados junto al contenido de ejemplo.
- **T2 — `check.EnsureStructure` (R2–R7).** Nueva función exportada que concentra
  los pasos 2–4 de `check.Run`. Para preservar exactamente la salida de
  `rei check --quiet`, se implementó un núcleo privado `ensureStructure(p, out,
  quiet)`; `EnsureStructure` es el envoltorio público (usado por `initwizard`) y
  `check.Run` llama al núcleo con su flag `quiet`. El fallo de `git init` es un
  aviso, no un error.
- **T3–T5 — `internal/initwizard` (R8, R9, R10, R11, R13, R14, R15, R16, R19).**
  `marker`, `Docs`, `Pending`, `Init` y `Status`. `Pending` considera pendiente un
  archivo inexistente para que el listado sea completo. `Init` no sobrescribe
  archivos existentes (idempotente); `Status` es de solo lectura.
- **T6 — Tests del paquete.** Detección con/sin marcador, `AGENTS.md` con un
  marcador, archivo inexistente, `Init` crea la estructura y no la sobrescribe,
  y `Status` con/sin pendientes.
- **T7–T9 — CLI (R1, R12, R14, R17, R36, R37).** `cmdInit` despacha `status` o el
  modo por defecto; argumentos no reconocidos → uso por `stderr` y código `2`.
  Entrada `init [status]` en la ayuda con propósito, subcomandos, códigos de
  salida y relación con el rol `initializer`.
- **T10 — `rei doctor` (R33–R35).** Reporta cada pendiente como
  `[WARN] Personalización pendiente: <ruta>` y, sin pendientes, `[OK]
  Documentación del proyecto personalizada.` No incrementa `failures`.
- **T11 — Rol `initializer` (R21–R32).** Frontmatter, precondición (hay
  marcadores), protocolo por paso, reanudabilidad, sugerencias marcadas, formato
  por documento, alcance de escritura, actualización de `.rei/config.json` y
  cierre con `rei init status` + `rei doctor`; prohibición de implementar código.
- **T12 — `RequiredFiles` (R41).** Añadido `.rei/agents/initializer.md`.
- **T13 — `AGENTS.md` (R39).** Instrucción al Leader de delegar en `initializer`
  cuando `rei init status` reporte pendientes.
- **T14 — Invariante `rei check` (R40).** `rei check --quiet` sigue devolviendo
  `0` con marcadores presentes; no se añadió ningún chequeo de marcadores a
  `check.Run`.

---

## Verificación y evidencia

Checkpoints del proyecto (`verification.md`) aún sin personalizar (`V1`, `V2`
son placeholders) y `.rei/config.json` sin checks declarados. Se ejecutó la
batería de `tasks.md` T15:

| Check | Comando | Resultado |
|-------|---------|-----------|
| Formato | `gofmt -l internal` | sin salida |
| Vet | `make vet` | `vet_exit=0` |
| Tests | `make test` | `test_exit=0` (todos los paquetes `ok`) |
| Build | `make build` | `build_exit=0` |
| Checkpoint rápido | `rei check --quiet` | `check_exit=0` |

### Verificación manual (T15)

- `rei init`: crea/verifica la estructura, informa `Repositorio git detectado.`,
  lista los 4 documentos como `[PENDIENTE]` y muestra los 6 pasos del plan.
  `init_exit=0`.
- `rei init` por segunda vez: misma salida, sin romper ni sobrescribir nada
  (`init2_exit=0`).
- `rei init status`: reporta los 4 pendientes y devuelve `1`
  (`status_exit=1`). En un proyecto temporal sin marcadores devuelve `0` y
  "Documentación del proyecto personalizada.".
- `rei doctor`: reporta los 4 pendientes como `[WARN]`, termina con
  `[OK] Sin fallos. Harness íntegro.` y devuelve `0` (`doctor_exit=0`).
- `rei help`: incluye `init [status]`. `rei help init`: documenta `status`, los
  códigos `0`/`1`/`2` y el rol `initializer` (`helpinit_exit=0`).
- `rei init bogus`: muestra `uso: rei init [status]` por `stderr` y devuelve `2`
  (`bogus_exit=2`).

---

## Observaciones

- **Diseño vs. detalle de implementación (T2).** El diseño da la firma
  `EnsureStructure(p, out)` y señala que «`check.Run` conserva su lógica de
  supresión si fuera necesario». Para no cambiar la salida de
  `rei check --quiet`, el núcleo compartido acepta un flag `quiet` interno y
  `EnsureStructure` es el envoltorio público sin supresión. No se modificó
  ninguna firma planificada ni el comportamiento de `check.Run`.
- **Repositorio plantilla.** Al añadir los marcadores, este repositorio queda
  como plantilla pendiente de personalización: `rei init status` devuelve `1` y
  `rei doctor` muestra WARN hasta que el rol `initializer` complete los pasos.
  Es el estado esperado según R18/D4.
- No se modificó `.rei/config.json` (no hay checkpoints rápidos definidos),
  `.rei/docs/harness/*` ni otros roles.
- `git add` de los archivos nuevos del Work Item (`internal/initwizard/` y
  `.rei/agents/initializer.md`). No se añadió `.rei/progress/` ni artefactos
  ajenos al Work Item.
