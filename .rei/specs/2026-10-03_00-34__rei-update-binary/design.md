# Design — `rei update` (auto-actualización del binario)

> Work Item: `2026-10-03_00-34__rei-update-binary`
> Tipo: feature

## Estrategia de implementación

Se añade un paquete nuevo `internal/update` con la lógica de actualización,
aislado de `internal/cli` para poder probarlo sin red ni sistema operativo real.
`internal/cli` solo se encarga del despacho, la validación de argumentos y la
presentación del resultado, coherente con el resto de comandos del proyecto.

Flujo de `rei update`:

```
cli.cmdUpdate(rest)
  └─ valida args (--check; otro -> 2)
     └─ update.Run(Options{CurrentVersion: version, Check: <bool>, Stdout: os.Stdout})
        ├─ version == dev  -> avisar + ofrecer última release (R10)
        ├─ GET releases/latest -> tag_name + assets (R2)
        ├─ comparar semver (R3)
        ├─ sin actualización -> informar, 0 (R4)
        ├─ --check -> informar, 0 (R9)
        ├─ Windows -> informar + openBrowser + URL, 0 (R8)
        └─ Linux/macOS
           ├─ elegir asset <os>_<arch> + checksums.txt (R5)
           ├─ descargar asset + checksums (R6)
           ├─ verificar sha256 -> si falla, abortar sin tocar el binario (R11)
           ├─ extraer binario del archive (R6)
           └─ reemplazar os.Executable(): temp + chmod + rename con backup (R7, R13)
```

`update.Run` devuelve un `int` (código de salida) y escribe por el `io.Writer`
recibido; los errores se traducen a mensajes accionables y a códigos
0/1/2, nunca a panic.

## Archivos involucrados

| Archivo | Cambio |
|---------|--------|
| `internal/update/update.go` | Nuevo. Orquestación, `Options`, `Run`, ramas dev/check/Windows/up-to-date y códigos de salida. |
| `internal/update/github.go` | Nuevo. Cliente de la API de GitHub: `GET /repos/k1wi777/my-harness-SDD/releases/latest`, tipos `release`/`asset`, selección de asset y de `checksums.txt`. |
| `internal/update/semver.go` | Nuevo. Normalización (`v` opcional) y comparación semver `major.minor.patch`; detección de `dev`. |
| `internal/update/archive.go` | Nuevo. Descarga con timeout, verificación SHA-256 contra `checksums.txt` y extracción de `.tar.gz`/`.zip`. |
| `internal/update/replace.go` | Nuevo (efectivo en Unix/macOS). Comprobación de escritura, escritura atómica y backup/restauración. |
| `internal/update/browser.go` | Nuevo. `openBrowser` best-effort (`xdg-open`/`open`/`cmd /c start`). |
| `internal/update/*_test.go` | Nuevos. Tests con `httptest` y `t.TempDir()`. |
| `internal/cli/cli.go` | Añadir `case "update"` y `cmdUpdate(rest)`. |
| `internal/cli/help.go` | Añadir la entrada `update` a la tabla `commands`. |

No se modifican `go.mod` ni dependencias: se usa solo la librería estándar.

## Componentes y decisiones técnicas

### Paquete `internal/update`

Tipos y costuras para poder testear sin red ni reemplazar el binario real:

```go
type Options struct {
    CurrentVersion string
    Check          bool
    Stdout         io.Writer
    // Costuras inyectables (con valores por defecto de producción):
    Client     *http.Client   // con timeout
    Executable func() (string, error) // os.Executable
    GOOS       string          // runtime.GOOS
    GOARCH     string          // runtime.GOARCH
    OpenBrowser func(url string) error
}
```

`Run(opts Options) int` aplica valores por defecto a los campos vacíos. Los
helpers de red/archivo/OS reciben estas dependencias, de modo que los tests
usen `httptest.NewServer` y un `Executable` apuntando a `t.TempDir()`.

### Fuente de la versión actual

La versión vive en `internal/cli` (`var version = "dev"`, inyectada con
`-ldflags`). Se mantiene sin exportar y `cmdUpdate` pasa su valor a
`update.Run`; así se evita un ciclo de importación y sigue habiendo una única
fuente de verdad. En el helper de CLI se importa `internal/update`.

### Consulta de la release y selección de asset

- Endpoint: `https://api.github.com/repos/k1wi777/my-harness-SDD/releases/latest`.
- Se envía cabecera `User-Agent` (GitHub la exige) y `Accept: application/vnd.github+json`.
- Se decodifica únicamente lo necesario: `tag_name` y `assets[].{name,browser_download_url}`.
- El asset se construye con el mismo esquema que goreleaser:
  - Unix/macOS: `rei_<version>_<os>_<arch>.tar.gz`
  - Windows: `rei_<version>_<os>_<arch>.zip`
  - checksums: `checksums.txt`
  - `<version>` es el `tag_name` sin la `v` inicial (p. ej. tag `v0.1.0` → `0.1.0`);
    `<os>`/`<arch>` son `linux`/`darwin`/`windows` y `amd64`/`arm64`.
