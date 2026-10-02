# Implementación — rei runtime adapters

> Work Item: `2026-10-02_15-24__rei-adapters` (feature).
> Implementer · estado final: `review`.

## Resumen

Se introduce un **formato canónico de rol** en `.rei/agents/<rol>.md`
(frontmatter YAML + `## Contrato` autosuficiente + `## Referencia` opcional) y se
construye un **adaptador de runtime** para **OpenCode V2** que genera
`.opencode/agents/<rol>.md` a partir del Contrato, con el comando
`rei init opencode [--check]`.

Los 5 roles (`leader`, `spec_author`, `implementer`, `reviewer`,
`initializer`) se refactorizaron moviendo **todo** el protocolo y **todas** las
reglas duras al `## Contrato`; la `## Referencia` queda como sección obligatoria
con una nota de que no hay detalle adicional (decisión D7 de `design.md`).

El runtime nativo y el fallback entregan el **mismo Contrato**: en un runtime sin
adaptador, el Leader transmite solo el `## Contrato`; en OpenCode, ese Contrato
es el system prompt del agente generado, con el frontmatter nativo `permission`
mapeado desde `tools`.

Se completaron T1–T21 según `tasks.md`. No se modificó ninguna decisión de
planificación ni se amplió el alcance (solo OpenCode; R37).

## Archivos

### Nuevos

| Archivo | Contenido |
|---------|-----------|
| `internal/adapter/adapter.go` | `Role`, `ParseRole`, `ExtractContract`, parser mínimo de frontmatter (sin dependencias) y lista canónica `Roles`. |
| `internal/adapter/opencode.go` | Mapeo a OpenCode V2 (`permission`), carga de `agent.tmpl`/`tools.json`, marca GENERATED y deduplicación `write`+`edit`. |
| `internal/adapter/generate.go` | `GenerateAll`, `Check` e `Install`. |
| `internal/adapter/adapter_test.go` | Parseo, extracción, validación de frontmatter y test estructural de los 5 Contratos. |
| `internal/adapter/generate_test.go` | Mapeo de permisos, generación, idempotencia, no sobrescritura y `--check`. |
| `.rei/adapters/opencode/agent.tmpl` | Plantilla `text/template` del agente nativo (frontmatter + Contrato). |
| `.rei/adapters/opencode/tools.json` | Mapa herramienta genérica → claves de permiso de OpenCode. |
| `.rei/docs/harness/adapters.md` | Documentación del formato canónico y del procedimiento de adaptación. |
| `.opencode/agents/<rol>.md` (×5) | Salida generada (dogfooding) con la marca GENERATED. |

### Modificados

| Archivo | Cambio |
|---------|--------|
| `.rei/agents/leader.md` | Formato canónico; regla de fallback (transmitir el `## Contrato`, nunca el archivo completo). |
| `.rei/agents/spec_author.md` | Formato canónico. |
| `.rei/agents/implementer.md` | Formato canónico. |
| `.rei/agents/reviewer.md` | Formato canónico. |
| `.rei/agents/initializer.md` | Formato canónico. |
| `AGENTS.md` | Menciona el formato de rol, el fallback y `.rei/adapters/`; enlaza `adapters.md`. |
| `.rei/docs/harness/README.md` | Añade `adapters.md` al índice. |
| `embed.go` | Incluye `all:.rei/adapters` en `Skeleton`. |
| `embed_test.go` | Exige `agent.tmpl` y `tools.json` en el esqueleto embebido. |
| `internal/cli/cli.go` | `cmdInit` acepta `opencode [--check]`; conserva `init` y `init status`. |
| `internal/cli/help.go` | Documenta `opencode`, `--check` y los códigos de salida. |
| `internal/cli/help_test.go` | Ajusta las aserciones de la ayuda general y de `init`. |

## Trazabilidad tareas → requisitos

- T1 → R1, R2, R3, R4, R6, R7, R8, R14, R15.
- T2–T5 → R1–R4, R6–R8, R10.
- T6 → R5.
- T7 → R16, R17, R18, R19.
- T8 → R1–R5, R20, R38.
- T9 → R11–R13, R18, R23, R24, R38.
- T10 → R19, R21, R22, R25–R28.
- T11–T13 → R32, R33.
- T14 → R21, R27, R29, R30.
- T15 → R35.
- T16 → R31.
- T17 → R14, R19, R36.
- T18 → R32, R34.
- T19 → R21–R29.
- T20 → R1–R38.
- T21 → trazabilidad.

## Checklist de cobertura por rol (R32)

`## Contrato` = la regla/paso es necesario y está en el Contrato.
`## Referencia` = detalle opcional. En esta Feature la Referencia solo contiene
la nota «sin detalle adicional», por lo que **todas** las reglas y pasos están
en el Contrato.

### `leader`

