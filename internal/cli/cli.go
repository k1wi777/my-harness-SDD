package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/k1wi777/my-harness-SDD/internal/check"
	"github.com/k1wi777/my-harness-SDD/internal/doctor"
	"github.com/k1wi777/my-harness-SDD/internal/gitx"
	"github.com/k1wi777/my-harness-SDD/internal/initwizard"
	"github.com/k1wi777/my-harness-SDD/internal/meta"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
	"github.com/k1wi777/my-harness-SDD/internal/show"
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
	cmd := args[0]
	rest := args[1:]
	for _, a := range rest {
		if a == "--help" || a == "-h" {
			return printCommandHelp(cmd)
		}
	}
	switch cmd {
	case "version", "--version", "-v":
		fmt.Printf("rei %s\n", version)
		return 0
	case "help", "--help", "-h":
		if len(rest) > 0 {
			return printCommandHelp(rest[0])
		}
		printHelp()
		return 0
	case "check":
		return cmdCheck(rest)
	case "init":
		return cmdInit(rest)
	case "doctor":
		return cmdDoctor(rest)
	case "new":
		return cmdNew(rest)
	case "session":
		return cmdSession(rest)
	case "items":
		return cmdItems(rest)
	case "status":
		return cmdStatus(rest)
	case "commit":
		return cmdCommit(rest)
	case "validate":
		return cmdValidate(rest)
	case "review-diff":
		return cmdReviewDiff(rest)
	case "test":
		return cmdTest(rest)
	case "item":
		return cmdItem(rest)
	default:
		fmt.Fprintf(os.Stderr, "comando desconocido: %s\n\n", cmd)
		printHelp()
		return 2
	}
}

func project() (*paths.Project, int) {
	p, err := paths.Find()
	if err != nil {
		fmt.Fprintln(os.Stderr, projectMessage(err))
		return nil, 1
	}
	return p, 0
}

// projectMessage traduce un error de resolución de proyecto a un mensaje útil.
func projectMessage(err error) string {
	if errors.Is(err, paths.ErrNotFound) {
		return "REI no está inicializado aquí; ejecuta `rei init` para inicializarlo."
	}
	return err.Error()
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

// cmdInit implementa `rei init [status]`: en modo instalador (sin argumentos)
// despliega el esqueleto embebido en el proyecto destino —resuelto con
// initwizard.ResolveRoot, sin exigir un .rei/ previo— y reporta la
// personalización; `status` es un reporte de solo lectura que sigue exigiendo
// un proyecto REI existente (project()).
func cmdInit(args []string) int {
	if len(args) > 1 || (len(args) == 1 && args[0] != "status") {
		fmt.Fprintln(os.Stderr, "uso: rei init [status]")
		return 2
	}
	if len(args) == 1 {
		p, code := project()
		if p == nil {
			return code
		}
		return initwizard.Status(p, os.Stdout)
	}
	p, err := initwizard.ResolveRoot()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return initwizard.Init(p, os.Stdout)
}

// cmdTest ejecuta solo los checks declarados en .rei/config.json, sin efectos
// secundarios. NO es la suite de tests de REI (esa es `make test`).
func cmdTest(args []string) int {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "uso: rei test")
		return 2
	}
	p, code := project()
	if p == nil {
		return code
	}
	return check.RunDeclared(p, os.Stdout)
}

// cmdItem muestra la ficha de un Work Item (solo lectura).
func cmdItem(args []string) int {
	if len(args) != 2 || args[0] != "show" {
		fmt.Fprintln(os.Stderr, "uso: rei item show <id>")
		return 2
	}
	p, code := project()
	if p == nil {
		return code
	}
	return show.WorkItem(p, args[1], os.Stdout)
}

func cmdDoctor(args []string) int {
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "uso: rei doctor")
		return 2
	}
	p, code := project()
	if p == nil {
		return code
	}
	return doctor.Run(p, os.Stdout)
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

func cmdStatus(args []string) int {
	if len(args) == 0 || args[0] != "set" {
		fmt.Fprintln(os.Stderr, "uso: rei status set <id> <status> [--force]")
		return 2
	}
	force := false
	var positional []string
	for _, a := range args[1:] {
		switch {
		case a == "--force":
			force = true
		case strings.HasPrefix(a, "-"):
			fmt.Fprintln(os.Stderr, "uso: rei status set <id> <status> [--force]")
			return 2
		default:
			positional = append(positional, a)
		}
	}
	if len(positional) < 2 {
		fmt.Fprintln(os.Stderr, "uso: rei status set <id> <status> [--force]")
		return 2
	}
	id, status := positional[0], positional[1]
	if !meta.IsValidStatus(status) {
		fmt.Fprintf(os.Stderr, "ERROR: estado inválido: %q\n", status)
		return 1
	}
	p, code := project()
	if p == nil {
		return code
	}
	path := p.MetaFile(id)
	if _, err := os.Stat(path); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: el Work Item %q no existe\n", id)
		return 1
	}
	previous, err := meta.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}
	if previous.Status == status {
		fmt.Printf("Work Item %s: ya está en '%s'.\n", id, status)
		return 0
	}
	forced := false
	if !meta.TransitionAllowed(previous.Status, status) {
		if !force {
			fmt.Fprintf(os.Stderr, "ERROR: transición no permitida: %s -> %s (ver workflow.md; usa --force para forzar)\n", previous.Status, status)
			return 1
		}
		forced = true
	}
	if _, err := meta.SetStatus(path, status); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		return 1
	}
	if forced {
		fmt.Printf("Work Item %s: %s -> %s (forzado).\n", id, previous.Status, status)
	} else {
		fmt.Printf("Work Item %s: %s -> %s.\n", id, previous.Status, status)
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
