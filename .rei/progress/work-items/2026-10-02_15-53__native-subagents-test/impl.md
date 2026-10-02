# Implementación — probar subagentes nativos

> Work Item: `2026-10-02_15-53__native-subagents-test` (task).
> Implementer · estado final: `review`.

## Resumen

Se documentó el uso de los **subagentes nativos de OpenCode** en la guía para
personas `.rei/docs/usage.md`, siguiendo exactamente `plan.md`. Se añadió una
sección nueva, "Subagentes nativos de OpenCode", redactada en el tono del resto
de la guía y sin modificar ninguna sección existente.

La sección cubre los cuatro puntos acordados:

- `rei init opencode` genera un archivo por rol en `.opencode/agents/<rol>.md` a
  partir del `## Contrato` canónico, con la marca `GENERATED` y comportamiento
  idempotente.
- Las dos formas de invocación: `@spec_author`, `@implementer`, `@reviewer` o
  `@initializer`, o bien cambiando al agente `leader`.
- `rei init opencode --check` verifica sin escribir y detecta deriva: código de
  salida `0` si todo coincide, `1` si falta o difiere.
- El fallback: sin adaptador nativo, el Leader transmite al subagente únicamente
  el `## Contrato` del rol, nunca el archivo completo.

No se mencionan otros runtimes ni se amplió el alcance.

## Archivos

### Modificados

| Archivo | Cambio |
|---------|--------|
| `.rei/docs/usage.md` | Nueva sección `# 7. Subagentes nativos de OpenCode` al final del documento. |

## Decisiones de implementación

- **Punto de inserción.** La nueva sección se añadió al final, como Sección 7,
  para respetar la restricción de `plan.md` de **no alterar otras secciones**;
  renumerar secciones existentes habría modificado encabezados fuera del
  alcance. La numeración resultante (1–7) es secuencial y la sección se separa
  con el mismo `---` que el resto.
- **Referencia cruzada.** Se enlaza a la Sección 3 ("El ciclo de trabajo") para
  situar las etapas en las que se invoca cada subagente.

## Pasos del plan

- [x] 1. Punto de inserción y numeración decididos.
- [x] 2. Sección redactada cubriendo generación, invocación, `--check` y fallback.
- [x] 3. Coherencia de estilo, encabezados y referencias cruzadas verificada.
- [x] 4. Sin menciones a otros runtimes ni cambios fuera de alcance.

## Verificación

- `rei check --quiet` → `exit 0`.
- `usage.md`: encabezados numerados 1–7 sin duplicados ni saltos; secciones
  existentes intactas.
- Alcance: el único cambio del Work Item es la sección nueva en
  `.rei/docs/usage.md` (más los artefactos de planificación y progreso). No hay
  referencias a otros runtimes.

## Estado

Trabajo listo para revisión. `status = review`.
