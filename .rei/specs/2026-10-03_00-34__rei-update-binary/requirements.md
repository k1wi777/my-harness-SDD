# Requirements — `rei update` (auto-actualización del binario)

> Work Item: `2026-10-03_00-34__rei-update-binary`
> Tipo: feature
> Sintaxis: EARS.

## Alcance

Comando `rei update [--check]` que actualiza el binario global de REI a la
última release publicada en GitHub Releases.

Fuera de alcance: firmas criptográficas (cosign), actualización del proyecto
(`rei init --update`) y cualquier mecanismo de auto-actualización silenciosa.

## Requisitos

### R1 — Comando disponible
El sistema DEBE exponer el comando `rei update [--check]`.

### R2 — Consulta de la última release
CUANDO se ejecuta `rei update` o `rei update --check`, el sistema DEBE consultar
la última release publicada del repositorio `k1wi777/my-harness-SDD` y obtener
su identificador de versión (`tag_name`).

### R3 — Comparación de versiones
CUANDO se obtiene la versión de la última release, el sistema DEBE compararla
(semver) con la versión del binario en ejecución.

### R4 — Sin actualización disponible
SI la versión del binario en ejecución es igual o mayor que la de la última
release, ENTONCES el sistema DEBE informarlo y terminar con código de salida 0
sin modificar el binario.

### R5 — Selección del asset
CUANDO existe una versión más reciente, el sistema DEBE seleccionar, de entre
los assets de la release, el correspondiente al sistema operativo y arquitectura
del binario en ejecución, y localizar el archivo de checksums.

### R6 — Descarga y verificación de integridad
CUANDO se selecciona el asset, el sistema DEBE descargar el asset y el archivo
de checksums, y verificar la integridad del asset descargado mediante SHA-256
antes de instalar.

### R7 — Instalación en Unix/macOS
DONDE el sistema operativo es Linux o macOS, CUANDO la verificación de
integridad es correcta, el sistema DEBE reemplazar el binario en ejecución por
la nueva versión y conservar el permiso de ejecución.

### R8 — Windows
DONDE el sistema operativo es Windows, CUANDO existe una versión más reciente,
el sistema DEBE informar al usuario y abrir la página de releases del
repositorio en el navegador, imprimiendo además la URL como alternativa.

### R9 — Modo `--check`
CUANDO se ejecuta `rei update --check`, el sistema DEBE informar del resultado
de la comprobación sin descargar ni reemplazar ningún archivo, en cualquier
sistema operativo, y terminar con código 0 si la comprobación pudo realizarse.

### R10 — Versión `dev`
CUANDO la versión del binario en ejecución es `dev` (compilación local), el
sistema DEBE advertir de que no corresponde a una release publicada y ofrecer
instalar la última versión disponible.

### R11 — Fallo de integridad
SI la verificación SHA-256 del asset descargado falla, ENTONCES el sistema DEBE
abortar la instalación sin modificar el binario e informar del error.

### R12 — Fallo de red o de permisos
SI no se puede consultar la release, descargar el asset, o escribir el binario
(sin red o sin permisos), ENTONCES el sistema DEBE informar del problema,
terminar con código distinto de cero y no dejar el binario en un estado
corrupto ni degradado.

### R13 — Restauración ante sustitución fallida
SI el reemplazo del binario falla una vez iniciado, ENTONCES el sistema DEBE
restaurar el binario original.

### R14 — Ayuda
El sistema DEBE documentar `rei update [--check]` en la ayuda general
(`rei help`) y en la ayuda específica del comando (`rei update --help`).

### R15 — Uso incorrecto
SI se pasan argumentos distintos de `--check`, ENTONCES el sistema DEBE mostrar
el uso del comando y terminar con código de salida 2.
