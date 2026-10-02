# Requirements — rei init wizard

> Work Item: `2026-10-02_12-18__rei-init-wizard` (feature).
>
> Los requisitos describen **qué** debe hacer el sistema, no **cómo**. Se redactan
> en EARS. El sujeto "el sistema" designa al CLI `rei` o al rol `initializer`,
> según el requisito.
>
> **Marcador canónico:** `<!-- REI:PENDIENTE -->`.
>
> **Documentos de personalización del proyecto:** `AGENTS.md` (secciones 2 y 3),
> `.rei/docs/project/architecture.md`, `.rei/docs/project/conventions.md` y
> `.rei/docs/project/verification.md`.

---

# 1. `rei init` — scaffold, detección y plan

**R1.** El sistema DEBE ofrecer el comando `rei init` sin argumentos obligatorios.

**R2.** CUANDO se ejecuta `rei init`, el sistema DEBE crear `.rei/specs/` si no existe.

**R3.** CUANDO se ejecuta `rei init`, el sistema DEBE crear `.rei/progress/work-items/` si no existe.

**R4.** CUANDO se ejecuta `rei init` y falta `.rei/progress/current.md`, el sistema DEBE crearlo desde la plantilla canónica.

**R5.** CUANDO se ejecuta `rei init` y falta `.rei/progress/history.md`, el sistema DEBE crearlo desde la plantilla canónica.

**R6.** CUANDO se ejecuta `rei init`, SI el directorio del proyecto no es un repositorio git y el binario `git` está disponible, ENTONCES el sistema DEBE inicializar un repositorio git.

**R7.** MIENTRAS `git` no esté disponible, `rei init` DEBE continuar, informar mediante un aviso y devolver el mismo código de salida que si git estuviera disponible.

**R8.** CUANDO se ejecuta `rei init`, el sistema DEBE detectar cada documento de personalización que contenga el marcador `<!-- REI:PENDIENTE -->`.

**R9.** CUANDO se ejecuta `rei init`, el sistema DEBE indicar, por cada documento de personalización, si está pendiente o completo.

**R10.** CUANDO se ejecuta `rei init`, el sistema DEBE mostrar el plan de pasos ordenado: (1) identidad/propósito, (2) stack y comandos, (3) arquitectura, (4) convenciones, (5) verificación y (6) cierre/validación, indicando para cada paso el documento destino.

**R11.** El sistema DEBE hacer `rei init` idempotente: ejecutarlo repetidamente NO DEBE sobrescribir archivos existentes ni fallar por estructura ya creada.

**R12.** SI `rei init` recibe argumentos que no correspondan a `status`, ENTONCES el sistema DEBE mostrar el uso por `stderr` y devolver código de salida `2`.

**R13.** `rei init` DEBE devolver `0` cuando completa el scaffold y el reporte, y `1` si no puede crear la estructura o los archivos base a partir de las plantillas.

---

# 2. `rei init status` — reporte determinista

**R14.** El sistema DEBE ofrecer el subcomando `rei init status`.

**R15.** CUANDO se ejecuta `rei init status`, el sistema DEBE reportar qué documentos de personalización siguen pendientes y cuáles están completos, SIN modificar ningún archivo.

**R16.** `rei init status` DEBE devolver `0` si no queda ningún documento pendiente y `1` si queda al menos uno.

**R17.** SI `rei init status` recibe argumentos adicionales no reconocidos, ENTONCES el sistema DEBE mostrar el uso y devolver código de salida `2`.

---

# 3. Marcadores de personalización

**R18.** Los documentos de personalización DEBEN contener el marcador `<!-- REI:PENDIENTE -->` mientras no estén personalizados. En concreto, `AGENTS.md` DEBE contener un marcador en la sección 2 y otro en la sección 3; `architecture.md`, `conventions.md` y `verification.md` DEBEN contener al menos un marcador cada uno.

**R19.** Un documento de personalización DEBE considerarse pendiente si y solo si contiene al menos un marcador `<!-- REI:PENDIENTE -->`.

