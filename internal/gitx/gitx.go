package gitx

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/k1wi777/my-harness-SDD/internal/meta"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

// ErrNoGit indica que no hay un repositorio git disponible.
var ErrNoGit = errors.New("sin git")

// ErrNoBase indica que no hay un commit base para calcular el diff.
var ErrNoBase = errors.New("sin base")

func gitOut(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), err
}

// IsRepo indica si dir está dentro de un repositorio git.
func IsRepo(dir string) bool {
	if _, err := exec.LookPath("git"); err != nil {
		return false
	}
	_, err := gitOut(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil
}

// HeadSHA devuelve el SHA de HEAD, o "" si no hay commits.
func HeadSHA(dir string) string {
	out, err := gitOut(dir, "rev-parse", "HEAD")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(out)
}

// SetCommit registra el HEAD actual en meta.json.
// Devuelve un mensaje descriptivo; no es un error si no hay git/commits.
func SetCommit(p *paths.Project, id, field string) (string, error) {
	if field != "base_commit" && field != "last_review_commit" {
		return "", fmt.Errorf("campo inválido: %s", field)
	}
	m, err := meta.Load(p.MetaFile(id))
	if err != nil {
		return "", err
	}
	if !IsRepo(p.Root) {
		return "Sin git: no se registra '" + field + "'.", nil
	}
	sha := HeadSHA(p.Root)
	if sha == "" {
		return "Sin commits todavía: no se registra '" + field + "'.", nil
	}
	if field == "base_commit" {
		m.BaseCommit = sha
	} else {
		m.LastReviewCommit = sha
	}
	if err := m.Save(p.MetaFile(id)); err != nil {
		return "", err
	}
	return field + "=" + sha, nil
}

// ReviewDiff genera el paquete de revisión de un Work Item.
func ReviewDiff(p *paths.Project, id string, full bool) (string, error) {
	m, err := meta.Load(p.MetaFile(id))
	if err != nil {
		return "", err
	}
	if !IsRepo(p.Root) {
		return "", ErrNoGit
	}

	base := ""
	origin := "base_commit"
	switch {
	case m.LastReviewCommit != "":
		base = m.LastReviewCommit
		origin = "last_review_commit"
	case m.BaseCommit != "":
		base = m.BaseCommit
	}
	if base == "" {
		return "", ErrNoBase
	}
	if _, err := gitOut(p.Root, "cat-file", "-e", base+"^{commit}"); err != nil {
		return "", ErrNoBase
	}

	tracked := filterPaths(gitLines(p.Root, "diff", "--name-only", base), id)
	untracked := filterPaths(gitLines(p.Root, "ls-files", "--others", "--exclude-standard"), id)

	var b strings.Builder
	fmt.Fprintf(&b, "== Paquete de revisión: %s ==\n", id)
	fmt.Fprintf(&b, "Base: %s (%s)\n\n", base, origin)

	codeTracked := filterCode(tracked)
	b.WriteString("--- Resumen (código) ---\n")
	if len(codeTracked) > 0 {
		args := append([]string{"diff", "--stat", base, "--"}, codeTracked...)
		stat, _ := gitOut(p.Root, args...)
		b.WriteString(stat)
	} else {
		b.WriteString("(sin cambios de código)\n")
	}

	if len(untracked) > 0 {
		b.WriteString("\n--- Nuevos sin rastrear ---\n")
		for _, u := range untracked {
			fmt.Fprintf(&b, "  %s\n", u)
		}
	}

	if full {
		b.WriteString("\n--- Diff completo (solo código; la spec se lista arriba) ---\n")
		code := filterCode(tracked)
		if len(code) > 0 {
			args := append([]string{"diff", base, "--"}, code...)
			diff, _ := gitOut(p.Root, args...)
			b.WriteString(diff)
		}
		for _, u := range untracked {
			if strings.HasPrefix(u, ".rei/") {
				continue
			}
			fmt.Fprintf(&b, "\n--- nuevo: %s ---\n", u)
			diff, _ := gitOut(p.Root, "diff", "--no-index", "--", "/dev/null", u)
			b.WriteString(diff)
		}
	}

	return b.String(), nil
}

func gitLines(dir string, args ...string) []string {
	out, _ := gitOut(dir, args...)
	var lines []string
	for _, line := range strings.Split(out, "\n") {
		if line = strings.TrimSpace(line); line != "" {
			lines = append(lines, line)
		}
	}
	sort.Strings(lines)
	return lines
}

// filterPaths conserva el código del proyecto y la spec del Work Item;
// descarta el resto del harness (.rei/progress, otros specs, docs...).
func filterPaths(paths []string, id string) []string {
	specPrefix := ".rei/specs/" + id + "/"
	var kept []string
	for _, p := range paths {
		switch {
		case strings.HasPrefix(p, specPrefix):
			kept = append(kept, p)
		case strings.HasPrefix(p, ".rei/"):
			// ruido del harness
		default:
			kept = append(kept, p)
		}
	}
	return kept
}

// filterCode conserva solo rutas fuera de .rei/.
func filterCode(paths []string) []string {
	var kept []string
	for _, p := range paths {
		if !strings.HasPrefix(p, ".rei/") {
			kept = append(kept, p)
		}
	}
	return kept
}
