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

EXIT_CODE=0

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
mkdir -p .rei/progress

ok ".rei/specs/"
ok ".rei/progress/"

CURRENT_FILE=".rei/progress/current.md"

# Plantilla definida en .rei/docs/harness/progress.md — si la modificas ahí, actualiza también este heredoc.
if [[ ! -f "$CURRENT_FILE" ]]; then
cat > "$CURRENT_FILE" <<'EOF'
# Sesión actual

> Estado vivo de la sesión.
> Se actualiza durante toda la ejecución.
> Al finalizar el Work Item su resumen se mueve a `history.md`
> y este archivo vuelve a su estado inicial.

- **Work Item:** _ninguno_
- **Tipo:** _—_
- **Estado:** _—_
- **Inicio:** _—_
- **Agente activo:** _—_

## Plan

_—_

## Bitácora

_—_

## Próximo paso

_—_
EOF

    ok ".rei/progress/current.md creado"
else
    ok ".rei/progress/current.md"
fi

HISTORY_FILE=".rei/progress/history.md"

# Plantilla definida en .rei/docs/harness/progress.md — si la modificas ahí, actualiza también este heredoc.
if [[ ! -f "$HISTORY_FILE" ]]; then
cat > "$HISTORY_FILE" <<'EOF'
# Bitácora histórica (append-only)

> Registro histórico de todas las sesiones completadas.
> Nunca modifiques entradas anteriores.
> Siempre añade nuevas entradas al final.

---
EOF

    ok ".rei/progress/history.md creado"
else
    ok ".rei/progress/history.md"
fi

echo

###########################################################
# 3. Validar invariantes de REI Harness
###########################################################

echo "── 3. Validando invariantes ───────────────"

# 3.1 — Sesión activa en .rei/progress/current.md
#
# .rei/progress/current.md es el único archivo que representa la sesión activa.
# Por diseño (AGENTS.md §9: "trabaja sobre un único Work Item por sesión"),
# la regla de .rei/docs/harness/workflow.md ("solo un Work Item en in_progress") queda
# garantizada por esta misma estructura: basta con leer este único archivo,
# no es necesario escanear .rei/specs/*/meta.json.
if grep -q "_ninguno_" "$CURRENT_FILE"; then
    ok "No existe ninguna sesión activa."
else
    warn "Existe una sesión registrada en .rei/progress/current.md."
    warn "Revisa si debe continuarse antes de iniciar un nuevo Work Item."
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
