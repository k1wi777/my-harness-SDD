# Plan — `rei init --update [--force]`

## Objetivo

Añadir el subcomando `rei init --update [--force]`, que **refresca el esqueleto
embebido del harness** en un proyecto **ya inicializado** sin pisar la
personalización del usuario ni su estado.

- `rei init --update` **no instala desde cero**: exige un proyecto REI (`.rei/`).
  Si no lo está, informa con el mensaje canónico de `rei init` y sale con 1.
- `--force` sobrescribe los archivos que difieran, incluidos los personalizados.
- **Nunca** se tocan `.rei/specs/**` ni `.rei/progress/**`.
- Se reporta cada archivo como **nuevo**, **actualizado**, **sin cambios** u
  **omitido** (conservado), y se regeneran los nativos si procede.

## Política de actualización (decisión)

Sin marcadores no se puede distinguir un archivo de harness «sin personalizar»
de uno «modificado a mano». Por eso la actualización combina tres mecanismos:

1. **Markers** para los documentos de personalización de proyecto.
2. **Manifiesto de hashes** para el resto del esqueleto.
3. **`--force`** como escape explícito.

### 1. Documentos de personalización (marker-aware)

Son `AGENTS.md` (§2/§3) y `.rei/docs/project/{architecture,conventions,verification}.md`
(la lista canónica `Docs` de `internal/initwizard`):

- Contiene `<!-- REI:PENDIENTE -->` → **no personalizado** → se actualiza desde
  el embebido (aunque el archivo difiera).
- No contiene el marcador → **personalizado** → se **conserva** con `[SKIP]` y
  aviso; `--force` lo sobrescribe.
- No existe → se crea desde el embebido.

> Límite conocido de `AGENTS.md`: la regla es a nivel de archivo, así que si el
> usuario editó secciones ajenas a §2/§3 **sin** haber personalizado §2/§3
> (marcador presente), `--update` reemplaza el archivo completo. Es el
> comportamiento que pide la política marker-aware; se documenta en la ayuda.

### 2. Manifiesto de hashes `.rei/install-manifest.json` (archivos de harness)

Alcance: `.rei/docs/harness/**`, `.rei/agents/**`, `.rei/adapters/**`,
`.rei/templates/**`, `.rei/config.json` y las partes de `AGENTS.md` fuera de
§2/§3 (a efectos prácticos, el archivo se gobierna por la regla del punto 1).

- Formato: `{"files": {"<ruta relativa>": "<sha256 hexadecimal>"}}`.
- Se escribe/refresca al desplegar (`InstallSkeleton`), registrando **solo** los
  archivos cuyo contenido en disco coincide **byte a byte** con el embebido
  (es decir, los que no han sido modificados). Se guarda ordenado y estable.
- En `--update`, para cada archivo de harness que difiere del embebido:
  - Tiene entrada en el manifiesto **y** `sha256(disco) == entrada` → no fue
    modificado desde la instalación → **se actualiza**.
  - No tiene entrada, o `sha256(disco) != entrada` → **modificado a mano** (o
    sin baseline) → **se conserva** con `[SKIP]` y aviso; `--force` lo sobrescribe.
- Tras la operación se reescribe el manifiesto con los archivos que quedan
  idénticos al embebido. Los archivos conservados no quedan registrados (para no
  marcar como baseline una modificación del usuario).

**Justificación (alternativas descartadas):**

- *Actualizar por defecto*: destruiría silenciosamente `.rei/config.json` (checks
  del usuario), `.rei/agents/*.md` adaptados o `.rei/adapters/**`. Descartada.
- *Solo marker-aware*: los archivos de harness no llevan marcador, así que nunca
  podrían actualizarse; el comando sería inútil. Descartada.
- *Sufijo/backup `.bak`*: ensucia el árbol, no es idempotente y pierde el
  baseline. Descartada.
- **Manifiesto de hashes**: mínimo, sin dependencias externas, idempotente y
  permite distinguir «sin tocar» de «modificado». Es la opción adoptada.

**Límites del manifiesto:**

- Proyectos inicializados **antes** de esta funcionalidad no tienen manifiesto:
  la primera `--update` conserva con aviso cualquier archivo de harness que
  difiera (no puede saber si es una versión antigua o una edición manual). Tras
  esa ejecución el manifiesto queda escrito para los archivos que sí coinciden,
  de modo que las siguientes ejecuciones ya actualizan con normalidad. Para
  forzar el primer refresco: `--force` o restaurar los archivos.
- El manifiesto no se embebe (no está bajo los patrones de `go:embed`) y vive en
  `.rei/install-manifest.json`, fuera de `specs/` y `progress/`.

### 3. `--force` y `--update`

- Sobrescribe cualquier archivo que difiera, incluidos los personalizados, y lo
  registra en el manifiesto con el hash embebido.
