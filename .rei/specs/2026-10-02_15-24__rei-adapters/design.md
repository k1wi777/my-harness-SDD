# Design — rei runtime adapters

> Work Item: `2026-10-02_15-24__rei-adapters` (feature).
>
> Este documento describe **cómo** se implementa la Feature. Las referencias `R<n>`
> apuntan a `requirements.md`.

---

# 1. Estrategia

La Feature separa dos conceptos que hoy están mezclados en `.rei/agents/<rol>.md`:

1. **Qué debe hacer el rol** — su contrato operativo. Vive en la sección
   `## Contrato` y es **autosuficiente**: contiene identidad, objetivo,
   precondiciones, protocolo completo, reglas duras, formato de salida,
   herramientas permitidas y punteros a la documentación del arnés. Es lo único
   que un agente necesita para ejecutar correctamente (R6–R8).
2. **Detalle de apoyo** — racional, ejemplos y casos límite. Vive en la sección
   `## Referencia` y es opcional: un agente que nunca la lea sigue ejecutando
   bien el rol (R9).

Sobre ese formato canónico se construye un **adaptador de runtime**. El adaptador
extrae el frontmatter y el Contrato, los mapea al formato nativo del runtime y
genera archivos de agente. En esta Feature el único runtime implementado es
**OpenCode** (R37), pero la infraestructura es genérica: añadir un runtime nuevo
consiste en crear un subdirectorio bajo `.rei/adapters/` (R19).

El mismo Contrato alimenta a ambos caminos (R11):

- **Runtime nativo (OpenCode):** `rei init opencode` genera
  `.opencode/agents/<rol>.md` con frontmatter mapeado + Contrato. OpenCode carga
  el agente y el Contrato es su system prompt.
- **Runtime no nativo (fallback):** el Leader pasa el Contrato del rol al
  subagente sin generar ningún archivo nativo (R14).

**Regla de oro de la refactorización:** si algo es necesario para ejecutar el rol
correctamente, va en `## Contrato`; si es solo justificación o ejemplo, va en
`## Referencia` (R10).

**Sin dependencias externas.** El módulo es `go 1.27` sin `require`. El parser de
frontmatter es el subconjunto mínimo de YAML descrito en §3; no se introduce
`gopkg.in/yaml.v3`.

---

# 2. Formato canónico de rol

Cada `.rei/agents/<rol>.md` queda así (R1–R5):

```md
---
name: implementer
description: <una frase; es el texto que OpenCode muestra al elegir agente>
mode: subagent            # primary | subagent
tools: [read, write, edit, search, shell]   # lista inline de herramientas genéricas
model: <provider/model>   # opcional
---

## Contrato

### Identidad
...

### Objetivo
...

### Precondiciones
...

### Protocolo
1. ...
2. ...

### Reglas duras
- ...

### Formato de salida
...
```

**Nota:** el bloque `## Referencia` se documenta aquí como sección canónica (R4),
pero este Work Item todavía no lo usa: los 5 roles se refactorizan con el
Contrato completo y la Referencia queda como sección opcional documentada
(ver Decisión D7).

## 2.1 Gramática de frontmatter aceptada

Para no añadir dependencias, el parser admite solo:

- delimitador de apertura y cierre `---` en la primera línea y en una posterior;
- líneas `clave: valor` con valores escalares sin comillas;
- `tools` como lista inline entre corchetes: `tools: [read, write, edit]`
  (separador `,`, espacios opcionales);
- `model` opcional.

Un frontmatter que no cumpla esta gramática DEBE producir un error de parseo
explícito, no una interpretación silenciosa.

## 2.2 Asignación de roles

| Rol | `mode` | `tools` genéricas |
|-----|--------|-------------------|
| `leader` | `primary` | `read, search, shell, write, edit, subagent` |
| `spec_author` | `subagent` | `read, write, edit, search, shell` |
| `implementer` | `subagent` | `read, write, edit, search, shell` |
| `reviewer` | `subagent` | `read, write, edit, search, shell` |
| `initializer` | `subagent` | `read, write, edit, search, shell` |

