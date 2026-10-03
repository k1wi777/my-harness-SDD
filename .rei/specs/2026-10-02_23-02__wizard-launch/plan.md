# Plan — lanzar el wizard de forma explícita

## Objetivo

Convertir la personalización inicial (`initializer`) en una acción **explícita**
del usuario, en vez de una auto-delegación disparada por la mera detección de
marcadores pendientes. Además, `rei init` debe explicar cómo iniciar la
entrevista y los adaptadores nativos deben exponer un comando `/personalize`
(OpenCode y Claude) y `CLAUDE.md` (Claude).

Puntos acordados con el usuario:

1. Quitar la auto-delegación de `AGENTS.md` §1.
2. Trigger explícito en el `## Contrato` de `leader`: si el usuario pide
   «personaliza mi proyecto» (o equivalente), delega en `initializer` pasándole
   su `## Contrato`; documentarlo también en `AGENTS.md`.
3. `rei init` y sus variantes (`opencode`/`claude`) terminan indicando cómo
   iniciar la entrevista.
4. El adaptador genera el comando de runtime `/personalize`:
   - OpenCode: `.opencode/commands/personalize.md` con `agent: initializer` y
     `subtask: true`.
   - Claude: `.claude/commands/personalize.md`.
5. El adaptador de Claude genera además `CLAUDE.md` con `@AGENTS.md`.

Fuera de alcance: Cursor, Codex y cualquier otro runtime.

## Archivos

### Documentación y roles (fuente de verdad)

- `AGENTS.md` — quitar el párrafo de auto-delegación de §1; documentar el
  trigger explícito (en §6).
- `.rei/agents/leader.md` — nuevo subcaso «Personalización del proyecto» dentro
  de `## Contrato`.
- `.rei/docs/harness/adapters.md` — documentar `/personalize` y `CLAUDE.md`.

### Plantillas de adaptador (nuevas, embebidas por el glob existente)

- `.rei/adapters/opencode/command.tmpl` — comando `/personalize` de OpenCode.
- `.rei/adapters/claude/command.tmpl` — comando `/personalize` de Claude.

### Código del CLI

- `internal/adapter/generate.go` — generar/verificar el comando (y `CLAUDE.md`
  en Claude), helper de escritura con marca, y mensaje final de la entrevista.
- `internal/adapter/opencode.go` — descriptor del comando de OpenCode.
- `internal/adapter/claude.go` — descriptor del comando de Claude y de
  `CLAUDE.md`.
- `internal/initwizard/initwizard.go` — mensaje final de `rei init`.
- `internal/cli/help.go` — ayuda de `init` mencionando `/personalize`.

### Tests

- `internal/adapter/generate_test.go` (OpenCode) y
  `internal/adapter/claude_test.go` (Claude) — generación, idempotencia,
  no-sobrescritura y `--check` del comando (y `CLAUDE.md`).
- `internal/initwizard/initwizard_test.go` — el mensaje final nombra
  `/personalize`.
- `embed_test.go` — el esqueleto contiene `command.tmpl` de ambos runtimes.
- `internal/cli/help_test.go` — ajustar aserciones si cambia el texto de ayuda.

### Nativos regenerados (en este repositorio, tras `make build`)

- `.opencode/agents/*.md` (cambia `leader`), `.opencode/commands/personalize.md`.
- `.claude/agents/*.md`, `.claude/commands/personalize.md`, `CLAUDE.md`.

## Cambios

### 1. `AGENTS.md`

- §1: eliminar el párrafo que obliga al Leader a delegar en `initializer` al
  detectar pendientes (`Si rei init status reporta ... antes de continuar con
  cualquier otro trabajo.`). No debe quedar ninguna auto-delegación.
- §6 (Rol del agente principal): añadir un párrafo que describa el trigger
  explícito: si el usuario pide personalizar el proyecto (p. ej. «personaliza mi
  proyecto»), el Leader delega en el rol `initializer`
  (`.rei/agents/initializer.md`) y le transmite su `## Contrato`, sin crear un
  Work Item.

### 2. `.rei/agents/leader.md` (dentro de `## Contrato`)

Insertar un subcaso en `#### Protocolo`, inmediatamente después de
`#### Arranque`:

```
#### Personalización del proyecto

Si el usuario pide personalizar el proyecto (p. ej. «personaliza mi proyecto»),
NO crees un Work Item: delega directamente en el subagente `initializer`
transmitiéndole únicamente su `## Contrato` (`.rei/agents/initializer.md`),
espera a que finalice y devuelve el control al usuario. En un runtime con
adaptador nativo, el comando `/personalize` ejecuta este mismo camino.
```

Debe quedar **dentro** de `## Contrato` (antes de `## Referencia`) para que el
adaptador lo incluya al regenerar. No alterar los encabezados `###` requeridos.

### 3. Adaptadores: comando `/personalize`

- Añadir `command.tmpl` por runtime en `.rei/adapters/<runtime>/` (se embebe
  automáticamente porque `embed.go` ya incluye `all:.rei/adapters`).
  - OpenCode: frontmatter `description`, `agent: initializer`, `subtask: true`;
    cuerpo con marca GENERATED y una instrucción breve de inicio.
  - Claude: frontmatter `description`; cuerpo con marca GENERATED e instrucción
    de invocar al subagente `initializer`.