| Regla / paso del rol | Ubicación |
|----------------------|-----------|
| Identidad: coordina, nunca implementa | Contrato · Identidad |
| Arranque: `rei session` + `rei items status`; caso D; campos de `meta.json` | Contrato · Protocolo |
| Antes de delegar: comprender, dudas, tipo, `id`, `rei new`, `title`/`description`, continuar | Contrato · Protocolo |
| Casos A–H completos | Contrato · Protocolo |
| Delegación: prompt, contexto, `rei validate`, anti teléfono descompuesto | Contrato · Protocolo |
| **Fallback**: transmitir solo el `## Contrato`, nunca el archivo completo | Contrato · Protocolo/Delegación y Reglas duras |
| Escalado de solicitudes grandes | Contrato · Protocolo |
| Reglas duras (no implementar, no escribir tests, no sustituir, no inventar, no modificar, no avanzar/omitir etapas, no asumir aprobación, no delegar desde el chat, no continuar sin aprobación, meta qué/no cómo, fallback, no releer) | Contrato · Reglas duras |
| Formato de salida | Contrato · Formato de salida |
| Herramientas permitidas | Contrato · Herramientas permitidas |
| Punteros a docs del arnés | Contrato · Documentos de referencia |

### `spec_author`

| Regla / paso del rol | Ubicación |
|----------------------|-----------|
| Identidad: planifica, nunca implementa | Contrato · Identidad |
| Protocolo 1–6 (docs por tipo, meta, status, sesión, docs de proyecto) | Contrato · Protocolo |
| Caso A (feature): requirements/design/tasks + `current.md` + `ready` | Contrato · Protocolo |
| Caso B (task): `plan.md` + `current.md` + `ready` | Contrato · Protocolo |
| Caso C (revisión sobre plan listo) | Contrato · Protocolo |
| Bloqueos (marcar `blocked`, `spec.md`) | Contrato · Protocolo |
| Reglas duras (no implementar, no tests, no cambiar alcance, no inventar, requisitos verificables, diseño justificado, tareas derivadas, sesión, `ready`, plantillas vía CLI, no releer) | Contrato · Reglas duras |
| Formato de salida (`ready`/`blocked`) | Contrato · Formato de salida |
| Herramientas permitidas | Contrato · Herramientas permitidas |
| Punteros a docs del arnés | Contrato · Documentos de referencia |

### `implementer`

| Regla / paso del rol | Ubicación |
|----------------------|-----------|
| Identidad: implementa, no replanifica | Contrato · Identidad |
| Precondiciones (`in_progress`, `current.md`, docs por tipo) | Contrato · Precondiciones |
| Protocolo 1–6 (docs, meta, `current.md`, rework) | Contrato · Protocolo |
| Caso A (feature): leer plan, orden, marcar al completar, cierre | Contrato · Protocolo |
| Caso B (task): plan, pasos, cierre | Contrato · Protocolo |
| Caso C (rework) | Contrato · Protocolo |
| Bloqueos | Contrato · Protocolo |
| Reglas duras (un Work Item, no cambiar planificación, no inicializar `current.md`, no inventar, no `done` sin Reviewer, convenciones, marcar tareas, `current.md`, `impl.md`, verificar, no releer, `git add` sin `.rei/progress/`) | Contrato · Reglas duras |
| Formato de salida (`review`/`blocked`) | Contrato · Formato de salida |
| Herramientas permitidas | Contrato · Herramientas permitidas |
| Punteros a docs del arnés | Contrato · Documentos de referencia |

### `reviewer`

| Regla / paso del rol | Ubicación |
|----------------------|-----------|
| Identidad: aprueba/rechaza, no implementa | Contrato · Identidad |
| Precondiciones (`status == review`) | Contrato · Precondiciones |
| Protocolo 1–4 (`verification.md`, meta, `review-diff`, modo lectura) | Contrato · Protocolo |
| Caso A (feature): comprobaciones, `review.md`, cierre/rechazo | Contrato · Protocolo |
| Caso B (task): comprobaciones, `review.md`, cierre/rechazo | Contrato · Protocolo |
| Bloqueos | Contrato · Protocolo |
| Reglas duras (no implementar, no cambiar planificación, no aprobar con fallos/`rei check`/desviación, `review-diff`, justificar rechazos, `review.md`, `session archive`, plantillas vía CLI, no releer) | Contrato · Reglas duras |
| Formato de salida (`done`/`changes_requested`/`blocked`) | Contrato · Formato de salida |
| Herramientas permitidas | Contrato · Herramientas permitidas |
| Punteros a docs del arnés | Contrato · Documentos de referencia |

### `initializer`