`leader` necesita `subagent` para lanzar a los demás roles (R2). En el fallback,
`subagent` se materializa como la capacidad de spawn del runtime; en OpenCode es
la clave de permiso `task`.

---

# 3. Extracción del Contrato

Se implementa en un paquete nuevo `internal/adapter`, sin dependencias externas.

```go
// Role es un rol canónico parseado desde .rei/agents/<rol>.md.
type Role struct {
    Name        string
    Description string
    Mode        string   // "primary" | "subagent"
    Tools       []string // genéricas: read, write, edit, search, shell, subagent
    Model       string   // opcional
    Contract    string   // cuerpo de ## Contrato, sin la sección ## Referencia
}

// ParseRole parsea un archivo de rol. Error si falta frontmatter, mode inválido,
// herramienta desconocida, o faltan las secciones ## Contrato / ## Referencia.
func ParseRole(data []byte) (Role, error)

// ExtractContract devuelve el texto entre el encabezado ## Contrato y el
// encabezado ## Referencia (excluyendo ambos encabezados).
func ExtractContract(body string) (string, error)
```

Reglas de extracción:

- El frontmatter es el primer bloque entre la primera línea `---` y la siguiente `---`.
- `## Contrato` y `## Referencia` son encabezados de nivel 2; se localizan por
  coincidencia exacta de la línea (con espacios de recorte).
- El Contrato es todo el texto entre el final de la línea `## Contrato` y el
  inicio de la línea `## Referencia`, recortado.
- Si el orden es incorrecto o falta alguna sección, `ParseRole` devuelve error (R4).
- La Referencia se excluye deliberadamente (R20).

---

# 4. Formato del adaptador OpenCode

## 4.1 Destino y versión objetivo

Se genera en **`.opencode/agents/<rol>.md`** (R22). Es la ubicación preferida por
OpenCode V2, que además sigue descubriendo el `agent/` singular de V1. La versión
instalada y verificada en el entorno de desarrollo es **v2.0.21**.

> La delegación mencionaba `.opencode/agent/<rol>.md` (singular) como conjetura y
> pedía basar el adaptador en el formato nativo real y documentarlo. La decisión
> aquí es usar la ruta preferida de V2 (plural) y el frontmatter nativo V2; ver
> D1 y D2.

## 4.2 Frontmatter nativo

OpenCode V2 define los agentes Markdown con `description`, `mode`, `model`
(opcional) y un mapa `permission` de entradas `clave: "allow" | "ask" | "deny"`;
el cuerpo del archivo es el *system prompt*. El campo `tools` es formato deprecado,
así que el adaptador usa `permission`. El mapeo desde el rol canónico es:

| Canónico | OpenCode V2 |
|----------|-------------|
| `description` | `description` (idéntico) |
| `mode` | `mode` (idéntico: `primary`/`subagent`) |
| `model` | `model` (idéntico, solo si está presente) |
| `tools` genéricas | `permission` (una entrada `clave: allow` por permiso) |
| `## Contrato` | cuerpo del archivo (precedido por la marca GENERATED) |

Mapa de herramientas genéricas → claves de permiso de OpenCode (R18):

| Genérica | Clave(s) de permiso OpenCode |
|----------|------------------------------|
| `read` | `read` |
| `write` | `edit` |
| `edit` | `edit` |
| `search` | `glob`, `grep` |
| `shell` | `bash` |
| `subagent` | `task` |

`write` y `edit` colapsan en la misma clave `edit` de OpenCode; si un rol declara
ambas, se emite una sola entrada `edit: allow` (deduplicación). El orden de las
entradas sigue el orden de declaración de `tools` en el rol.

No se emiten entradas `deny` para las herramientas ausentes: OpenCode aplica sus
propios permisos por defecto y el Contrato ya restringe el uso. Esto mantiene el
mapeo simple y evita conceder o revocar más de lo declarado.