- Los archivos nativos regenerados (`rei init opencode` / `claude`) siguen su
  propia regla de marca `GENERATED`; `--force` no elimina archivos nativos
  escritos a mano (se conservan con aviso).

### 4. Regeneración de nativos («si procede»)

Tras refrescar el esqueleto, si el proyecto tiene `.opencode/` se regeneran los
nativos de OpenCode y si tiene `.claude/` los de Claude Code. Si ninguno existe,
no se instala nada (eso sería `rei init opencode`/`claude`, fuera del alcance de
`--update`). Se reutiliza `adapter.Install`/`adapter.InstallClaude` (principio
«reutilizar antes que duplicar»); su `EnsureStructure` es idempotente y no
modifica estado existente en un proyecto ya inicializado.

## Tabla de decisión (por archivo del esqueleto)

| Situación | Sin `--force` | Con `--force` |
|-----------|---------------|---------------|
| No existe | crear `[NEW]` | crear `[NEW]` |
| Igual al embebido | sin cambios `[OK]` | sin cambios `[OK]` |
| Difiere · personalización con marcador | actualizar `[UPD]` | actualizar `[UPD]` |
| Difiere · personalización sin marcador | conservar `[SKIP]` | actualizar `[UPD]` |
| Difiere · harness con hash de manifiesto == disco | actualizar `[UPD]` | actualizar `[UPD]` |
| Difiere · harness modificado o sin entrada | conservar `[SKIP]` | actualizar `[UPD]` |

Resumen final: `N nuevo(s), M actualizado(s), K sin cambios, J omitido(s).`
Códigos de salida: `0` correcto (aunque haya omitidos), `1` si falla una
lectura/escritura del esqueleto o un fallo no recuperable, `2` uso incorrecto.

## Archivos

| Archivo | Cambio |
|---------|--------|
| `internal/initwizard/manifest.go` (nuevo) | `ManifestPath = ".rei/install-manifest.json"`, `loadManifest`, `saveManifest`, `hashBytes` (SHA-256, `crypto/sha256`, `encoding/hex`, `encoding/json`). |
| `internal/initwizard/installer.go` | `InstallSkeleton` refresca el manifiesto al terminar (solo archivos idénticos al embebido). |
| `internal/initwizard/update.go` (nuevo) | `Update(p *paths.Project, force bool, out io.Writer) int` con la política completa y el reporte. |
| `internal/initwizard/update_test.go` (nuevo) | Tests de la política, idempotencia y no-interferencia con `specs/`/`progress/`. |
| `internal/initwizard/manifest_test.go` (nuevo) | Tests de carga/guardado/hashing del manifiesto. |
| `internal/cli/cli.go` | `cmdInit` reconoce `--update [--force]`; exige `project()`; llama a `Update` y regenera nativos si `.opencode`/`.claude` existen. |
| `internal/cli/help.go` | Uso y detalle de `init` documentan `--update [--force]`, la política marker-aware y el manifiesto. |
| `internal/cli/help_test.go` | Ajustar la aserción del uso de `init` y añadir test de que la ayuda menciona `--update`, `--force` y el manifiesto. |
| `internal/cli/cli_test.go` | Tests de uso correcto/incorrecto de `cmdInit --update` (2 con args inválidos). |

> `internal/initwizard/initwizard.go` no cambia: `Init` sigue igual y obtiene el
> manifiesto de forma transparente vía `InstallSkeleton`.

## Cambios

1. **Manifiesto** (`manifest.go`):
   - `const ManifestPath = ".rei/install-manifest.json"`.
   - `type manifest struct { Files map[string]string \`json:"files"\` }`.
   - `hashBytes(data []byte) string` con `sha256` + `hex.EncodeToString`.
   - `loadManifest(p) map[string]string`: devuelve `nil` si no existe o es
     inválido (no fatal); error solo ante fallos de lectura distintos de
     «no existe».
   - `saveManifest(p, entries) error`: `MkdirAll`, serialización con claves
     ordenadas e indentación, escritura `0644`.
2. **`InstallSkeleton`**: tras el resumen, recorrer de nuevo los archivos del
   `Skeleton` (o acumular en el primer walk), calcular `hashBytes(embebido)` y
   `hashBytes(disco)`; registrar los que coincidan y guardar el manifiesto. Si el
   guardado falla, emitir `[WARN]` y **no** cambiar el código de salida (la
   instalación sigue siendo válida).
3. **`Update`**: recorrer `reiskel.Skeleton` con `fs.WalkDir`, ignorando
   directorios salvo para crearlos. Por cada archivo:
   - Leer embebido y disco.
   - Clasificar con `Docs` (personalización) o harness.
   - Aplicar la tabla de decisión.
   - Escribir con `MkdirAll` + `os.WriteFile` cuando corresponda (`0644`).
   - Acumular contadores y entradas de manifiesto para archivos que quedan
     idénticos al embebido.
   - `[FAIL]` y `exit = 1` ante errores de E/S.
   - Al final, `saveManifest` y resumen. **No** llamar a `check.EnsureStructure`.
