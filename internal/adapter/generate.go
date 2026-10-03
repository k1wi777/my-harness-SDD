package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	harnesscheck "github.com/k1wi777/my-harness-SDD/internal/check"
	"github.com/k1wi777/my-harness-SDD/internal/initwizard"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

// runtime describe cómo adaptar el formato canónico de rol a un runtime nativo.
// Es el punto de extensión: para añadir un runtime solo hay que declarar su
// descriptor y registrar su plantilla y su mapa de herramientas en
// .rei/adapters/<name>/.
type runtime struct {
	// name es el identificador del runtime (nombre del directorio del
	// adaptador y del subcomando de `rei init`).
	name string
	// hostedDir es el directorio destino de los agentes nativos.
	hostedDir string
	// mark es la marca GENERATED que autoriza a sobrescribir un archivo.
	mark string
	// nativeKey es la clave del mapa de herramientas que contiene los destinos
	// nativos de ese runtime (p. ej. "permission" o "tools").
	nativeKey string
	// validTargets son los destinos nativos admitidos por el runtime.
	validTargets map[string]bool
	// buildData construye los datos que consume agent.tmpl a partir del rol y de
	// los destinos nativos ya resueltos y deduplicados.
	buildData func(role Role, targets []string) any
}

// toolEntry es una entrada de tools.json: una clave nativa por runtime
// (p. ej. "permission" o "tools") con la lista de destinos nativos.
type toolEntry map[string][]string

// toolMap es el contenido de tools.json: herramienta genérica -> entrada.
type toolMap map[string]toolEntry

// sortedKeys devuelve las claves de un toolMap ordenadas.
func sortedKeys(tm toolMap) []string {
	keys := make([]string, 0, len(tm))
	for key := range tm {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// readRole lee y parsea .rei/agents/<rol>.md.
func readRole(p *paths.Project, name string) (Role, error) {
	data, err := os.ReadFile(filepath.Join(p.ReiDir(), "agents", name+".md"))
	if err != nil {
		return Role{}, fmt.Errorf("no se pudo leer el rol %s: %w", name, err)
	}
	return ParseRole(data)
}

// adapterDir devuelve el directorio del adaptador de un runtime dentro del
// proyecto destino.
func adapterDir(p *paths.Project, rt runtime) string {
	return filepath.Join(p.ReiDir(), "adapters", rt.name)
}

// agentPath devuelve la ruta destino del agente nativo de un rol.
func agentPath(p *paths.Project, rt runtime, name string) string {
	return filepath.Join(p.Root, filepath.FromSlash(rt.hostedDir), name+".md")
}

// loadToolMap lee y valida .rei/adapters/<rt.name>/tools.json.
func loadToolMap(p *paths.Project, rt runtime) (toolMap, error) {
	data, err := os.ReadFile(filepath.Join(adapterDir(p, rt), "tools.json"))
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer tools.json del adaptador: %w", err)
	}
	var tm toolMap
	if err := json.Unmarshal(data, &tm); err != nil {
		return nil, fmt.Errorf("tools.json inválido: %w", err)
	}
	if err := validateToolMap(tm, rt); err != nil {
		return nil, err
	}
	return tm, nil
}

// validateToolMap exige que el mapa cubra todas las herramientas genéricas y
// que cada entrada declare al menos un destino nativo válido para el runtime.
func validateToolMap(tm toolMap, rt runtime) error {
	generics := make([]string, 0, len(genericTools))
	for tool := range genericTools {
		generics = append(generics, tool)
	}
	sort.Strings(generics)
	for _, tool := range generics {
		if _, ok := tm[tool]; !ok {
			return fmt.Errorf("tools.json no cubre la herramienta genérica %q", tool)
		}
	}
	for _, tool := range sortedKeys(tm) {
		targets := tm[tool][rt.nativeKey]
		if len(targets) == 0 {
			return fmt.Errorf("tools.json: %q no declara ninguna clave de %s", tool, rt.nativeKey)
		}
		for _, target := range targets {
			if !rt.validTargets[target] {
				return fmt.Errorf("tools.json: destino nativo desconocido %q para %q", target, tool)
			}
		}
	}
	return nil
}

// resolveTargets traduce las herramientas genéricas del rol a destinos nativos
// del runtime, en el orden de declaración de `tools` y deduplicando destinos
// (p. ej. write y edit colapsan en Write, Edit).
func resolveTargets(tools []string, tm toolMap, rt runtime) ([]string, error) {
	var out []string
	seen := map[string]bool{}
	for _, tool := range tools {
		entry, ok := tm[tool]
		if !ok {
			return nil, fmt.Errorf("herramienta genérica sin mapeo en tools.json: %q", tool)
		}
		for _, target := range entry[rt.nativeKey] {
			if seen[target] {
				continue
			}
			seen[target] = true
			out = append(out, target)
		}
	}
	return out, nil
}

// buildAgent construye el contenido nativo canónico de un rol.
func buildAgent(p *paths.Project, rt runtime, name string) (string, error) {
	role, err := readRole(p, name)
	if err != nil {
		return "", err
	}
	return renderAgent(p, rt, role)
}

// renderAgent construye el archivo nativo de un rol: frontmatter mapeado, marca
// GENERATED y el ## Contrato como cuerpo.
func renderAgent(p *paths.Project, rt runtime, role Role) (string, error) {
	raw, err := os.ReadFile(filepath.Join(adapterDir(p, rt), "agent.tmpl"))
	if err != nil {
		return "", fmt.Errorf("no se pudo leer agent.tmpl del adaptador: %w", err)
	}
	tmpl, err := template.New("agent").Parse(string(raw))
	if err != nil {
		return "", fmt.Errorf("agent.tmpl inválido: %w", err)
	}
	tm, err := loadToolMap(p, rt)
	if err != nil {
		return "", err
	}
	targets, err := resolveTargets(role.Tools, tm, rt)
	if err != nil {
		return "", err
	}
	data := rt.buildData(role, targets)
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("no se pudo aplicar agent.tmpl: %w", err)
	}
	return buf.String(), nil
}

