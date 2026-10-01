#!/usr/bin/env bash
#
# .rei/scripts/reset-current.sh
#
# Restablece .rei/progress/current.md a su estado inicial ("ninguno").
# Se usa al cerrar una sesión o para recuperar un estado limpio.

set -u

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
source "$SCRIPT_DIR/_lib.sh"

mkdir -p .rei/progress

if ! reset_current; then
    echo "ERROR: no se pudo restablecer current.md." >&2
    exit 1
fi

echo "current.md restablecido."
