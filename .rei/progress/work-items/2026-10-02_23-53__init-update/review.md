# Revisión — `rei init --update [--force]`

## Resultado

Aprobado (`done`).

## Verificación (independiente)

- V1/V2/V3: `gofmt -l .` limpio, `go vet ./...` 0, `make test` 0, `make build` 0.
- V5: `rei check --quiet` 0, `rei init opencode --check` 0, `rei init claude --check` 0.
- V4 en proyecto temporal (`/tmp/opencode/rev-update-*`) con `bin/rei` recién compilado:
  - Personalizado (`conventions.md` sin marcador), `leader.md` y `config.json`
    modificados → `[SKIP]`, contenido intacto.
  - Doc con marcador (`architecture.md`) → `[UPD]` al embebido.
  - `--force` → `[UPD]` de los 3 (0 omitidos), restaurados al embebido.
  - `.rei/specs/**` y `.rei/progress/**` sin cambios (hashes antes/después, también con `--force`).
  - Segunda `--update` → `0 nuevo(s), 0 actualizado(s), 28 sin cambios, 0 omitido(s)` (idempotente).
  - Fuera de proyecto REI → mensaje canónico y exit 1.
  - `.opencode/` presente → nativos regenerados (0 si no existen).
  - `.rei/install-manifest.json` escrito por `init`, coherente (28 entradas, 0 mismatch) y no embebido.
- Trazabilidad: los 8 pasos del plan están `[x]`; el diff son los 9 archivos
  previstos, sin tocar `embed.go`, `initwizard.go` ni el esqueleto.
- `rei help init` documenta `--update [--force]`, la política marker-aware, el
  manifiesto y que no toca `specs/**`/`progress/**`.

## Observaciones

- `package-lock.json` sin rastrear en la raíz es preexistente (10-sep), ajeno al Work Item.
- Sin dependencias externas nuevas.

## Acciones requeridas

Ninguna.