## 4.3 Ejemplo de salida generada

```md
---
description: Implementa un único Work Item siguiendo la planificación aprobada.
mode: subagent
permission:
  read: allow
  edit: allow
  glob: allow
  grep: allow
  bash: allow
---

<!-- GENERATED by rei — no editar; regenerar con `rei init opencode` -->

## Contrato
...
```

La marca GENERATED es una constante del paquete adapter (R24); su presencia en el
archivo es lo que autoriza a sobrescribirlo (R25, R26).

---

# 5. Layout de `.rei/adapters/`

```text
.rei/adapters/
└── opencode/
    ├── agent.tmpl   # plantilla del agente nativo (text/template)
    └── tools.json   # mapa herramienta genérica -> claves de permiso OpenCode
```

- **`agent.tmpl`** (R17): plantilla con los campos `Description`, `Mode`, `Model`,
  `Permission` y `Contract`. El marcador GENERATED lo aporta el paquete adapter
  (no la plantilla) para que sea una constante verificable y no se pueda olvidar.
  Los datos de la plantilla se resuelven desde `.rei/adapters/opencode/` en tiempo
  de ejecución (leídos del proyecto destino).
- **`tools.json`** (R17, R18): mapa declarativo para que el adaptador no incruste
  el mapeo en Go. Forma propuesta:

```json
{
  "read":     { "permission": ["read"] },
  "write":    { "permission": ["edit"] },
  "edit":     { "permission": ["edit"] },
  "search":   { "permission": ["glob", "grep"] },
  "shell":    { "permission": ["bash"] },
  "subagent": { "permission": ["task"] }
}
```

El paquete adapter valida que el mapa cubra todas las genéricas de R3, que cada
valor sea una clave de permiso válida de OpenCode (`read`, `edit`, `glob`, `grep`,
`bash`, `task`) y falla con error explícito si falta alguna o es desconocida.

> Decisión abierta a revisión (D6): leer la plantilla y el mapa desde el
> directorio del proyecto destino (`.rei/adapters/opencode/`) permite
> personalizarlos y es coherente con el resto del arnés, que vive en `.rei/`. La
> alternativa es embeberlos y usarlos como única fuente; se descarta porque
> impediría ajustar el adaptador por proyecto.

---

# 6. Generación y `--check`

## 6.1 API del paquete `internal/adapter`

```go
// GenerateAll genera .opencode/agents/<rol>.md para cada rol de p.Root.
// No sobrescribe archivos que no lleven la marca GENERATED. Devuelve 0 si
// completa (aunque omita algún archivo) y 1 si no pudo escribir un archivo
// generado.
func GenerateAll(p *paths.Project, out io.Writer) int

// Check comprueba, sin escribir, que cada archivo generado existe y coincide
// byte a byte con la generación canónica. Devuelve 0 si todo coincide y 1 si
// falta o difiere algún archivo.
func Check(p *paths.Project, out io.Writer) int

// Install instala el esqueleto, asegura la estructura de estado y genera los
// agentes nativos. Devuelve 1 si alguna fase falla.
func Install(p *paths.Project, out io.Writer) int
```

`GenerateAll` recorre los roles canónicos `leader`, `spec_author`, `implementer`,
`reviewer`, `initializer` (lista ordenada y estable; misma que
`check.RequiredFiles`) y, por cada uno:

1. Lee `.rei/agents/<rol>.md` y lo parsea con `ParseRole`.
2. Construye el contenido con la plantilla + mapa + marca GENERATED.
3. Si el destino no existe → escribe `[OK]`.
4. Si el destino existe y contiene la marca → reescribe solo si el contenido
   difiere; reporta `[OK]`/`[UPD]`.
5. Si el destino existe y **no** contiene la marca → `[WARN]` y no lo toca
   (R25).

`Check` reutiliza exactamente la misma función de construcción que
`GenerateAll`, de modo que "generación canónica" y "verificación" no pueden
divergir (R27, R28). Si un archivo existe sin marca, `Check` lo reporta como
deriva.

