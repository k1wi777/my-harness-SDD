#!/usr/bin/env bash
#
# .rei/scripts/start-session.sh
#
# Inicia la sesión actual en .rei/progress/current.md a partir de la plantilla.
#
# Uso:
#   start-session.sh <work-item-id> <feature|task>

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

if (( $# < 2 )); then
    echo "Uso: start-session.sh <work-item-id> <feature|task>" >&2
    exit 1
fi

id="$1"
type="$2"
inicio="$(date -Iseconds)"

mkdir -p .rei/progress

if ! render_template "$TEMPLATES_DIR/current.md" "$CURRENT_FILE" \
    "WORK_ITEM=$id" \
    "TYPE=$type" \
    "STATE=pending" \
    "START=$inicio" \
    "AGENT=spec_author" \
    "PLAN=$REI_EMPTY" \
    "LOG=$REI_EMPTY" \
    "NEXT=Planificación pendiente."; then
    echo "ERROR: no se pudo iniciar la sesión." >&2
    exit 1
fi

echo "Sesión iniciada: $id ($type) -> pending (agente: spec_author)"
