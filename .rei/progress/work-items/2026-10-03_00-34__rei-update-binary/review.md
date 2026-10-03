# Revisión — `rei update [--check]`

> Work Item: `2026-10-03_00-34__rei-update-binary`
> Tipo: feature · Resultado: **APROBADO**

## Resultado

Implementación conforme a la planificación aprobada (T1–T12). Cobertura
completa de R1–R15, verificación reproducible y sin desviaciones.

## Verificaciones

- `gofmt -l internal/` → sin salida.
- `go vet ./...` → OK.
- `make test` → OK (todos los paquetes; `internal/update` y `internal/cli` incluidos).
- `make build` → OK.
- `rei check --quiet` → exit 0.
- `rei help` incluye `update [--check]`; `rei update --help` documenta `--check` (R14).
- `rei update --check` real → 404 del repo sin releases publicadas, exit 1 con
  mensaje claro (degradación aceptable: la comprobación no pudo realizarse).
- `rei update bogus` → uso + exit 2 (R15).
- Tests herméticos: `httptest` + `t.TempDir()`, sin red externa ni sustitución
  del binario real.

## Cobertura

R1 (`case "update"`), R2 (`fetchLatestRelease`), R3/R4 (`parseSemver`/
`compareSemver`, up-to-date → 0), R5 (`assetName`/`findAsset`/`findChecksums`),
R6/R11 (SHA-256 antes de extraer/instalar), R7/R13 (`replaceBinary`: temp junto
al destino, chmod 0755, rename con backup y restauración), R8 (Windows abre
navegador + imprime URL), R9 (`--check` no escribe), R10 (`dev` avisa y ofrece
instalar), R12 (red/permisos → 1 sin corromper), R14/R15.

## Seguridad

- No se escribe el binario antes de verificar el SHA-256 (`update.go`:
  descarga → verificacion → extrae → reemplaza).
- Sin `InsecureSkipVerify` (TLS por defecto de `net/http`).
- `safeMember` rechaza rutas absolutas y `..` en tar/zip (sin path traversal).
- El temporal de sustitución se crea en el directorio destino (rename atómico).

## Observaciones

- Ni `requirements.md`, `design.md` ni `tasks.md` presentan anexos o cambios
  posteriores a la planificación.
- Artefacto ajeno al Work Item: `package-lock.json` vacío sin rastrear en la
  raíz (no afecta a R1–R15; conviene limpiarlo en otro item).

## Acciones requeridas

Ninguna.
