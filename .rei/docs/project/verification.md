# Verification

> **Documento del proyecto** — personaliza este archivo al adaptar REI Harness a tu repositorio.
>
> Este documento define cómo se verifica que un Work Item ha sido implementado correctamente.
>
> El objetivo no es indicar **qué** se implementó, sino **cómo demostrar que funciona**.

## Estrategia de verificación

_Describe aquí cómo debe verificarse el proyecto._

Ejemplos:

- Tests automatizados.
- Build del proyecto.
- Linter.
- Type checking.
- Validación manual.
- Smoke tests.
- Otro mecanismo propio del proyecto.

---

## Checkpoints de verificación

_Indica los comandos que normalmente deben ejecutarse antes de considerar un Work Item listo para revisión._

Cada checkpoint posee un identificador estable (`V1`, `V2`, ...). El Implementer y el Reviewer deben referenciar estos IDs en sus reportes en vez de describir la verificación de forma libre.

| ID | Comando | Descripción |
|----|---------|-------------|
| `V1` | `<comando 1>` | <qué valida> |
| `V2` | `<comando 2>` | <qué valida> |
| `V3` | `<comando 3>` | <qué valida> |

> Estos mismos comandos deben estar cableados en `.rei/init.sh` (Sección 4, `run_check`) para que su resultado se refleje en el código de salida del script. Si añades o modificas un checkpoint aquí, actualiza también `.rei/init.sh`.

---

## Evidencia

_Describe qué evidencia debe dejar el Implementer para demostrar que el trabajo fue verificado._

Como mínimo, por cada checkpoint ejecutado:

- ID del checkpoint (`V1`, `V2`, ...);
- resultado (pasa / falla);
- comando ejecutado;
- observaciones relevantes (logs, capturas, validaciones manuales) cuando el checkpoint no sea automatizable.

La evidencia debe documentarse en:

```
.rei/progress/work-items/<work-item-id>/impl.md
```

---

## Bloqueos

Un Work Item pasa a `blocked` cuando alguna verificación falla y no puede resolverse dentro del alcance de la sesión.

El protocolo exacto para declarar un bloqueo (qué documentar, dónde, y cómo detenerse) está definido en el archivo de rol correspondiente (`.rei/agents/implementer.md`, `.rei/agents/reviewer.md`).

Nunca continúes la implementación ni solicites revisión ignorando una verificación fallida.

---

## Qué NO hacer

- No asumir que un cambio funciona sin verificarlo contra los checkpoints definidos.
- No omitir checkpoints definidos para el proyecto.
- No solicitar revisión sin evidencia suficiente por cada checkpoint.
- No describir una verificación de forma libre si ya existe un ID (`V1`, `V2`, ...) para ella.
