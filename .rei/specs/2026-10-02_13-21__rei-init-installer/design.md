# Design — rei init installer

> Work Item: `2026-10-02_13-21__rei-init-installer` (feature).
>
> Este documento describe **cómo** se implementa la Feature. Las referencias `R<n>`
> apuntan a `requirements.md`.

---

# 1. Estrategia

El esqueleto del harness se **embebe en el binario** en tiempo de compilación y
`rei init` lo **despliega** en el proyecto destino. Así el binario es
autosuficiente: no necesita el repositorio fuente, ni red, ni un paso previo de
copia manual.

La Feature se apoya en tres piezas ya existentes:

1. **`go:embed`** para empaquetar el esqueleto. El archivo Go con la directiva vive
   en la **raíz del repositorio** (mismo directorio que `.rei/` y `AGENTS.md`),
   porque `go:embed` no admite `..` en los patrones. Expone un `embed.FS` como API
   pública para que lo consuma `internal/initwizard`.
2. **`check.EnsureStructure`** (ya extraída en la Feature anterior) para crear
   `.rei/specs/`, `.rei/progress/work-items/`, `current.md`, `history.md` e
   inicializar git. No se duplica el scaffold.
3. **`initwizard.Pending`/`Status`** para el reporte de personalización, que se
   mantiene sin cambios.

**Orden de ejecución de `rei init`:**

1. Resolver el proyecto destino (R5).
2. Desplegar el esqueleto embebido sin sobrescribir (R1–R10).
3. `check.EnsureStructure`: estructura de estado + git, sin ejecutar checks
   (R11–R17).
4. `initwizard.Pending`: reporte de pendientes y plan de pasos (R18–R19).

El despliegue precede a `EnsureStructure` porque este último lee
`.rei/templates/current.md` y `.rei/templates/history.md` para crear los archivos
de estado.

**Estado del wizard = marcadores.** No se introduce ningún archivo de estado nuevo;
se mantiene el comportamiento actual.

---

# 2. Archivos afectados

## Nuevos

| Archivo | Propósito |
|---------|-----------|
| `embed.go` (raíz, `package rei`) | Directiva `//go:embed` y variable `Skeleton embed.FS`. |
| `embed_test.go` (raíz, `package rei`) | Verifica el contenido embebido y sus exclusiones. |
| `internal/initwizard/installer.go` | `ResolveRootFrom`, `ResolveRoot` e `InstallSkeleton`. |
| `internal/initwizard/installer_test.go` | Tests del despliegue, idempotencia y resolución de destino. |

## Modificados

| Archivo | Cambio |
|---------|--------|
| `internal/initwizard/initwizard.go` | `Init` despliega el esqueleto antes de `EnsureStructure`; mantiene el reporte. |
| `internal/cli/cli.go` | `cmdInit` usa `initwizard.ResolveRoot()` en el modo instalador (no exige `.rei/` previo); `status` sigue con `project()`. |
| `internal/cli/help.go` | Actualiza el detalle de `init` (despliegue del esqueleto, estructura, git). |
| `internal/cli/help_test.go` | Ajusta/amplía las aserciones de la ayuda de `init`. |

No se modifican `.rei/docs/harness/*`, otros roles ni las plantillas.

---

# 3. Componentes

## 3.1 Paquete raíz `rei` (`embed.go`)

```go
// Package rei expone el esqueleto de REI Harness embebido en el binario.
package rei

import "embed"

// Skeleton contiene el esqueleto desplegable por `rei init`.
// Excluye deliberadamente .rei/specs/ y .rei/progress/ (estado del proyecto).
//
//go:embed all:.rei/docs all:.rei/agents all:.rei/templates all:.rei/config.json AGENTS.md
var Skeleton embed.FS
```

- El patrón con prefijo `all:` incluye subárboles completos bajo `.rei/` sin arrastrar
  `specs`/`progress` (R1, R2). `all:` es necesario para recorrer directorios cuyo
  nombre empieza por `.`.
- El paquete no importa nada: no puede crear ciclos. `internal/initwizard` lo importa
  con alias (`reiskel "github.com/k1wi777/my-harness-SDD"`).
