# Requirements — rei init installer

> Work Item: `2026-10-02_13-21__rei-init-installer` (feature).
>
> Los requisitos describen **qué** debe hacer el sistema, no **cómo**. Se redactan
> en EARS. El sujeto "el sistema" designa al CLI `rei`.
>
> **Objetivo:** convertir `rei init` en un instalador completo autosuficiente:
> desplegar el esqueleto de REI Harness en cualquier proyecto **sin clonar el
> repositorio fuente**.
>
> **Esqueleto desplegable:** `AGENTS.md` + `.rei/docs/**` + `.rei/agents/**` +
> `.rei/templates/**` + `.rei/config.json`.
>
> **Excluido del esqueleto:** `.rei/specs/` y `.rei/progress/` (estado del proyecto
> destino, no del arnés).
>
> **Fuera de alcance:** actualización del esqueleto (`rei init --update`) y
> adaptadores de runtime.

---

# 1. Empaquetado del esqueleto

**R1.** El binario `rei` DEBE incorporar en tiempo de compilación el esqueleto del
harness formado por `AGENTS.md`, todo `.rei/docs/**`, todo `.rei/agents/**`, todo
`.rei/templates/**` y `.rei/config.json`.

**R2.** El esqueleto embebido NO DEBE incluir `.rei/specs/` ni `.rei/progress/`.

**R3.** El contenido de cada archivo desplegado DEBE ser idéntico, byte a byte, al
del archivo correspondiente del esqueleto embebido.

**R4.** El binario `rei` DEBE poder instalar el esqueleto sin acceder al repositorio
fuente ni a la red.

---

# 2. Destino de la instalación

**R5.** CUANDO se ejecuta `rei init`, el sistema DEBE determinar el proyecto destino
como el ancestro más cercano que contenga un directorio `.rei/`; SI no existe
ninguno, ENTONCES DEBE usar el directorio de trabajo actual como raíz del destino.

---

# 3. Despliegue del esqueleto

**R6.** CUANDO se ejecuta `rei init`, el sistema DEBE desplegar todos los archivos
del esqueleto embebido dentro del proyecto destino, creando los directorios que
falten.

**R7.** SI un archivo del esqueleto ya existe en el destino, ENTONCES el sistema
DEBE omitirlo y NO DEBE modificar su contenido.

**R8.** SI `AGENTS.md` ya existe en el destino, ENTONCES el sistema DEBE emitir un
aviso indicándolo y omitirlo.

**R9.** CUANDO `rei init` termina el despliegue, el sistema DEBE reportar de forma
diferenciada cada archivo creado y cada archivo omitido, e incluir un resumen con
el número de archivos creados y de archivos omitidos.

**R10.** El sistema DEBE hacer `rei init` idempotente: ejecutarlo repetidamente NO
DEBE sobrescribir archivos existentes ni fallar por estructura ya creada.

---

# 4. Estructura de estado y repositorio git

**R11.** CUANDO se ejecuta `rei init`, el sistema DEBE crear `.rei/specs/` si no
existe.

**R12.** CUANDO se ejecuta `rei init`, el sistema DEBE crear
`.rei/progress/work-items/` si no existe.

**R13.** CUANDO se ejecuta `rei init` y falta `.rei/progress/current.md`, el sistema
DEBE crearlo desde la plantilla canónica.

**R14.** CUANDO se ejecuta `rei init` y falta `.rei/progress/history.md`, el sistema
DEBE crearlo desde la plantilla canónica.

**R15.** CUANDO se ejecuta `rei init`, SI el destino no es un repositorio git y el
binario `git` está disponible, ENTONCES el sistema DEBE inicializar un repositorio
git.

**R16.** MIENTRAS `git` no esté disponible, `rei init` DEBE continuar, informar
mediante un aviso y devolver el mismo código de salida que si git estuviera
disponible.

**R17.** `rei init` NO DEBE ejecutar los checks del proyecto definidos en
`.rei/config.json`.

---

# 5. Reporte de personalización

**R18.** TRAS el despliegue y la creación de estructura, `rei init` DEBE detectar los
documentos de personalización que contengan el marcador `<!-- REI:PENDIENTE -->`.

**R19.** CUANDO se ejecuta `rei init`, el sistema DEBE indicar, por cada documento de
personalización, si está pendiente o completo, y DEBE mostrar el plan de pasos, tal
como hacía antes de esta Feature.

---

# 6. `rei init status`

**R20.** El sistema DEBE seguir ofreciendo `rei init status` con el comportamiento
actual: reportar qué documentos de personalización están pendientes y cuáles
completos SIN modificar ningún archivo, devolviendo `0` si no queda ninguno y `1` si
queda al menos uno.

**R21.** SI `rei init status` recibe argumentos adicionales no reconocidos, ENTONCES
el sistema DEBE mostrar el uso por `stderr` y devolver código de salida `2`.

---

# 7. Códigos de salida

**R22.** `rei init` DEBE devolver `0` cuando completa el despliegue y la creación de
estructura, y `1` si no puede crear algún archivo del esqueleto o algún archivo o
carpeta base.

**R23.** SI `rei init` recibe argumentos que no correspondan a `status`, ENTONCES el
sistema DEBE mostrar el uso por `stderr` y devolver código de salida `2`.

---

# 8. Integridad tras instalar

**R24.** TRAS ejecutar `rei init` en un directorio sin `.rei/`, `rei check` NO DEBE
reportar archivos faltantes del harness ni fallos de integridad.

**R25.** Los documentos desplegados DEBEN conservar sus marcadores
`<!-- REI:PENDIENTE -->`, de modo que `rei init status` reporte la personalización
como pendiente y el wizard siga operativo.

---

# 9. Ayuda

**R26.** `rei help init` DEBE documentar que `rei init` despliega el esqueleto del
harness, crea la estructura de estado y, si falta, inicializa git, además de sus
códigos de salida.

---

# 10. Fuera de alcance

**R27.** `rei init` NO DEBE ofrecer la actualización del esqueleto (`--update`) ni
adaptadores de runtime en esta Feature.
