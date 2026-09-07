# meta.json

> `meta.json` describe un Work Item y permite a los agentes conocer su tipo, estado e información básica.

Cada carpeta dentro de `.rei/specs/` debe contener un único archivo `meta.json`. El nombre de la carpeta es el `id` inmutable del Work Item.

```text
.rei/specs/
└── <work-item-id>/
    ├── meta.json
    └── ...
```

---

# Estructura

```json
{
  "id": "",
  "title": "",
  "description": "",
  "type": "",
  "status": "",
  "created_at": ""
}
```

---

# Campos

## id

Identificador inmutable del Work Item. Deberá coincidir exactamente con el mismo nombre de su carpeta en `.rei/specs/` y `.rei/progress/work-items/`.

Formato obligatorio:

```text
YYYY-MM-DD_HH-mm__slug-en-kebab-case
```

La fecha y hora corresponden al momento de creación del Work Item. El identificador no cambia aunque el trabajo se reanude o cambie su título.

Si ya existe un identificador igual, añade un sufijo numérico al slug (`-2`, `-3`, ...) para evitar colisiones sin modificar la fecha y hora originales.

## title

Nombre corto y descriptivo del Work Item.
Debe resumir claramente el objetivo del trabajo.

---

## description

Descripción técnica del trabajo.
Debe resumir el acuerdo alcanzado con el usuario y proporcionar el contexto necesario para continuar el workflow.
No debe copiar literalmente la conversación.

---

## type

Define el tipo de Work Item.
Valores válidos:

- `feature`
- `task`

---

## status

Define la etapa actual del workflow.
Los estados válidos y sus transiciones se definen en `workflow.md`.

---

## created_at

Fecha y hora exactas de creación del Work Item en formato ISO 8601, incluyendo la zona horaria. Debe corresponder al momento utilizado para construir el prefijo del campo `id`.
