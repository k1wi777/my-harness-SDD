# Uso de REI Harness

> Esta guía está dirigida a **personas**, no a agentes.
> Explica cómo operar REI Harness en el día a día: cómo arrancarlo, qué esperar en cada etapa y cómo participar en los puntos donde se requiere tu aprobación.
>
> Si eres un agente de IA trabajando en este repositorio, tu punto de entrada es `AGENTS.md`, no este documento.

---

# Tu responsabilidad en este proceso

REI Harness no reemplaza tu criterio — lo agiliza. Los agentes redactan planificación, implementan y verifican, pero la responsabilidad del proyecto sigue siendo tuya en dos momentos concretos:

- **Al negociar el Work Item.** Tú decides qué se construye, por qué y con qué alcance. El Leader no actúa hasta llegar a un acuerdo explícito contigo — su trabajo es comprender tu solicitud y formalizarla, no interpretarla por su cuenta.
- **Al aprobar la planificación.** Antes de que se escriba una sola línea de código, tú confirmas que la lógica, las decisiones de diseño y los requisitos son correctos. No se te pide revisar código — esa es justamente la parte que REI Harness te ahorra —, pero sí se te pide revisar las decisiones que determinan qué se va a construir.

Ninguna de las dos es un trámite. Son los dos puntos donde el proceso se detiene y te devuelve el control antes de avanzar.

---

# 1. Primera vez en el repositorio

1. Obtén el binario `rei` para tu sistema operativo (compílalo con `make build` si tienes Go, o descarga una release publicada). El binario **no necesita Go ni ninguna dependencia** para ejecutarse.
2. Copia `AGENTS.md` y la carpeta `.rei/` dentro de tu proyecto.
3. **Personaliza la documentación de tu proyecto** (ver Sección 1.1).
4. Ejecuta:
   ```bash
   rei check
   ```
5. `rei check` verificará que REI Harness esté completo y preparará las carpetas `.rei/specs/` y `.rei/progress/work-items/` si no existen todavía.

Si `rei check` reporta `[FAIL]`, falta algún archivo de REI Harness — revisa el listado que imprime antes de continuar. No inicies trabajo hasta que el resumen final diga "REI Harness listo para trabajar".

### Permisos de ejecución

El CLI `rei` es un binario compilado; no necesita permisos de ejecución especiales más allá de estar disponible en el `PATH`.

---

## 1.1 Qué personalizar y qué no

La documentación en `docs/` está dividida en dos carpetas. Solo una de ellas debes adaptarla a tu proyecto.

### Personaliza (documentación del proyecto)

| Archivo | Qué debes definir |
|---------|-------------------|
| `.rei/docs/project/architecture.md` | Decisiones arquitectónicas de tu repositorio. |
| `.rei/docs/project/conventions.md` | Estilo de código y convenciones de desarrollo. |
| `.rei/docs/project/verification.md` | Checkpoints (`V1`, `V2`, ...) y cómo demostrar que el trabajo funciona. |
| `AGENTS.md` — Sección 2 | Propósito del producto, usuarios, principios y límites. |
| `AGENTS.md` — Sección 3 | Stack técnico y comandos principales del proyecto. |
| `.rei/config.json` | Checks de verificación (`checks`) alineados con `verification.md`. |

### No personalices (documentación de REI Harness)

| Carpeta / archivo | Motivo |
|-------------------|--------|
| `.rei/docs/harness/` | Define el funcionamiento del arnés — es igual en todos los proyectos. |
| `.rei/agents/` | Roles y protocolos de los subagentes. |
| `.rei/docs/usage.md` | Esta guía. |

Si modificas un archivo de `.rei/docs/harness/` o `.rei/agents/`, estarías alterando REI Harness en sí, no la configuración de tu proyecto.

---

# 2. Negociar el Work Item con el Leader

Descríbele al agente, en lenguaje natural, lo que necesitas. El agente actúa como **Leader**, y su primer trabajo no es ejecutar — es entender.

Si tu solicitud tiene cualquier ambigüedad de alcance, el Leader **debe** detenerse y preguntarte antes de continuar — no está autorizado a asumir nada por su cuenta. Esto normalmente significa una conversación de ida y vuelta: el Leader propone una interpretación, tú la ajustas, hasta llegar a un acuerdo claro sobre qué se va a construir.

Esta negociación es tu responsabilidad, no un paso automático. El `meta.json` que el Leader crea al final **representa el acuerdo alcanzado contigo** — es tu criterio quedando registrado, no una decisión que el agente tomó solo. Vale la pena tomarte este paso en serio: cuanto más claro quede el acuerdo aquí, menos ambigüedad tendrá que resolver el Spec Author después.

