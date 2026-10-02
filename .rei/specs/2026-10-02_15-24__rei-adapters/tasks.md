# Tasks — rei runtime adapters

> Work Item: `2026-10-02_15-24__rei-adapters` (feature).
>
> Orden de ejecución. Cada tarea indica los requisitos que implementa.
> El Implementer marca `[x]` **inmediatamente** al terminar cada tarea.

---

## Fase 1 — Formato canónico de rol

- [x] **T1 — Refactorizar `.rei/agents/leader.md`** al formato canónico: frontmatter con `name`, `description`, `mode: primary`, `tools: [read, search, shell, write, edit, subagent]`; sección `## Contrato` autosuficiente (identidad, objetivo, precondiciones, protocolo completo por casos A–H, reglas duras, formato de salida, herramientas permitidas, punteros a `.rei/docs/harness/`) y sección `## Referencia`. El Contrato DEBE incorporar la regla de fallback: en runtime no nativo, transmitir al subagente únicamente el `## Contrato` del rol, nunca el archivo completo. (R1, R2, R3, R4, R6, R7, R8, R14, R15)

- [x] **T2 — Refactorizar `.rei/agents/spec_author.md`** al formato canónico con `mode: subagent`, `tools: [read, write, edit, search, shell]`, protocolo de los casos A/B/C completo dentro del Contrato y reglas duras (nunca implementa código) en el Contrato. (R1–R4, R6–R8, R10)

- [x] **T3 — Refactorizar `.rei/agents/implementer.md`** al formato canónico con `mode: subagent`, `tools: [read, write, edit, search, shell]`, precondiciones, casos A/B/C, reglas duras y formato de salida dentro del Contrato. (R1–R4, R6–R8, R10)

- [x] **T4 — Refactorizar `.rei/agents/reviewer.md`** al formato canónico con `mode: subagent`, `tools: [read, write, edit, search, shell]`, precondiciones, casos A/B, reglas duras y formato de salida dentro del Contrato. (R1–R4, R6–R8, R10)

- [x] **T5 — Refactorizar `.rei/agents/initializer.md`** al formato canónico con `mode: subagent`, `tools: [read, write, edit, search, shell]`, precondiciones, alcance de escritura, plan de pasos, reglas duras y formato de salida dentro del Contrato. (R1–R4, R6–R8, R10)

- [x] **T6 — Verificar la ausencia de contenido normativo fuera de las secciones** en los 5 roles: el frontmatter es válido y todo el texto del rol está bajo `## Contrato` o `## Referencia`. (R5)

## Fase 2 — Infraestructura de adaptadores

- [x] **T7 — Crear `.rei/adapters/opencode/agent.tmpl`** (plantilla `text/template` del agente nativo con `Description`, `Mode`, `Model`, `Permission`, `Contract`) y **`.rei/adapters/opencode/tools.json`** (mapa herramienta genérica → claves de permiso OpenCode: `read→[read]`, `write→[edit]`, `edit→[edit]`, `search→[glob,grep]`, `shell→[bash]`, `subagent→[task]`). (R16, R17, R18, R19)

- [x] **T8 — Implementar el parseo de roles en `internal/adapter/adapter.go`**: `Role`, `ParseRole(data []byte) (Role, error)` (parser mínimo de frontmatter sin dependencias: `clave: valor` y `tools` como lista inline) y `ExtractContract(body string) (string, error)`. Validar `mode ∈ {primary, subagent}`, herramientas del conjunto genérico y presencia/orden de `## Contrato` y `## Referencia`. (R1, R2, R3, R4, R5, R20, R38)

- [x] **T9 — Implementar el mapeo OpenCode en `internal/adapter/opencode.go`**: cargar `tools.json`, mapear las herramientas del rol a un mapa `permission` V2 (deduplicando `write`+`edit`), aplicar `agent.tmpl`, añadir la marca GENERATED y validar que el mapa cubre todas las genéricas y usa claves de permiso válidas (`read`, `edit`, `glob`, `grep`, `bash`, `task`). (R11, R12, R13, R18, R23, R24, R38)