- Si el asset exacto no aparece en la release, se degrada con aviso y se
  sugiere la URL de releases (no se inventa un nombre alternativo).

### Comparación de versiones

Implementación propia mínima (sin dependencias): parsea `v?major.minor.patch`,
compara numéricamente campo a campo. Una versión `dev` o no parseable se trata
como "siempre actualizable" (ver rama `dev`). No se comparan pre-releases
(las releases de goreleaser son `vX.Y.Z`).

### Descarga y verificación

- `http.Client` con timeout (p. ej. 30 s) y `io.Copy` a archivo temporal.
- `checksums.txt` (formato goreleaser: `<sha256>  <nombre>`) se descarga y se
  busca la línea cuyo nombre coincide con el asset. Si el nombre no está,
  es error (R11/R12).
- Se calcula `sha256` del archivo descargado y se compara en tiempo constante
  (`crypto/subtle.ConstantTimeCompare` o `hmac.Equal`). Si no coincide: borrar
  el temporal, informar y salir 1 sin tocar el binario.

### Extracción

- `.tar.gz`: `compress/gzip` + `archive/tar`.
- `.zip`: `archive/zip`.
- Solo se extrae la entrada cuyo nombre base es `rei` (o `rei.exe` en Windows).
- Protección contra path traversal: se ignoran entradas con rutas absolutas o
  con `..`; nunca se escribe fuera del directorio temporal.
- Se copia el contenido al archivo temporal de destino, no al binario final.

### Reemplazo en Unix/macOS

1. Localizar el destino con `os.Executable()`; aplicar `filepath.EvalSymlinks`
   best-effort para soportar instalaciones enlazadas desde el `PATH`.
2. Pre-comprobar escritura creando el temporal en **el mismo directorio** que el
   destino (necesario para que `os.Rename` sea atómico en el mismo filesystem).
   Si falla, degradar y salir 1 sin modificar nada (R12).
3. Escribir el binario nuevo en ese temporal con modo `0600` y luego `chmod 0755`.
4. `os.Rename(destino, destino+".bak")`; `os.Rename(tmp, destino)`.
   - Si el segundo rename falla, restaurar el backup (R13) y borrar el temporal.
   - Si tiene éxito, borrar el backup pesimistamente-ignorado (best-effort).
5. El binario nuevo queda en la ruta de `os.Executable()`, por lo que la
   siguiente ejecución usa la versión actualizada.

### Windows y apertura del navegador

No se intenta reemplazar un `.exe` en ejecución: la rama Windows (R8) solo
informa e intenta abrir el navegador. `openBrowser(url)` es best-effort y
prueba, según `GOOS`:

- `windows`: `cmd /c start "" <url>`
- `darwin`: `open <url>`
- resto (linux): `xdg-open <url>`

Cualquier fallo del comando se ignora; la URL se imprime siempre por stdout como
alternativa, de modo que la función nunca rompe el comando.

### `--check` y rama `dev`

- `--check`: imprime versión actual, versión de la última release y si hay
  actualización; no descarga, no escribe. Sale 0 si la consulta funcionó
  (independientemente de si hay versión nueva) y 1 si no pudo consultar.
- `dev`: avisa de que no es una release publicada, muestra la última versión
  disponible y continúa con el flujo normal de instalación (ofrece instalarla),
  respetando `--check` y la rama Windows.

### Seguridad y límites

- Confianza en el canal de release: el `checksums.txt` se descarga del mismo
  release. Las firmas (cosign) quedan fuera de alcance por decisión explícita.
- TLS por defecto de `net/http`; sin `InsecureSkipVerify`.
- Archivos temporales con permisos restrictivos; limpieza con `defer` en
  cualquier salida.
- Sin escritura en el binario hasta haber verificado la integridad.
- Límite/prudencia: no se descomprimen entradas ajenas al binario esperado.

### Códigos de salida

| Situación | Código |
|-----------|--------|
| Actualizado, ya al día, `--check` correcto, Windows guiado | 0 |
| Red/permisos/verificación/extracción fallidos | 1 |
| Argumentos inválidos | 2 |

## Alternativas descartadas

- **Librería semver externa** (`Masterminds/semver`, `golang.org/x/mod/semver`):
  descartada para no introducir dependencias; el formato de releases es simple.
- **`go install ...@latest` / recompilar**: descartado porque los usuarios del
  binario distribuido no tienen Go instalado.
- **Auto-reemplazo en Windows con borrado diferido / re-exec**: descartado por
  frágil; se opta por guiar al usuario a la página de releases.
- **Descargar en `os.TempDir()` y renombrar al destino**: descartado porque
  `os.Rename` entre filesystems no es atómico; el temporal se crea junto al
  destino.
- **Prompt interactivo de confirmación**: descartado; el comando es explícito y
  no debe bloquearse en contextos no interactivos.
- **Actualización silenciosa al arrancar cualquier comando**: descartada por
  transparencia y control del usuario.
