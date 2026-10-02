# Implementación — rei init installer

> Work Item: `2026-10-02_13-21__rei-init-installer` (feature).
> Implementer · estado final: `review`.

## Resumen

`rei init` deja de depender del repositorio fuente. El esqueleto del harness se
embebe en el binario con `//go:embed` (patrón selectivo que excluye
`.rei/specs` y `.rei/progress`) y `rei init` lo despliega en el proyecto
destino sin sobrescribir archivos existentes; después crea la estructura de
estado, inicializa git si falta y mantiene el reporte de personalización.

El destino se resuelve como el ancestro más cercano con `.rei/` o, si no
existe, el directorio de trabajo actual, de modo que se puede inicializar un
proyecto nuevo y re-ejecutar `rei init` desde un subdirectorio sin anidar
`.rei/`.

Implementación completa de T1–T11 según `design.md`. No se modificó ninguna
decisión de planificación y no se amplió el alcance (sin `--update` ni
adaptadores de runtime).

## Archivos

### Nuevos

| Archivo | Contenido |
|---------|-----------|
| `embed.go` (raíz, `package rei`) | `//go:embed all:.rei/docs all:.rei/agents all:.rei/templates all:.rei/config.json AGENTS.md` → `var Skeleton embed.FS`. |
| `embed_test.go` (raíz, `package rei`) | Contenido embebido (AGENTS.md, `.rei/config.json`, subárboles docs/agents/templates) y exclusión de `.rei/specs`/`.rei/progress`. |
| `internal/initwizard/installer.go` | `ResolveRootFrom`, `ResolveRoot` e `InstallSkeleton`. |
| `internal/initwizard/installer_test.go` | Despliegue, idempotencia, aviso de `AGENTS.md`, resumen creados/omitidos, resolución de destino e integridad tras instalación limpia. |

### Modificados

| Archivo | Cambio |
|---------|--------|
| `internal/initwizard/initwizard.go` | `Init` llama a `InstallSkeleton` antes de `check.EnsureStructure`; devuelve `1` si falla cualquiera de los dos (o `Pending`). El resto del reporte no cambia. |
| `internal/cli/cli.go` | `cmdInit`: el modo instalador obtiene el destino con `initwizard.ResolveRoot()` (ya no exige `.rei/` previo); `status` sigue con `project()`. Validación de argumentos y código `2` intactos. |
| `internal/cli/help.go` | La ayuda de `init` documenta despliegue del esqueleto, estructura de estado, git y códigos de salida. |
| `internal/cli/help_test.go` | Nueva comprobación de que la ayuda menciona el instalador (`esqueleto`, `git`, `Códigos de salida`, `.rei/config.json`). |

## Trazabilidad tareas → requisitos

- T1–T2 → R1, R2, R4.
- T3 → R5.
- T4 → R3, R6, R7, R8, R9, R10.
- T5 → R11–R19, R22.
- T6 → R3, R5–R12.
- T7 → R5, R20, R21, R23.
- T8 → R26.
- T9 → R24, R25.
- T10 → R1–R27 (validación).

## Verificación

### Automática

```text
$ gofmt -l .            # sin salida
$ make vet              # exit 0
$ make test             # exit 0 (todos los paquetes ok)
$ make build            # exit 0 → bin/rei
$ rei check --quiet     # exit 0
```

### Manual (directorios temporales reales)

a) Instalación limpia (`mktemp -d`, sin `.rei/`): `rei init` = `0`; crea los
   21 archivos del esqueleto (incluidos `AGENTS.md` y `.rei/config.json`),
   reporta `21 archivo(s) creado(s), 0 omitido(s).`, crea `.rei/specs/`,
   `.rei/progress/work-items/`, `current.md`, `history.md` e inicializa git.

b) `rei check --quiet` dentro del temporal = `0` (sin faltantes).

c) Segunda ejecución de `rei init`: reporta `0 archivo(s) creado(s), 21
   omitido(s).`; tras modificar `.rei/config.json` manualmente, el contenido del
   usuario se conserva (`{"checks":[{"id":"user",...}]}`).

d) `AGENTS.md` preexistente: emite `[WARN]  AGENTS.md ya existe; se conserva sin
   cambios.` junto al `[SKIP]`; el contenido propio (`# Mi AGENTS propio`) se
   conserva.

e) `rei init status` con marcadores = `1`; `rei init` muestra el "Plan de pasos"
   (6 pasos). `rei help init` menciona despliegue del esqueleto, estructura,
   git y códigos de salida.

f) Subdirectorio de un proyecto existente: `rei init` resuelve el root al
   ancestro con `.rei/` y no crea un `.rei/` anidado (`nested=NO`,
   `root_exists=SI`).

Los directorios temporales se eliminaron al terminar.

## Observaciones

- `.rei/progress/current.md` está en estado `review` con agente activo
  `implementer`.
- `base_commit` verificado: `c9a11cd0dfd5636e7aea4cdb3e75fa564731ee33` (coincide
  con `git rev-parse HEAD` al iniciar; no se hizo `git commit`).
- No se modificaron `.rei/docs/**`, `.rei/agents/**` ni `.rei/templates/**`.
- `status` de `initwizard` no se tocó: conserva su contrato (0/1 por marcadores),
  según R20/R21 y D10.

## Rework (2026-10-02)

El Reviewer devolvió el Work Item como `changes_requested` (detalle en
`review.md`): la implementación estaba funcionalmente completa y R1–R27
verificados, pero **T5** figuraba como `[ ]` en `tasks.md` pese a estar
implementada, contradiciendo a `T11` y a este documento. No hubo defectos de
código ni de alcance.

Acción aplicada (única requerida):

- En `.rei/specs/2026-10-02_13-21__rei-init-installer/tasks.md`, marcar **T5**
  como `[x]`; ahora T1–T11 están completadas (`- [ ]` no aparece en el archivo).
  Sin cambio de código.

Verificación del rework:

```text
$ rei check --quiet     # exit 0
$ make test             # exit 0
```

Estado final: `current.md` en `review`, agente activo `implementer`.
`base_commit` y `last_review_commit` se conservan
(`c9a11cd0dfd5636e7aea4cdb3e75fa564731ee33`).