// generateAll genera <rt.hostedDir>/<rol>.md para cada rol canónico. No
// sobrescribe archivos que no lleven la marca GENERATED: los reporta con
// [WARN] y los conserva. Devuelve 0 si completa (aunque omita algún archivo) y
// 1 si no pudo escribir un archivo generado.
func generateAll(p *paths.Project, rt runtime, out io.Writer) int {
	exit := 0
	for _, name := range Roles {
		rel := rt.hostedDir + "/" + name + ".md"
		content, err := buildAgent(p, rt, name)
		if err != nil {
			fmt.Fprintf(out, "[FAIL]  %s: %v\n", rel, err)
			exit = 1
			continue
		}

		dest := agentPath(p, rt, name)
		existing, err := os.ReadFile(dest)
		switch {
		case err == nil:
			if !strings.Contains(string(existing), rt.mark) {
				fmt.Fprintf(out, "[WARN]  %s no fue generado por rei; se conserva sin cambios. Elimínalo o regenera el archivo a mano.\n", rel)
				continue
			}
			if string(existing) == content {
				fmt.Fprintf(out, "[OK]    %s\n", rel)
				continue
			}
			if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
				fmt.Fprintf(out, "[FAIL]  No se pudo escribir %s: %v\n", rel, err)
				exit = 1
				continue
			}
			fmt.Fprintf(out, "[UPD]   %s\n", rel)
		case os.IsNotExist(err):
			if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
				fmt.Fprintf(out, "[FAIL]  No se pudo crear el directorio de %s: %v\n", rel, err)
				exit = 1
				continue
			}
			if err := os.WriteFile(dest, []byte(content), 0o644); err != nil {
				fmt.Fprintf(out, "[FAIL]  No se pudo escribir %s: %v\n", rel, err)
				exit = 1
				continue
			}
			fmt.Fprintf(out, "[OK]    %s\n", rel)
		default:
			fmt.Fprintf(out, "[FAIL]  No se pudo comprobar %s: %v\n", rel, err)
			exit = 1
		}
	}
	return exit
}

// check comprueba, sin escribir, que cada archivo generado existe y coincide
// byte a byte con la generación canónica. Devuelve 0 si todo coincide y 1 si
// falta, no lleva la marca o difiere algún archivo.
func check(p *paths.Project, rt runtime, out io.Writer) int {
	exit := 0
	for _, name := range Roles {
		rel := rt.hostedDir + "/" + name + ".md"
		content, err := buildAgent(p, rt, name)
		if err != nil {
			fmt.Fprintf(out, "[FAIL]  %s: %v\n", rel, err)
			exit = 1
			continue
		}

		existing, err := os.ReadFile(agentPath(p, rt, name))
		switch {
		case os.IsNotExist(err):
			fmt.Fprintf(out, "[FAIL]  Falta %s\n", rel)
			exit = 1
			continue
		case err != nil:
			fmt.Fprintf(out, "[FAIL]  No se pudo leer %s: %v\n", rel, err)
			exit = 1
			continue
		}
		if !strings.Contains(string(existing), rt.mark) {
			fmt.Fprintf(out, "[FAIL]  %s existe pero no fue generado por rei (sin la marca GENERATED)\n", rel)
			exit = 1
			continue
		}
		if string(existing) != content {
			fmt.Fprintf(out, "[FAIL]  %s difiere de la generación canónica\n", rel)
			exit = 1
			continue
		}
		fmt.Fprintf(out, "[OK]    %s\n", rel)
	}
	return exit
}

// install despliega el esqueleto (incluido .rei/adapters/**), asegura la
// estructura de estado y genera los agentes nativos. Reutiliza
// initwizard.InstallSkeleton y check.EnsureStructure para no duplicar la
// instalación. Devuelve 1 si alguna fase falla.
func install(p *paths.Project, rt runtime, out io.Writer) int {
	fmt.Fprintf(out, "== rei init %s ==\n", rt.name)
	fmt.Fprintln(out)

	installed := initwizard.InstallSkeleton(p, out)
	fmt.Fprintln(out)

	scaffold := harnesscheck.EnsureStructure(p, out)
	fmt.Fprintln(out)

	generated := generateAll(p, rt, out)

	if installed != 0 || scaffold != 0 || generated != 0 {
		return 1
	}
	return 0
}
