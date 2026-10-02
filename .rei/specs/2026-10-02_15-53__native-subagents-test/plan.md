# Plan — probar subagentes nativos

Objetivo

Documentar en `.rei/docs/usage.md` el uso de los **subagentes nativos de
OpenCode**. La guía debe explicar, en lenguaje dirigido a personas:

- que `rei init opencode` genera `.opencode/agents/<rol>.md`;
- cómo invocarlos: con `@spec_author`, `@implementer`, `@reviewer` o
  `@initializer`, o cambiando al agente `leader`;
- que `rei init opencode --check` detecta deriva de los archivos generados;
- que, sin adaptador nativo (fallback), el Leader transmite el mismo
  `## Contrato` del rol.

Fuera de alcance: cualquier otro runtime (no mencionar ni documentar otros
adaptadores).

Archivos

- `.rei/docs/usage.md` — único archivo a modificar. Añadir una sección nueva
  sobre subagentes nativos de OpenCode e integrarla con la numeración y el
  estilo existentes.

Cambios

- Insertar una sección dedicada (p. ej. "Subagentes nativos de OpenCode") en
  `.rei/docs/usage.md`, redactada para personas y coherente con el tono actual.
- Explicar que `rei init opencode` genera un archivo por rol en
  `.opencode/agents/<rol>.md` a partir del `## Contrato` canónico.
- Explicar las dos formas de invocación: mencionar `@spec_author`,
  `@implementer`, `@reviewer` o `@initializer`, o cambiar al agente `leader`.
- Explicar que `rei init opencode --check` verifica sin escribir y detecta
  deriva (código de salida 0 si coincide, 1 si falta o difiere).
- Aclarar el fallback: sin adaptador nativo, el Leader transmite al subagente
  únicamente el `## Contrato` del rol (nunca el archivo completo).
- No documentar otros runtimes.

Restricciones

- Lenguaje dirigido a personas; no convertir la guía en documentación de agente.
- No alterar otras secciones de `usage.md` ni otros documentos.
- Tono y formato consistentes con el resto de la guía (encabezados `#`,
  tablas/lisas, bloques de código cuando aplique).
- No inventar comportamiento: ceñirse a `meta.json` y a lo descrito en
  `.rei/docs/harness/adapters.md`.

Pasos

- [x] 1. Localizar el punto de inserción en `.rei/docs/usage.md` y decidir la numeración de la nueva sección.
- [x] 2. Redactar la sección sobre subagentes nativos de OpenCode cubriendo: generación de `.opencode/agents/<rol>.md`, formas de invocación (`@...` o agente `leader`), `--check` y fallback.
- [x] 3. Verificar coherencia de estilo, encabezados y referencias cruzadas con el resto de `usage.md`.
- [x] 4. Revisar que no se introduzcan menciones a otros runtimes ni cambios fuera de alcance.
