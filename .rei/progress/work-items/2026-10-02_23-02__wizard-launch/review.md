# Revisión — lanzar el wizard de forma explícita

> Work Item: `2026-10-02_23-02__wizard-launch` (task).
> Reviewer · resultado: **aprobado**.

## Alcance revisado

Revisión contra `plan.md` (10 pasos), `impl.md` y el paquete de
`rei review-diff 2026-10-02_23-02__wizard-launch` (base
`1b81195833bcb5b49635fd69261660125b961f8f`). Tipo `task`: sin
`requirements.md`/`design.md`/`tasks.md`, conforme a la planificación.

## Comprobaciones

### Objetivo y planificación

- [x] **§1 AGENTS.md sin auto-delegación.** El párrafo que obligaba al Leader a
      delegar en `initializer` al detectar pendientes fue eliminado. La única
      mención a `initializer` en `AGENTS.md` está en §6 (trigger explícito).
- [x] **Trigger explícito en el `## Contrato` de `leader`.** Subcaso
      «Personalización del proyecto» insertado tras `#### Arranque` (línea 43),
      dentro de `## Contrato` y antes de `## Referencia` (línea 231). Regenerado
      en `.opencode/agents/leader.md` y `.claude/agents/leader.md`.
- [x] **`rei init` finaliza indicando la entrevista.** `initwizard.Init` y
      `adapter.install` cierran mencionando `/personalize` y el rol
      `initializer`; verificado en proyecto temporal.
- [x] **Comando `/personalize` en ambos adaptadores.**
      - OpenCode: `.opencode/commands/personalize.md` con `agent: initializer` y
        `subtask: true`.
      - Claude: `.claude/commands/personalize.md`.
      Plantillas embebidas `.rei/adapters/{opencode,claude}/command.tmpl`.
- [x] **Claude genera `CLAUDE.md`** con la marca GENERATED y `@AGENTS.md`.
- [x] **Ayuda y documentación del arnés** actualizadas: `internal/cli/help.go`
      (menciona `/personalize`, `.opencode/commands/`, `.claude/commands/`,
      `CLAUDE.md`) y `.rei/docs/harness/adapters.md` (comando por adaptador,
      `CLAUDE.md`, `command.tmpl` en «Añadir un runtime nuevo»).

### Restricciones

- [x] Solo se tocaron OpenCode y Claude; Cursor/Codex intactos.
- [x] No se modificó el formato canónico de rol ni el mapa de herramientas.
- [x] Formato correcto de OpenCode: `.opencode/commands/` (plural) con
      frontmatter `agent` y `subtask`.
- [x] No-sobrescritura de archivos sin marca GENERATED: comprobado con un
      archivo de prueba (`.opencode/commands/personalize.md` y `CLAUDE.md`); se
      emite `[WARN]` y se conserva; `--check` devuelve 1.
- [x] Idempotencia: reejecutar `rei init opencode|claude` deja los generados
      idénticos (sin `[UPD]`).
- [x] `rei init` (sin args) y `rei init status` conservan su comportamiento.
- [x] La lógica de escritura/marca/idempotencia se comparte en helpers
      (`writeGenerated`, `checkGenerated`, `renderTemplate`), sin duplicación.

### Verificación (checkpoints)

`verification.md` no define aún checkpoints `V*` (documento pendiente de
personalización); se ejecutaron las verificaciones declaradas en `plan.md`:

- [x] `go test ./...` → exit 0 (todos los paquetes `ok`).
- [x] `make vet` → exit 0.
- [x] `rei check --quiet` → exit 0.
- [x] `rei init opencode --check` → exit 0 (6 `[OK]`: 5 agentes + comando).
- [x] `rei init claude --check` → exit 0 (7 `[OK]`: 5 agentes + comando +
      `CLAUDE.md`).
- [x] Proyecto temporal limpio: `rei init opencode` + `--check` → 0;
      `rei init claude` + `--check` → 0; `rei init` (sin args) termina indicando
      `/personalize`.

Evidencia en `.rei/progress/work-items/2026-10-02_23-02__wizard-launch/impl.md`.

## Desviaciones

Ninguna. La implementación corresponde a lo planificado, sin ampliar alcance.
Nota: `package-lock.json` figura como archivo sin rastrear, pero es previo a
este Work Item (fecha sep 10; no forma parte del diff ni de los archivos
registrados) y no afecta a la revisión.

## Conclusión

El Work Item cumple la planificación y todas las verificaciones pasan.
**Aprobado.**
