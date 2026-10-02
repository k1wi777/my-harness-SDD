# Revisión — slimming de contratos de rol

Work Item: `2026-10-02_16-12__contrato-slimming` (type: `task`)
Plan: `.rei/specs/2026-10-02_16-12__contrato-slimming/plan.md`
Estado final: **done**

## Alcance revisado

Tipo `task`. Se revisa el objetivo del plan: reducir el `## Contrato` de los 5
roles moviendo a `## Referencia` el detalle no esencial, **sin perder capacidad
ni cambiar el flujo**, regenerar los agentes nativos y medir el ahorro.

Verificación mediante `rei review-diff 2026-10-02_16-12__contrato-slimming`
(exit 0, base `7195524…`), que lista los 5 `.rei/agents/*.md`, los 5
`.opencode/agents/*.md` regenerados y la spec del Work Item
(`meta.json`/`plan.md`).

## Verificaciones

| Check | Comando | Resultado |
|-------|---------|-----------|
| Paquete de revisión | `rei review-diff 2026-10-02_16-12__contrato-slimming` | exit 0 |
| Test estructural | `go test ./...` | exit 0, todos los paquetes en verde |
| Harness | `rei check --quiet` | exit 0 |
| Adaptador nativo | `rei init opencode --check` | exit 0 (`[OK]` en los 5 agentes) |
| Validación del Work Item | `rei validate 2026-10-02_16-12__contrato-slimming` | `Resultado: OK`, exit 0 |

`.rei/docs/project/verification.md` sigue siendo plantilla (sin `V1`/`V2`
reales) y `.rei/config.json` tiene `checks: []`, por lo que no hay checkpoints
rápidos o lentos adicionales que ejecutar. Los checkpoints efectivos del plan
§6 (`go test ./...`, `rei init opencode --check`, `rei check --quiet`) se
verificaron en verde.

## Sin pérdida de capacidad (checklist §4 vs. roles resultantes)

Se contrastó cada ítem del **checklist de cobertura de `plan.md` §4** con los
Contratos resultantes; todas las reglas operativas conservan su ubicación
prevista (C = Contrato, R = Referencia, A = cubierto por `AGENTS.md` §9):

- **leader** — Identidad, Precondiciones, Arranque (→ Caso D), Antes de delegar
  (comprender/dudar/`feature|task`/`id`/`rei new`/`title`/`description`), casos
  **A–H** completos, Delegación (incl. fallback «solo `## Contrato`»), «Tras cada
  subagente» (`rei validate` + corregir `current.md` ante `WARN`) y Escalado
  presentes en el Contrato; el ejemplo de aprobación, los campos de `meta.json`,
  la aclaración del Caso B, el racional del fallback y la plantilla del prompt
  viven en Referencia. La regla genérica §9 queda en Referencia y en
  `AGENTS.md` §9 (A).
- **spec_author** — los 8 encabezados, casos A/B/C y Bloqueos presentes. Único
  recorte: regla genérica §9 (A).
- **implementer** — los 8 encabezados, casos A/B/C, `#### Cierre (común a Caso A
  y Caso B)` y Bloqueos presentes; `rei check --quiet`, marcado de tareas/pasos
  y la regla `git add` (nunca `.rei/progress/`) se conservan. Ejemplos de
  verificación/nota de `verification.md` y regla genérica §9 en Referencia.
- **reviewer** — los 8 encabezados, casos A/B, `#### Cierre` y Bloqueos;
  `rei review-diff` con fallback a modo lectura, `rei validate`,
  `rei session archive` y `rei commit set … last_review_commit` presentes. Los
  racionales entre paréntesis y la regla genérica §9 en Referencia.
- **initializer** — Identidad, Objetivo, Precondiciones, Alcance de escritura,
  Protocolo 1–5, Plan de pasos, Entrevista guiada, Criterio de finalización y
  Bloqueos presentes. Nota de los marcadores y regla genérica §9 en Referencia.

