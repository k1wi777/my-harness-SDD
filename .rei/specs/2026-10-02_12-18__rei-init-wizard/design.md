# Design — rei init wizard

> Work Item: `2026-10-02_12-18__rei-init-wizard` (feature).
>
> Este documento describe **cómo** se implementa la Feature. Las referencias
> `R<n>` apuntan a `requirements.md`.

---

# 1. Estrategia

La Feature tiene dos niveles con responsabilidades separadas y una única fuente
de verdad: los marcadores en los documentos.

1. **Nivel determinista (CLI `rei init`).** Nuevo comando que:
   - reutiliza la creación de estructura y la inicialización de git que ya vive
     en `check.Run` (extraída a una función compartida);
   - detecta marcadores `<!-- REI:PENDIENTE -->` en los documentos de
     personalización;
   - imprime el plan de pasos y el estado de cada documento;
   - añade `rei init status` como reporte de solo lectura.
   Todo esto es Go puro, sin IA, sin red.

2. **Nivel IA (rol `initializer`).** Nuevo archivo `.rei/agents/initializer.md`
   que describe una entrevista guiada por paso, la interpretación de respuestas y
   la escritura en el formato de cada documento. El rol elimina los marcadores al
   completar cada paso. No es código: es un protocolo de agente.

**Estado del wizard = marcadores en los archivos.** No hay archivo de estado
adicional. Un documento está pendiente si contiene el marcador. Esto hace el
wizard reanudable y consistente con el principio "el repositorio es la memoria".

---

# 2. Archivos afectados

## Nuevos

| Archivo | Propósito |
|---------|-----------|
| `internal/initwizard/initwizard.go` | Detección de pendientes, scaffold y plan de pasos (`Init`, `Status`, `Pending`). |
| `internal/initwizard/initwizard_test.go` | Tests del paquete. |
| `.rei/agents/initializer.md` | Rol de personalización asistida por IA. |

## Modificados

| Archivo | Cambio |
|---------|--------|
| `internal/check/check.go` | Extraer `EnsureStructure` (carpetas, `current.md`/`history.md`, git init) y añadir `.rei/agents/initializer.md` a `RequiredFiles`. |
| `internal/check/check_test.go` | Tests de `EnsureStructure` (si aplica) y de `RequiredFiles`. |
| `internal/cli/cli.go` | Dispatch de `init` y `init status`. |
| `internal/cli/help.go` | Entrada `init` en la tabla `commands`. |
| `internal/cli/help_test.go` | Tests de la ayuda de `init` y del dispatch. |
| `internal/doctor/doctor.go` | Reporte de documentos pendientes (WARN). |
| `internal/doctor/doctor_test.go` | Tests del reporte de pendientes. |
| `AGENTS.md` | Marcadores en §2 y §3; indicación al Leader de delegar en `initializer` cuando haya pendientes. |
| `.rei/docs/project/architecture.md` | Marcador `<!-- REI:PENDIENTE -->`. |
| `.rei/docs/project/conventions.md` | Marcador `<!-- REI:PENDIENTE -->`. |
| `.rei/docs/project/verification.md` | Marcador `<!-- REI:PENDIENTE -->`. |

No se modifican `.rei/docs/harness/*` (documentación del arnés) ni otros roles,
salvo lo indicado.

---

# 3. Componentes

## 3.1 `internal/check.EnsureStructure`

Función exportada que concentra los pasos 2–4 de `check.Run`:

```go
// EnsureStructure crea .rei/specs/, .rei/progress/work-items/,
// current.md e history.md si faltan, e inicializa git si no existe.
// Devuelve un código de salida (0 correcto; 1 si no pudo crear los archivos
// base). No ejecuta los checks de .rei/config.json.
func EnsureStructure(p *paths.Project, out io.Writer) int
```

- `check.Run` pasa a llamarla para sus pasos 2–4, preservando la salida actual.
- `initwizard.Init` la usa para el scaffold (R2–R7).

