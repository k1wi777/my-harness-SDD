#!/usr/bin/env bash
#
# .rei/scripts/work-items-status.sh
#
# Recorre todos los .rei/specs/*/meta.json y lista su estado.
# Marca el Work Item que esté en `in_progress`.
#
# Salida:
#   - Lista "id [status]" (exit 0 si no hay ninguno en in_progress).
#   - Añade "AVISO" y devuelve exit 1 si existe algún in_progress.

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

shopt -s nullglob
files=("$SPECS_DIR"/*/meta.json)

if (( ${#files[@]} == 0 )); then
    echo "No hay Work Items registrados."
    exit 0
fi

extract() {
    grep -m1 -o "\"$2\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" "$1" 2>/dev/null \
        | sed -E "s/.*\"$2\"[[:space:]]*:[[:space:]]*\"//; s/\"$//"
}

active=0
echo "Work Items:"
for f in "${files[@]}"; do
    id="$(extract "$f" id)"
    [[ -z "$id" ]] && id="$(basename "$(dirname "$f")")"

    status="$(extract "$f" status)"
    [[ -z "$status" ]] && status="desconocido"

    marker=""
    if [[ "$status" == "in_progress" ]]; then
        marker="   << ACTIVO"
        active=$((active + 1))
    fi

    printf '  - %s [%s]%s\n' "$id" "$status" "$marker"
done

if (( active > 0 )); then
    echo "AVISO: existe un Work Item en in_progress."
    exit 1
fi

exit 0
