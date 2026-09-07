# Spec Driven Development

> Este repositorio utiliza Spec Driven Development (SDD) para planificar el trabajo antes de escribir código.

El objetivo del proceso no es producir documentación extensa, sino proporcionar al Implementer el contexto necesario para realizar un cambio correctamente, minimizando ambigüedades y evitando decisiones improvisadas durante la implementación.

---

# Objetivo

Todo Work Item debe ser planificado antes de implementarse.

La planificación siempre es responsabilidad del **Spec Author** y debe ser aprobada explícitamente por el usuario antes de comenzar la implementación.

El tipo de planificación depende del tipo de Work Item.

---

# Estructura

Cada Work Item dispone de una carpeta propia dentro de `.rei/specs/`. El nombre de la carpeta es su `id` inmutable.

```
.rei/specs/
└── <work-item-id>/
    ├── meta.json
    ├── requirements.md   (solo Features)
    ├── design.md         (solo Features)
    ├── tasks.md          (solo Features)
    └── plan.md           (solo Tasks)
```

---

# Tipos de planificación

## Feature

Las Features requieren una planificación completa.

El Spec Author debe generar:

- `requirements.md`
- `design.md`
- `tasks.md`

Cada documento cumple un propósito distinto.

---

### requirements.md

Describe **qué debe hacer** el sistema.
No contiene decisiones técnicas.
Los requisitos deben redactarse utilizando **EARS (Easy Approach to Requirements Syntax)**.

Reglas:

- [ ] Cada requisito tiene un identificador estable (`R1`, `R2`, ...)
- [ ] Cada requisito representa una única responsabilidad
- [ ] Cada requisito es verificable
- [ ] Ningún requisito describe decisiones de implementación

---

### EARS

Los requisitos pueden escribirse utilizando cualquiera de los siguientes patrones.

| Tipo | Plantilla |
|------|-----------|
| Ubicuo | El sistema DEBE... |
| Evento | CUANDO... el sistema DEBE... |
| Estado | MIENTRAS... el sistema DEBE... |
| Opcional | DONDE... el sistema DEBE... |
| Error | SI... ENTONCES el sistema DEBE... |

---

### design.md

Describe **cómo se implementará** la Feature.

Debe documentar únicamente las decisiones necesarias para implementar el cambio.

Puede incluir:

- archivos afectados;
- componentes implicados;
- decisiones técnicas;
- restricciones;
- alternativas descartadas.

No debe convertirse en documentación de arquitectura general.

---

### tasks.md

Divide la implementación en pasos concretos y ordenados.

Cada tarea representa una unidad de trabajo que el Implementer ejecutará de forma secuencial.

Formato recomendado:

```md
- [ ] T1 — ...
- [ ] T2 — ...
- [ ] T3 — ...
```

El Implementer marcará cada tarea como completada (`[x]`) en `tasks.md` **inmediatamente** al terminarla, antes de continuar con la siguiente.

---

## Task

Las Tasks utilizan una planificación simplificada.

El Spec Author únicamente genera `plan.md`.

---

### plan.md

Resume el trabajo que debe realizar el Implementer sin generar documentación innecesaria.

Formato recomendado:

```text
Objetivo

...

Archivos

...

Cambios

...

Restricciones

...

Pasos

1.
2.
3.
```

Debe contener únicamente la información necesaria para implementar correctamente el cambio.

---

# Estados

Todos los Work Items utilizan los mismos estados.

Los estados válidos y sus transiciones se definen en `workflow.md`.
---

# Aprobación humana

La planificación nunca se implementa automáticamente.

Cuando el Spec Author finaliza su trabajo, el Leader debe detener el workflow y esperar una aprobación explícita del usuario.

Solo después de esa aprobación el Work Item puede pasar a `in_progress`.

---

# Reglas generales

- Todo Work Item comienza en estado `pending`.
- Todo Work Item debe pasar por el Spec Author.
- Ningún Work Item puede implementarse sin planificación aprobada.
- El Implementer nunca modifica la planificación, salvo marcar tareas o pasos como completados en `tasks.md` o `plan.md`.
- El Reviewer verifica la implementación contra la planificación aprobada.