El diff confirma que los recortes son exactamente los previstos en `plan.md` §3
(mover duplicación/racional/ejemplos y consolidar cierres), sin eliminar reglas,
casos, comandos ni estados. **El flujo y el comportamiento no cambian.**

## Invariante estructural

- Los 5 Contratos contienen las 8 subsecciones obligatorias (`### Identidad`,
  `### Objetivo`, `### Precondiciones`, `### Protocolo`, `### Reglas duras`,
  `### Formato de salida`, `### Herramientas permitidas`,
  `### Documentos de referencia`). `initializer` conserva además su subsección
  previa `### Alcance de escritura`, permitida.
- `## Contrato` precede a `## Referencia` y ambos aparecen exactamente una vez
  por archivo, sin duplicados.
- Ningún Contrato contiene la cadena `Sin detalle adicional`.
- `go test ./...` verde, incluido `TestContratosDeLosRolesSonEstructurados`
  (`internal/adapter/adapter_test.go`). El test lee el esqueleto embebido, que
  `embed.go` apunta directamente a `.rei/agents/**`, por lo que valida los
  archivos editados.

## Nativos sincronizados

`rei init opencode --check` devuelve exit 0 con `[OK]` en los 5 agentes, lo que
confirma que `.opencode/agents/*.md` está sincronizado byte a byte con la
fuente y no fue editado a mano.

## Medición (antes/después)

Contrato total reproducido con el método de `plan.md` §5
(`awk '/^## Contrato/{c=1;next}/^## Referencia/{c=0}c' | wc -c/-w`):

| Rol | bytes antes | bytes después | % red. | palabras antes | palabras después |
|-----|------------:|--------------:|-------:|---------------:|-----------------:|
| `leader.md` | 8704 | 7791 | −10.49 % | 1228 | 1102 |
| `spec_author.md` | 5290 | 5149 | −2.67 % | 702 | 677 |
| `implementer.md` | 5744 | 4991 | −13.11 % | 709 | 610 |
| `reviewer.md` | 5835 | 4917 | −15.73 % | 740 | 626 |
| `initializer.md` | 5699 | 5449 | −4.39 % | 803 | 764 |
| **Total Contrato** | **31272** | **28297** | **−9.51 %** | **4182** | **3779** |

Nativos `.opencode/agents/*.md`: **33070 → 30095 bytes** (−2975, −9.00 %).
Ningún rol creció. La tabla de `impl.md` coincide exactamente con la medición
reproducida por el Reviewer y con el baseline de `plan.md`.

## Restricciones

- No se modificó frontmatter (`name`, `description`, `mode`, `tools`) en ningún
  rol (diff sin cambios en esas líneas).
- No se tocó `.rei/adapters/` ni el código del CLI (`git status` limpio en
  `cmd/`, `internal/`, `.rei/adapters/`).
- No se editó `.opencode/agents/*.md` a mano (confirmado por `--check`).
- `.rei/progress/work-items/<id>/plan.md` solo lleva marcas de pasos; el alcance
  de `meta.json` no cambió.

## Observaciones

- Menor: en `leader.md`, el ítem «Tras cada subagente, ejecuta `rei validate`…»
  quedó como un bullet dentro de la lista «El prompt de delegación debe
  incluir:». Es una ubicación mejorable, pero la instrucción sigue presente en
  el Contrato y es accionable; no altera el comportamiento ni pierde capacidad.
- `package-lock.json` figura como no rastreado (fechado el 10-sep, previo a este
  Work Item) y es ajeno a su alcance.
- `rei check --quiet` pasa trivialmente al no haber checks configurados.

## Resultado

Cumple la planificación: objetivo alcanzado, restricciones respetadas,
sin pérdida de capacidad, invariante estructural intacto, nativos sincronizados,
medición reproducible y todas las verificaciones en verde. Se aprueba.