Impresión de git: el fallo al inicializar git es un **aviso**, no un error
(R7). `EnsureStructure` no distingue `--quiet`; `check.Run` conserva su lógica
de supresión si fuera necesario.

## 3.2 `internal/initwizard`

```go
// Doc describe un documento de personalización.
type Doc struct {
    ID       string // "agentes", "arquitectura", "convenciones", "verificacion"
    Path     string // ruta relativa a la raíz del proyecto
    Label    string // descripción legible
    Steps    []int  // pasos del plan asociados
}

// Docs es la lista canónica de documentos de personalización.
var Docs = []Doc{...}

// Pending devuelve los documentos que contienen el marcador.
func Pending(p *paths.Project) ([]Doc, error)

// Init ejecuta el scaffold, reporta pendientes y muestra el plan.
func Init(p *paths.Project, out io.Writer) int

// Status reporta pendientes sin modificar nada.
func Status(p *paths.Project, out io.Writer) int
```

- `marker = "<!-- REI:PENDIENTE -->"`. `Pending` lee cada `Doc.Path` y cuenta
  ocurrencias del marcador (`strings.Count > 0`). Un archivo inexistente se
  reporta como pendiente (no como error fatal) para que el listado sea completo
  (R19).
- `Init`: llama a `check.EnsureStructure`; luego `Pending`; imprime cabecera,
  resultado del scaffold, estado por documento y plan de pasos (R8–R10). Es
  idempotente porque nunca sobrescribe archivos existentes (R11).
- `Status`: solo `Pending` + impresión. Sin escrituras (R15).
- Ambos usan los mismos `Docs` y la misma tabla de pasos, evitando duplicar el
  plan.

**Plan de pasos (tabla canónica):**

| Paso | Título | Destino |
|------|--------|---------|
| 1 | Identidad y propósito | `AGENTS.md` §2 |
| 2 | Stack y comandos | `AGENTS.md` §3 |
| 3 | Arquitectura | `.rei/docs/project/architecture.md` |
| 4 | Convenciones | `.rei/docs/project/conventions.md` |
| 5 | Verificación | `.rei/docs/project/verification.md` + `.rei/config.json` |
| 6 | Cierre y validación | `rei init status` / `rei doctor` |

## 3.3 `internal/doctor`

Tras el bloque de integridad y antes del resultado final, añade una sección de
personalización: por cada documento pendiente, `[WARN] Personalización pendiente:
<ruta>`; si no hay, `[OK] Documentación del proyecto personalizada.` (R33–R35).
No incrementa `failures`.

## 3.4 CLI

- `internal/cli/cli.go`: `case "init": return cmdInit(rest)`, con
  `cmdInit` que despacha `status` o el modo por defecto y valida argumentos
  (R12, R17).
- `internal/cli/help.go`: nueva entrada `init` con `usage: "init [status]"`,
  `short` y `detail` (R36, R37).

## 3.5 Rol `initializer`

`.rei/agents/initializer.md` sigue la estructura del resto de roles
(frontmatter `name`/`description`/`tools`, secciones de reglas, protocolo y
comunicación) y especifica:

- **Precondición:** hay marcadores pendientes (`rei init status` → 1).
- **Protocolo por paso:** leer el documento destino; inferir del repositorio lo
  verificable (lenguajes, estructura, `Makefile`, manifiestos, comandos); hacer
  preguntas concretas; ofrecer sugerencias marcadas cuando el usuario no sepa;
  interpretar y escribir en el formato del documento; eliminar los marcadores de
  la sección; confirmar con el usuario; pasar al siguiente paso.
- **Reanudabilidad:** empezar por el primer documento pendiente; no rehacer los
  completos.
- **No inventar:** solo información confirmada o verificable.
- **Alcance de escritura:** los documentos de personalización + `.rei/config.json`.
- **Cierre:** `rei init status` y `rei doctor` sin pendientes; devolver control al
  Leader.
- **Prohibición:** no implementa código del proyecto.

