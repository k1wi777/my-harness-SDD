# Implementación — rei doctor

## Resumen

Se implementó el comando `rei doctor`, un diagnóstico de **solo lectura** del
harness que verifica archivos/directorios requeridos, valida todos los Work Items
existentes, reporta la sesión activa y el estado de git, informa cuántos checks
hay configurados (sin ejecutarlos) y emite un veredicto con el código de salida
correspondiente. Se añadió el despacho del subcomando y su línea en `rei help`.

El trabajo siguió `plan.md` sin desviaciones de alcance.

## Archivos

- `internal/doctor/doctor.go` — **nuevo**. Función `Run(p *paths.Project, out io.Writer) int` con las seis secciones del reporte y helpers `ok`/`warn`/`info`/`fail`.
- `internal/doctor/doctor_test.go` — **nuevo**. Tres tests: harness completo → `0`, harness incompleto → `1`, Work Item con `meta.json` inválido → `1`.
- `internal/cli/cli.go` — **modificado**. Import de `doctor`, `case "doctor"`, función `cmdDoctor` y línea en `printHelp`.
- `.rei/specs/2026-10-02_00-14__rei-doctor/plan.md` — pasos marcados como completados.

## Cambios

### `internal/doctor/doctor.go`

- `Run` define un contador `failures` y los helpers de etiquetas con el mismo
  estilo que `internal/check` (`[OK]`, `[WARN]`, `[INFO]`, `[FAIL]`).
- **Integridad**: recorre `check.RequiredFiles` (variable reutilizada, sin
  duplicar la lista) y comprueba `.rei/specs/`,
  `.rei/progress/work-items/`, `.rei/progress/current.md` y
  `.rei/progress/history.md` con `os.Stat` (nada de crear archivos).
- **Work Items**: `state.ListWorkItems` + `validate.WorkItem`; cada issue se
  imprime como `[<Level>] <id>: <Message>` y cada `validate.LevelFail`
  incrementa el contador global.
- **Sesión**: `state.ReadSession` informa sesión activa o ausencia; un error de
  lectura se degrada a `WARN`.
- **Git**: `exec.LookPath` detecta la ausencia de `git` (`WARN`); si es repo,
  usa `gitx.IsRepo`/`gitx.HeadSHA` e informa `HEAD=<sha>` o
  `Sin commits todavía.`. Nunca inicializa git.
- **Checks**: `config.Load(.rei/config.json)`; error → `FAIL`, si no →
  `[INFO] <n> checks ... (no ejecutados)`. Nunca se ejecutan.
- **Resultado**: sin fallos → `[OK] Sin fallos. Harness íntegro.` y `return 0`;
  con fallos → `[FAIL] <n> fallo(s) detectado(s).` y `return 1`.

### `internal/cli/cli.go`

- Nuevo `case "doctor": return cmdDoctor(args[1:])`.
- `cmdDoctor`: valida ausencia de argumentos extra (`uso: rei doctor` en
  `stderr`, código `2`), obtiene el proyecto con `project()` y delega en
  `doctor.Run(p, os.Stdout)`.
- `printHelp` incluye la línea del comando en el bloque `Comandos:`.

### `internal/doctor/doctor_test.go`

- `setupValidHarness` crea un harness válido en `t.TempDir()` (archivos de
  `check.RequiredFiles`, directorios requeridos y `current.md`/`history.md`).
- `TestRunOK` espera `0`; `TestRunMissingFile` y `TestRunInvalidWorkItem`
  esperan `1`. Los tests no dependen de que el entorno sea un repo git (git solo
  produce `WARN`).

## Verificación

Ejecutado desde la raíz del repositorio con
`PATH="$HOME/.local/bin:/tmp/opencode/go/bin:$PATH"`:

- `gofmt -l internal/doctor internal/cli` → sin salida (formato correcto).
- `rei check --quiet` → exit `0`.
- `make build` → OK.
- `make test` → OK (`internal/doctor` pasa; resto de paquetes sin regresiones).
- `make vet` → OK.
- `rei doctor` → exit `0`; reporta integridad OK, Work Item
  `2026-10-02_00-14__rei-doctor` OK, sesión activa, git detectado con `HEAD` y
  `0 checks` configurados.
- `rei help` → lista `doctor`.

## Observaciones

- `~/.local/bin/rei` es un symlink a `bin/rei`, por lo que la verificación con
  el binario del PATH usa la build recién generada.
- No se tocó `.rei/docs/usage.md` ni `package-lock.json` (aparecen como cambios
  preexistentes ajenos a este Work Item).
- Se registraron únicamente los archivos nuevos con
  `git add internal/doctor/`; no se añadió `.rei/progress/`.