- [x] **T10 — Implementar `GenerateAll`, `Check` e `Install` en `internal/adapter/generate.go`**: `GenerateAll` recorre los 5 roles, genera `.opencode/agents/<rol>.md`, no sobrescribe archivos sin marca GENERATED y reporta creados/actualizados/omitidos; `Check` compara la generación canónica con el disco sin escribir; `Install` reutiliza `initwizard.InstallSkeleton` + `check.EnsureStructure` y luego `GenerateAll`. (R19, R21, R22, R25, R26, R27, R28)

## Fase 3 — Tests de la infraestructura

- [x] **T11 — Tests de parseo y extracción (`internal/adapter/adapter_test.go`)**: rol válido; falta de frontmatter; `mode` inválido; herramienta desconocida; falta de `## Contrato`; falta de `## Referencia`; orden invertido; `ExtractContract` excluye la Referencia y recorta límites. (R33)

- [x] **T12 — Tests de mapeo y generación (`internal/adapter/generate_test.go`)**: cada genérica produce las claves de permiso esperadas; `write`+`edit` deduplican; mapa incompleto ⇒ error; `GenerateAll` crea un archivo por rol con frontmatter y Contrato; segunda ejecución idéntica; archivo preexistente sin marca no se sobrescribe y avisa; `Check` devuelve 0 con todo generado, 1 si falta un archivo y 1 si el contenido difiere. (R33)

- [x] **T13 — Test estructural de los Contratos**: para cada uno de los 5 roles, verificar que el `## Contrato` contiene las subsecciones mínimas (identidad, objetivo, precondiciones, protocolo, reglas, formato de salida, herramientas y punteros a documentos) y que la referencia queda fuera. (R32, R33)

## Fase 4 — CLI, empaquetado y documentación

- [x] **T14 — Integrar el subcomando en `cmdInit` (`internal/cli/cli.go`)**: `rei init opencode` → `adapter.Install`; `rei init opencode --check` → `adapter.Check`; conservar `rei init` y `rei init status`; argumentos no reconocidos ⇒ uso + `2`. (R21, R27, R29, R30)

- [x] **T15 — Actualizar la ayuda (`internal/cli/help.go` y `internal/cli/help_test.go`)**: documentar `opencode`, `--check` y los códigos de salida de `rei init`. (R35)

- [x] **T16 — Empaquetar los adaptadores en `embed.go`**: añadir `all:.rei/adapters` a la directiva `//go:embed` y ampliar `embed_test.go` para exigir `agent.tmpl` y `tools.json` en `Skeleton`. (R31)

- [x] **T17 — Documentar el formato canónico y la adaptación** en `.rei/docs/harness/` (sección nueva "Roles y adaptadores" en `workflow.md` o documento `adapters.md` enlazado desde `AGENTS.md`): frontmatter, Contrato vs Referencia, extracción/generación, cómo añadir un runtime y el comando `rei init opencode [--check]`. Actualizar `AGENTS.md` (carga de contexto y rol del Leader) para reflejar el formato y el fallback. (R14, R19, R36)

## Fase 5 — Verificación y evidencia

- [x] **T18 — Redactar la checklist de cobertura por rol y la medición de tokens** en `.rei/progress/work-items/2026-10-02_15-24__rei-adapters/impl.md`: por cada rol, cada regla dura y cada paso del protocolo actual con su ubicación (Contrato/Referencia), y el tamaño en caracteres/tokens del Contrato frente al archivo completo. (R32, R34)

- [x] **T19 — Verificación manual del flujo**: `rei init opencode` genera los 5 archivos con la marca GENERATED; editar a mano un archivo no generado y confirmar que no se sobrescribe; `rei init opencode --check` devuelve 0 con todo generado y 1 tras alterar un archivo; `rei init` y `rei init status` conservan su comportamiento. (R21–R29)

- [x] **T20 — Ejecutar la batería de validación** y corregir lo que falle: `gofmt -l internal` (sin salida), `make vet`, `make test`, `make build`, `rei check --quiet` (`exit 0`). (R1–R38)

- [x] **T21 — Marcar como completadas (`[x]`) todas las tareas** de este archivo a medida que se terminan, dejando `tasks.md` al día. (trazabilidad)