4. **`cmdInit`**:
   - Detectar `--update` (independiente del orden respecto a `--force`) y
     rechazar cualquier otro argumento (`2`).
   - Resolver con `project()` (exige `.rei/`; mensaje de no inicializado).
   - `rc := initwizard.Update(p, force, os.Stdout)`.
   - Si `os.Stat(p.Root+"/.opencode")` es dir → `rc |= adapter.Install(p, os.Stdout)`;
     si `.claude` → `rc |= adapter.InstallClaude(p, os.Stdout)`.
   - Devolver `rc`.
5. **`help.go`**: actualizar `usage` de `init` a
   `init [status|opencode|claude [--check]|--update [--force]]` y añadir líneas de
   detalle (qué actualiza, marker-aware, manifiesto, `--force`, no toca
   `specs/`/`progress/`, no instala desde cero).
6. **Tests**: unitarios de `Update`/manifiesto con `t.TempDir()` y `reiskel.Skeleton`;
   ajuste de `help_test.go` y `cli_test.go`.

## Restricciones

- No añadir dependencias externas (solo `crypto/sha256`, `encoding/hex`,
  `encoding/json`, `io/fs`, `os`, `path/filepath`).
- No modificar el esqueleto embebido (`embed.go` / `.rei/**` desplegable) ni el
  alcance de `meta.json`.
- No tocar `.rei/specs/**` ni `.rei/progress/**` (el walk solo recorre
  `reiskel.Skeleton`, que ya los excluye; `--update` no llama a `EnsureStructure`).
- Mantener la idempotencia: una segunda ejecución no debe escribir nada.
- Respetar el estilo existente (comentarios en español, prefijos `[OK]`/`[UPD]`/
  `[WARN]`/`[FAIL]`/`[SKIP]`).
- `rei init` sin argumentos, `status`, `opencode` y `claude` deben seguir
  comportándose igual (regresión cubierta por los tests existentes).

## Verificación

- `V1` — `go test ./...` (incluye los tests nuevos de manifiesto, `Update` y CLI).
- `V2` — `go vet ./...`.
- `V3` — `make build` (compila `bin/rei`).
- `V4` — Prueba de aceptación en **proyecto temporal** con el binario:
  1. `rei init` en un dir vacío; personalizar un doc (quitar marcador) y
     modificar `.rei/agents/leader.md` y `.rei/config.json`.
  2. `rei init --update`: el doc personalizado, `leader.md` y `config.json` se
     conservan (`[SKIP]`); los docs con marcador se actualizan; un archivo de
     harness intacto cambia si se alteró el embebido.
  3. `rei init --update --force`: se sobrescriben los personalizados/modificados.
  4. Verificar que `.rei/specs/**` y `.rei/progress/**` (incluido `current.md` y
     `history.md`) no cambian (comparar hashes antes/después).
  5. Ejecutar `rei init --update` dos veces seguidas: la segunda todo «sin
     cambios» (idempotencia).
  6. `rei init --update` fuera de un proyecto REI: mensaje de no inicializado y
     exit 1.
- `V5` — `rei check --quiet` en el repo (exit 0) y `rei init opencode --check` /
  `rei init claude --check` (exit 0) para confirmar que no hay deriva tras
  regenerar nativos.

## Pasos

- [x] 1. Crear `internal/initwizard/manifest.go` (`ManifestPath`, `hashBytes`, `loadManifest`, `saveManifest`).
- [x] 2. Modificar `InstallSkeleton` (`internal/initwizard/installer.go`) para refrescar el manifiesto tras el despliegue (solo archivos idénticos al embebido; fallo de guardado = `[WARN]`).
- [x] 3. Crear `internal/initwizard/update.go` con `Update(p, force, out) int` implementando la tabla de decisión, el reporte y la reescritura del manifiesto.
- [x] 4. Integrar `--update [--force]` en `cmdInit` (`internal/cli/cli.go`): validación de args, `project()`, llamada a `Update` y regeneración condicional de nativos.
- [x] 5. Actualizar `internal/cli/help.go` (uso y detalle de `init`).
- [x] 6. Escribir tests: `manifest_test.go`, `update_test.go`; ajustar `help_test.go` y `cli_test.go`.
- [x] 7. Ejecutar la prueba de aceptación en proyecto temporal (V4) y documentar la evidencia.
- [x] 8. Pasar V1–V3 y V5 (tests, vet, build, `rei check` y `--check` de nativos) y registrar resultados.