**R20.** CUANDO el rol `initializer` termina de personalizar una sección o documento, el sistema DEBE eliminar todos los marcadores `<!-- REI:PENDIENTE -->` correspondientes a esa sección o documento.

---

# 4. Rol `initializer`

**R21.** CUANDO un proyecto contiene marcadores `<!-- REI:PENDIENTE -->`, el Leader DEBE poder delegar la personalización en el rol `initializer`.

**R22.** El rol `initializer` DEBE estar definido en `.rei/agents/initializer.md`, con sus precondiciones, protocolo, reglas y criterio de finalización.

**R23.** CUANDO el rol `initializer` inicia, DEBE determinar el primer paso pendiente mediante `rei init status` y reanudar desde ahí, sin repetir los documentos ya completos.

**R24.** El rol `initializer` DEBE conducir una entrevista guiada por cada paso pendiente, respetando el orden del plan de pasos.

**R25.** DONDE el usuario no sepa qué responder, el rol `initializer` DEBE ofrecer sugerencias concretas y marcarlas explícitamente como sugerencias.

**R26.** El rol `initializer` DEBE interpretar las respuestas del usuario y redactar el contenido en el formato esperado por cada documento: secciones 2 y 3 de `AGENTS.md`, `architecture.md`, `conventions.md` y `verification.md`.

**R27.** CUANDO el rol `initializer` completa un paso, DEBE escribir el documento correspondiente, eliminar sus marcadores y pedir confirmación al usuario antes de continuar con el siguiente paso.

**R28.** El rol `initializer` NO DEBE inventar hechos: solo DEBE registrar información confirmada por el usuario o verificable directamente en el repositorio.

**R29.** El rol `initializer` DEBE limitar sus escrituras a `AGENTS.md` (secciones 2 y 3), `.rei/docs/project/architecture.md`, `.rei/docs/project/conventions.md`, `.rei/docs/project/verification.md` y `.rei/config.json`.

**R30.** CUANDO el rol `initializer` completa el paso de verificación, DEBE actualizar `.rei/config.json` para alinear sus `checks` con los checkpoints rápidos declarados en `verification.md`.

**R31.** CUANDO no quedan marcadores pendientes, el rol `initializer` DEBE ejecutar `rei init status` y `rei doctor`, confirmar que no hay pendientes y devolver el control al Leader.

**R32.** El rol `initializer` NO DEBE implementar código del proyecto.

---

# 5. `rei doctor`

**R33.** CUANDO se ejecuta `rei doctor`, el sistema DEBE reportar los documentos de personalización que contienen el marcador `<!-- REI:PENDIENTE -->`.

**R34.** Un documento de personalización pendiente DEBE reportarse como advertencia y NO DEBE incrementar el número de fallos de `rei doctor`.

**R35.** CUANDO no hay documentos de personalización pendientes, `rei doctor` DEBE indicar que la documentación del proyecto está personalizada.

---

# 6. Ayuda

**R36.** `rei help` DEBE listar el comando `init` con su uso (`init [status]`) y una descripción breve.

**R37.** `rei help init` DEBE documentar el propósito, los subcomandos, los códigos de salida y la relación con el rol `initializer`.

---

# 7. Funcionamiento sin IA

**R38.** `rei init` y `rei init status` DEBEN funcionar sin IA, sin acceso a red y sin dependencias externas distintas de `git` (opcional).

---

# 8. Integración e invariantes

**R39.** CUANDO existen marcadores `<!-- REI:PENDIENTE -->` pendientes, `AGENTS.md` DEBE indicar que el Leader delegue en el rol `initializer` antes de continuar con otro trabajo.

**R40.** `rei check` NO DEBE devolver fallo por la presencia de marcadores `<!-- REI:PENDIENTE -->` (la personalización pendiente no es un error de integridad).

**R41.** El arnés DEBE considerar `.rei/agents/initializer.md` como parte de su integridad, de modo que `rei check` y `rei doctor` reporten un fallo si el archivo falta.