## 6.2 Integración con `rei init`

`cmdInit` (`internal/cli/cli.go`) pasa a aceptar:

```text
rei init [status]
rei init opencode [--check]
```

- Sin argumentos: comportamiento actual (`initwizard.Init`).
- `status`: comportamiento actual (`initwizard.Status`).
- `opencode`: `adapter.Install(p, out)` (R21).
- `opencode --check`: `adapter.Check(p, out)` (R27).
- Cualquier otro argumento o combinación no reconocida: uso por `stderr` + `2`
  (R30).

`adapter.Install` reutiliza piezas existentes (`initwizard.InstallSkeleton`,
`check.EnsureStructure`) para no duplicar la instalación (R21, R29). No modifica
`rei init` ni `rei init status`. La instalación del esqueleto precede a la
generación porque la plantilla y el mapa deben existir en el destino.

> Orden: `InstallSkeleton` → `EnsureStructure` → `GenerateAll`. El mismo orden que
> `initwizard.Init`, con la generación al final.

## 6.3 Empaquetado

`embed.go` añade `all:.rei/adapters` a la directiva `//go:embed` (R31):

```go
//go:embed all:.rei/docs all:.rei/agents all:.rei/templates all:.rei/adapters all:.rei/config.json AGENTS.md
var Skeleton embed.FS
```

No se añaden los archivos de `.rei/adapters/` a `check.RequiredFiles`: `rei check`
no debe fallar en proyectos ya inicializados con una versión anterior del binario
hasta que se ejecute `rei init opencode`. La completitud del adaptador es
responsabilidad de `--check`.

---

# 7. Cambios en el fallback (roles y `AGENTS.md`)

- **`leader.md`**: su Contrato incorpora la regla de fallback (R14, R15): cuando
  el runtime no tenga adaptador nativo, transmitir al subagente únicamente la
  sección `## Contrato` del rol, nunca el archivo completo. La plantilla de
  delegación existente se mantiene; se le añade que la entrada es el Contrato.
- **`AGENTS.md`**: la sección de carga de contexto y el mapa de documentación
  mencionan que `.rei/agents/<rol>.md` tiene frontmatter + `## Contrato` +
  `## Referencia`, que el Contrato es lo que se transmite a los subagentes y que
  existe la infraestructura de adaptadores en `.rei/adapters/`.
- **`check.RequiredFiles`**: sin cambios (siguen siendo los mismos 5 roles y los
  mismos documentos). El formato interno de los roles cambia, pero sus rutas no.

---

# 8. Verificación

## 8.1 Cobertura del Contrato (R32)

Por cada rol, una checklist que enumere:

- cada regla dura del rol actual, y
- cada paso del protocolo actual,

indicando si está en `## Contrato` o si es detalle opcional en `## Referencia`.
Se documenta en `.rei/progress/work-items/2026-10-02_15-24__rei-adapters/impl.md`.
El Reviewer la usa como evidencia de que no se perdió ningún paso.

Además, una prueba automatizada recorre los 5 roles y exige que su `## Contrato`
contenga las subsecciones mínimas (identidad, objetivo, precondiciones,
protocolo, reglas, salida, herramientas y punteros a documentos) (R33).

## 8.2 Tests (R33)

`internal/adapter/*_test.go`, con `t.TempDir()`:

- `ParseRole`: rol válido; falta de frontmatter; `mode` inválido; herramienta
  desconocida; falta de `## Contrato`; falta de `## Referencia`; orden invertido.
- `ExtractContract`: excluye la Referencia; recorta límites.
- Mapeo: cada genérica produce las claves de permiso esperadas; `write`+`edit`
  deduplican; falta de cobertura del mapa ⇒ error.
- `GenerateAll`: crea un archivo por rol con frontmatter y Contrato; segunda
  ejecución idéntica (idempotencia); archivo preexistente sin marca no se
  sobrescribe y avisa.
