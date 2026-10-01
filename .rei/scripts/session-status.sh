#!/usr/bin/env bash
#
# .rei/scripts/session-status.sh
#
# Indica si existe una sesión activa y en qué estado se encuentra.
# Solo lectura: nunca modifica archivos.
#
# Salida:
#   - "Sin sesión activa."                       (exit 0)
#   - "Sesión activa: <id> (estado: ..., ...)"   (exit 1)

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

if [[ ! -f "$CURRENT_FILE" ]]; then
    echo "Sin sesión activa (no existe .rei/progress/current.md)."
    exit 0
fi

field() {
    grep -m1 -E "^- \*\*$1:\*\*" "$CURRENT_FILE" \
        | sed -E "s/^- \*\*$1:\*\*[[:space:]]*//"
}

work_item="$(field 'Work Item')"
estado="$(field 'Estado')"
agente="$(field 'Agente activo')"

if [[ -z "$work_item" || "$work_item" == "_ninguno_" ]]; then
    echo "Sin sesión activa."
    exit 0
fi

echo "Sesión activa: ${work_item} (estado: ${estado:-?}, agente: ${agente:-?})"
exit 1
