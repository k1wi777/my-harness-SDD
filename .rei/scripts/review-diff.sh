#!/usr/bin/env bash
#
# .rei/scripts/review-diff.sh
#
# Genera el "paquete de revisión" de un Work Item a partir de git:
#   - resumen (git diff --stat)
#   - lista de archivos modificados / añadidos / eliminados / sin rastrear
#   - con --full, el diff completo
#
# Uso:
#   review-diff.sh <work-item-id> [--full]
#
# Códigos de salida:
#   0  paquete generado
#   2  sin git o sin base → el Reviewer debe caer al modo lectura

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

if (( $# < 1 )); then
    echo "Uso: review-diff.sh <work-item-id> [--full]" >&2
    exit 1
fi

id="$1"
shift

full=0
for arg in "$@"; do
    case "$arg" in
        --full) full=1 ;;
    esac
done

meta="$SPECS_DIR/$id/meta.json"
if [[ ! -f "$meta" ]]; then
    echo "ERROR: no existe $meta" >&2
    exit 1
fi

if ! command -v git >/dev/null 2>&1 || ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
    echo "SIN_GIT"
    exit 2
fi

base_commit="$(meta_field "$meta" base_commit)"
last_review_commit="$(meta_field "$meta" last_review_commit)"

base=""
origin="base_commit"
if [[ -n "$last_review_commit" ]]; then
    base="$last_review_commit"
    origin="last_review_commit"
elif [[ -n "$base_commit" ]]; then
    base="$base_commit"
fi

if [[ -z "$base" ]] || ! git cat-file -e "${base}^{commit}" 2>/dev/null; then
    echo "SIN_BASE"
    exit 2
fi

# Filtra el ruido del harness: conserva el código del proyecto y la spec del
# Work Item; descarta el resto de `.rei/` (progress, otros specs, docs...).
filter_paths() {
    local p
    while IFS= read -r p; do
        [[ -z "$p" ]] && continue
        case "$p" in
            ".rei/specs/$id/"*) printf '%s\n' "$p" ;;
            .rei/*) ;;
            *) printf '%s\n' "$p" ;;
        esac
    done
}

tracked=()
while IFS= read -r p; do
    tracked+=("$p")
done < <(git diff --name-only "$base" 2>/dev/null | filter_paths)

untracked=()
while IFS= read -r p; do
    untracked+=("$p")
done < <(git ls-files --others --exclude-standard 2>/dev/null | filter_paths)

echo "== Paquete de revisión: $id =="
echo "Base: $base ($origin)"
echo

echo "--- Resumen ---"
if (( ${#tracked[@]} > 0 )); then
    git diff --stat "$base" -- "${tracked[@]}"
else
    echo "(sin cambios en archivos rastreados)"
fi

if (( ${#untracked[@]} > 0 )); then
    echo
    echo "--- Nuevos sin rastrear ---"
    printf '  %s\n' "${untracked[@]}"
fi

if (( full )); then
    echo
    echo "--- Diff completo ---"
    if (( ${#tracked[@]} > 0 )); then
        git diff "$base" -- "${tracked[@]}"
    fi
    for p in "${untracked[@]}"; do
        echo
        echo "--- nuevo: $p ---"
        git diff --no-index -- /dev/null "$p" 2>/dev/null || true
    done
fi

exit 0