- `Check`: 0 con todo generado; 1 si falta un archivo; 1 si el contenido difiere.
- `embed_test.go`: el `embed.FS` contiene `.rei/adapters/opencode/agent.tmpl` y
  `.rei/adapters/opencode/tools.json`.
- CLI: `rei init opencode` y `--check` despachan a `adapter`; argumentos inválidos
  ⇒ `2`; `rei help init` documenta `opencode` y `--check` (R35).

## 8.3 Medición de tokens (R34)

El Implementer registra en `impl.md`, por rol, el tamaño del `## Contrato` (lo que
se transmite por spawn) frente al archivo completo, en caracteres y en una
estimación de tokens (caracteres/4). No es un gate de CI: es evidencia del ahorro.

---

# 9. Archivos afectados

## Nuevos

| Archivo | Propósito |
|---------|-----------|
| `internal/adapter/adapter.go` | `Role`, `ParseRole`, `ExtractContract`, construcción de contenido. |
| `internal/adapter/opencode.go` | Mapeo a OpenCode V2 y lectura de `agent.tmpl`/`tools.json`. |
| `internal/adapter/generate.go` | `GenerateAll`, `Check`, `Install`. |
| `internal/adapter/adapter_test.go` | Tests de parseo, extracción, mapeo y generación. |
| `internal/adapter/generate_test.go` | Tests de generación, idempotencia, no sobrescritura y `--check`. |
| `.rei/adapters/opencode/agent.tmpl` | Plantilla del agente nativo. |
| `.rei/adapters/opencode/tools.json` | Mapa de herramientas genéricas → OpenCode. |

## Modificados

| Archivo | Cambio |
|---------|--------|
| `.rei/agents/leader.md` | Formato canónico; regla de fallback (Contrato, no archivo completo). |
| `.rei/agents/spec_author.md` | Formato canónico. |
| `.rei/agents/implementer.md` | Formato canónico. |
| `.rei/agents/reviewer.md` | Formato canónico. |
| `.rei/agents/initializer.md` | Formato canónico. |
| `AGENTS.md` | Documenta formato de rol, fallback y adaptadores. |
| `embed.go` | Incluye `all:.rei/adapters`. |
| `embed_test.go` | Verifica el nuevo subárbol embebido. |
| `internal/cli/cli.go` | `cmdInit` acepta `opencode [--check]`. |
| `internal/cli/help.go` | Documenta `opencode` y `--check`. |
| `internal/cli/help_test.go` | Ajusta las aserciones de la ayuda de `init`. |

No se modifican `.rei/docs/harness/*` salvo lo necesario para documentar el
formato (§10) ni las plantillas de `.rei/templates/`.

---

# 10. Documentación del formato (R36)

La descripción canónica del formato de rol y del procedimiento de adaptación se
añade a `.rei/docs/harness/`. Opción preferida: una sección nueva "Roles y
adaptadores" dentro de `workflow.md` o un documento nuevo
`.rei/docs/harness/adapters.md` enlazado desde `AGENTS.md`. Describe:

- frontmatter canónico y valores válidos;
- contrato vs referencia y la regla de oro;
- extracción y generación;
- cómo añadir un runtime (nuevo subdirectorio en `.rei/adapters/`);
- el comando `rei init opencode [--check]`.

La decisión entre sección o documento nuevo se deja al Implementer, siempre que
quede enlazado desde `AGENTS.md` (R36). `RequiredFiles` no cambia, así que un
documento nuevo no es obligatorio para `rei check`.

---

# 11. Decisiones de diseño

**D1 — Target OpenCode V2 y ruta `.opencode/agents/` (plural).** La versión
instalada verificada es v2.0.21 y la documentación de V2 es la fuente de verdad.
V2 descubre además el `agent/` singular de V1, pero la ruta preferida y no
deprecada es `agents/`. Se elige la preferida.