Solo una vez que existe ese acuerdo explícito, el Leader crea el Work Item (`.rei/specs/<work-item-id>/meta.json`) con estado `pending`, genera su `id` con fecha y hora de creación y decide si corresponde a una Feature o una Task (ver `.rei/docs/harness/workflow.md`).

---

# 3. El ciclo de trabajo

Cada Work Item pasa por las mismas etapas, sin importar si es una Feature o una Task:

```
pending → spec_author → ready → ⏸ tu aprobación → in_progress → review → reviewer → done
```

En la práctica, esto es lo que vas a ver:

1. **Planificación.** El Leader delega en el Spec Author, que redacta la planificación (`requirements.md` + `design.md` + `tasks.md` para una Feature, o `plan.md` para una Task) dentro de `.rei/specs/<work-item-id>/`.
2. **Tu aprobación.** El Leader se detiene y te pide revisar la planificación. Esto **nunca se salta** — es el punto donde tienes control total antes de que se escriba una sola línea de código.
3. **Implementación.** Solo después de tu aprobación explícita, el Implementer ejecuta la planificación tal como fue aprobada.
4. **Revisión.** El Reviewer valida el trabajo contra la planificación y los checkpoints de `.rei/docs/project/verification.md`.
5. **Finalización.** Si todo pasa, el Work Item queda `done` y su resumen se archiva en `.rei/progress/history.md`.

---

# 4. Cómo aprobar o pedir cambios

Cuando el Leader te presente una planificación lista (`ready`), te toca ejercer la parte de este proceso que nadie más puede hacer por ti: revisar que la lógica, las decisiones de diseño y los requisitos sean los correctos **antes** de que exista código. No se trata de leer por leer — es la comprobación real de que lo que se va a construir es lo que realmente acordaste en la negociación del Work Item (Sección 2).

- Si estás de acuerdo, responde algo como **"aprobado"**. El Leader avanzará el Work Item a `in_progress`.
- Si algo no refleja lo acordado, o encuentras una decisión de diseño con la que no estás de acuerdo, dilo en lenguaje natural. El Leader relanzará al Spec Author y volverá a presentarte la planificación actualizada.

REI Harness nunca avanza sin esta aprobación explícita — no por burocracia, sino porque la responsabilidad de qué se construye es tuya, no del agente.

Lo mismo aplica después de la revisión: si el Reviewer rechaza el trabajo (`changes_requested`), el Leader te lo informará y relanzará al Implementer con las correcciones necesarias — no necesitas hacer nada salvo revisar el resultado otra vez al final.

---

# 5. Continuar una sesión interrumpida

Si cierras la conversación con un Work Item a medias, no se pierde nada — vive en `.rei/specs/<work-item-id>/` y `.rei/progress/work-items/<work-item-id>/`, no en el historial del chat.

La próxima vez que ejecutes `rei check` y hables con el agente, el Leader consultará el estado mediante el CLI (`rei session`, `rei items status`) y te preguntará si quieres continuar, reiniciar o cancelar ese Work Item. Tú decides; el Leader nunca lo asume por su cuenta.

---

# 6. Preguntas frecuentes

**¿Puedo tener varios Work Items en paralelo?**
No al mismo tiempo en estado `in_progress` — REI Harness solo permite uno activo a la vez (ver `.rei/docs/harness/workflow.md`). Sí puedes tener varios `pending` esperando su turno.

**¿Dónde declaro los checks de mi proyecto?**
En `.rei/config.json` (ver `.rei/docs/project/verification.md`). `rei check` los ejecuta y refleja el resultado en su código de salida.

**¿Qué hago si `rei check` falla por un check del proyecto?**
Significa que algún comando del proyecto (tests, build, lint) falló. Revisa el detalle en la salida y corrige antes de solicitar revisión de un Work Item.

**¿Dónde veo el historial de todo lo que se ha hecho?**
En `.rei/progress/history.md` — es un registro permanente, nunca se sobrescribe.

**¿Necesito leer toda la documentación de `.rei/docs/` para usar REI Harness?**
No. Esta guía es suficiente para el uso diario. Los agentes cargan solo lo necesario según la etapa — como humano, consulta `.rei/docs/project/` si quieres revisar las reglas de tu proyecto, o `.rei/docs/harness/` si quieres entender el detalle del arnés.

**¿REI Harness necesita git?**
No, pero funciona mejor con él. Si el repositorio usa git, el Reviewer revisa solo los cambios reales del Work Item (`rei review-diff`) en lugar de releer toda la especificación, y en `changes_requested` puede acotar la revisión a lo que cambió. Si no hay git, el flujo sigue funcionando en modo lectura. `rei check` inicializa git automáticamente si no existe.
