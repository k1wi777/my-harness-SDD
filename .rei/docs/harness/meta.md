# meta.json

> `meta.json` describe un Work Item y permite a los agentes conocer su tipo, estado e información básica.

Cada carpeta dentro de `.rei/specs/` debe contener un único archivo `meta.json`.

```text
.rei/specs/
└── <work-item>/
    ├── meta.json
    └── ...
```

---

# Estructura

```json
{
  "title": "",
  "description": "",
  "type": "",
  "status": "",
  "created_at": ""
}
```

---

# Campos

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
