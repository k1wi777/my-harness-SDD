# Revisión — corregir hallazgos del dogfooding

Work Item: `2026-10-02_16-07__dogfood-fixes` (task)
Plan: `.rei/specs/2026-10-02_16-07__dogfood-fixes/plan.md`
Estado final: **done**

## Alcance revisado

Tipo `task`. Se revisan los dos objetivos declarados en el plan:

1. **Filtro de `rei review-diff`.** Debe excluir solo `.rei/progress/**` y
   conservar código, `AGENTS.md`, la spec del Work Item y el resto del harness
   (`.rei/docs`, `.rei/agents`, `.rei/adapters`, `.rei/templates`,
   `.rei/config.json`).
2. **Handoff coherente del Leader.** El `## Contrato` de `.rei/agents/leader.md`
   debe documentar que, ante el `WARN` de incoherencia `current.md`↔`meta.json`
   tras `rei validate`, el Leader corrige `current.md`; y
   `.opencode/agents/leader.md` debe quedar regenerado y coherente.

## Verificaciones

| Check | Comando | Resultado |
|-------|---------|-----------|
| Paquete de revisión | `rei review-diff 2026-10-02_16-07__dogfood-fixes` | exit 0 (ver abajo) |
| Tests | `go test ./...` | exit 0, todos los paquetes en verde |
| Vet | `make vet` | exit 0 |
| Harness | `rei check --quiet` | exit 0 |
| Adaptador nativo | `rei init opencode --check` | exit 0 (`[OK]` en los 5 agentes) |

No hay checkpoints `V1`/`V2` declarados en `.rei/docs/project/verification.md`
(documento aún plantilla) ni en `.rei/config.json` (`checks: []`), por lo que no
aplica ningún checkpoint rápido o lento adicional.

## Objective 1 — filtro de `rei review-diff`

El paquete (`exit 0`, base `e635e8f…`) lista:

- `.rei/agents/leader.md` (antes oculto) — **incluido**;
- `internal/gitx/gitx.go` y `internal/gitx/gitx_test.go` — **incluidos**;
- `.opencode/agents/leader.md` (generado) — incluido;
- `.rei/specs/2026-10-02_16-07__dogfood-fixes/{meta.json,plan.md}` — incluidos.

No aparece ninguna ruta bajo `.rei/progress/**`, que sigue excluida.

El diff de `internal/gitx/gitx.go` confirma el criterio único:

- `filterPaths` conserva toda ruta salvo las que empiezan por `.rei/progress/`;
  se eliminaron el parámetro `id` y `filterCode`.
- El resumen usa la lista filtrada (rótulo `--- Resumen ---`), por lo que
  `.rei/docs`/`.rei/agents` dejan de quedar ocultos.
- El salto de no rastreados en `--full` omite solo `.rei/progress/**`.
- No se alteraron base/`last_review_commit`, formato general ni exit codes.

`internal/gitx/gitx_test.go` cubre la conservación de `.rei/docs/`, `.rei/agents/`,
`.rei/adapters/`, `.rei/templates/`, `.rei/config.json`, specs de otros Work
Items y `AGENTS.md`, y el descarte de `.rei/progress/**`; `TestFilterCode` fue
eliminado junto con el filtro.

## Objective 2 — handoff coherente del Leader

`.rei/agents/leader.md` (`## Contrato`) añade, en **Delegación**, la corrección
del campo **Estado** de `.rei/progress/current.md` cuando `rei validate` emite el
`WARN` de incoherencia con `meta.json`; y refleja la misma obligación en
**Reglas duras** (no avanzar mientras estén incoherentes).

`.opencode/agents/leader.md` contiene la misma redacción y conserva la marca
`GENERATED`; `rei init opencode --check` devuelve exit 0, confirmando que el
adaptador está sincronizado con la fuente. No hubo edición manual del archivo
generado, en línea con la restricción del plan.

## Observaciones

- La implementación no se desvía de la planificación. El `.opencode/agents/leader.md`
  aparece en el diff por regeneración con el CLI (mecanismo previsto), no por
  edición manual.
- `package-lock.json` figura como no rastreado y es ajeno al alcance de este Work
  Item; no afecta la revisión.
- `rei check --quiet` pasa trivialmente porque no hay checkpoints configurados.

## Resultado

Cumple la planificación: objetivos 1 y 2 verificados, restricciones respetadas y
verificaciones en verde. Se aprueba.