- El módulo es `github.com/k1wi777/my-harness-SDD`; el paquete raíz puede llamarse
  `rei` (el nombre del paquete no tiene que coincidir con el último segmento de la
  ruta).

## 3.2 Resolución del destino (`internal/initwizard`)

```go
// ResolveRootFrom determina el proyecto destino partiendo de dir: el ancestro
// más cercano con .rei/ o, si no existe, dir.
func ResolveRootFrom(dir string) (*paths.Project, error)

// ResolveRoot aplica ResolveRootFrom sobre el directorio de trabajo actual.
func ResolveRoot() (*paths.Project, error)
```

- Reutiliza `paths.FindFrom(dir)`: si encuentra `.rei/`, devuelve ese proyecto; si
  devuelve `paths.ErrNotFound`, construye `&paths.Project{Root: dir}` (R5).
- La variante `From` permite testear la resolución sin cambiar el directorio de
  trabajo global.

## 3.3 Despliegue del esqueleto (`internal/initwizard`)

```go
// InstallSkeleton copia el esqueleto embebido en p.Root sin sobrescribir
// archivos existentes. Devuelve 0 si completa y 1 si no pudo crear algún
// archivo o directorio.
func InstallSkeleton(p *paths.Project, out io.Writer) int
```

Algoritmo con `fs.WalkDir(reiskel.Skeleton, ".", ...)`:

1. **Directorios:** `os.MkdirAll(destino, 0o755)`.
2. **Archivos:** si `destino` ya existe (`os.Stat` sin `IsNotExist`), se **omite**;
   si no, se asegura el directorio padre, se lee del `embed.FS` y se escribe con
   `0o644` (R3, R6, R7).
3. **`AGENTS.md` existente:** se omite y se emite un `[WARN]` específico, además del
   `[SKIP]` genérico (R8).
4. **Reporte:** `[OK]    <ruta>` por archivo creado, `[SKIP]  <ruta> (ya existe)` por
   omitido y una línea de resumen `N archivo(s) creado(s), M omitido(s).` (R9).

La omisión basada en existencia (en lugar de comparar contenido) es la que garantiza
la idempotencia y la no destrucción de contenido del usuario, incluido
`.rei/config.json` y `AGENTS.md` (R10).

## 3.4 Integración en `Init` (`internal/initwizard`)

`Init(p, out) int` mantiene su firma y pasa a:

1. Llamar `InstallSkeleton(p, out)` (paso nuevo).
2. Llamar `check.EnsureStructure(p, out)` (existente).
3. Detectar pendientes con `Pending(p)` y llamar a los `printStatus`/`printPlan`
   existentes (R18, R19).
4. Devolver `1` si `InstallSkeleton` o `EnsureStructure` fallan, o si `Pending`
   devuelve error; `0` en caso contrario (R22).

`Status(p, out) int` no se toca (R20, R21).

## 3.5 CLI y ayuda

- `internal/cli/cli.go`: en `cmdInit`, el modo `status` sigue usando `project()`
  (`paths.Find`). El modo instalador obtiene el destino con
  `initwizard.ResolveRoot()`; si falla la resolución devuelve `1`. La validación de
  argumentos (`init [status]`, resto ⇒ `2`) no cambia (R23).
- `internal/cli/help.go`: la entrada `init` describe el despliegue del esqueleto, la
  creación de estructura y git, y los códigos de salida (R26). `help_test.go` se
  ajusta a las nuevas cadenas.

## 3.6 Integridad del esqueleto

El conjunto embebido coincide con `check.RequiredFiles` más `.rei/config.json`
(leído por `config.Load`) y `.rei/docs/usage.md`. Tras una instalación limpia todos
los archivos requeridos existen, por lo que `rei check` no reporta faltantes (R24).
`.rei/config.json` embebido trae `"checks": []`, de modo que `rei check` no ejecuta
checks del proyecto al terminar de instalar.

---

# 4. Decisiones de diseño

