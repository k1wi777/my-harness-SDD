package cli

import (
	"fmt"
	"os"
	"strings"
)

// commandInfo describe un comando para la ayuda general y por comando.
type commandInfo struct {
	name   string
	usage  string
	short  string
	detail []string
}

// commands es la tabla única que alimenta la ayuda general y la de cada
// comando (evita duplicar el texto entre ambas).
var commands = []commandInfo{
	{
		name:  "check",
		usage: "check [--quiet]",
		short: "Verifica el harness, inicializa la estructura y corre los checks",
		detail: []string{
			"--quiet, -q    omite los [OK] y el resumen.",
			"Además de los checks, crea la estructura y, si falta, inicializa git.",
			"Sale 0 si todo pasa; 1 si hay algún [FAIL].",
		},
	},
	{
		name:  "init",
		usage: "init [status]",
		short: "Instala el esqueleto del harness y reporta la personalización",
		detail: []string{
			"Sin argumentos: despliega el esqueleto embebido (AGENTS.md, .rei/docs/, .rei/agents/, .rei/templates/ y .rei/config.json) sin sobrescribir archivos existentes, crea la estructura de estado (.rei/specs/, .rei/progress/work-items/, current.md, history.md) y, si falta y git está disponible, inicializa el repositorio.",
			"status: reporta qué documentos siguen pendientes, sin modificar archivos (0 si no queda ninguno; 1 si queda alguno).",
			"Códigos de salida: 0 si completa la instalación; 1 si no puede crear algún archivo del esqueleto o de la estructura; 2 uso incorrecto.",
			"Si AGENTS.md ya existe se avisa y se conserva sin cambios.",
			"La personalización guiada la conduce el rol `initializer` (.rei/agents/initializer.md).",
		},
	},
	{
		name:   "doctor",
		usage:  "doctor",
		short:  "Diagnóstico de solo lectura del harness (no modifica nada)",
		detail: []string{"No escribe archivos ni ejecuta comandos que modifiquen el proyecto."},
	},
	{
		name:  "new",
		usage: "new <id> <feature|task> [title]",
		short: "Crea un Work Item nuevo",
		detail: []string{
			"Crea .rei/specs/<id>/meta.json y la carpeta de progreso.",
			"El estado inicial es 'pending'.",
		},
	},
	{
		name:  "session",
		usage: "session [status]",
		short: "Muestra o gestiona la sesión activa",
		detail: []string{
			"start <id> <feature|task>    inicia la sesión de un Work Item.",
			"archive                       archiva la sesión en history.md y resetea current.md.",
			"reset                         restablece current.md.",
		},
	},
	{
		name:   "items",
		usage:  "items status",
		short:  "Lista el estado de todos los Work Items",
		detail: []string{"Marca con '<< ACTIVO' el Work Item en in_progress."},
	},
	{
		name:  "status",
		usage: "status set <id> <status> [--force]",
		short: "Cambia el estado de un Work Item en meta.json",
		detail: []string{
			"Valida la transición según workflow.md.",
			"--force permite transiciones no estándar.",
		},
	},
	{
		name:   "commit",
		usage:  "commit set <id> <base_commit|last_review_commit>",
		short:  "Registra el commit actual de git en meta.json",
		detail: []string{"Guarda HEAD en el campo indicado del Work Item."},
	},
	{
		name:  "validate",
		usage: "validate [<id>]",
		short: "Comprueba la consistencia interna de un Work Item",
		detail: []string{
			"Sin <id>, usa la sesión activa.",
			"Sale 0 si no hay [FAIL]; 1 si hay al menos uno.",
		},
	},
	{
		name:   "review-diff",
		usage:  "review-diff <id> [--full]",
		short:  "Genera el paquete de revisión de un Work Item",
		detail: []string{"--full muestra el diff completo, no solo el resumen."},
	},
	{
		name:  "test",
		usage: "test",
		short: "Ejecuta los checks de .rei/config.json",
		detail: []string{
			"Ejecuta SOLO los checks declarados en .rei/config.json, sin efectos secundarios.",
			"NO es la suite de tests de REI; esos son `make test` (dev, requiere Go).",
			"Sin checks configurados informa y sale 0; sale 1 si algún check falla.",
		},
	},
	{
		name:  "item",
		usage: "item show <id>",
		short: "Muestra la ficha de un Work Item",
		detail: []string{
			"Ficha de solo lectura: metadatos, documentos, sesión y validación.",
			"Sale 1 si hay algún [FAIL] de validación.",
		},
	},
	{
		name:   "version",
		usage:  "version",
		short:  "Muestra la versión",
		detail: []string{"También disponible como --version o -v."},
	},
	{
		name:   "help",
		usage:  "help [<comando>]",
		short:  "Muestra la ayuda general o de un comando",
		detail: []string{"Sin argumentos muestra la ayuda general."},
	},
}

// helpText construye la ayuda general desde la tabla.
func helpText() string {
	var b strings.Builder
	b.WriteString("rei — CLI de REI Harness\n")
	b.WriteString("\nUso:\n")
	b.WriteString("  rei <comando> [argumentos]\n")
	b.WriteString("\nComandos:\n")
	for _, c := range commands {
		fmt.Fprintf(&b, "  %s\n", c.usage)
		fmt.Fprintf(&b, "      %s\n", c.short)
	}
	b.WriteString("\nUsa 'rei <comando> --help' para ver el detalle de un comando.\n")
	return b.String()
}

func printHelp() {
	fmt.Print(helpText())
}

// commandHelpText devuelve el texto de ayuda de un comando a partir de la
// tabla. No escribe a stdout, para poder testearlo.
func commandHelpText(name string) (string, bool) {
	for _, c := range commands {
		if c.name != name {
			continue
		}
		var b strings.Builder
		fmt.Fprintf(&b, "rei %s — %s\n", c.name, c.short)
		fmt.Fprintf(&b, "Uso: %s\n", c.usage)
		for _, d := range c.detail {
			fmt.Fprintf(&b, "  %s\n", d)
		}
		return b.String(), true
	}
	return "", false
}

// printCommandHelp imprime la ayuda de un comando. Si el comando no existe,
// informa por stderr, muestra la ayuda general y devuelve 2.
func printCommandHelp(name string) int {
	text, ok := commandHelpText(name)
	if !ok {
		fmt.Fprintf(os.Stderr, "comando desconocido: %s\n\n", name)
		printHelp()
		return 2
	}
	fmt.Print(text)
	return 0
}
