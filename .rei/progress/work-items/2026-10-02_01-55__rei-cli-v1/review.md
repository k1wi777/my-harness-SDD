# Revisión — `rei cli v1`

- **Work Item:** `2026-10-02_01-55__rei-cli-v1`
- **Tipo:** `task`
- **Estado final:** `done`
- **Base revisada:** `e4680010c952f682b6830408be4396b4b5864665`
- **Fecha:** 2026-10-02

## Alcance revisado

Paquete de `rei review-diff 2026-10-02_01-55__rei-cli-v1` (exit 0):

- `internal/check/check.go` — modificado (+33/-2: `runCheck` → `Exec`; nuevo `RunDeclared`).
- `internal/check/check_test.go` — modificado (+75: `TestRunDeclaredSinChecks`, `ChecksVacios`, `CheckOK`, `CheckFail`, `SinEfectos`).
- `internal/cli/cli.go` — modificado (+91/-40: dispatch de `--help`/`-h`, `help <cmd>`, `test`, `item`, `cmdTest`, `cmdItem`; `printHelp` retirado).
- `internal/cli/help.go` — nuevo (+172: tabla `commands`, `helpText`, `printHelp`, `commandHelpText`, `printCommandHelp`).
- `internal/cli/help_test.go` — nuevo (+53).
- `internal/show/show.go` — nuevo (+121: `WorkItem`, `docPath`, `present`, `orEmpty`).
- `internal/show/show_test.go` — nuevo (+118: 4 casos).
- `package-lock.json` — sin rastrear, preexistente (sep 10) y ajeno a este Work Item.

## Objetivo cumplido

Las tres capacidades de la CLI v1 están implementadas según `plan.md`:

1. **`rei item show <id>`** — `internal/show.WorkItem` imprime metadatos
   (`id`, `title`, `type`, `status`, `created_at`, `base_commit`,
   `last_review_commit`, con vacíos como `(vacío)`), los siete documentos en
   orden fijo con `sí`/`no`, la sesión activa (estado y agente) y los hallazgos
   de `validate.WorkItem`. Work Item inexistente o `meta.json` ilegible →
   `ERROR` por `stderr` y `1`; ≥1 `LevelFail` → `Resultado: FALLO` y `1`.
2. **`rei test`** — `check.RunDeclared` usa solo `config.Load` + `Exec`: config
   inválida → `[FAIL]` y `1`; sin checks → `[INFO]` y `0`; por check `[OK]`/
   `[FAIL]` y `exit 1` si alguno falla. No llama a `check.Run`. La ayuda aclara
   que **no** es la suite de REI (`make test`).
3. **Ayuda por comando** — `internal/cli/help.go` centraliza la tabla única
   (13 comandos, incluidos `test` e `item`); `rei <cmd> --help`/`-h` devuelve `0`
   para comandos conocidos y `2` para desconocidos; `rei help <cmd>` funciona;
   `rei help` general lista `test` e `item`.

## Restricciones respetadas

- **`item show` es solo lectura:** `show.go` solo usa `os.Stat`, `meta.Load`,
  `state.ReadSession` y `validate.WorkItem`; no hay `os.WriteFile`,
  `os.MkdirAll`, git mutante ni ejecución de checks. Confirmado por `grep` en
  `internal/show`, `internal/state` y `internal/validate` (solo aparecen
  escrituras en archivos `_test.go`).
- **`test` sin efectos secundarios ni `check.Run`:** `RunDeclared` solo invoca
  `config.Load` + `Exec`; el diff no contiene `os.MkdirAll`, `template.*`,
  `gitx.*`, `RequiredFiles`, `validate.*` ni `check.Run`. Test unitario
  `TestRunDeclaredSinEfectos` verifica que no crea `.rei/specs/` ni
  `current.md`.
- **Comandos existentes sin cambios de comportamiento:** el refactor es
  `runCheck` → `Exec` (misma semántica y llamada desde `Run`); el único cambio
  visible es la ayuda general, explícitamente permitida en el plan.
- **Sin dependencias externas nuevas** (stdlib + paquetes internos).
- **Mensajes y ayuda en español**, coherentes con el CLI.
- **`--help`/`-h` solo tras el comando:** `rei check --quiet` sigue funcionando
  (verificado, exit 0).

## Verificaciones (checkpoints y comandos solicitados)

`.rei/config.json` no declara checkpoints (`"checks": []`), por lo que el
checkpoint rápido efectivo es `rei check`; las validaciones del `plan.md`
(§Validaciones) son la referencia.

| Verificación | Comando | Resultado |
|--------------|---------|-----------|
| Formato | `gofmt -l internal` | sin salida ✅ |
| Vet | `make vet` | exit `0` ✅ |
| Tests | `make test` | exit `0`; `internal/show`, `internal/check`, `internal/cli` ok ✅ |
| Checkpoint rápido | `rei check --quiet` | exit `0` ✅ |

### Verificación independiente

**(a) `item show`**

```text
rei item show 2026-10-02_01-55__rei-cli-v1  -> exit 0
  metadatos: title/type/status/created_at/base_commit correctos; last_review_commit (vacío)
  Documentos: plan.md sí; requirements/design/tasks no; impl.md sí; review/spec no
  Sesión: activa (estado: review, agente: implementer)
  Validación: (sin hallazgos) / Resultado: OK
rei item show no-existe                      -> exit 1
  ERROR: el Work Item "no-existe" no existe
```

**(b) `rei test`** (config original `{ "checks": [] }`, respaldado y restaurado)

```text
rei test (vacío)          -> exit 0  [INFO]  Sin checks configurados (.rei/config.json). Nada que ejecutar.
rei test (check temporal) -> exit 1  [FAIL]  check temporal de revision
rei test (restaurado)     -> exit 0  [INFO]  Sin checks configurados (.rei/config.json). Nada que ejecutar.
```

`.rei/config.json` quedó idéntico al original (`{ "checks": [] }`) tras la
prueba.

**(c) Sin efectos secundarios:** instantánea de `find .rei -type f` antes y
después de toda la secuencia de `rei test` → sin archivos nuevos ni eliminados;
`.rei/specs/` y `current.md` no fueron creados ni modificados.

**(d) Ayuda**

```text
rei <cmd> --help  para check, doctor, new, session, items, status, commit,
  validate, review-diff, test, item, version, help  -> todos exit 0
rei no-existe --help                                -> exit 2
rei help test / rei help item                       -> exit 0 (contenido correcto)
rei help                                            -> exit 0; lista test e item
rei -h                                              -> exit 0
```

**(e) Regresión:** `rei check --quiet` → exit `0`.

## Observaciones

- `printCommandHelp` muestra la ayuda general ante un comando desconocido
  (`exit 2`), consistente con el `default` del dispatch.
- `cmdTest`/`cmdItem` devuelven `2` ante uso incorrecto, según los códigos de
  salida de `plan.md`.
- Los tests de `RunDeclared` asumen entorno POSIX (`sh -c`); el plan lo previó y
  el cross-build a Windows no ejecuta tests.
- `package-lock.json` figura como cambio sin rastrear ajeno al Work Item; no se
  le atribuye impacto.

## Acciones requeridas

Ninguna.

## Resultado

Aprobado. El objetivo se cumple, las restricciones se respetan, la verificación
independiente confirma el comportamiento esperado y todas las verificaciones
(`gofmt`, `make vet`, `make test`, `rei check --quiet`) pasan.