- Extender el descriptor `runtime` (`internal/adapter/generate.go`) con la ruta
  de destino del comando (`.opencode/commands/personalize.md` /
  `.claude/commands/personalize.md`) y los datos de plantilla correspondientes.
- Reutilizar la lógica de escritura (marca GENERATED, idempotencia,
  no-sobrescritura) en un helper compartido entre agentes y comando, en lugar de
  duplicarla.
- `install` genera también el comando; `check` lo verifica (si falta o difiere,
  `--check` devuelve 1).

### 4. Claude: `CLAUDE.md`

- `rei init claude` genera además `CLAUDE.md` en la raíz con la marca GENERATED
  y el import `@AGENTS.md`.
- Aplicar las mismas reglas de marca/idempotencia/no-sobrescritura.
- `rei init claude --check` verifica también `CLAUDE.md`.

### 5. Mensaje final de la entrevista

- `internal/initwizard/initwizard.go` (`Init`): reemplazar el bloque de
  personalización pendiente para indicar cómo iniciar la entrevista:
  `/personalize` (OpenCode/Claude Code) o invocando el rol `initializer`
  (`.rei/agents/initializer.md`); mantener la referencia a `rei init status`.
- `internal/adapter/generate.go` (`install`): tras generar, imprimir el cierre
  indicando `/personalize` para el runtime correspondiente (y el fallback al rol
  `initializer`). No imprimirlo en `--check`.
- `internal/cli/help.go`: actualizar el detalle de `init` para mencionar
  `/personalize`, `.opencode/commands/`, `.claude/commands/` y `CLAUDE.md`.

### 6. Documentación del arnés

- `.rei/docs/harness/adapters.md`: documentar que cada adaptador genera el
  comando `/personalize` y que el de Claude genera además `CLAUDE.md`
  (`@AGENTS.md`); ampliar la nota de «Añadir un runtime nuevo» con
  `command.tmpl`.

### 7. Regenerar nativos

Con el binario nuevo (`bin/rei`), desde la raíz del repo:

- `bin/rei init opencode` → actualiza `.opencode/agents/leader.md` y crea
  `.opencode/commands/personalize.md`.
- `bin/rei init claude` → crea `.claude/agents/*.md`,
  `.claude/commands/personalize.md` y `CLAUDE.md`.

## Restricciones

- Es una `task`: no crear `requirements.md`/`design.md`/`tasks.md`.
- No tocar Cursor, Codex ni otros runtimes (sin adaptador nativo).
- No cambiar el formato canónico de rol ni el mapa de herramientas.
- El comando debe usar el formato verificado de OpenCode:
  `.opencode/commands/<name>.md` (plural), frontmatter con `agent` y `subtask`.
- Nunca sobrescribir archivos propios del usuario sin la marca GENERATED:
  avisar (`[WARN]`) y conservar.
- Idempotencia: reejecutar `rei init <runtime>` deja los generados idénticos.
- Mantener el comportamiento existente de `rei init` (sin args) y
  `rei init status`.
- No romper la verificación: `rei init opencode --check` y
  `rei init claude --check` deben devolver 0 tras generar en un proyecto
  temporal limpio.

## Pasos

- [x] 1. Editar `AGENTS.md`: quitar la auto-delegación de §1 y documentar el
      trigger explícito en §6.
- [x] 2. Añadir el subcaso «Personalización del proyecto» en el `## Contrato` de
      `.rei/agents/leader.md`.
- [x] 3. Crear `.rei/adapters/opencode/command.tmpl` y
      `.rei/adapters/claude/command.tmpl` con su frontmatter y marca GENERATED.
- [x] 4. Extender `internal/adapter` (descriptor `runtime`, helper de escritura
      con marca) para generar/verificar `.opencode/commands/personalize.md` y
      `.claude/commands/personalize.md`.
- [x] 5. Generar/verificar `CLAUDE.md` (`@AGENTS.md`) en el adaptador de Claude.
- [x] 6. Actualizar el mensaje final de `rei init` (`initwizard.Init`) y de
      `rei init opencode|claude` (`adapter.install`) para indicar `/personalize`.
- [x] 7. Actualizar `internal/cli/help.go` y `.rei/docs/harness/adapters.md`.
- [x] 8. Añadir/ajustar tests (adapters, initwizard, embed, help) y ejecutar
      `make test` y `make vet`.
- [x] 9. Verificar en un proyecto temporal: `rei init opencode` + `--check` = 0
      y `rei init claude` + `--check` = 0.
- [x] 10. Reconstruir (`make build`) y regenerar los nativos del repositorio
      (`.opencode/**`, `.claude/**`, `CLAUDE.md`).

## Verificación

- `make test` y `make vet` en verde.
- Proyecto temporal limpio:
  - `rei init opencode` y `rei init opencode --check` → código 0;
    existe `.opencode/commands/personalize.md` con `agent: initializer` y
    `subtask: true`.
  - `rei init claude` y `rei init claude --check` → código 0;
    existen `.claude/commands/personalize.md` y `CLAUDE.md` con `@AGENTS.md`.
- `rei init` (sin argumentos) termina indicando `/personalize` /
  rol `initializer`.
- Reejecutar los adaptadores es idempotente y no sobrescribe archivos sin la
  marca GENERATED.
- `git status` del repo muestra solo los nativos regenerados esperados.
