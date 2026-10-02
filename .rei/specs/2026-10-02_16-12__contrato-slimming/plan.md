# plan.md — slimming de contratos de rol

## Objetivo

Reducir el tamaño del `## Contrato` de los 5 roles canónicos
(`.rei/agents/{leader,spec_author,implementer,reviewer,initializer}.md`) para que
cada spawn de subagente consuma menos tokens, **sin perder capacidad operativa**.
El recorte se logra moviendo a `## Referencia` el detalle no esencial
(ejemplos, racional, aclaraciones no operativas y duplicaciones) y dejando en el
Contrato solo lo necesario para ejecutar el rol. El Contrato debe seguir siendo
autosuficiente y regenerar los agentes nativos de OpenCode sin diferencias.

Resultado esperado: menor tamaño de `.rei/agents/*.md` (sección Contrato) y, por
ende, de `.opencode/agents/*.md`, con `rei init opencode --check` en exit 0.

---

## Archivos

- `.rei/agents/leader.md`
- `.rei/agents/spec_author.md`
- `.rei/agents/implementer.md`
- `.rei/agents/reviewer.md`
- `.rei/agents/initializer.md`
- `.opencode/agents/*.md` — regenerados por `rei init opencode` (no se editan a mano).
- `.rei/progress/work-items/2026-10-02_16-12__contrato-slimming/impl.md` — informe de medición y evidencia.
- `.rei/progress/current.md` y `.rei/specs/2026-10-02_16-12__contrato-slimming/meta.json` — estado de sesión.

No se prevé tocar frontmatter (`name`, `description`, `mode`, `tools`), ni
`.rei/adapters/`, ni el código del CLI.

---

## Cambios

### 1. Criterio de recorte

Regla de oro (de `.rei/docs/harness/adapters.md`): **si una regla o un paso es
necesario para ejecutar el rol, permanece en `## Contrato`; si es justificación,
ejemplo o duplicación, se mueve a `## Referencia`.**

Se admiten únicamente estas técnicas, en este orden de preferencia:

1. **Mover a `## Referencia`**: ejemplos reproducibles, plantillas, racional,
   aclaraciones no operativas y párrafos que repiten una regla ya presente en
   el Contrato o en `AGENTS.md` §9.
2. **Consolidar duplicación interna** dentro del Contrato: cuando dos casos
   comparten literalmente una misma secuencia de pasos, extraerla una sola vez
   (p. ej. un «Cierre común») y referenciarla. Esto reduce texto **sin cambiar
   pasos ni comportamiento**.
3. **Comprimir redacción** sin eliminar contenido normativo (frases más cortas,
   misma instrucción).

No se permite eliminar reglas operativas, casos, comandos ni estados.

### 2. Invariante: qué permanece en el Contrato

En los 5 roles deben permanecer, intactas y operativas:

- Las 8 subsecciones obligatorias requeridas por el test estructural:
  `### Identidad`, `### Objetivo`, `### Precondiciones`, `### Protocolo`,
  `### Reglas duras`, `### Formato de salida`, `### Herramientas permitidas`,
  `### Documentos de referencia`.
- Identidad y objetivo del rol.
- Precondiciones y condiciones de parada (DETENTE).
- El protocolo completo, paso a paso, con **todos** sus casos.
- Las reglas duras (NUNCA/SIEMPRE) operativas.
- El formato de salida exacto (cadenas `estado -> ruta`).
- La lista de herramientas permitidas.
- Los punteros a documentación (Documentos de referencia).

La sección `## Contrato` debe seguir precediendo a `## Referencia`, sin
duplicados de esos encabezados. La Referencia puede contener subsecciones `###`
propias sin romper `ExtractContract`.

### 3. Qué se mueve a Referencia (por rol)

Se listan los candidatos concretos por rol. El Implementer redacta la Referencia
correspondiente y deja en el Contrato una instrucción breve y accionable con
puntero cuando sea útil.

**leader.md**
- El ejemplo del mensaje de aprobación («Tu mensaje deberá ser similar a: …»):
  mover el bloque textual. En el Contrato queda la instrucción operativa
  «solicita una aprobación explícita indicando la ruta del spec» (ajustar
  `### Formato de salida` para no depender del ejemplo).
- La enumeración descriptiva de campos de `meta.json` (los rellena `rei new`):
  mover el detalle. En el Contrato queda «solo editas `title`/`description`;
  registras `base_commit` con `rei commit set`».
