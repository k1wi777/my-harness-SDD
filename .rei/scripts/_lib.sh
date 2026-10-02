#!/usr/bin/env bash
#
# .rei/scripts/_lib.sh
#
# Librería compartida por los scripts de REI Harness.
# No está pensada para ejecutarse directamente: se carga con `source` desde
# cada script, que previamente define $SCRIPT_DIR.

PROJECT_ROOT="$(cd -- "${SCRIPT_DIR}/../.." && pwd)"
cd "$PROJECT_ROOT" || exit 1

CURRENT_FILE=".rei/progress/current.md"
HISTORY_FILE=".rei/progress/history.md"
SPECS_DIR=".rei/specs"
WORK_ITEMS_DIR=".rei/progress/work-items"
TEMPLATES_DIR=".rei/templates"

# Valor canónico para un campo vacío: "_—_".
# Se construye con printf para no depender de literales multibyte frágiles.
REI_EMPTY="_$(printf '\xe2\x80\x94')_"

# meta_field <meta.json> <clave> — imprime el valor string de una clave.
meta_field() {
    grep -m1 -o "\"$2\"[[:space:]]*:[[:space:]]*\"[^\"]*\"" "$1" 2>/dev/null \
        | sed -E "s/.*\"$2\"[[:space:]]*:[[:space:]]*\"//; s/\"$//"
}

# meta_set <meta.json> <clave> <valor> — reemplaza el valor string de una clave.
meta_set() {
    local file="$1" key="$2" value="$3"
    local tmp="${file}.tmp.$$"
    local value_esc

    value_esc="$(printf '%s' "$value" | sed 's/[&\\]/\\&/g')"

    if sed -E "s/(\"$key\"[[:space:]]*:[[:space:]]*\")[^\"]*(\")/\1${value_esc}\2/" "$file" > "$tmp"; then
        mv "$tmp" "$file"
    else
        rm -f "$tmp"
        return 1
    fi
}

# render_template <plantilla> <destino> CLAVE=valor [CLAVE=valor ...]
#
# Sustituye cada {{CLAVE}} por su valor. La sustitución es literal (bash), así
# que los valores pueden contener "/", "&" o "\" sin necesidad de escaparlos.
render_template() {
    local template="$1" target="$2"
    shift 2

    local -a pairs=("$@")
    local line pair key value

    if [[ ! -f "$template" ]]; then
        printf 'ERROR: no existe la plantilla %s\n' "$template" >&2
        return 1
    fi

    : > "$target" || return 1

    while IFS= read -r line || [[ -n "$line" ]]; do
        for pair in "${pairs[@]}"; do
            key="${pair%%=*}"
            value="${pair#*=}"
            line="${line//\{\{${key}\}\}/$value}"
        done
        printf '%s\n' "$line" >> "$target"
    done < "$template"
}

# reset_current — restablece .rei/progress/current.md a su estado inicial.
reset_current() {
    render_template "$TEMPLATES_DIR/current.md" "$CURRENT_FILE" \
        "WORK_ITEM=_ninguno_" \
        "TYPE=$REI_EMPTY" \
        "STATE=$REI_EMPTY" \
        "START=$REI_EMPTY" \
        "AGENT=$REI_EMPTY" \
        "PLAN=$REI_EMPTY" \
        "LOG=$REI_EMPTY" \
        "NEXT=$REI_EMPTY"
}
