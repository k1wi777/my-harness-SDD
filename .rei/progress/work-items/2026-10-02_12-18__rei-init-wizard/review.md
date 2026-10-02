# Revisión — rei init wizard

> Work Item: `2026-10-02_12-18__rei-init-wizard` (feature).
> Revisor: rol `reviewer`.
> Paquete: `rei review-diff 2026-10-02_12-18__rei-init-wizard` (exit 0, `base_commit`
> `dae4ad3a1b4c6784c7904c8d9d17e3b32390f5e2`).

---

## Estado final

**Aprobado.** Todos los requisitos (`R1`–`R41`) están cubiertos por la
implementación y verificados. `rei check --quiet` devuelve `0` con los marcadores
presentes (R40) y `rei validate` devuelve `OK`.

---

## Verificaciones ejecutadas

### Checkpoints de la batería (T15)

| Check | Comando | Resultado |
|-------|---------|-----------|
| Formato | `gofmt -l internal` | sin salida (`exit 0`) |
| Vet | `make vet` | `exit 0` |
| Tests | `make test` | `exit 0` (todos los paquetes `ok`) |
| Build | `make build` | `exit 0` |
| Checkpoint rápido | `rei check --quiet` | `exit 0` **con marcadores presentes** (R40) |
| Consistencia | `rei validate 2026-10-02_12-18__rei-init-wizard` | `Resultado: OK` (`exit 0`) |

### Verificación funcional (R1–R41)

- **`rei init` (R1–R13).** En este repositorio (con `.rei/` y plantillas) crea
  `.rei/specs/`, `.rei/progress/work-items/`, `current.md` e `history.md`,
  informa `Repositorio git detectado.`, lista los 4 documentos como `[PENDIENTE]`
  y muestra el plan de 6 pasos con su destino. `exit 0`. R12/R17: `rei init bogus`
  imprime `uso: rei init [status]` por `stderr` y devuelve `2`.
- **Idempotencia (R11).** Segunda ejecución `exit 0`; en un proyecto temporal se
  añadió `SENTINEL` a `current.md` y `rei init` no lo sobrescribió (se conservó).
- **Sin git (R7).** Ejecución con `PATH` sin `git`: emite `[WARN] git no
  disponible…`, continúa y devuelve `0`. R6: con `git` disponible en un proyecto
  sin repositorio, lo inicializa.
- **`rei init status` (R14–R17).** En este repo reporta 4 `[PENDIENTE]` y devuelve
  `1`; en un proyecto temporal con los 4 documentos sin marcador reporta
  `[COMPLETO]` y devuelve `0`, sin crear ni modificar archivos.
- **Marcadores (R18/R19).** `<!-- REI:PENDIENTE -->` presente en `AGENTS.md` §2
  (línea 33), `AGENTS.md` §3 (línea 68), `.rei/docs/project/architecture.md:8`,
  `.rei/docs/project/conventions.md:8` y `.rei/docs/project/verification.md:9`.
  `Pending` marca pendiente ⇔ existe el marcador.
- **`rei doctor` (R33–R35).** Reporta los 4 pendientes como
  `[WARN] Personalización pendiente: <ruta>`, termina en
  `[OK] Sin fallos. Harness íntegro.` y devuelve `0` (los WARN no incrementan
  fallos).
- **Ayuda (R36/R37).** `rei help` incluye `init [status]`; `rei help init`
  documenta propósito, subcomandos, códigos `0`/`1`/`2` y el rol `initializer`
  (`exit 0`).
- **Rol `initializer` (R20–R32).** `.rei/agents/initializer.md` existe y tiene
  frontmatter (`name: initializer`, `description`, `tools`). Su protocolo cubre:
  precondición vía `rei init status` (R21/R23), entrevista guiada en el orden del
  plan (R24), sugerencias marcadas explícitamente (R25), redacción en el formato
  de cada documento (R26), escritura + borrado de marcadores + confirmación
  (R20/R27), regla de no inventar (R28), alcance de escritura acotado (R29),
  actualización de `.rei/config.json` en verificación (R30), cierre con
  `rei init status` + `rei doctor` (R31) y prohibición de implementar código
  (R32).
- **Integridad (R41).** `.rei/agents/initializer.md` está en
  `check.RequiredFiles`. En una copia temporal sin ese archivo, `rei check`
  devuelve `1` (`[FAIL] Falta .rei/agents/initializer.md`) y `rei doctor` también
  falla.
- **Delegación (R39).** `AGENTS.md` §1 indica al Leader delegar primero en
  `initializer` cuando `rei init status` reporte pendientes.
- **Sin IA/red (R38).** El paquete `initwizard` usa solo librería estándar (más
  `git` opcional); el nivel IA es documental.

### Trazabilidad

`tasks.md` marca `[x]` T1–T16. El diff de `rei review-diff` corresponde a lo
planificado: `internal/initwizard/*`, `.rei/agents/initializer.md`,
`internal/{check,cli,doctor}`, `AGENTS.md` y los marcadores en los tres
documentos de proyecto. `internal/check` extrae `ensureStructure`/`EnsureStructure`
sin cambiar la salida de `check.Run` (el flag `quiet` se conserva interno, tal y
como documenta `impl.md`).

---

## Observaciones (no bloqueantes)

- **Archivo ajeno al Work Item.** `package-lock.json` aparece como no rastreado
  en el paquete de revisión, pero es previo al Work Item (fecha `sep 10 20:46`,
  contenido vacío `"packages": {}`) y no fue introducido por la implementación. No
  afecta a la aprobación.
- **`rei init` requiere `.rei/` preexistente.** En un directorio sin `.rei/`,
  `rei init` devuelve `1` con «no se encontró un proyecto REI (.rei/)». Es
  consistente con la plantilla (que siempre incluye `.rei/` y `.rei/templates/`),
  y `R2`/`R3` solo exigen crear `.rei/specs/` y `.rei/progress/work-items/` si
  faltan, lo que sí cumple.

---

## Acciones requeridas

Ninguna.