- La aclaración del Caso B («el estado ya lo reportaste… reconfirmarlo»): mover.
- El párrafo extenso «**Fallback sin adaptador nativo.**» de la sección
  Delegación: consolidar con la regla dura equivalente («transmite solo el
  `## Contrato` del rol; nunca el archivo completo») y mover el racional.
- El bloque plantilla del prompt de delegación (fenced `text`): mover a
  Referencia. En el Contrato queda la lista operativa de lo que debe incluir el
  prompt (objetivo, ruta del Work Item, tipo, acuerdo clave, restricciones).

**spec_author.md**
- La regla genérica «No releas un archivo que ya esté en tu contexto…»: mover a
  Referencia (duplicada en `AGENTS.md` §9).
- Sin más recortes relevantes: el resto del Contrato es operativo.

**implementer.md**
- Consolidar el cierre idéntico de Caso A y Caso B (pasos «Al finalizar»,
  actualización de `current.md` a `review`, cambio de `status` y DETENTE) en una
  única subsección de cierre compartida.
- Mover a Referencia los ejemplos/alternativas parentéticas de verificación
  («por ejemplo … o solo los tests del módulo afectado si el proyecto lo
  permite»); en el Contrato queda `rei check --quiet`.
- Mover a Referencia la nota explicativa sobre `verification.md` («NO leas …
  por defecto …»), conservando en el Contrato la instrucción de cuándo leerla.
- La regla genérica «No releas un archivo que ya esté en tu contexto…»: mover a
  Referencia (duplicada en `AGENTS.md` §9).

**reviewer.md**
- Consolidar el cierre idéntico de Caso A y Caso B (pasos de aprobación:
  `rei validate`, `rei session archive`, `status = done`; y pasos de rechazo:
  `status = changes_requested`, `rei commit set …`, actualización de
  `current.md`) en una única subsección de cierre compartida.
- Mover a Referencia los racionales entre paréntesis («(así la próxima revisión
  verá solo los cambios pedidos)», «(archiva el resumen en `history.md`…)»,
  «(consistencia interna)», «(trazabilidad tarea ↔ `R-id`)»), conservando los
  comandos y pasos.
- La regla genérica «No releas un archivo que ya esté en tu contexto…»: mover a
  Referencia (duplicada en `AGENTS.md` §9).

**initializer.md**
- Mover a Referencia la nota de la sección «Plan de pasos» sobre los dos
  marcadores de `AGENTS.md`, por estar ya cubierta por las Reglas duras.
- Mover a Referencia cualquier racional de la «Entrevista guiada» que no agregue
  una instrucción ejecutable; conservar las reglas operativas de la entrevista.
- La regla genérica «No releas un archivo que ya esté en tu contexto…»: mover a
  Referencia (duplicada en `AGENTS.md` §9).

### 4. Checklist de cobertura por rol

Verificación de que cada regla/paso necesario sigue disponible. Marcar cada
ítem; ubicación esperada: **C** = Contrato, **R** = Referencia (no esencial),
**A** = cubierto por `AGENTS.md` §9.

- [x] **leader**
  - [x] Identidad: recibe solicitudes, selecciona workflow, coordina; no implementa (C)
  - [x] Precondiciones: no repetir `rei check`; usar `rei session`/`rei items status`, no escanear `meta.json` (C)
  - [x] Arranque: consultar sesión e ítems; `in_progress` → Caso D (C)
  - [x] Antes de delegar: comprender solicitud, dudar→preguntar, decidir feature/task, generar `id`, `rei new`, completar `title`/`description` (C)
  - [x] Caso A: lanzar `spec_author` y pedir aprobación (C; ejemplo del mensaje en R)
  - [x] Caso B: sin aprobación no continuar; con cambios relanzar `spec_author` (C)
  - [x] Caso C: `status=in_progress`, `rei commit set base_commit`, lanzar `implementer` (C)
  - [x] Caso D: interrupción → preguntar al usuario (C)
  - [x] Caso E: lanzar `reviewer` (C)
  - [x] Caso F: leer `review.md`, `status=in_progress`, lanzar `implementer` (C)
  - [x] Caso G: `done` → no continuar (C)
  - [x] Caso H: `blocked` → leer motivo, informar, no continuar (C)
  - [x] Delegación: transmitir solo lo relevante; fallback sin adaptador usa solo el Contrato (C; racional en R)
  - [x] Tras cada subagente: `rei validate` y corregir `current.md` ante `WARN` (C)
  - [x] Escalado: dividir solicitudes grandes (C)
  - [x] Reglas duras, formato de salida, herramientas, documentos de referencia (C)
  - [x] Regla «no releas archivo ya en contexto» (A)
