# impl.md — slimming de contratos de rol

Work Item: `2026-10-02_16-12__contrato-slimming` (type: `task`)
Plan: `.rei/specs/2026-10-02_16-12__contrato-slimming/plan.md`

## Resumen

Se recortó el `## Contrato` de los 5 roles canónicos moviendo a `## Referencia`
el detalle no esencial (ejemplos, racional, aclaraciones no operativas,
duplicaciones) y consolidando cierres idénticos, **sin cambiar el flujo ni el
comportamiento** de ningún rol. Los agentes nativos de OpenCode se regeneraron
con `rei init opencode` y quedaron sincronizados (`--check` exit 0).

Resultado: Contrato **31272 → 28297 bytes** (−2975, **−9.51 %**) y
`.opencode/agents/*.md` **33070 → 30095 bytes** (−2975, **−9.00 %**). Ningún rol
creció.

## Cambios aplicados

- **leader.md**: ejemplos/plantillas y campos de `meta.json` a Referencia;
  Caso A sin bloque de ejemplo (instrucción operativa de aprobación en su
  lugar); aclaración del Caso B a Referencia; párrafo del fallback sin adaptador
  consolidado con la regla dura; lista operativa de qué debe incluir el prompt
  de delegación. Regla genérica de AGENTS §9 retirada del Contrato.
- **spec_author.md**: regla genérica de AGENTS §9 movida a Referencia.
- **implementer.md**: cierre A/B consolidado en `#### Cierre (común a Caso A y
  Caso B)`; ejemplos de verificación y nota de `verification.md` a Referencia
  (el Contrato conserva `rei check --quiet` y cuándo leer `verification.md`);
  regla genérica a Referencia.
- **reviewer.md**: cierre A/B consolidado en `#### Cierre (común a Caso A y Caso
  B)`; racionales entre paréntesis (`rei validate`, `rei session archive`,
  `last_review_commit`, trazabilidad) a Referencia; regla genérica a Referencia.
- **initializer.md**: nota de los dos marcadores de `AGENTS.md` a Referencia;
  regla genérica a Referencia.

En todos: los 8 encabezados `###` obligatorios permanecen en el Contrato y
`## Contrato` sigue precediendo a `## Referencia` sin duplicados.

## Paso 1 — Baseline (antes)

Método de `plan.md` §5 (texto entre el primer `## Contrato` y el primer
`## Referencia`) y `wc -c .opencode/agents/*.md`. Baseline tomado de
`git show HEAD:<path>`.

| Rol | Contrato bytes | Contrato palabras |
|-----|---------------:|------------------:|
| `leader.md` | 8704 | 1228 |
| `spec_author.md` | 5290 | 702 |
| `implementer.md` | 5744 | 709 |
| `reviewer.md` | 5835 | 740 |
| `initializer.md` | 5699 | 803 |
| **Total** | **31272** | **4182** |

Nativos (`.opencode/agents/*.md`): **33070 bytes**.

## Paso 9 — Medición después

| Rol | bytes antes | bytes después | Δ bytes | % reducción | palabras antes | palabras después |
|-----|------------:|--------------:|--------:|------------:|---------------:|-----------------:|
| `leader.md` | 8704 | 7791 | −913 | −10.49 % | 1228 | 1102 |
| `spec_author.md` | 5290 | 5149 | −141 | −2.67 % | 702 | 677 |
| `implementer.md` | 5744 | 4991 | −753 | −13.11 % | 709 | 610 |
| `reviewer.md` | 5835 | 4917 | −918 | −15.73 % | 740 | 626 |
| `initializer.md` | 5699 | 5449 | −250 | −4.39 % | 803 | 764 |
| **Total Contrato** | **31272** | **28297** | **−2975** | **−9.51 %** | **4182** | **3779** |

Ahorro en palabras del Contrato: **4182 → 3779** (−403, −9.64 %).

Agentes nativos regenerados:

| Rol (`.opencode/agents/`) | antes bytes | después bytes | Δ bytes | % reducción |
|---------------------------|------------:|--------------:|--------:|------------:|
| `leader.md` | 9071 | 8158 | −913 | −10.06 % |
| `spec_author.md` | 5589 | 5448 | −141 | −2.52 % |
| `implementer.md` | 6102 | 5349 | −753 | −12.34 % |
| `reviewer.md` | 6170 | 5252 | −918 | −14.88 % |
| `initializer.md` | 6138 | 5888 | −250 | −4.07 % |
| **Total nativos** | **33070** | **30095** | **−2975** | **−9.00 %** |

## Paso 7 — Checklist de cobertura

Verificado que cada ítem de `plan.md` §4 sigue disponible. Resultado por rol:

- **leader**: los 8 encabezados `###` presentes; casos A–H presentes;
  `#### Delegación` y `#### Escalado` presentes; `rei new`, `rei commit set`,
  `rei validate` presentes. Regla genérica de §9 cubierta por `AGENTS.md` (A).
- **spec_author**: los 8 encabezados presentes; casos A, B, C y Bloqueos
  presentes. Regla genérica cubierta por `AGENTS.md` (A).
- **implementer**: los 8 encabezados presentes; casos A, B, C, `#### Cierre` y
  `#### Bloqueos` presentes; `git add` (sin `.rei/progress/`) presente. Regla
  genérica cubierta por `AGENTS.md` (A).
- **reviewer**: los 8 encabezados presentes; casos A, B, `#### Cierre` y
  `#### Bloqueos` presentes; `rei review-diff` con fallback al modo lectura
  presente. Regla genérica cubierta por `AGENTS.md` (A).
- **initializer**: los 8 encabezados presentes; plan de pasos, entrevista
  guiada, criterio de finalización y bloqueos presentes. Regla genérica cubierta
  por `AGENTS.md` (A).

Comprobación mecánica: cada Contrato contiene los 8 encabezados `###` y
**ninguno** contiene la cadena `Sin detalle adicional` (cuerpo de la Referencia
por defecto).

## Paso 8 — Regeneración y verificación

- `rei init opencode` → `[UPD]` de los 5 `.opencode/agents/*.md`, exit 0.
- `rei init opencode --check` → `[OK]` los 5, exit 0 (sincronía byte a byte).

## Paso 10 — Validaciones

- `rei check --quiet` → exit 0.
- `go test ./...` → todos los paquetes `ok`; en particular
  `TestContratosDeLosRolesSonEstructurados`, `TestGenerateAll`, `TestCheck`.
- `rei validate 2026-10-02_16-12__contrato-slimming` → `Resultado: OK`, exit 0.

## Archivos tocados

- `.rei/agents/{leader,spec_author,implementer,reviewer,initializer}.md`
- `.opencode/agents/*.md` (regenerados por el CLI, no editados a mano)
- `.rei/specs/2026-10-02_16-12__contrato-slimming/plan.md` (marcas de pasos)
- `.rei/progress/current.md` (estado de sesión)

No se tocó frontmatter, `.rei/adapters/` ni código del CLI.

## Incidencias

Ninguna. La implementación se completó sin bloqueos.
