# Plan — `rei status set <id> <status>`

## Objetivo

Agregar el comando `rei status set <id> <status> [--force]` que cambia el
`status` de un Work Item en `.rei/specs/<id>/meta.json` de forma segura,
preservando **todos** los campos restantes del archivo (en particular
`base_commit` y `last_review_commit`). El comando valida que el Work Item exista,
que el estado destino sea válido y que la transición esté permitida por
`workflow.md`; `--force` permite aplicar transiciones no estándar. También
actualiza `rei help`.

Motivo: durante el dogfooding de `rei doctor`, editar `meta.json` a mano borró
`base_commit`. El comando evita ese error al leer/reescribir la estructura
completa en lugar de editar texto.

## Archivos

- `internal/meta/meta.go` — lógica de estados y transiciones + actualización segura.
- `internal/meta/meta_test.go` — pruebas unitarias de transiciones y preservación.
- `internal/cli/cli.go` — subcomando `status set`, dispatch y ayuda.

No se crea un paquete nuevo: `internal/meta` ya es dueño de los estados
(`ValidStatuses`, `IsValidStatus`) y del tipo `Meta`. Reutilizarlo evita
duplicar lógica (convenciones: reutilizar antes que duplicar).

## Cambios

### 1. `internal/meta/meta.go`

Añadir la tabla de transiciones estándar derivada de `.rei/docs/harness/workflow.md`:

| Desde | Hacia | Origen en `workflow.md` |
|-------|-------|--------------------------|
| `pending` | `ready` | Spec Author termina la planificación |
| `ready` | `in_progress` | Aprobación humana |
| `in_progress` | `review` | Implementer termina |
| `review` | `done` | Reviewer aprueba |
| `review` | `changes_requested` | Reviewer rechaza (ramificación) |
| `changes_requested` | `in_progress` | Leader relanza al Implementer |
| cualquiera | `blocked` | `* → blocked` (el trabajo no puede continuar) |
| mismo → mismo | — | No-op idempotente (se trata en el CLI, ver abajo) |

Cualquier otra combinación **no** es estándar (p. ej. desbloquear
`blocked → X`, reabrir `done → X` u omitir etapas como `pending → in_progress`).
Esas solo se permiten con `--force`.

```go
// allowedTransitions describe las transiciones estándar de workflow.md.
// Cualquier estado puede pasar a blocked.
var allowedTransitions = map[string][]string{
    StatusPending:          {StatusReady},
    StatusReady:            {StatusInProgress},
    StatusInProgress:       {StatusReview},
    StatusReview:           {StatusDone, StatusChangesRequested},
    StatusChangesRequested: {StatusInProgress},
}

// TransitionAllowed indica si from -> to es una transición estándar.
func TransitionAllowed(from, to string) bool {
    if to == StatusBlocked {
        return true
    }
    for _, s := range allowedTransitions[from] {
        if s == to {
            return true
        }
    }
    return false
}

// SetStatus actualiza el status de un meta.json preservando el resto de campos.
// Devuelve el estado anterior.
func SetStatus(path, status string) (string, error) {
    m, err := Load(path)
    if err != nil {
        return "", err
    }
    previous := m.Status
    m.Status = status
    if err := m.Save(path); err != nil {
        return "", err
    }
    return previous, nil
}
```

`SetStatus` usa `Load` + `Save` (struct completo), garantizando que
`id`, `title`, `description`, `type`, `created_at`, `base_commit` y
`last_review_commit` se conservan. **Prohibido** editar `meta.json` como texto
o con expresiones regulares.

### 2. `internal/meta/meta_test.go`

- `TestTransitionAllowed`: verificar cada transición permitida de la tabla y
  rechazar casos no estándar (`pending → in_progress`, `done → in_progress`,
  `blocked → in_progress`) y verificar que `* → blocked` siempre se permite.