- [x] **spec_author**
  - [x] Identidad/objetivo: planificar sin implementar (C)
  - [x] Precondiciones y parada ante `status` inválido (C)
  - [x] Protocolo 1–6: leer `specs.md`/`task.md`, `meta.json`, verificar `status`, iniciar sesión solo si `pending`, consultar docs del proyecto solo si aplica (C)
  - [x] Caso A (feature): `requirements.md` con R-id/EARS/verificable, `design.md`, `tasks.md` con trazabilidad y checkboxes, dejar `ready`, DETENTE (C)
  - [x] Caso B (task): `plan.md` simple, dejar `ready`, DETENTE (C)
  - [x] Caso C (revisión): no reiniciar sesión; aplicar cambios; dejar `ready` (C)
  - [x] Bloqueos: documentar en `spec.md`, `status=blocked`, DETENTE (C)
  - [x] Reglas duras, formato de salida, herramientas, documentos de referencia (C)
  - [x] Regla «no releas archivo ya en contexto» (A)
- [x] **implementer**
  - [x] Identidad/objetivo y precondiciones (`in_progress`, `current.md`, spec/plan presentes) (C)
  - [x] Protocolo 1–6: leer docs, `meta.json`, `current.md`, `type`, detectar rework → Caso C (C)
  - [x] Caso A (feature): implementar tareas en orden; marcar `[x]` y actualizar `current.md` tras cada una; verificar; cierre (C)
  - [x] Caso B (task): implementar pasos de `plan.md`; marcar y actualizar tras cada uno; cierre (C)
  - [x] Caso C (rework): aplicar `review.md`, desmarcar/remarcar tareas, continuar flujo (C)
  - [x] Cierre: `rei check --quiet`, `impl.md`, `current.md` a `review`, `status=review`, DETENTE (C)
  - [x] Bloqueos: documentar en `impl.md`, `status=blocked`, DETENTE (C)
  - [x] Reglas duras (incl. `git add` sin `.rei/progress/`), formato de salida, herramientas, documentos de referencia (C)
  - [x] Regla «no releas archivo ya en contexto» (A)
- [x] **reviewer**
  - [x] Identidad/objetivo y precondición `status == review` (C)
  - [x] Protocolo: leer `verification.md` y `meta.json`, `type`, `rei review-diff` con fallback a modo lectura si exit 2 (C)
  - [x] Caso A (feature): leer `tasks.md`/`impl.md`/diff, comprobar tareas, R-ids, arquitectura, convenciones, checkpoints, `rei check --quiet` (C)
  - [x] Caso B (task): leer `plan.md`/`impl.md`/diff, comprobar objetivo, restricciones, arquitectura, convenciones, checkpoints, `rei check --quiet` (C)
  - [x] Cierre aprobación: `review.md`, `rei validate`, `rei session archive`, `status=done`, DETENTE (C)
  - [x] Cierre rechazo: `status=changes_requested`, `rei commit set last_review_commit`, `current.md`, `review.md`, DETENTE (C)
  - [x] Bloqueos: `status=blocked`, documentar, DETENTE (C)
  - [x] Reglas duras, formato de salida, herramientas, documentos de referencia (C)
  - [x] Regla «no releas archivo ya en contexto» (A)
- [x] **initializer**
  - [x] Identidad/objetivo: entrevista guiada, sin implementar (C)
  - [x] Precondiciones: `rei init status == 1`; si `0`, informar y devolver control (C)
  - [x] Alcance de escritura (AGENTS §2/§3, docs de proyecto, `config.json`) (C)
  - [x] Protocolo 1–5: estado, reanudabilidad, inferir/entrevistar/redactar/escribir/confirmar, alinear `config.json`, cierre con `rei init status`/`rei doctor` (C)
  - [x] Plan de pasos (tabla) y regla de los dos marcadores de `AGENTS.md` (C; nota explicativa en R)
  - [x] Entrevista guiada operativa y criterio de finalización (C)
  - [x] Bloqueos: documentar en `current.md`, DETENTE (C)
  - [x] Reglas duras, formato de salida, herramientas, documentos de referencia (C)
  - [x] Regla «no releas archivo ya en contexto» (A)

### 5. Medición del ahorro

**Baseline (Contrato, medido antes de editar):**

