// Package rei expone el esqueleto de REI Harness embebido en el binario.
//
// El paquete raíz no tiene dependencias: solo declara el sistema de archivos
// embebido que consume `internal/initwizard` para desplegar el esqueleto con
// `rei init`. Debe vivir en la raíz del módulo porque `go:embed` no admite
// patrones con `..`.
package rei

import "embed"

// Skeleton contiene el esqueleto desplegable por `rei init`: AGENTS.md,
// .rei/docs/**, .rei/agents/**, .rei/templates/** y .rei/config.json.
//
// El patrón selectivo con prefijo `all:` incluye los subárboles completos bajo
// `.rei/` (incluidos nombres que empiezan por `.`) y, a la vez, excluye
// deliberadamente .rei/specs/ y .rei/progress/ (estado del proyecto destino,
// no del arnés).
//
//go:embed all:.rei/docs all:.rei/agents all:.rei/templates all:.rei/config.json AGENTS.md
var Skeleton embed.FS
