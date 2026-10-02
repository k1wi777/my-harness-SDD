#!/usr/bin/env bash
#
# .rei/scripts/validate.sh [<work-item-id>]
#
# Comprueba, solo con código, la consistencia interna de un Work Item:
#   - campos de meta.json
#   - documentos requeridos según estado/tipo
#   - evidencia (impl.md / review.md)
#   - coherencia entre meta.json y current.md
#
# Sin id, valida el Work Item de la sesión activa (si existe).
# No consume tokens. Salida: lista de [FAIL]/[WARN]; exit 1 si hay algún FAIL.

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

VALID_STATUS="pending ready in_progress review done blocked changes_requested"

issues=0
fail() { printf '[FAIL] %s\n' "$1"; issues=1; }
warn() { printf '[WARN] %s\n' "$1"; }

current_work_item() {
    [[ -f "$CURRENT_FILE" ]] || return 1
    grep -m1 -E '^- \*\*Work Item:\*\*' "$CURRENT_FILE" \
        | sed -E 's/^- \*\*Work Item:\*\*[[:space:]]*//'
}

current_state() {
    [[ -f "$CURRENT_FILE" ]] || return 1
    grep -m1 -E '^- \*\*Estado:\*\*' "$CURRENT_FILE" \
        | sed -E 's/^- \*\*Estado:\*\*[[:space:]]*//'
}

id="${1:-}"
if [[ -z "$id" ]]; then
    id="$(current_work_item || true)"
    if [[ -z "$id" || "$id" == "_ninguno_" ]]; then
        echo "Sin sesión activa."
        exit 0
    fi
fi

echo "== Validando Work Item: $id =="

meta="$SPECS_DIR/$id/meta.json"
if [[ ! -f "$meta" ]]; then
    fail "no existe $meta"
    echo
    echo "Resultado: FALLO"
    exit 1
fi

# --- campos obligatorios ---
for key in id title description type status created_at; do
    if [[ -z "$(meta_field "$meta" "$key")" ]]; then
        fail "campo vacío en meta.json: $key"
    fi
done

m_id="$(meta_field "$meta" id)"
m_type="$(meta_field "$meta" type)"
m_status="$(meta_field "$meta" status)"
m_created="$(meta_field "$meta" created_at)"

[[ "$m_id" == "$id" ]] || fail "id de meta.json ('$m_id') no coincide con la carpeta ('$id')"

case "$m_type" in
    feature|task) ;;
    *) fail "type inválido: '$m_type'" ;;
esac

case " $VALID_STATUS " in
    *" $m_status "*) ;;
    *) fail "status inválido: '$m_status'" ;;
esac

# created_at coherente con el prefijo del id
prefix="${id%%__*}"
date_part="${prefix%%_*}"
time_part="${prefix##*_}"
time_part="${time_part/-/:}"
case "$m_created" in
    "${date_part}T${time_part}"*) ;;
    *) warn "created_at ('$m_created') no coincide con el prefijo del id ('${date_part}T${time_part}...')" ;;
esac

# --- estructura ---
[[ -d "$WORK_ITEMS_DIR/$id" ]] || warn "no existe $WORK_ITEMS_DIR/$id"

# --- documentos de planificación según estado ---
needs_planning=0
case "$m_status" in
    ready|in_progress|review|changes_requested|done) needs_planning=1 ;;
esac
if (( needs_planning )); then
    if [[ "$m_type" == "feature" ]]; then
        for f in requirements design tasks; do
            [[ -f "$SPECS_DIR/$id/$f.md" ]] || fail "falta $f.md (feature en estado $m_status)"
        done
    else
        [[ -f "$SPECS_DIR/$id/plan.md" ]] || fail "falta plan.md (task en estado $m_status)"
    fi
fi

# --- evidencia ---
case "$m_status" in
    review|done|changes_requested)
        [[ -f "$WORK_ITEMS_DIR/$id/impl.md" ]] || fail "falta impl.md (estado $m_status)" ;;
esac
if [[ "$m_status" == "done" && ! -f "$WORK_ITEMS_DIR/$id/review.md" ]]; then
    fail "falta review.md (estado done)"
fi

# --- coherencia con current.md (si es la sesión de este Work Item) ---
if [[ -f "$CURRENT_FILE" ]]; then
    cur_wi="$(current_work_item || true)"
    if [[ "$cur_wi" == "$id" ]]; then
        cur_state="$(current_state || true)"
        if [[ -n "$cur_state" && "$cur_state" != "$m_status" ]]; then
            warn "current.md dice Estado='$cur_state' pero meta.json dice status='$m_status'"
        fi
    fi
fi

echo
if (( issues )); then
    echo "Resultado: FALLO"
    exit 1
fi
echo "Resultado: OK"
exit 0