**D2 — Frontmatter nativo V2 (`permission`), no el V1 `tools`.** La documentación
de V2 pide no usar campos legacy (`tools`) en configuración nueva; el campo vigente
es `permission`, un mapa `clave: "allow" | "ask" | "deny"`. El frontmatter canónico
sigue siendo genérico (`tools`), y el adaptador traduce a la forma nativa V2. Esto
satisface "el nativo impone tools/model" (R12) sin exponer OpenCode en el formato
canónico (R38).

**D3 — Parser de frontmatter propio y mínimo.** Evita dependencias externas
(go.mod sin `require`) y hace explícita la gramática soportada. `tools` en lista
inline es suficiente para los 5 roles.

**D4 — El Contrato es la única fuente del cuerpo del agente nativo.** Se extrae
por encabezados `## Contrato`/`## Referencia`, que quedan fijados como contrato
de extracción. Así nativos y no nativos reciben texto idéntico (R11, R13, R20).

**D5 — Marca GENERATED como autorización de sobrescritura.** Comparar solo el
contenido no distingue "archivo generado y desactualizado" de "archivo editado a
mano". La marca permite regenerar de forma segura y no destruir trabajo manual
(R24–R26).

**D6 — Plantilla y mapa leídos del proyecto destino (`.rei/adapters/`).** Permite
personalizar el adaptador por proyecto y mantiene `.rei/` como fuente única. La
alternativa (embeberlos y usarlos directamente) se descarta porque impediría
personalizarlos.

**D7 — Referencia documentada pero aún vacía.** R4/R9 fijan el contrato de las
dos secciones. En esta Feature la refactorización mueve el protocolo y las reglas
al Contrato y no conserva prosa de referencia en los 5 roles; la Referencia queda
como sección obligatoria (posiblemente con una nota de que no hay detalle
adicional). Así el Contrato es autosuficiente desde el primer momento.

**D8 — No añadir los adaptadores a `check.RequiredFiles`.** Añadirlos rompería
`rei check` en proyectos inicializados con binarios previos. La verificación del
adaptador se hace con `rei init opencode --check`.

**D9 — `rei init opencode` reutiliza la instalación existente.** No duplica el
despliegue del esqueleto ni la estructura; llama a `InstallSkeleton` y
`EnsureStructure`. Mantiene `rei init` y `rei init status` intactos (R29).

**D10 — El fallback se resuelve por instrucción en `leader.md`, no por CLI.**
Añadir un comando para imprimir el Contrato sería una superficie nueva no pedida;
en un runtime no nativo el Leader ya tiene o puede leer el rol y transmite su
Contrato. Queda como alternativa A3 para el futuro.

---

# 12. Alternativas descartadas

**A1 — Mantener el formato actual de roles y generar el agente nativo con el
archivo completo (Contrato + Referencia).** Descartada: la Referencia es detalle
opcional y transmitirla derrocha tokens en cada spawn; el objetivo es enviar solo
el Contrato (R14).

**A2 — Emitir el frontmatter legacy (`tools` como mapa de bool) por
compatibilidad amplia.** Descartada: V2 la traduce automáticamente, pero es
formato deprecado y la documentación de V2 pide no usarlo en archivos nuevos; se
prefiere el mapa nativo actual `permission` (D2).

**A3 — Comando `rei agent contract <rol>` para el fallback.** Descartada en esta
Feature por no estar en el alcance acordado; el Leader transmite el Contrato
directamente. Queda como extensión natural.

**A4 — Guardar el Contrato extraído en un archivo intermedio.** Descartada:
introduciría una segunda fuente de verdad que podría desincronizarse del rol.

**A5 — Un único adaptador monolítico sin `.rei/adapters/`.** Descartada: incrustar
el mapeo en Go dificulta añadir runtimes y personalizarlos; el directorio por
runtime es el punto de extensión (R19).

**A6 — Detectar deriva comparando solo un hash.** Descartada: no permitiría
distinguir un archivo no generado de uno generado y modificado; la marca
GENERATED es más informativa y ya es necesaria para la no sobrescritura.
