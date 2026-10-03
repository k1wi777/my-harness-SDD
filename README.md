# REI Harness

REI Harness es un arnés de desarrollo para proyectos asistidos por IA basado en **Spec Driven Development (SDD)** y **orquestación multiagente**.

Su objetivo no es generar código automáticamente, sino proporcionar una estructura donde distintos agentes colaboran de forma controlada, verificable y siempre bajo supervisión humana.

El repositorio actúa como la fuente de verdad del sistema: toda la planificación, implementación y revisión queda documentada dentro del propio proyecto.

---

# Instalación

El CLI `rei` se distribuye como **binario precompilado** en las
[releases de GitHub](https://github.com/k1wi777/my-harness-SDD/releases); no
necesita Go ni ninguna dependencia para ejecutarse.

1. Entra en la página de releases:
   <https://github.com/k1wi777/my-harness-SDD/releases>.
2. Descarga el archive correspondiente a tu sistema operativo y arquitectura:
   `rei_<versión>_<os>_<arch>.tar.gz` (en Windows, `rei_<versión>_windows_<arch>.zip`).
3. Descomprímelo y mueve el binario `rei` a un directorio de tu `PATH`
   (por ejemplo `/usr/local/bin` o `~/.local/bin`):

   ```bash
   tar -xzf rei_<versión>_<os>_<arch>.tar.gz
   mv rei /usr/local/bin/
   ```

4. Verifica la instalación:

   ```bash
   rei version
   ```

Si desarrollas sobre el propio repositorio y tienes Go instalado, puedes compilar
el binario localmente con `make build`; quedará en `bin/rei`.

## Actualizar el CLI

El binario `rei` puede actualizarse por sí mismo desde la línea de comandos:

- `rei update` descarga la última release publicada, verifica su `checksums.txt`
  (SHA-256) y reemplaza el ejecutable en uso de forma atómica (con backup y
  restauración si algo falla). En Windows, donde no se puede reemplazar el `.exe`
  en ejecución, abre la página de releases para descargarla a mano.
- `rei update --check` solo informa de si hay una versión nueva, sin descargar ni
  modificar nada.

Ambos comandos devuelven `0` cuando no hay nada que actualizar (ya estás en la
última versión o todavía no hay releases publicadas). Un código `1` se reserva
para fallos reales: red, HTTP 5xx, checksum inválido o asset ausente.

---

# Primeros pasos

Para entender cómo se usa REI Harness en el día a día (aprobar planificaciones, continuar una sesión interrumpida, preguntas frecuentes), consulta **`.rei/docs/usage.md`**.

1. Copia `AGENTS.md` y la carpeta `.rei/` dentro de tu proyecto.
2. Ejecuta `rei check` para verificar REI Harness e inicializar `.rei/specs/` y `.rei/progress/`.
3. Empieza a hablar con el agente describiendo lo que necesitas — actuará como Leader y coordinará el resto.

---

# Principios

REI Harness se construye sobre cuatro principios fundamentales.

## 1. El repositorio es la memoria

Los agentes no dependen del historial del chat.

Toda la información relevante vive dentro del proyecto:

- documentación
- especificaciones
- progreso
- historial
- estados

Una conversación puede perderse.
El repositorio no.

---

## 2. Un agente, una responsabilidad

Cada agente posee un único objetivo.

| Agente | Responsabilidad |
|---------|-----------------|
| Leader | Comprender la solicitud, coordinar el workflow y delegar. |
| Spec Author | Transformar un Work Item en una planificación técnica. |
| Implementer | Implementar únicamente la planificación aprobada. |
| Reviewer | Validar que el trabajo cumple la planificación y las reglas del proyecto. |

Ningún agente sustituye el trabajo de otro.

---

## 3. Spec Driven Development

Todo trabajo sigue el mismo flujo base.

```
Usuario
      │
      ▼
 Leader
      │
      ▼
 meta.json
      │
      ▼
Spec Author
      │
      ▼
Planificación
      │
      ▼
Aprobación humana
      │
      ▼
Implementer
      │
      ▼
Reviewer
      │
      ▼
Finalización
```

El código nunca se implementa antes de existir una planificación aprobada.

> Este diagrama muestra la ruta principal. Las ramificaciones (rechazo en revisión, bloqueos) están documentadas en `.rei/docs/harness/workflow.md`.

> El tramo `Usuario → Leader → meta.json` representa una negociación explícita, no una conversión automática — ver `.rei/docs/usage.md` (para el usuario) o `.rei/agents/leader.md` (protocolo del agente).


---

## 4. El humano siempre mantiene el control

Ningún agente puede avanzar automáticamente entre etapas críticas.

Toda planificación debe ser aprobada explícitamente antes de comenzar la implementación.

La aprobación humana forma parte del workflow y nunca puede omitirse.

---

# Organización del repositorio
## Estructura

```text
.
├── AGENTS.md                     # Punto de entrada para los agentes
├── README.md                     # Descripción del repositorio
├── go.mod                        # Módulo Go del CLI
├── cmd/rei/                      # Punto de entrada del CLI `rei`
├── internal/                     # Paquetes internos del CLI
├── .rei/                         # REI Harness
│   ├── agents/                   # Roles y protocolos de los agentes
│   ├── templates/                # Plantillas canónicas (fuente única de verdad)
│   ├── config.json               # Checks de verificación del proyecto
│   ├── docs/                     # Documentación de REI Harness y del proyecto
│   ├── specs/                    # Work Items y planificaciones por fecha
│   └── progress/                 # Estado e historial del trabajo
│       ├── current.md            # Estado de la sesión actual
│       ├── history.md            # Resumen histórico
│       └── work-items/           # Detalle por Work Item
│
├── src/                          # Código fuente del proyecto
├── tests/                        # Pruebas del proyecto
└── ...
```
---

# Work Items

REI Harness trabaja sobre **Work Items**. Existen dos tipos, **Feature** y **Task**, según el nivel de planificación que requiere el cambio.

Los criterios para elegir entre uno y otro, así como los estados y transiciones que sigue cada Work Item, están definidos en `.rei/docs/harness/workflow.md`.

---

# Documentación

La documentación está dividida en dos grupos según quién debe modificarla y desacoplada por responsabilidad.

## Documentación de REI Harness (`.rei/docs/harness/`)

Define el funcionamiento del arnés. **No requiere personalización** al adaptarlo a un proyecto.

| Documento | Propósito |
|-----------|-----------|
| `workflow.md` | Flujo completo: tipos de Work Item, estados y transiciones. |
| `specs.md` | Cómo se construyen las especificaciones (Features). |
| `task.md` | Formato del plan de una Task. |
| `progress.md` | Funcionamiento del sistema de progreso. |
| `meta.md` | Estructura y significado de `meta.json`. |

## Documentación del proyecto (`.rei/docs/project/`)

Describe las reglas de **tu repositorio**. **Debes personalizarla** al implementar REI Harness por primera vez.

| Documento | Propósito |
|-----------|-----------|
| `architecture.md` | Principios de arquitectura del proyecto. |
| `conventions.md` | Convenciones de desarrollo. |
| `verification.md` | Checkpoints y reglas de validación. |

También personaliza las secciones 2 y 3 de `AGENTS.md` (propósito y stack) y los checks en `.rei/config.json` (verificación).

Los agentes cargan únicamente la documentación necesaria para su etapa.

---

# Sistema de progreso

El progreso del proyecto también vive dentro del repositorio.

```
.rei/progress/
```

contiene:

- sesión actual (`current.md`)
- historial (`history.md`)
- carpeta por Work Item (`work-items/<work-item-id>/`) con reportes de implementación, revisión y bloqueos

Esto permite:

- continuar sesiones interrumpidas
- mantener trazabilidad
- conservar un historial permanente del proyecto

---

# Revisión por diff (git)

REI Harness funciona sin git, pero está optimizado para usarlo. Cuando el proyecto es un repositorio git, el Reviewer recibe un **paquete de revisión** (`rei review-diff`) con los cambios reales del Work Item en lugar de releer toda la especificación. Al aprobar se registra un `base_commit`, y en un `changes_requested` se registra un `last_review_commit` para que la próxima revisión vea solo los cambios pedidos.

Si no hay git, el flujo sigue funcionando en modo lectura. `rei check` inicializa git automáticamente si no existe.

---

# Filosofía del proyecto

Este repositorio no pretende construir un asistente autónomo.

Pretende construir un sistema donde:

- los agentes tienen responsabilidades claras;
- el contexto está distribuido;
- la documentación es la fuente de verdad;
- el humano conserva siempre el control del proceso.

La IA no reemplaza el proceso de desarrollo.

Lo sigue.
