# Plan — rei doctor

## Objetivo

Agregar el comando `rei doctor`: un diagnóstico **de solo lectura** que verifica la
integridad del harness y devuelve un veredicto.

Debe:

- verificar los archivos y directorios requeridos del harness **sin crear nada**;
- validar **todos** los Work Items existentes usando la validación ya existente;
- reportar la sesión activa (`current.md`);
- reportar el estado de git (`repositorio detectado` / `HEAD`, o ausencia de git) **sin inicializar git**;
- informar cuántos checks hay en `.rei/config.json` **sin ejecutarlos**;
- devolver `exit 0` si no hay ningún `FAIL`, `exit 1` si hay al menos uno;
- añadirse a la ayuda (`rei help`).

Fuera de alcance (no hacer): crear archivos o directorios, inicializar git, ejecutar
los checks del proyecto, modificar el comportamiento de otros comandos.

## Archivos

- `internal/doctor/doctor.go` — **nuevo**. Lógica del diagnóstico y reporte.
- `internal/doctor/doctor_test.go` — **nuevo**. Tests del paquete.
- `internal/cli/cli.go` — **modificado**. Despacho del subcomando y ayuda.
- `.rei/specs/2026-10-02_00-14__rei-doctor/plan.md` — este archivo (el Implementer marca los pasos).

No modificar otros archivos.

## Cambios

### `internal/doctor/doctor.go`

Nuevo paquete `doctor` con una única función pública:

```go
func Run(p *paths.Project, out io.Writer) int
```

Devuelve el código de salida (`0` si no hay `FAIL`, `1` en caso contrario) y escribe
el reporte en `out`. **No** debe escribir en disco ni invocar git mutante.

Secciones del reporte, en este orden, usando el mismo estilo de etiquetas que
`internal/check` (`[OK]`, `[WARN]`, `[FAIL]`, `[INFO]`):

1. **Integridad del harness** (solo lectura, `os.Stat`). Cada ausencia es `FAIL`:
   - cada archivo de `check.RequiredFiles` (reutilizar la variable exportada, **no** duplicar la lista);
   - directorios `.rei/specs/` y `.rei/progress/work-items/`;
   - `.rei/progress/current.md` y `.rei/progress/history.md`.
   Si existe, imprimir `[OK] <ruta>`.

2. **Work Items**. Obtener los IDs con `state.ListWorkItems(p)`.
   - Si falla: `FAIL`.
   - Si no hay ninguno: `[INFO] No hay Work Items.`.
   - Para cada item, llamar `validate.WorkItem(p, it.ID)` y:
     - si no devuelve issues: `[OK] <id>`;
     - si devuelve issues: imprimir cada uno como `[<Level>] <id>: <Message>`;
       cada `validate.LevelFail` incrementa el contador global de fallos.

3. **Sesión**. `state.ReadSession(p)`:
   - sesión activa → `[OK] Sesión activa: <id> (estado: <estado>, agente: <agente>)`;
   - sin sesión → `[OK] Sin sesión activa.`.
   Un error de lectura se trata como `WARN` (no debe impedir el diagnóstico).

4. **Git** (solo consulta, nunca `git init`):
   - si `git` no está en el `PATH` (`exec.LookPath`) → `[WARN] git no disponible; ...`;
   - si `gitx.IsRepo(p.Root)` → `[OK] Repositorio git detectado. HEAD=<sha>` usando
     `gitx.HeadSHA(p.Root)`; si el SHA es vacío, reportar `Sin commits todavía.`;
   - si no es repo → `[WARN] Sin repositorio git; ...`.

5. **Checks**. `config.Load(filepath.Join(p.ReiDir(), "config.json"))`:
   - error → `FAIL` con el detalle;
   - sin error → `[INFO] <n> checks configurados en .rei/config.json (no ejecutados).`
     donde `<n>` es `len(cfg.Checks)`. **Nunca** ejecutar los comandos.