| Rol | bytes | palabras |
|-----|------:|---------:|
| `.rei/agents/leader.md` | 8704 | 1228 |
| `.rei/agents/spec_author.md` | 5290 | 702 |
| `.rei/agents/implementer.md` | 5744 | 709 |
| `.rei/agents/reviewer.md` | 5835 | 740 |
| `.rei/agents/initializer.md` | 5699 | 803 |
| **Total Contrato** | **31272** | **4182** |

Baseline de agentes nativos: **33070 bytes** (`wc -c .opencode/agents/*.md`).

**Método reproducible** (idéntico antes y después; el Contrato es el texto entre
el primer `## Contrato` y el primer `## Referencia`):

```bash
for f in .rei/agents/*.md; do
  echo "== $f =="
  awk '/^## Contrato/{c=1;next}/^## Referencia/{c=0}c' "$f" | wc -c
  awk '/^## Contrato/{c=1;next}/^## Referencia/{c=0}c' "$f" | wc -w
done
wc -c .opencode/agents/*.md
```

- Capturar el baseline al **inicio** de la implementación (o usar la tabla de
  arriba) antes de aplicar cualquier edición.
- Recalcular tras editar y regenerar.
- Reportar en `impl.md` una tabla por rol con `bytes` y `palabras` antes/después
  y el **% de reducción**, más el total del Contrato y el total de
  `.opencode/agents/*.md`. El ahorro real es la diferencia total; ninguna rol
  debería crecer en Contrato.

### 6. Regeneración y verificación

- Regenerar agentes nativos: `rei init opencode` (debe terminar con exit 0 y
  actualizar `.opencode/agents/*.md`).
- Confirmar sincronía: `rei init opencode --check` debe salir con **exit 0**.
- Confirmar estructura del esqueleto: `rei check --quiet` exit 0.
- Confirmar test estructural: `go test ./...` en verde, en particular
  `TestContratosDeLosRolesSonEstructurados` (exige las 8 subsecciones `###` en
  cada Contrato y prohíbe filtrar «Sin detalle adicional») y
  `TestGenerateAll`/`TestCheck` (los nativos contienen `## Contrato` y
  `### Identidad`).
- Verificar `rei validate 2026-10-02_16-12__contrato-slimming`.

---

## Restricciones

- **No cambiar el flujo ni el comportamiento de los roles**: mismos pasos,
  mismos comandos, mismos estados y mismas condiciones de parada. Solo cambia
  dónde vive el detalle (Contrato vs Referencia) y la redacción.
- **No cambiar el alcance de `meta.json`.**
- No tocar frontmatter de los roles (salvo necesidad justificada; `tools` y
  `mode` son inmutables en la práctica).
- Mantener el orden y la unicidad de los encabezados `## Contrato` y
  `## Referencia`; no introducir esas cadenas exactas dentro de otro contenido.
- No editar `.opencode/agents/*.md` a mano: solo regenerarlos con el CLI.
- No modificar `.rei/adapters/` ni código del CLI.
- La Referencia solo contiene detalle opcional; su omisión no debe impedir
  ejecutar el rol.

---

## Pasos

- [x] 1. Capturar y registrar el baseline de la sección Contrato (bytes/palabras) y de `.opencode/agents/*.md`.
- [x] 2. Recortar `.rei/agents/leader.md` según §3 (mover ejemplos, lista de campos de `meta.json`, aclaración Caso B, racional del fallback y plantilla de delegación a `## Referencia`).
- [x] 3. Recortar `.rei/agents/spec_author.md` según §3 (mover la regla genérica «no releas…»).
- [x] 4. Recortar `.rei/agents/implementer.md` según §3 (consolidar cierre A/B; mover ejemplos/nota y regla genérica).
- [x] 5. Recortar `.rei/agents/reviewer.md` según §3 (consolidar cierre A/B; mover racionales y regla genérica).
- [x] 6. Recortar `.rei/agents/initializer.md` según §3 (mover nota de `AGENTS.md` y racional; mover regla genérica).
- [x] 7. Verificar el **checklist de cobertura por rol** (§4): cada ítem marcado C/R/A según lo previsto; ninguna regla operativa se perdió.
- [x] 8. Regenerar nativos con `rei init opencode` y confirmar `rei init opencode --check` (exit 0).
- [x] 9. Medir el después (mismo método de §5) y escribir en `impl.md` la tabla antes/después por rol + totales + % de ahorro.
- [x] 10. Ejecutar `rei check --quiet`, `go test ./...` y `rei validate 2026-10-02_16-12__contrato-slimming`; documentar resultados en `impl.md`.
- [x] 11. Actualizar `current.md` (Estado `review`, Próximo paso «esperar revisión») y cambiar `status` a `review`.
