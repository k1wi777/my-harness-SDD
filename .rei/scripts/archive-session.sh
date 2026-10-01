#!/usr/bin/env bash
#
# .rei/scripts/archive-session.sh
#
# Cierra la sesión activa:
#   1. Añade el contenido de .rei/progress/current.md al final de history.md.
#   2. Restablece current.md a su estado inicial.
#
# history.md es append-only: nunca se sobrescribe.

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

if [[ ! -f "$CURRENT_FILE" ]]; then
    echo "No hay sesión que archivar."
    exit 0
fi

work_item="$(grep -m1 -E '^- \*\*Work Item:\*\*' "$CURRENT_FILE" \
    | sed -E 's/^- \*\*Work Item:\*\*[[:space:]]*//')"

if [[ -z "$work_item" || "$work_item" == "_ninguno_" ]]; then
    echo "No hay sesión activa que archivar."
    exit 0
fi

mkdir -p "$(dirname "$HISTORY_FILE")"

{
    printf '\n## %s — %s\n\n' "$(date '+%Y-%m-%d %H:%M')" "$work_item"
    cat "$CURRENT_FILE"
    printf '\n'
} >> "$HISTORY_FILE"

if ! reset_current; then
    echo "ERROR: no se pudo restablecer current.md tras archivar." >&2
    exit 1
fi

echo "Sesión '$work_item' archivada en history.md. current.md restablecido."