6. **Resultado**. Si el contador de fallos es 0:
   `[OK] Sin fallos. Harness íntegro.` y `return 0`.
   Si es > 0: `[FAIL] <n> fallo(s) detectado(s).` y `return 1`.

Imports previstos: `fmt`, `io`, `os`, `os/exec`, `path/filepath` y los paquetes
internos `check`, `config`, `gitx`, `paths`, `state`, `validate`.

### `internal/cli/cli.go`

- Añadir el caso al `switch` de `Run`:
  ```go
  case "doctor":
      return cmdDoctor(args[1:])
  ```
- Implementar `cmdDoctor`:
  - obtener el proyecto con `project()` (si falla, devolver su código);
  - si hay argumentos extra, imprimir `uso: rei doctor` en `stderr` y devolver `2`;
  - devolver `doctor.Run(p, os.Stdout)`.
- Añadir a `printHelp` la línea del comando, en el bloque `Comandos:`:
  `doctor                       Diagnóstico de solo lectura del harness (no modifica nada)`

### `internal/doctor/doctor_test.go`

Tests con `t.TempDir()` y `paths.Project{Root: ...}`, siguiendo el estilo de
`internal/check/check_test.go`:

- un harness completo y válido (archivos de `check.RequiredFiles`, directorios
  requeridos, `current.md`/`history.md`, sin Work Items y sin git) → `Run` devuelve `0`;
- harness incompleto (falta un archivo requerido) → devuelve `1`;
- un Work Item con `meta.json` inválido/incompleto → devuelve `1`.

Los tests no deben depender de que el entorno sea un repositorio git (git solo
produce `WARN`, nunca `FAIL`).

## Restricciones

- **Solo lectura**: prohibido `os.Create`/`os.WriteFile`/`os.MkdirAll`, `git init`,
  `git add` o cualquier comando de git que modifique el repositorio.
- **No ejecutar** los checks de `.rei/config.json`.
- No llamar a `check.Run` (crea archivos e inicializa git); reutilizar únicamente
  `check.RequiredFiles`.
- No duplicar lógica existente: apoyarse en `state`, `validate`, `gitx`, `config`.
- No alterar el comportamiento de los comandos existentes ni su salida.
- Mantener el estilo y los mensajes en español, coherentes con el resto del CLI.
- Sin dependencias externas nuevas (solo stdlib y paquetes internos).

## Pasos

- [x] 1. Crear `internal/doctor/doctor.go` con la firma `Run(p *paths.Project, out io.Writer) int` y el esqueleto de contadores/helpers (`ok`, `warn`, `fail`) siguiendo el estilo de `internal/check/check.go`.
- [x] 2. Implementar la sección de integridad (archivos de `check.RequiredFiles`, directorios requeridos y `current.md`/`history.md`) usando solo `os.Stat` (sin crear nada).
- [x] 3. Implementar la sección de Work Items con `state.ListWorkItems` + `validate.WorkItem`, agregando el contador de fallos.
- [x] 4. Implementar la sección de sesión con `state.ReadSession` (activa / sin sesión).
- [x] 5. Implementar la sección de git con `gitx.IsRepo` y `gitx.HeadSHA` (sin inicializar git) y `exec.LookPath` para detectar la ausencia de `git`.
- [x] 6. Implementar la sección de checks con `config.Load` reportando `len(cfg.Checks)` sin ejecutarlos.
- [x] 7. Implementar el veredicto final (`[OK]`/`[FAIL]`) y devolver `0`/`1` según el contador de fallos.
- [x] 8. En `internal/cli/cli.go`, añadir `case "doctor"`, la función `cmdDoctor` y la línea correspondiente en `printHelp`.
- [x] 9. Crear `internal/doctor/doctor_test.go` con los tres casos descritos.
- [x] 10. Verificar: `make build`, `make test`, `make vet`. Comprobar manualmente que `rei doctor` devuelve exit `0` en el repo actual y que `rei help` lista el comando.
- [x] 11. Marcar cada paso como completado (`[x]`) en este archivo a medida que se termina.
