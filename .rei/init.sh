#!/usr/bin/env bash
#
# .rei/init.sh
#
# Inicializa y verifica REI Harness antes de comenzar una sesión.
#
# Ejecución recomendada (no requiere permisos de ejecución):
#   bash .rei/init.sh
#
# Alternativa, si el script tiene permiso de ejecución:
#   chmod +x .rei/init.sh && ./.rei/init.sh
#
# Si algún elemento crítico de REI Harness falta, la sesión no debe comenzar.
#

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd -- "$SCRIPT_DIR/.." && pwd)"
cd "$PROJECT_ROOT" || exit 1

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

ok()   { printf "${GREEN}[OK]${NC}    %s\n" "$1"; }
warn() { printf "${YELLOW}[WARN]${NC}  %s\n" "$1"; }
fail() { printf "${RED}[FAIL]${NC}  %s\n" "$1"; }

center_line() {
    local line="$1"
    local terminal_width="$2"
    local color="${3:-}"
    local reset="${4:-}"
    local line_width=${#line}
    local padding=$(( (terminal_width - line_width) / 2 ))

    (( padding < 0 )) && padding=0
    printf '%*s%s%s%s\n' "$padding" '' "$color" "$line" "$reset"
}

print_banner() {
    local terminal_width=80
    local detected_width=''

    if command -v tput >/dev/null 2>&1; then
        detected_width="$(tput cols 2>/dev/null || true)"
    fi

    if [[ "$detected_width" =~ ^[0-9]+$ ]] && (( detected_width > 0 )); then
        terminal_width="$detected_width"
    fi

    local banner_white=''
    local banner_reset=''
    if [[ -t 1 ]]; then
        banner_white=$'\033[97m'
        banner_reset=$'\033[0m'
    fi

    printf '\n'

    if (( terminal_width < 28 )); then
        center_line 'REI HARNESS' "$terminal_width" "$banner_white" "$banner_reset"
        printf '\n'
        return
    fi

    local art_line
    local rei_art=(
        '██████╗ ███████╗██╗'
        '██╔══██╗██╔════╝██║'
        '██████╔╝█████╗  ██║'
        '██╔══██╗██╔══╝  ██║'
        '██║  ██║███████╗██║'
        '╚═╝  ╚═╝╚══════╝╚═╝'
    )

    for art_line in "${rei_art[@]}"; do
        center_line "$art_line" "$terminal_width" "$banner_white" "$banner_reset"
    done

    printf '\n'
    center_line 'H A R N E S S' "$terminal_width" "$banner_white" "$banner_reset"
    printf '\n'

    local box_text='SDD and Multi-agent Orchestration'
    if (( terminal_width < ${#box_text} + 4 )); then
        box_text='REI Harness'
    fi

    local dashes
    printf -v dashes '%*s' "$(( ${#box_text} + 2 ))" ''
    dashes=${dashes// /-}

    center_line "+${dashes}+" "$terminal_width" "$banner_white" "$banner_reset"
    center_line "| ${box_text} |" "$terminal_width" "$banner_white" "$banner_reset"
    center_line "+${dashes}+" "$terminal_width" "$banner_white" "$banner_reset"
    printf '\n'
}

EXIT_CODE=0

print_banner
echo "────────────────────────────────────────────"
echo " REI Harness Initialization"
echo "────────────────────────────────────────────"
echo

###########################################################
# 1. Verificar archivos críticos de REI Harness
###########################################################

echo "── 1. Verificando REI Harness ─────────────"

REQUIRED_FILES=(
    "AGENTS.md"

    ".rei/docs/harness/workflow.md"
    ".rei/docs/harness/specs.md"
    ".rei/docs/harness/meta.md"
    ".rei/docs/harness/progress.md"
    ".rei/docs/project/architecture.md"
    ".rei/docs/project/conventions.md"
    ".rei/docs/project/verification.md"

    ".rei/agents/leader.md"
    ".rei/agents/spec_author.md"
    ".rei/agents/implementer.md"
    ".rei/agents/reviewer.md"

    ".rei/scripts/_lib.sh"
    ".rei/scripts/session-status.sh"
    ".rei/scripts/work-items-status.sh"
    ".rei/scripts/reset-current.sh"
    ".rei/scripts/start-session.sh"
    ".rei/scripts/archive-session.sh"
    ".rei/scripts/new-work-item.sh"

    ".rei/templates/current.md"
    ".rei/templates/history.md"
    ".rei/templates/meta.json"
)

FAILED_CHECKS=()

for file in "${REQUIRED_FILES[@]}"; do
    if [[ ! -f "$file" ]]; then
        fail "Falta $file"
        EXIT_CODE=1
        FAILED_CHECKS+=("Falta $file")
    else
        ok "$file"
    fi
done

echo

###########################################################
# 2. Inicializar estructura auxiliar de REI Harness
###########################################################

echo "── 2. Inicializando estructura ────────────"

mkdir -p .rei/specs
mkdir -p .rei/progress/work-items

ok ".rei/specs/"
ok ".rei/progress/work-items/"

CURRENT_FILE=".rei/progress/current.md"
HISTORY_FILE=".rei/progress/history.md"

# Las plantillas viven en .rei/templates/ (fuente única de verdad).
# No se duplican aquí: se copian desde el archivo canónico.
if [[ ! -f "$CURRENT_FILE" ]]; then
    bash .rei/scripts/reset-current.sh > /dev/null
    ok ".rei/progress/current.md creado desde plantilla"
else
    ok ".rei/progress/current.md"
fi

if [[ ! -f "$HISTORY_FILE" ]]; then
    cp .rei/templates/history.md "$HISTORY_FILE"
    ok ".rei/progress/history.md creado desde plantilla"
else
    ok ".rei/progress/history.md"
fi

echo

###########################################################
# 3. Validar invariantes de REI Harness
###########################################################

echo "── 3. Validando invariantes ───────────────"

# El estado se resuelve con código, nunca leyendo archivos desde la IA.
# Ver .rei/scripts/session-status.sh y .rei/scripts/work-items-status.sh

session_report="$(bash .rei/scripts/session-status.sh)"
session_rc=$?
if (( session_rc == 0 )); then
    ok "$session_report"
else
    warn "$session_report"
    warn "Revisa si debe continuarse antes de iniciar un nuevo Work Item."
fi

work_items_report="$(bash .rei/scripts/work-items-status.sh)"
work_items_rc=$?
echo "$work_items_report"
if (( work_items_rc != 0 )); then
    warn "Resuelve el Work Item activo antes de iniciar uno nuevo."
fi

echo

###########################################################
# 4. Verificación del proyecto
###########################################################

echo "── 4. Verificando proyecto ───────────────"

echo "[INFO] Personaliza esta sección según tu proyecto."
echo "[INFO] Los comandos deben coincidir con .rei/docs/project/verification.md."

# Usa run_check para que cualquier fallo se propague correctamente a $EXIT_CODE.
# El Reviewer exige que .rei/init.sh finalice sin errores antes de aprobar un Work Item,
# así que un comando de verificación que falle DEBE marcar REI Harness como fallido.
run_with_spinner() {
    local pid="$1"
    local message="$2"
    local frames='|/-\'
    local frame_index=0
    local status=0

    while kill -0 "$pid" 2>/dev/null; do
        printf '\r\033[2K[%s] %s' "${frames:frame_index:1}" "$message"
        frame_index=$(( (frame_index + 1) % ${#frames} ))
        sleep 0.1
    done

    wait "$pid" || status=$?
    printf '\r\033[2K'
    return "$status"
}

run_check() {
    local desc="$1"; shift

    if [[ -t 1 ]]; then
        "$@" > /dev/null 2>&1 &
        local pid=$!

        if run_with_spinner "$pid" "$desc"; then
            ok "$desc"
        else
            fail "$desc"
            EXIT_CODE=1
            FAILED_CHECKS+=("$desc")
        fi
    elif "$@" > /dev/null 2>&1; then
        ok "$desc"
    else
        fail "$desc"
        EXIT_CODE=1
        FAILED_CHECKS+=("$desc")
    fi
}

#
# Ejemplos (descomenta y adapta a tu stack):
#
# run_check "Lint"  npm run lint
# run_check "Tests" npm test
# run_check "Build" npm run build
#
# run_check "Tests" pytest
#
# run_check "Tests" cargo test
#
echo

###########################################################
# 5. Resumen
###########################################################

echo "── 5. Resumen ─────────────────────────────"

if [[ $EXIT_CODE -eq 0 ]]; then
    ok "REI Harness listo para trabajar."
else
    fail "REI Harness contiene errores:"
    for issue in "${FAILED_CHECKS[@]}"; do
        fail "  - $issue"
    done
fi

exit $EXIT_CODE
