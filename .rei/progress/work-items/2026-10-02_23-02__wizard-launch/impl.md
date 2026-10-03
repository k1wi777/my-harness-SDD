# Implementación — lanzar el wizard de forma explícita

> Work Item: `2026-10-02_23-02__wizard-launch` (task).
> Implementer · estado final: `review`.

## Resumen

Se convierte la personalización inicial (`initializer`) en una acción
**explícita** del usuario, en vez de una auto-delegación disparada por la mera
detección de marcadores pendientes. Se implementaron los 10 pasos de `plan.md`
sin modificar la planificación ni ampliar el alcance.

Cambios de comportamiento:

1. `AGENTS.md` §1 ya no obliga al Leader a delegar en `initializer` al detectar
   pendientes; §6 documenta el trigger explícito («personaliza mi proyecto»).
2. El `## Contrato` de `.rei/agents/leader.md` incorpora el subcaso
   «Personalización del proyecto»: delega en `initializer` sin crear Work Item.
3. `rei init` (sin args) y `rei init opencode|claude` terminan indicando
   `/personalize` (y el fallback al rol `initializer`).
4. Cada adaptador genera el comando `/personalize`:
   `.opencode/commands/personalize.md` (`agent: initializer`, `subtask: true`) y
   `.claude/commands/personalize.md`.
5. El adaptador de Claude genera además `CLAUDE.md` con `@AGENTS.md`.

## Archivos

### Nuevos

| Archivo | Contenido |
|---------|-----------|
| `.rei/adapters/opencode/command.tmpl` | Comando `/personalize` de OpenCode: frontmatter `description`/`agent: initializer`/`subtask: true` + marca `{{.Mark}}`. |
| `.rei/adapters/claude/command.tmpl` | Comando `/personalize` de Claude Code: frontmatter `description` + marca `{{.Mark}}`. |
| `.opencode/commands/personalize.md` | Nativo regenerado (OpenCode). |
| `.claude/agents/*.md`, `.claude/commands/personalize.md`, `CLAUDE.md` | Nativos regenerados (Claude Code). |

### Modificados

| Archivo | Cambio |
|---------|--------|
| `AGENTS.md` | §1: eliminado el párrafo de auto-delegación en `initializer`. §6: documentado el trigger explícito y el comando `/personalize`. |
| `.rei/agents/leader.md` | Nuevo subcaso «Personalización del proyecto» dentro de `## Contrato`, tras `#### Arranque`. |
| `.rei/docs/harness/adapters.md` | Documentados el comando `/personalize` de cada adaptador, `CLAUDE.md` (`@AGENTS.md`) y `command.tmpl` en «Añadir un runtime nuevo». |
| `internal/adapter/generate.go` | Descriptor `runtime.artifacts []artFile`; helpers `renderTemplate`, `writeGenerated` y `checkGenerated` (marca GENERATED, idempotencia, no-sobrescritura compartidas); `generateAll`/`check` recorren agentes + artefactos; `install` imprime el cierre con `/personalize`. |
| `internal/adapter/opencode.go` | `opencodeRuntime.artifacts` incluye `.opencode/commands/personalize.md`. |
| `internal/adapter/claude.go` | `claudeRuntime.artifacts` incluye `.claude/commands/personalize.md` y `CLAUDE.md` (`claudeMDArtifact`). |
| `internal/initwizard/initwizard.go` | `initializerPath` → exportado `InitializerPath`; mensaje final de `Init` con `/personalize`. |
| `internal/cli/help.go` | Ayuda de `init` menciona `/personalize`, `.opencode/commands/`, `.claude/commands/` y `CLAUDE.md`. |
| `internal/adapter/generate_test.go`, `claude_test.go` | `writeAdapter` escribe `command.tmpl`; tests de generación/`--check`/no-sobrescritura del comando y de `CLAUDE.md`; test de `install` con el esqueleto real. |
| `internal/initwizard/initwizard_test.go` | `TestInitMencionaPersonalize`. |
| `internal/cli/help_test.go` | `TestHelpInitDocumentaPersonalize`. |
| `embed_test.go` | El esqueleto incluye `command.tmpl` de ambos runtimes; contenido del de OpenCode. |
| `.opencode/agents/leader.md` | Nativo regenerado (incluye el nuevo subcaso). |
| `.rei/specs/2026-10-02_23-02__wizard-launch/plan.md` | Pasos 1–10 marcados `[x]`. |