**D1 — Paquete raíz con `go:embed` en lugar de generar archivos desde constantes Go.**
`go:embed` mantiene una única fuente de verdad (los archivos reales del repo) y evita
que el esqueleto se desincronice del código. Es la única ubicación válida, porque
`go:embed` no admite `..` y el esqueleto debe vivir junto a los archivos que embebe.

**D2 — `embed.FS` como API pública y despliegue en `initwizard`.**
El paquete raíz expone solo el `FS`; la lógica de caminar, crear sin sobrescribir y
reportar vive en `internal/initwizard`, junto al resto del comportamiento de
`rei init`. Se mantiene el paquete raíz sin dependencias.

**D3 — Exclusión por patrón selectivo, no por filtrado en tiempo de ejecución.**
`all:.rei/docs all:.rei/agents all:.rei/templates all:.rei/config.json AGENTS.md`
garantiza que `specs`/`progress` ni siquiera formen parte del binario (R2). Filtrado
posterior sería más código y permitiría fugas accidentales.

**D4 — No sobrescribir por existencia, no por comparación.**
Comparar contenido podría reescribir archivos que el usuario modificó
intencionalmente, y no aporta a la idempotencia. La existencia es la señal correcta
de "ya desplegado" (R7, R10).

**D5 — `AGENTS.md` como caso especial con aviso.**
Es el archivo que con más probabilidad tiene contenido propio; omitirlo en silencio
confundiría al usuario. Se avisa explícitamente (R8).

**D6 — Destino = ancestro con `.rei/`, o el directorio actual.**
Permite inicializar un proyecto nuevo (sin `.rei/`) y, a la vez, re-ejecutar
`rei init` desde un subdirectorio de un proyecto existente sin crear un `.rei/`
anidado (R5).

**D7 — El tamaño del binario es un coste aceptable.**
El esqueleto son documentos de texto de pocos cientos de KB; el aumento de tamaño es
irrelevante frente al objetivo de autosuficiencia. Se descarta cualquier esquema de
descarga en ejecución.

**D8 — Orden despliegue → `EnsureStructure`.**
`EnsureStructure` crea `current.md`/`history.md` a partir de
`.rei/templates/current.md` y `.rei/templates/history.md`. Las plantillas deben
existir antes; de ahí el orden.

**D9 — `rei init` no ejecuta los checks del proyecto.**
Se conserva la decisión de la Feature anterior: la inicialización no debe acoplarse
a la verificación. `EnsureStructure` no ejecuta `.rei/config.json` (R17).

**D10 — `rei init status` intacto.**
Es un contrato observable (detección de marcadores, códigos 0/1). No se modifica para
no romper al rol `initializer` (R20, R21).

---

# 5. Alternativas descartadas

**A1 — Generar los archivos del esqueleto desde literales en Go.**
Descartada: duplicaría cada documento dentro del código y se desincronizaría con los
archivos reales; además dificultaría su edición.

**A2 — Embeber en un paquete `internal/skeleton` con copias de los archivos.**
Descartada: mantendría dos copias del esqueleto y `go:embed` no puede apuntar con `..`
a `.rei/`, obligando a sincronización manual.

**A3 — Desplegar con `os.CopyFS`.**
Descartada: no ofrece el control por archivo (skip/reporte) ni la garantía explícita
de no sobrescribir que exige R7/R10; el walk manual es explícito y testeable.

**A4 — `git clone` o descarga del repositorio fuente en tiempo de ejecución.**
Descartada: reintroduce la dependencia de red y del repo fuente, justo lo que la
Feature elimina (R4).

**A5 — Sobrescribir los archivos existentes para "actualizar" el esqueleto.**
Descartada: destruiría personalizaciones del usuario (`AGENTS.md`,
`.rei/config.json`, `architecture.md`, …). La actualización del esqueleto queda fuera
de alcance (R27).

**A6 — Sembrar `.rei/specs/` o `.rei/progress/` desde el esqueleto.**
Descartada: son estado del proyecto destino, no parte del arnés. Se crean vacíos vía
`EnsureStructure` (R2, R11–R14).

**A7 — Requerir un `.rei/` existente para ejecutar `rei init`.**
Descartada: impide la inicialización de un proyecto nuevo desde cero. La resolución
de destino (D6) cubre ambos casos.