## 3.6 Integración en `AGENTS.md`

Además de los marcadores de §2 y §3, se añade una instrucción en §1 (o §6) para
que el Leader, cuando `rei init status` reporte pendientes, delegue primero en el
rol `initializer` (R39). Así el wizard no requiere que el usuario conozca el
nombre del rol.

---

# 4. Decisiones de diseño

**D1 — Detección por marcador literal a nivel de documento.**
`R19` define pendiente ⇔ existe el marcador. `AGENTS.md` puede contener dos
marcadores (uno por sección); la detección es por archivo, por lo que el estado
de `AGENTS.md` se reporta a nivel de documento y el rol `initializer` gestiona
las dos secciones. Es simple, robusto y no depende de heurísticas.

**D2 — Sin archivo de estado del wizard.**
El estado vive en los propios documentos. Reanudar = volver a ejecutar
`rei init status`. Evita desincronización y respeta "el repositorio es la
memoria" (R23).

**D3 — `rei init` no ejecuta los checks del proyecto.**
`rei check` crea estructura *y* corre `.rei/config.json`. `rei init` solo hace el
scaffold y el reporte de personalización; ejecutar los checks en la
inicialización (cuando aún no hay comandos declarados) no aporta y acopla init a
la verificación. Por eso se extrae `EnsureStructure` y no se reutiliza
`check.Run` completo.

**D4 — Marcadores como WARN, no como FAIL.**
La personalización pendiente es el estado normal de un proyecto recién copiado.
`rei doctor` la advierte (R33–R35) y `rei check` no falla por ella (R40), para no
bloquear la inicialización ni el trabajo normal.

**D5 — Dos niveles explícitos.**
El CLI no intenta interpretar lenguaje natural; la IA no intenta crear
estructura. El CLI es la interfaz estable y verificable; el rol `initializer` es
el que interpreta y redacta. `rei init` imprime el plan y la invocación al rol.

**D6 — `initializer.md` entra en `RequiredFiles`.**
Es parte del arnés (como los demás roles), así que su ausencia debe detectarse
(R41). El implementador lo añade al arreglo existente.

**D7 — `.rei/config.json` se trata como parte del paso 5.**
No admite marcadores HTML (es JSON), así que su estado no se detecta por marcador:
la finalización del paso de verificación lo actualiza (R30). `verification.md` es
la señal de pendiente del paso.

---

# 5. Alternativas descartadas

**A1 — Archivo de estado `.rei/init.json` con el progreso del wizard.**
Descartada: duplica información, puede desincronizarse con los documentos y
contradice que el repositorio sea la memoria. Los marcadores ya codifican el
estado.

**A2 — Marcador con clave (`<!-- REI:PENDIENTE:arquitectura -->`).**
Descartada: la detección a nivel de documento es suficiente para `rei init` y
`rei init status`; una clave añade complejidad sin cambiar el comportamiento. Los
dos marcadores de `AGENTS.md` se resuelven a nivel de sección en el rol.

**A3 — Heurística sobre placeholders (`<...>`, `_..._`, `TODO`).**
Descartada: frágil y propensa a falsos positivos (los documentos ya usan `<...>`
en ejemplos). El marcador explícito es inequívoco.

**A4 — `rei init` llamando a `check.Run`.**
Descartada: ejecutaría los checks del proyecto durante la inicialización y
mezclaría dos responsabilidades. Se extrae `EnsureStructure` (D3).

**A5 — Wizard enteramente determinista (sin rol IA).**
Descartada: un CLI no puede interpretar respuestas abiertas ni redactar los
documentos en su formato; el nivel 2 es lo que resuelve "el usuario no sabe".

**A6 — Hacer que `rei check` falle con marcadores pendientes.**
Descartada: convertiría el estado de onboarding en un error que bloquea
`rei check` (que se ejecuta al inicio de cada sesión) e impediría arrancar el
propio wizard. Se reporta como advertencia (D4).