- `TestSetStatusPreservesFields`: crear un `Meta` con `BaseCommit` y
  `LastReviewCommit` no vacíos, guardarlo, llamar a `SetStatus` con otro estado
  y comprobar que `Load` devuelve el nuevo `Status` y que los demás campos
  (incluidos los dos commits) quedan intactos.

### 3. `internal/cli/cli.go`

- En `Run`, añadir `case "status": return cmdStatus(args[1:])`.
- Implementar `cmdStatus`:

  1. Parsear: debe empezar por `set` y recibir `<id>` y `<status>`. Flags
     admitidos tras ellos: `--force`. Un flag desconocido o argumentos
     insuficientes → uso + salida `2`.
  2. Validar que `status` pertenezca a `meta.ValidStatuses`; si no, error y
     salida `1`.
  3. Obtener el proyecto (`project()`). Verificar que exista
     `p.MetaFile(id)`; si no, error indicando que el Work Item no existe y
     salida `1`.
  4. `meta.Load` para conocer el estado anterior.
  5. Si el estado anterior ya es igual al destino: informar
     `Work Item <id>: ya está en '<status>'.` y salida `0` (no-op, no requiere
     `--force`).
  6. Si `--force` no está presente y `meta.TransitionAllowed(anterior, destino)`
     es falso: error indicando la transición rechazada, que se consulte
     `workflow.md` y que puede usarse `--force`; salida `1`.
  7. `meta.SetStatus(path, destino)`. En éxito imprimir
     `Work Item <id>: <anterior> -> <destino>.` (añadir `(forzado)` si se usó
     `--force` sobre una transición no estándar); salida `0`. En error de
     escritura, error y salida `1`.

- Añadir los imports necesarios: `internal/meta` y `strings`.
- `--force` **solo** omite la validación de transición: nunca omite la
  validación de estado válido ni la de existencia del Work Item.

### 4. Ayuda (`printHelp`)

Añadir, tras `items status`:

```text
  status set <id> <status> [--force]
                               Cambia el estado de un Work Item en meta.json
```

## Validaciones y códigos de salida

| Situación | Código |
|-----------|--------|
| Éxito (cambio aplicado o no-op mismo estado) | `0` |
| Subcomando distinto de `set`, faltan `<id>`/`<status>`, flag desconocido | `2` |
| `status` destino inválido (no está en `meta.ValidStatuses`) | `1` |
| El Work Item (su `meta.json`) no existe | `1` |
| Transición no estándar sin `--force` | `1` |
| Error al cargar/guardar `meta.json` | `1` |

## Restricciones

- **No** editar `.rei/progress/current.md` desde el comando.
- **No** ejecutar git ni checks.
- **No** tocar otros Work Items ni otros archivos fuera de `meta.json` del id.
- **No** añadir dependencias nuevas ni paquetes nuevos.
- **Siempre** preservar todos los campos de `meta.json` (lectura/escritura con
  struct).
- Mantener el estilo existente del CLI (mensajes `ERROR:` a stderr, `uso:` para
  errores de uso).

## Pasos

- [x] 1. Añadir a `internal/meta/meta.go` `allowedTransitions`, `TransitionAllowed` y `SetStatus`.
- [x] 2. Añadir pruebas en `internal/meta/meta_test.go` (`TestTransitionAllowed`, `TestSetStatusPreservesFields`).
- [x] 3. Añadir `cmdStatus` y su dispatch (`case "status"`) en `internal/cli/cli.go`.
- [x] 4. Actualizar `printHelp` con la entrada de `status set`.
- [x] 5. Ejecutar `go build ./...` y `go test ./...` (o al menos `./internal/meta/...` y `./internal/cli/...`) y corregir lo que falle.
- [x] 6. Verificación manual: usar un Work Item de prueba (o uno existente), ejecutar `rei status set <id> <status>` y comprobar con `rei doctor`/`cat meta.json` que `base_commit` y `last_review_commit` se conservan.
