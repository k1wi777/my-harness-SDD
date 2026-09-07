# Progress

> Este documento define cómo se documenta el progreso del trabajo dentro del proyecto.
>
> La carpeta `.rei/progress/` constituye el registro vivo del estado de desarrollo y permite que cualquier agente o persona pueda comprender qué ocurrió durante una sesión, retomarla si fue interrumpida y consultar el historial del proyecto.
>
> Ningún agente debe inventar nuevos formatos. Todos los archivos de `.rei/progress/` deben respetar las plantillas descritas en este documento.
>
> **Este documento es la única fuente de verdad** para plantillas y estructura. Los agentes deben leerlo y aplicarlo directamente — nunca dupliques formatos en otros archivos.

---

# Estructura

```text
.rei/progress/
│
├── current.md              # Sesión activa (único Work Item en curso)
├── history.md              # Bitácora histórica (append-only)
│
└── <work-item>/            # Documentos generados por los agentes
    ├── impl.md             # Reporte del Implementer
    ├── review.md           # Reporte del Reviewer
    └── spec.md             # Bloqueo del Spec Author (solo si aplica)
```

Los archivos de sesión (`current.md`, `history.md`) viven en la raíz de `.rei/progress/`.

Cada Work Item tiene su propia carpeta `.rei/progress/<work-item>/`, análoga a `.rei/specs/<work-item>/`, donde se concentran todos los documentos generados durante su ciclo de vida.

---

# current.md

Representa el estado **actual** de la sesión.

> Esta plantilla también está embebida como heredoc en `.rei/init.sh` (Sección 2).
> Si la modificas aquí, actualiza también el script para mantenerlas sincronizadas.

El **Spec Author** lo inicializa al comenzar la planificación (`pending`) y lo deja en `ready` al terminar, esperando aprobación humana.

El **Implementer** lo actualiza al iniciar la implementación (`in_progress`) y lo mantiene durante toda la ejecución hasta pasar a `review`.

Debe mantenerse actualizado durante todo el ciclo de vida del Work Item — desde `pending` hasta `review`.

No debe rellenarse únicamente al finalizar el trabajo.

Su propósito es permitir que una sesión pueda retomarse en cualquier momento, incluso durante la planificación o la espera de aprobación.

Debe utilizar siempre la siguiente plantilla:

```md
# Sesión actual

> Estado vivo de la sesión.
> Se actualiza durante toda la ejecución.
> Al finalizar el Work Item su resumen se mueve a `history.md` y este archivo vuelve a su estado inicial.

- **Work Item:** _ninguno_
- **Tipo:** _—_
- **Estado:** _—_
- **Inicio:** _—_
- **Agente activo:** _—_

## Plan

_—_

## Bitácora

_—_

## Próximo paso

_—_
```

---

# history.md

Es la bitácora histórica del proyecto.

> Esta plantilla también está embebida como heredoc en `.rei/init.sh` (Sección 2).
> Si la modificas aquí, actualiza también el script para mantenerlas sincronizadas.

Su contenido es **append-only**.

Nunca deben modificarse entradas anteriores.

Al finalizar correctamente un Work Item, el resumen de `current.md` debe añadirse al final de este archivo.

La plantilla base es:

```md
# Bitácora histórica (append-only)

> Registro histórico de todas las sesiones completadas.
> Nunca modifiques entradas anteriores.
> Siempre añade nuevas entradas al final.

---
```

Cada entrada debe resumir:

- fecha;
- Work Item;
- agente;
- trabajo realizado;
- archivos modificados;
- resultado de la verificación;
- estado final.

---

# .rei/progress/<work-item>/impl.md

Documento generado por el **Implementer** en `.rei/progress/<work-item>/impl.md`.

Describe el trabajo realizado durante la implementación.

Su estructura puede adaptarse según la naturaleza del Work Item, pero siempre debe incluir como mínimo:

- resumen de la implementación;
- archivos modificados;
- cambios realizados;
- proceso de verificación;
- observaciones relevantes.

Este archivo constituye la evidencia principal para la revisión.

---

# .rei/progress/<work-item>/review.md

Documento generado por el **Reviewer** en `.rei/progress/<work-item>/review.md`.

Resume el resultado de la revisión realizada sobre el Work Item.

Debe incluir como mínimo:

- estado final (`done`, `changes_requested` o `blocked`);
- verificaciones realizadas;
- observaciones;
- acciones requeridas si la revisión fue rechazada.

---

# .rei/progress/<work-item>/spec.md

Documento generado únicamente cuando el **Spec Author** no puede completar correctamente la planificación.

Debe explicar claramente:

- motivo del bloqueo;
- información faltante;
- decisiones que requieren intervención humana.

Si no existen bloqueos, este archivo no debe crearse.

---

# Responsabilidades

| Archivo | Responsable |
|----------|-------------|
| `current.md` (inicialización: `pending` → `ready`) | Spec Author |
| `current.md` (implementación: `in_progress` → `review`) | Implementer |
| `history.md` | Reviewer |
| `<work-item>/impl.md` | Implementer |
| `<work-item>/review.md` | Reviewer |
| `<work-item>/spec.md` | Spec Author |

---

# Reglas

- Mantén `current.md` actualizado durante todo el ciclo de vida del Work Item.
- Crea `.rei/progress/<work-item>/` al iniciar el trabajo sobre un Work Item.
- Nunca sobrescribas el historial.
- No elimines documentación existente.
- Utiliza siempre las plantillas definidas en este documento.
- Todo Work Item debe dejar evidencia suficiente para poder comprender qué ocurrió sin depender del historial del chat.
