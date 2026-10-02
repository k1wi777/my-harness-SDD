# Task — plan.md

> Las Tasks utilizan una planificación simplificada. El Spec Author genera
> únicamente `plan.md`; no hay `requirements.md`, `design.md` ni `tasks.md`.
>
> Este documento es la referencia del Spec Author cuando el Work Item es de
> `type == task`. Para Features, consulta `specs.md`.

---

# plan.md

Resume el trabajo que debe realizar el Implementer sin generar documentación innecesaria.

Debe contener únicamente la información necesaria para implementar correctamente el cambio.

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

- [ ] 1. ...
- [ ] 2. ...
- [ ] 3. ...
```

El Implementer marca cada paso como completado (`[x]`) en `plan.md` **inmediatamente** al terminarlo, antes de continuar con el siguiente.
