# Roles y adaptadores

> Documento del arnés. Define el formato canónico de rol y cómo se adapta a un
> runtime nativo. No requiere personalización por proyecto, salvo
> `.rei/adapters/`, que puede ajustarse.
>
> Consulta este documento si necesitas crear o modificar un rol, entender qué
> se transmite a un subagente o añadir un adaptador de runtime.

---

# Objetivo

Separar dos cosas que antes vivían mezcladas en `.rei/agents/<rol>.md`:

1. **Qué debe hacer el rol** — su contrato operativo, en `## Contrato`.
2. **Detalle de apoyo** — racional y ejemplos, en `## Referencia`.

Sobre ese formato canónico se construye un adaptador que traduce el rol al
formato nativo de un runtime concreto.

---

# Formato canónico de rol

Cada `.rei/agents/<rol>.md` consta de, en este orden:

1. Un bloque de frontmatter YAML delimitado por `---`.
2. La sección `## Contrato`.
3. La sección `## Referencia`.

No debe haber contenido normativo fuera de `## Contrato` y `## Referencia`.

## Frontmatter

| Clave | Obligatoria | Valor |
|-------|-------------|-------|
| `name` | Sí | Identificador del rol (nombre del archivo). |
| `description` | Sí | Una frase; es el texto que muestra el runtime al elegir agente. |
| `mode` | Sí | `primary` para el rol `leader`; `subagent` para el resto. |
| `tools` | Sí | Lista inline de herramientas **genéricas**: `read`, `write`, `edit`, `search`, `shell`, `subagent`. |
| `model` | No | Modelo del runtime, p. ej. `provider/modelo`. |

```md
---
name: implementer
description: Implementa un único Work Item siguiendo la planificación aprobada.
mode: subagent
tools: [read, write, edit, search, shell]
---
```

El campo `tools` usa nombres **genéricos**: el formato canónico nunca menciona
herramientas de un runtime concreto. La traducción al runtime la hace el
adaptador.

## Contrato vs Referencia

- **`## Contrato`** es autosuficiente. Contiene identidad, objetivo,
  precondiciones, el protocolo completo paso a paso, las reglas duras, el
  formato de salida, las herramientas permitidas y punteros a la documentación
  del arnés. Es lo único que un agente necesita para ejecutar el rol y lo único
  que se transmite a un subagente.
- **`## Referencia`** contiene solo detalle opcional (racional, ejemplos, casos
  límite). Su omisión no debe impedir ejecutar el rol.

**Regla de oro:** si una regla o un paso es necesario para ejecutar el rol, vive
en `## Contrato`; si es justificación o ejemplo, vive en `## Referencia`.

---

# Extracción y generación

El adaptador extrae el frontmatter y el `## Contrato` (excluyendo `## Referencia`)
y los mapea al formato nativo. El mismo Contrato alimenta a los dos caminos:

- **Runtime sin adaptador nativo (fallback):** el Leader transmite al subagente
  únicamente la sección `## Contrato` del rol, nunca el archivo completo. Esto
  ahorra tokens por spawn.
- **Runtime con adaptador nativo (OpenCode V2):** `rei init opencode` genera
  `.opencode/agents/<rol>.md` con el frontmatter nativo (`permission`), la marca
  `GENERATED` y el `## Contrato` como cuerpo.

## Mapa de herramientas de OpenCode V2

| Genérica canónica | Clave(s) de permiso |
|-------------------|---------------------|
| `read` | `read` |
| `write` | `edit` |
| `edit` | `edit` |
| `search` | `glob`, `grep` |
| `shell` | `bash` |
| `subagent` | `task` |

`write` y `edit` colapsan en la misma clave `edit` (se deduplica). El orden del
mapa `permission` sigue el orden de declaración de `tools` en el rol.

---

# Añadir un runtime nuevo

El formato canónico no cambia al añadir un runtime; solo se añade su adaptador:

1. Crea `.rei/adapters/<runtime>/` con su plantilla y su mapa de herramientas.
2. Añade la construcción nativa correspondiente en `internal/adapter`.
3. Registra el subcomando de instalación/verificación en el CLI.

`.rei/adapters/` es el punto de extensión: no expone ningún runtime en el
formato canónico de rol.

---

# Comando `rei init opencode`

```text
rei init opencode           # instala el esqueleto (incluye .rei/adapters/) y genera los agentes nativos
rei init opencode --check   # verifica sin escribir: 0 si todo coincide; 1 si falta o difiere
```

- Los archivos generados viven en `.opencode/agents/<rol>.md` e incluyen la
  marca `GENERATED`.
- `rei init opencode` **no sobrescribe** archivos que no lleven esa marca; los
  avisa.
- Es idempotente: reejecutarlo sobre archivos generados los deja idénticos.
- `rei init` y `rei init status` conservan su comportamiento.
