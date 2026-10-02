#!/usr/bin/env bash
#
# .rei/scripts/set-commit.sh
#
# Registra en meta.json el commit actual de git.
#
# Uso:
#   set-commit.sh <work-item-id> <base_commit|last_review_commit>
#
#   base_commit         — punto de partida de la implementación (al aprobar).
#   last_review_commit  — punto revisado (al rechazar), para que la próxima
#                         revisión vea solo los cambios solicitados.
#
# Sin git o sin commits, no hace nada y termina con 0 (el review por diff
# simplemente se deshabilita y el Reviewer usa el modo lectura).

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

if (( $# < 2 )); then
    echo "Uso: set-commit.sh <work-item-id> <base_commit|last_review_commit>" >&2
    exit 1
fi

id="$1"
field="$2"

case "$field" in
    base_commit|last_review_commit) ;;
    *)
        echo "ERROR: campo inválido '$field'." >&2
        exit 1
        ;;
esac

meta="$SPECS_DIR/$id/meta.json"
if [[ ! -f "$meta" ]]; then
    echo "ERROR: no existe $meta" >&2
    exit 1
fi

if ! command -v git >/dev/null 2>&1 || ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "Sin git: no se registra '$field'."
    exit 0
fi

sha="$(git rev-parse HEAD 2>/dev/null || true)"
if [[ -z "$sha" ]]; then
    echo "Sin commits todavía: no se registra '$field'."
    exit 0
fi

if meta_set "$meta" "$field" "$sha"; then
    echo "$field=$sha"
else
    echo "ERROR: no se pudo actualizar '$field'." >&2
    exit 1
fi
