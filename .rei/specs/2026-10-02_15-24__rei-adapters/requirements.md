# Requirements — rei runtime adapters

> Work Item: `2026-10-02_15-24__rei-adapters` (feature).
>
> Los requisitos describen **qué** debe hacer el sistema, no **cómo**. Se redactan
> en EARS. Salvo indicación contraria, el sujeto "el sistema" designa al CLI `rei`
> y a los artefactos del arnés que este genera.
>
> **Objetivo:** introducir un formato canónico de rol en `.rei/agents/<rol>.md`
> (frontmatter + `## Contrato` autosuficiente + `## Referencia` opcional) e
> implementar un adaptador de runtime para **OpenCode** que genere agentes nativos
> a partir del Contrato, consumibles tanto por runtimes nativos como por el
> fallback sin adaptador.
>
> **Roles a cubrir:** `leader` (primary), `spec_author`, `implementer`, `reviewer`,
> `initializer` (subagent).
>
> **Fuera de alcance:** adaptadores para otros runtimes distintos de OpenCode. El
> formato debe ser extensible, pero en esta Feature solo se implementa OpenCode.

---

# 1. Formato canónico de rol

**R1.** Cada archivo `.rei/agents/<rol>.md` DEBE comenzar con un bloque de
frontmatter YAML delimitado por `---` que contenga, como mínimo, las claves
`name`, `description`, `mode` y `tools`; la clave `model` DEBE ser opcional.

**R2.** El valor de `mode` DEBE ser `primary` para el rol `leader` y `subagent`
para los roles `spec_author`, `implementer`, `reviewer` e `initializer`.

**R3.** El valor de `tools` DEBE ser una lista de herramientas **genéricas**
cuyos elementos pertenezcan al conjunto `{read, write, edit, search, shell,
subagent}`; el sistema NO DEBE usar nombres de herramientas específicas de un
runtime en este campo.

**R4.** Cada archivo `.rei/agents/<rol>.md` DEBE contener exactamente una sección
`## Contrato` y exactamente una sección `## Referencia`, en ese orden.

**R5.** El archivo de rol NO DEBE contener contenido normativo fuera de las
secciones `## Contrato` y `## Referencia`.

---

# 2. Autosuficiencia del Contrato

**R6.** La sección `## Contrato` DEBE incluir, de forma explícita: identidad del
rol, objetivo, precondiciones, el protocolo completo paso a paso, las reglas
duras, el formato de salida y la lista de herramientas permitidas.

**R7.** La sección `## Contrato` DEBE incluir punteros a los documentos del
arnés que el rol necesita (por ejemplo, `AGENTS.md`, `.rei/docs/harness/...` y
`.rei/docs/project/...`).

**R8.** El Contrato DEBE ser autosuficiente: toda regla, precondición o paso
necesario para ejecutar correctamente el rol DEBE estar contenido en `## Contrato`.

**R9.** La sección `## Referencia` DEBE contener únicamente detalle opcional
(racional, ejemplos, casos límite) y su omisión NO DEBE impedir que un agente
ejecute correctamente el rol.

**R10.** SI una regla o un paso del rol actual es necesario para ejecutarlo,
ENTONCES el sistema DEBE ubicarlo en `## Contrato` y NO DEBE dejarlo únicamente
en `## Referencia`.

---

# 3. Equivalencia entre runtime nativo y no nativo

**R11.** CUANDO un rol se ejecuta en un runtime con adaptador nativo, el sistema
DEBE entregar al agente exactamente el mismo `## Contrato` que se entregaría en
un runtime no nativo, preservando el mismo flujo y las mismas capacidades.

**R12.** CUANDO un rol se ejecuta en un runtime con adaptador nativo, el sistema
DEBE además imponer el `tools` y, si está presente, el `model` declarados en el
frontmatter del rol.

**R13.** El sistema NO DEBE alterar el texto del Contrato al adaptarlo a un
runtime nativo, salvo para añadir la marca de archivo generado.

---

# 4. Fallback sin adaptador nativo

**R14.** CUANDO el runtime no disponga de un adaptador nativo para un rol, el
Leader DEBE transmitir al subagente únicamente la sección `## Contrato` del rol
correspondiente, y NO el archivo de rol completo.