No se tocó Cursor, Codex ni el formato canónico de rol / mapa de herramientas.

## Verificación

### Automática

```text
$ gofmt -l .            # sin salida
$ go vet ./...          # exit 0 (sin salida)
$ make test             # exit 0 (todos los paquetes ok)
$ make build            # exit 0 -> bin/rei
$ rei check --quiet     # exit 0
$ ./bin/rei init opencode --check   # exit 0
$ ./bin/rei init claude --check     # exit 0
```

`make test` cubre `internal/adapter` (generación del comando y de `CLAUDE.md`,
idempotencia, no sobrescritura, `--check` con deriva, `install` con el esqueleto
real), `internal/initwizard` (mensaje con `/personalize`), `internal/cli` (ayuda
de `init`) y el `Skeleton` (plantillas `command.tmpl` embebidas).

### Manual (proyecto temporal limpio)

```text
$ rei init
...
Personalización pendiente. Inicia la entrevista con el comando
`/personalize` (OpenCode o Claude Code) o invocando el rol `initializer`
(.rei/agents/initializer.md). Consulta el detalle con `rei init status`.
exit 0

$ rei init opencode
...
[OK]    .opencode/commands/personalize.md
Para personalizar el proyecto ejecuta el comando `/personalize` (runtime opencode)
o invoca el rol `initializer` (.rei/agents/initializer.md).
exit 0
$ rei init opencode --check     # exit 0 (6 [OK]: 5 agentes + comando)

$ rei init claude
...
[OK]    .claude/commands/personalize.md
[OK]    CLAUDE.md
exit 0
$ rei init claude --check       # exit 0 (7 [OK]: 5 agentes + comando + CLAUDE.md)
```

Artefactos generados (evidencia):

```text
$ cat .opencode/commands/personalize.md
---
description: Lanza la entrevista de personalización del proyecto (rol initializer).
agent: initializer
subtask: true
---

<!-- GENERATED by rei — no editar; regenerar con `rei init opencode` -->

Conduce la entrevista guiada de personalización de este proyecto siguiendo tu
rol `initializer`.

$ cat .claude/commands/personalize.md
---
description: Lanza la entrevista de personalización del proyecto (rol initializer).
---

<!-- GENERATED by rei — no editar; regenerar con `rei init claude` -->

Invoca al subagente `initializer` para conducir la entrevista guiada de
personalización de este proyecto.

$ cat CLAUDE.md
<!-- GENERATED by rei — no editar; regenerar con `rei init claude` -->

@AGENTS.md
```

Idempotencia y no-sobrescritura (proyecto temporal):

```text
$ rei init opencode        # 2.ª ejecución: sin [UPD] en generados
$ rei init claude          # 2.ª ejecución: sin [UPD] en generados
$ echo manual > .opencode/commands/personalize.md
$ rei init opencode
[WARN]  .opencode/commands/personalize.md no fue generado por rei; se conserva sin cambios. ...
$ cat .opencode/commands/personalize.md   # manual (conservado)
$ rei init opencode --check               # exit 1 (detecta el archivo sin marca)
```

Regeneración de nativos en este repositorio (`bin/rei`):

```text
$ ./bin/rei init opencode
[UPD]   .opencode/agents/leader.md
[OK]    .opencode/commands/personalize.md
$ ./bin/rei init claude
[OK]    .claude/agents/{leader,spec_author,implementer,reviewer,initializer}.md
[OK]    .claude/commands/personalize.md
[OK]    CLAUDE.md
```

## Observaciones

- `base_commit` se conserva (`1b81195833bcb5b49635fd69261660125b961f8f`); no se
  hizo `git commit`.
- Se registraron con `git add` los archivos nuevos del Work Item (plantillas
  `command.tmpl`, nativos `.claude/**`, `.opencode/commands/**`, `CLAUDE.md` y
  `.rei/specs/...`); no se tocó `.rei/progress/`.
- Sin dependencias externas: solo biblioteca estándar (`text/template`, `embed`).
- `install` imprime el cierre con `rt.name`; `--check` no imprime el mensaje.
- `.rei/progress/current.md` queda en `review` con agente activo `implementer`.
