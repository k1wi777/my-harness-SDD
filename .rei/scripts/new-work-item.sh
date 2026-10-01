#!/usr/bin/env bash
#
# .rei/scripts/new-work-item.sh
#
# Crea la estructura de un Work Item nuevo:
#   .rei/specs/<id>/meta.json
#   .rei/progress/work-items/<id>/
#
# Uso:
#   new-work-item.sh <work-item-id> <feature|task> [title]
#
# El `description` queda vacío: el Leader debe completarlo con `Edit`.

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

if (( $# < 2 )); then
    echo "Uso: new-work-item.sh <work-item-id> <feature|task> [title]" >&2
    exit 1
fi

id="$1"
type="$2"
title="${3:-}"

case "$type" in
    feature|task) ;;
    *)
        echo "ERROR: type debe ser 'feature' o 'task' (recibido: '$type')." >&2
        exit 1
        ;;
esac

if [[ ! "$id" =~ ^[0-9]{4}-[0-9]{2}-[0-9]{2}_[0-9]{2}-[0-9]{2}__[a-z0-9]+(-[a-z0-9]+)*$ ]]; then
    echo "ERROR: id inválido '$id'. Formato: YYYY-MM-DD_HH-mm__slug-en-kebab-case." >&2
    exit 1
fi

if [[ -e "$SPECS_DIR/$id" || -e "$WORK_ITEMS_DIR/$id" ]]; then
    echo "ERROR: ya existe un Work Item con id '$id'." >&2
    exit 1
fi

# created_at se deriva del propio id para que coincida con su prefijo.
prefix="${id%%__*}"
date_part="${prefix%%_*}"
time_part="${prefix##*_}"
time_part="${time_part/-/:}"
tz="$(date +%:z)"
created_at="${date_part}T${time_part}:00${tz}"

json_escape() {
    printf '%s' "$1" | sed 's/\\/\\\\/g; s/"/\\"/g'
}
title_esc="$(json_escape "$title")"

mkdir -p "$SPECS_DIR/$id" "$WORK_ITEMS_DIR/$id"

if ! render_template "$TEMPLATES_DIR/meta.json" "$SPECS_DIR/$id/meta.json" \
    "ID=$id" \
    "TITLE=$title_esc" \
    "DESCRIPTION=" \
    "TYPE=$type" \
    "STATUS=pending" \
    "CREATED_AT=$created_at"; then
    echo "ERROR: no se pudo crear meta.json." >&2
    exit 1
fi

echo "Work Item creado: $SPECS_DIR/$id/ (type: $type, status: pending)"
echo "Completa 'description' en $SPECS_DIR/$id/meta.json."