**R15.** El Contrato del rol `leader` DEBE documentar explícitamente la regla de
fallback de R14.

---

# 5. Infraestructura de adaptadores

**R16.** El sistema DEBE disponer de un directorio `.rei/adapters/` con un
subdirectorio por runtime; en esta Feature DEBE existir `.rei/adapters/opencode/`.

**R17.** El directorio `.rei/adapters/opencode/` DEBE contener la plantilla del
agente nativo y el mapa de herramientas genéricas a herramientas de OpenCode.

**R18.** El mapa de herramientas de OpenCode DEBE cubrir todas las herramientas
genéricas definidas en R3.

**R19.** La infraestructura de adaptadores DEBE permitir añadir un runtime nuevo
creando su subdirectorio bajo `.rei/adapters/` sin modificar el formato canónico
de rol.

**R20.** El adaptador DEBE extraer del archivo de rol únicamente el frontmatter y
la sección `## Contrato`, excluyendo la sección `## Referencia`.

---

# 6. Comando `rei init opencode`

**R21.** El sistema DEBE ofrecer el comando `rei init opencode` que instala el
esqueleto del arnés y genera los archivos de agente nativos de OpenCode.

**R22.** La generación DEBE escribir un archivo por rol en el directorio de
agentes de OpenCode, con nombre `<rol>.md`.

**R23.** Cada archivo generado DEBE contener el frontmatter nativo mapeado a
partir del frontmatter canónico y, como cuerpo, el `## Contrato` del rol.

**R24.** Cada archivo generado DEBE incluir una marca que lo identifique como
archivo generado y prohíba su edición manual, indicando cómo regenerarlo.

**R25.** SI el archivo destino ya existe y NO contiene la marca de archivo
generado, ENTONCES el sistema NO DEBE sobrescribirlo y DEBE informarlo.

**R26.** El comando `rei init opencode` DEBE ser idempotente: ejecutarlo de nuevo
sobre archivos generados por el propio sistema DEBE dejarlos idénticos.

**R27.** El sistema DEBE ofrecer `rei init opencode --check`, que verifica que los
archivos generados existen y coinciden con la generación canónica, SIN escribir
ningún archivo.

**R28.** `rei init opencode --check` DEBE devolver `0` si todos los archivos
generados coinciden y `1` si falta alguno o difiere.

**R29.** El sistema DEBE seguir ofreciendo `rei init` y `rei init status` con su
comportamiento actual.

**R30.** SI `rei init` recibe argumentos no reconocidos, ENTONCES el sistema DEBE
mostrar el uso por `stderr` y devolver código de salida `2`.

---

# 7. Empaquetado en el esqueleto

**R31.** El esqueleto embebido en el binario `rei` DEBE incluir todo el árbol
`.rei/adapters/**`, de modo que `rei init` (y `rei init opencode`) lo despliegue
sin acceso al repositorio fuente ni a la red.

---

# 8. Verificación de la cobertura

**R32.** El sistema DEBE incluir una checklist de cobertura por cada rol que
enumere cada regla dura y cada paso del protocolo del rol y señale si está
presente en el `## Contrato` o, por ser detalle opcional, en `## Referencia`.

**R33.** El sistema DEBE incluir pruebas automatizadas que verifiquen: el parseo
del formato canónico, la extracción del Contrato, el mapeo de herramientas, la
generación de archivos nativos, la idempotencia, la detección de deriva con
`--check` y la no sobrescritura de archivos no generados.

**R34.** El sistema DEBE registrar una medición de tokens del `## Contrato`
frente al archivo de rol completo para cada rol, correspondiente a lo que se
transmite por spawn.

---

# 9. Ayuda y documentación

**R35.** `rei help init` DEBE documentar el subcomando `opencode`, su flag
`--check` y sus códigos de salida.

**R36.** El sistema DEBE documentar el formato canónico de rol (frontmatter,
Contrato y Referencia) y el procedimiento de adaptación, de modo que sea
extensible a otros runtimes.

---

# 10. Fuera de alcance

**R37.** El sistema NO DEBE implementar adaptadores para runtimes distintos de
OpenCode en esta Feature.

**R38.** El formato canónico de rol NO DEBE incluir nombres de herramientas ni
campos específicos de OpenCode.
