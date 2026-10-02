package adapter

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/k1wi777/my-harness-SDD/internal/check"
	"github.com/k1wi777/my-harness-SDD/internal/initwizard"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

// hostedDir es el directorio destino de los agentes nativos de OpenCode V2.
const hostedDir = ".opencode/agents"

// readRole lee y parsea .rei/agents/<rol>.md.
func readRole(p *paths.Project, name string) (Role, error) {
	data, err := os.ReadFile(filepath.Join(p.ReiDir(), "agents", name+".md"))
	if err != nil {
		return Role{}, fmt.Errorf("no se pudo leer el rol %s: %w", name, err)
	}
	return ParseRole(data)
}

// agentPath devuelve la ruta destino del agente nativo de un rol.
func agentPath(p *paths.Project, name string) string {
	return filepath.Join(p.Root, filepath.FromSlash(hostedDir), name+".md")
}

// buildAgent construye el contenido nativo canónico de un rol.
func buildAgent(p *paths.Project, name string) (string, error) {
	role, err := readRole(p, name)
	if err != nil {
		return "", err
	}
	return renderAgent(p, role)
}

// GenerateAll genera .opencode/agents/<rol>.md para cada rol canónico. No
// sobrescribe archivos que no lleven la marca GENERATED: los reporta con
// [WARN] y los conserva. Devuelve 0 si completa (aunque omita algún archivo) y
// 1 si no pudo escribir un archivo generado.
func GenerateAll(p *paths.Project, out io.Writer) int {
	exit := 0
	for _, name := range Roles {
		rel := hostedDir + "/" + name + ".md"
		content, err := buildAgent(p, name)
		if err != nil {
			fmt.Fprintf(out, "[FAIL]  %s: %v\n", rel, err)
			exit = 1
			continue
		}

		dest := agentPath(p, name)
		existing, err := os.ReadFile(dest)
		switch {
		case err == nil:
			if !strings.Contains(string(existing), generatedMark) {
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

// Check comprueba, sin escribir, que cada archivo generado existe y coincide
// byte a byte con la generación canónica. Devuelve 0 si todo coincide y 1 si
// falta, no lleva la marca o difiere algún archivo.
func Check(p *paths.Project, out io.Writer) int {
	exit := 0
	for _, name := range Roles {
		rel := hostedDir + "/" + name + ".md"
		content, err := buildAgent(p, name)
		if err != nil {
			fmt.Fprintf(out, "[FAIL]  %s: %v\n", rel, err)
			exit = 1
			continue
		}

		existing, err := os.ReadFile(agentPath(p, name))
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
		if !strings.Contains(string(existing), generatedMark) {
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

// Install despliega el esqueleto (incluido .rei/adapters/**), asegura la
// estructura de estado y genera los agentes nativos. Reutiliza
// initwizard.InstallSkeleton y check.EnsureStructure para no duplicar la
// instalación. Devuelve 1 si alguna fase falla.
func Install(p *paths.Project, out io.Writer) int {
	fmt.Fprintln(out, "== rei init opencode ==")
	fmt.Fprintln(out)

	installed := initwizard.InstallSkeleton(p, out)
	fmt.Fprintln(out)

	scaffold := check.EnsureStructure(p, out)
	fmt.Fprintln(out)

	generated := GenerateAll(p, out)

	if installed != 0 || scaffold != 0 || generated != 0 {
		return 1
	}
	return 0
}
