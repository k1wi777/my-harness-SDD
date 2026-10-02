package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/k1wi777/my-harness-SDD/internal/check"
	"github.com/k1wi777/my-harness-SDD/internal/gitx"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
	"github.com/k1wi777/my-harness-SDD/internal/state"
	"github.com/k1wi777/my-harness-SDD/internal/template"
	"github.com/k1wi777/my-harness-SDD/internal/validate"
	"github.com/k1wi777/my-harness-SDD/internal/workitem"
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
	case "check":
		return cmdCheck(args[1:])
	case "new":
		return cmdNew(args[1:])
	case "session":
		return cmdSession(args[1:])
	case "items":
		return cmdItems(args[1:])
	case "commit":
		return cmdCommit(args[1:])
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
  check [--quiet]              Verifica el harness, inicializa la estructura y corre los checks
  new <id> <feature|task> [title]
                               Crea un Work Item nuevo
  session                      Muestra la sesión activa
  session start <id> <type>    Inicia la sesión de un Work Item
  session archive              Archiva la sesión en history.md y resetea current.md
  session reset                Restablece current.md
  items status                 Lista el estado de todos los Work Items
  commit set <id> <base_commit|last_review_commit>
                               Registra el commit actual de git en meta.json
  validate [<id>]              Comprueba la consistencia interna de un Work Item
  review-diff <id> [--full]    Genera el paquete de revisión de un Work Item
  version                      Muestra la versión
  help                         Muestra esta ayuda
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

func cmdCheck(args []string) int {
	quiet := false
	for _, a := range args {
		if a == "--quiet" || a == "-q" {
			quiet = true
		}
	}
	p, code := project()
	if p == nil {
		return code
	}
	return check.Run(p, quiet, os.Stdout)
}

func cmdNew(args []string) int {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "uso: rei new <id> <feature|task> [title]")
		return 2
	}
	id, typ := args[0], args[1]
	title := ""
	if len(args) > 2 {
		title = args[2]
	}
	p, code := project()
	if p == nil {
		return code
	}
	if err := workitem.New(p, id, typ, title); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}
	fmt.Printf("Work Item creado: .rei/specs/%s/ (type: %s, status: pending)\n", id, typ)
	fmt.Printf("Completa 'description' en .rei/specs/%s/meta.json.\n", id)
	return 0
}

func cmdSession(args []string) int {
	if len(args) == 0 || args[0] == "status" {
		return cmdSessionStatus()
	}
	p, code := project()
	if p == nil {
		return code
	}
	switch args[0] {
	case "start":
		if len(args) < 3 {
			fmt.Fprintln(os.Stderr, "uso: rei session start <id> <feature|task>")
			return 2
		}
		if err := template.StartSession(p, args[1], args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			return 1
		}
		fmt.Printf("Sesión iniciada: %s (%s) -> pending (agente: spec_author)\n", args[1], args[2])
		return 0
	case "archive":
		id, err := template.ArchiveSession(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			return 1
		}
		if id == "" {
			fmt.Println("No hay sesión activa que archivar.")
			return 0
		}
		fmt.Printf("Sesión '%s' archivada en history.md. current.md restablecido.\n", id)
		return 0
	case "reset":
		if err := template.ResetCurrent(p); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
			return 1
		}
		fmt.Println("current.md restablecido.")
		return 0
	default:
		fmt.Fprintf(os.Stderr, "subcomando desconocido: session %s\n", args[0])
		return 2
	}
}

func cmdSessionStatus() int {
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

func cmdCommit(args []string) int {
	if len(args) < 3 || args[0] != "set" {
		fmt.Fprintln(os.Stderr, "uso: rei commit set <id> <base_commit|last_review_commit>")
		return 2
	}
	p, code := project()
	if p == nil {
		return code
	}
	msg, err := gitx.SetCommit(p, args[1], args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}
	fmt.Println(msg)
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
