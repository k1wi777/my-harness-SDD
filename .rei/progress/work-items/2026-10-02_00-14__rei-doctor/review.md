# Revisión — rei doctor

- **Work Item:** `2026-10-02_00-14__rei-doctor`
- **Tipo:** `task`
- **Estado final:** `done`
- **Base revisada:** `25bfe9c1f34fa1536b8a20d37c9973bde2ed80d9`
- **Fecha:** 2026-10-02

## Alcance revisado

Paquete de `rei review-diff 2026-10-02_00-14__rei-doctor` (exit 0):

- `internal/cli/cli.go` — modificado (16 inserciones: import, `case "doctor"`, `cmdDoctor`, línea en `printHelp`).
- `internal/doctor/doctor.go` — nuevo (sin rastrear en el paquete).
- `internal/doctor/doctor_test.go` — nuevo (sin rastrear en el paquete).
- `.rei/specs/2026-10-02_00-14__rei-doctor/{meta.json,plan.md}` — nuevos (spec).
- `package-lock.json` — sin rastrear, preexistente y ajeno a este Work Item.

## Objetivo cumplido

`internal/doctor/doctor.go` implementa `Run(p *paths.Project, out io.Writer) int` con las
seis secciones planificadas:

1. **Integridad** — recorre `check.RequiredFiles`, directorios `.rei/specs/` y
   `.rei/progress/work-items/`, y `current.md`/`history.md` con solo `os.Stat`.
2. **Work Items** — `state.ListWorkItems` + `validate.WorkItem`; los `LevelFail`
   incrementan el contador.
3. **Sesión** — `state.ReadSession` (activa / sin sesión).
4. **Git** — `exec.LookPath("git")`, `gitx.IsRepo` y `gitx.HeadSHA`; nunca inicializa git.
5. **Checks** — `config.Load` reporta `len(cfg.Checks)` sin ejecutarlos.
6. **Resultado** — `0` sin fallos, `1` con al menos uno.

`internal/cli/cli.go` añade el `case "doctor"`, `cmdDoctor` (uso en `stderr` + código 2
ante argumentos extra) y la línea correspondiente en `rei help`. El objetivo del plan
queda cumplido sin desviaciones de alcance.

## Restricciones respetadas

- **Solo lectura:** `internal/doctor/doctor.go` no contiene `os.WriteFile`, `os.MkdirAll`,
  `os.Create`, `git init` ni `check.Run` (grep sin coincidencias). Solo usa `os.Stat` y
  consultas de git (`IsRepo`/`HeadSHA`).
- **No ejecuta los checks** de `.rei/config.json`; únicamente informa su número.
- **No duplica lógica** existente: reutiliza `check.RequiredFiles`, `state`, `validate`,
  `gitx` y `config`.
- Sin dependencias externas nuevas (stdlib + paquetes internos).
- No se alteró el comportamiento de otros comandos.

## Verificaciones (V1/V2 y comandos solicitados)

| Verificación | Comando | Resultado |
|--------------|---------|-----------|
| Checkpoint rápido (`rei check`) | `rei check --quiet` | exit `0` ✅ |
| Build | `make build` | exit `0` ✅ |
| Tests | `make test` | exit `0`; `internal/doctor` ok ✅ |
| Vet | `make vet` | exit `0` ✅ |
| Smoke manual | `rei doctor` | exit `0`, `[OK] Sin fallos. Harness íntegro.` ✅ |
| Ayuda | `rei help` | lista `doctor` ✅ |

`rei doctor` en el repo actual reporta integridad OK (archivos, directorios,
`current.md`/`history.md`), Work Item `2026-10-02_00-14__rei-doctor` OK, sesión activa
(`review`, agente `implementer`), git detectado (`HEAD=25bfe9c1...`) y `0 checks`
configurados (no ejecutados).

## Observaciones

- `internal/doctor/doctor_test.go` usa `os.WriteFile`/`os.MkdirAll`, pero solo para crear
  el harness temporal de las pruebas (`t.TempDir()`), tal como contempla `plan.md`; el
  código de producción (`doctor.go`) permanece de solo lectura.
- `package-lock.json` figura como cambio sin rastrear ajeno al Work Item; no se le
  atribuye impacto.

## Acciones requeridas

Ninguna.

## Resultado

Aprobado. El objetivo se cumple, las restricciones se respetan y todas las
verificaciones pasan.