| Regla / paso del rol | Ubicación |
|----------------------|-----------|
| Identidad: personaliza documentación, no implementa | Contrato · Identidad |
| Precondiciones (`rei init status == 1`) | Contrato · Precondiciones |
| Alcance de escritura (5 destinos) | Contrato · Alcance de escritura |
| Protocolo 1–5 (status, reanudabilidad, entrevista, `config.json`, cierre) | Contrato · Protocolo |
| Plan de pasos (6 pasos) | Contrato · Protocolo |
| Entrevista guiada | Contrato · Protocolo |
| Criterio de finalización | Contrato · Protocolo |
| Bloqueos | Contrato · Protocolo |
| Reglas duras (no implementar, no inventar, no salirse del alcance, no dejar marcadores, orden, sugerencias, confirmación, `config.json`, `init status`/`doctor`, no releer) | Contrato · Reglas duras |
| Formato de salida (`listo`/`blocked`) | Contrato · Formato de salida |
| Herramientas permitidas | Contrato · Herramientas permitidas |
| Punteros a docs del arnés | Contrato · Documentos de referencia |

Además, `TestContratosDeLosRolesSonEstructurados` exige en los 5 Contratos las
subsecciones mínimas (identidad, objetivo, precondiciones, protocolo, reglas,
salida, herramientas y documentos de referencia) y que el cuerpo de la
Referencia quede fuera.

## Medición de tokens (R34)

Tamaño del archivo de rol frente a la sección `## Contrato` (lo que se transmite
por spawn). Estimación de tokens ≈ caracteres / 4.

| Rol | Rol completo (chars) | Contrato (chars) | Ahorro (chars) | Tokens completo | Tokens Contrato | Ahorro (tokens) |
|-----|----------------------|------------------|----------------|-----------------|-----------------|-----------------|
| `leader` | 8736 | 8358 | 378 | ~2184 | ~2089 | ~94 |
| `spec_author` | 5609 | 5290 | 319 | ~1402 | ~1322 | ~79 |
| `implementer` | 6122 | 5744 | 378 | ~1530 | ~1436 | ~94 |
| `reviewer` | 6187 | 5835 | 352 | ~1546 | ~1458 | ~88 |
| `initializer` | 6158 | 5699 | 459 | ~1539 | ~1424 | ~114 |

En esta Feature la `## Referencia` solo contiene la nota de «sin detalle
adicional» (D7), por lo que el ahorro actual es modesto. El valor del mecanismo
es que, cuando la Referencia acumule racional y ejemplos, el spawn seguirá
transmitiendo únicamente el Contrato sin tocar el formato.

## Verificación

### Automática

```text
$ gofmt -l .            # sin salida
$ make vet              # exit 0
$ make test             # exit 0 (todos los paquetes ok)
$ make build            # exit 0 → bin/rei
$ rei check --quiet     # exit 0
```

`make test` incluye `internal/adapter` (parseo, extracción, mapeo, generación,
idempotencia, no sobrescritura, `--check` y estructura de los Contratos), la
verificación de `Skeleton` (adaptadores embebidos) y la ayuda del CLI.

### Manual (T19)

Sobre el propio repositorio, tras `make build`:

```text
$ rei init opencode
== rei init opencode ==
...
[OK]    .opencode/agents/leader.md
[OK]    .opencode/agents/spec_author.md
[OK]    .opencode/agents/implementer.md
[OK]    .opencode/agents/reviewer.md
[OK]    .opencode/agents/initializer.md
exit 0

$ rei init opencode --check      # 0 con todo generado
$ rei init opencode --check      # 1 tras añadir una línea manual a implementer.md
[FAIL]  .opencode/agents/implementer.md difiere de la generación canónica
```

No sobrescritura de archivos no generados (se sustituyó temporalmente
`leader.md` por contenido manual sin marca y se restauró después):

```text
$ rei init opencode
[WARN]  .opencode/agents/leader.md no fue generado por rei; se conserva sin cambios. Elimínalo o regenera el archivo a mano.
# el contenido manual se conservó intacto

$ rei init opencode --check
[FAIL]  .opencode/agents/leader.md existe pero no fue generado por rei (sin la marca GENERATED)
exit 1
```

Comportamiento conservado de `rei init` y `rei init status`:

```text
$ rei init                # exit 0, instalador + reporte de personalización + plan de pasos
$ rei init status         # exit 1 (4 documentos pendientes), sin escribir
$ rei init bogus          # exit 2, uso por stderr
$ rei init opencode foo   # exit 2, uso por stderr
```

## Observaciones

- `.rei/progress/current.md` queda en `review` con agente activo `implementer`.
- `base_commit` se conserva: `05dcdae6b06d72822be55a43adfe14efca555b5e`; no se
  hizo `git commit`.
- Se conservan los 5 archivos generados en `.opencode/agents/` como dogfooding
  del adaptador (paso 4 de la delegación).
- El parser de frontmatter es propio y mínimo (D3): no se añaden dependencias;
  `go.mod` sigue sin `require`.
- `check.RequiredFiles` no cambia (D8): `rei check` no falla en proyectos
  inicializados sin adaptadores; la completitud se verifica con
  `rei init opencode --check`.
- La `## Referencia` se documenta como sección obligatoria pero vacía (D7);
  toda regla y paso vive en el `## Contrato`.
