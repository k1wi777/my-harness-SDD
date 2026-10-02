package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/k1wi777/my-harness-SDD/internal/gitx"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
	"github.com/k1wi777/my-harness-SDD/internal/state"
	"github.com/k1wi777/my-harness-SDD/internal/validate"
)

const version = "0.1.0-dev"

// Run despacha un comando y devuelve el código de salida.
func Run(args []string) int {
	if len(args) == 0 {
		printHelp()
		return 0
	}
	switch args[0] {
	case "version", "--version", "-v":
		fmt.Printf("rei %s\n", version)
		return 0
	case "help", "--help", "-h":
		printHelp()
		return 0
	case "session":
		return cmdSession()
	case "items":
		return cmdItems(args[1:])
	case "validate":
		return cmdValidate(args[1:])
	case "review-diff":
		return cmdReviewDiff(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "comando desconocido: %s\n\n", args[0])
		printHelp()
		return 2
	}
}

func printHelp() {
	fmt.Print(`rei — CLI de REI Harness

Uso:
  rei <comando> [argumentos]

Comandos:
  session                 Muestra la sesión activa
  items status            Lista el estado de todos los Work Items
  validate [<id>]         Comprueba la consistencia interna de un Work Item
  review-diff <id> [--full]
                          Genera el paquete de revisión de un Work Item
  version                 Muestra la versión
  help                    Muestra esta ayuda
`)
}

func project() (*paths.Project, int) {
	p, err := paths.Find()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return nil, 1
	}
	return p, 0
}

func cmdSession() int {
	p, code := project()
	if p == nil {
		return code
	}
	s, err := state.ReadSession(p)
	if err != nil || !s.Active() {
		fmt.Println("Sin sesión activa.")
		return 0
	}
	fmt.Printf("Sesión activa: %s (estado: %s, agente: %s)\n", s.WorkItem, orQ(s.State), orQ(s.Agent))
	return 0
}

func cmdItems(args []string) int {
	if len(args) == 0 || args[0] != "status" {
		fmt.Fprintln(os.Stderr, "uso: rei items status")
		return 2
	}
	p, code := project()
	if p == nil {
		return code
	}
	items, err := state.ListWorkItems(p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if len(items) == 0 {
		fmt.Println("No hay Work Items registrados.")
		return 0
	}
	fmt.Println("Work Items:")
	active := 0
	for _, it := range items {
		marker := ""
		if it.Status == "in_progress" {
			marker = "   << ACTIVO"
			active++
		}
		fmt.Printf("  - %s [%s]%s\n", it.ID, it.Status, marker)
	}
	if active > 0 {
		fmt.Println("AVISO: existe un Work Item en in_progress.")
	}
	return 0
}

func cmdValidate(args []string) int {
	id := ""
	if len(args) > 0 {
		id = args[0]
	}
	p, code := project()
	if p == nil {
		return code
	}
	if id == "" {
		s, err := state.ReadSession(p)
		if err != nil || !s.Active() {
			fmt.Println("Sin sesión activa.")
			return 0
		}
		id = s.WorkItem
	}
	issues, err := validate.WorkItem(p, id)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Printf("== Validando Work Item: %s ==\n", id)
	fails := 0
	for _, is := range issues {
		fmt.Printf("[%s] %s\n", is.Level, is.Message)
		if is.Level == validate.LevelFail {
			fails++
		}
	}
	fmt.Println()
	if fails > 0 {
		fmt.Println("Resultado: FALLO")
		return 1
	}
	fmt.Println("Resultado: OK")
	return 0
}

func cmdReviewDiff(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "uso: rei review-diff <work-item-id> [--full]")
		return 2
	}
	id := args[0]
	full := false
	for _, a := range args[1:] {
		if a == "--full" {
			full = true
		}
	}
	p, code := project()
	if p == nil {
		return code
	}
	out, err := gitx.ReviewDiff(p, id, full)
	switch {
	case errors.Is(err, gitx.ErrNoGit):
		fmt.Println("SIN_GIT")
		return 2
	case errors.Is(err, gitx.ErrNoBase):
		fmt.Println("SIN_BASE")
		return 2
	case err != nil:
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	fmt.Print(out)
	return 0
}

func orQ(s string) string {
	if s == "" {
		return "?"
	}
	return s
}
