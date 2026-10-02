# Revisión — probar subagentes nativos

> Work Item: `2026-10-02_15-53__native-subagents-test` (task).
> Reviewer · veredicto: **`done`**.
> Base revisada: `5e9bd7d7921901317ec078d3e2dfc1d7a025f44c` (`base_commit`).

## Estado final

**`done`** — el Work Item cumple su planificación aprobada (`plan.md`) sin
desviaciones de alcance.

Registro consultado:

- `rei review-diff 2026-10-02_15-53__native-subagents-test` → exit `0`, pero
  **paquete vacío**: `(sin cambios de código)`. Como se anticipó, `review-diff`
  no expone cambios bajo `.rei/` salvo la propia spec del Work Item, así que la
  revisión se realizó en **modo lectura complementaria**, contrastando
  directamente `.rei/docs/usage.md` con `plan.md` e `impl.md` y con
  `git diff 5e9bd7d -- .rei/docs/usage.md`.
- `rei validate 2026-10-02_15-53__native-subagents-test` → `Resultado: OK`.

## Verificaciones realizadas

### Checkpoints del proyecto

| ID | Comando | Resultado |
|----|---------|-----------|
| — | `rei check --quiet` | `exit 0` |

> `.rei/docs/project/verification.md` conserva los ejemplos de plantilla
> (`V1`–`V3` sin definir) porque el proyecto no está personalizado, y
> `.rei/config.json` declara `"checks": []`. Por eso el único checkpoint
> aplicable es `rei check --quiet`, que pasa con `exit 0`, igual que registra
> `impl.md`.

### Objetivo: los 4 puntos documentados

Contrastado contra el diff real (`git diff 5e9bd7d -- .rei/docs/usage.md`, 29
líneas añadidas al final del documento) y contra `.rei/docs/harness/adapters.md`:

1. **Generación de agentes.** La Sección 7 explica que `rei init opencode` crea
   `.opencode/agents/<rol>.md` (p. ej. `.opencode/agents/implementer.md`) a
   partir del `## Contrato` canónico de `.rei/agents/<rol>.md`, con la marca
   `GENERATED`, comportamiento idempotente y sin sobrescribir archivos ajenos.
   Coincide con `adapters.md` (líneas 80–82 y 120–125).
2. **Formas de invocación.** Documenta invocar `@spec_author`, `@implementer`,
   `@reviewer` o `@initializer`, o cambiar al agente `leader`. Coincide con
   `plan.md`.
3. **Detección de deriva.** Documenta `rei init opencode --check` como
   verificación sin escritura, con `exit 0` si todo coincide y `exit 1` si falta
   o difiere. Coincide con `adapters.md` (líneas 117, 122–124).
4. **Fallback.** Aclara que, sin adaptador nativo, el Leader transmite al
   subagente únicamente el `## Contrato` del rol, nunca el archivo completo.
   Coincide con `adapters.md` (líneas 77–79).

### Restricciones y alcance

- **Un solo archivo modificado.** El único cambio de contenido del Work Item es
  la sección nueva en `.rei/docs/usage.md`. El resto del diff respecto a
  `base_commit` son artefactos de planificación/progreso (`meta.json`,
  `plan.md`, `current.md`) y el directorio de `work-items/`.
- **Sin tocar secciones existentes.** El diff es puramente aditivo al final del
  documento: no se modificó ninguna línea de las Secciones 1–6 ni de los
  encabezados previos. Los encabezados numerados resultan `1`–`7`, secuenciales
  y sin duplicados.
- **Sin otros runtimes.** Búsqueda de `claude|cursor|gemini|aider|copilot|codex|windsurf`
  en `usage.md`: sin coincidencias. La sección solo menciona OpenCode y el
  fallback genérico.
- **Sin inventar comportamiento.** Todo lo afirmado (marca `GENERATED`,
  idempotencia, no sobrescritura, códigos `0`/`1`, fallback) está respaldado por
  `.rei/docs/harness/adapters.md`.
- **Tono para personas.** La sección está redactada de forma narrativa, con
  `---`, encabezado `# 7.` y bloques de código, consistente con el resto de la
  guía.

## Trazabilidad

Los 4 pasos de `plan.md` están marcados `[x]` y se corresponden con la sección
redactada. El único archivo planificado (`.rei/docs/usage.md`) es el único
archivo de contenido modificado.

## Observaciones

- `package-lock.json` (raíz) aparece como archivo no rastreado en `review-diff`.
  Es un residuo preexistente, ajeno a este Work Item y no incorporado por él. No
  bloquea la aprobación.
- `current.md` figuraba como `in_progress`/`implementer` pese a que `meta.json`
  ya estaba en `review`; se corrige a `review`/`reviewer` antes de archivar la
  sesión, de modo que el asiento de `history.md` sea coherente.

## Acciones requeridas

Ninguna. El Work Item se aprueba.
