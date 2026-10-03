package rei

import (
	"errors"
	"io/fs"
	"strings"
	"testing"
)

func TestSkeletonContieneArchivosRequeridos(t *testing.T) {
	for _, path := range []string{"AGENTS.md", ".rei/config.json"} {
		if _, err := fs.Stat(Skeleton, path); err != nil {
			t.Errorf("Skeleton no contiene %s: %v", path, err)
		}
	}
}

func TestSkeletonContieneAdaptadores(t *testing.T) {
	for _, path := range []string{
		".rei/adapters/opencode/agent.tmpl",
		".rei/adapters/opencode/tools.json",
		".rei/adapters/claude/agent.tmpl",
		".rei/adapters/claude/tools.json",
	} {
		if _, err := fs.Stat(Skeleton, path); err != nil {
			t.Errorf("Skeleton no contiene %s: %v", path, err)
		}
	}
}

func TestSkeletonContieneSubarbolesDelHarness(t *testing.T) {
	prefixes := []string{".rei/docs/", ".rei/agents/", ".rei/adapters/", ".rei/templates/"}
	found := map[string]bool{}
	err := fs.WalkDir(Skeleton, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		for _, p := range prefixes {
			if strings.HasPrefix(path, p) {
				found[p] = true
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WalkDir del esqueleto: %v", err)
	}
	for _, p := range prefixes {
		if !found[p] {
			t.Errorf("Skeleton no contiene ningún archivo bajo %s", p)
		}
	}
}

func TestSkeletonExcluyeEstadoDelProyecto(t *testing.T) {
	for _, path := range []string{".rei/specs", ".rei/progress"} {
		_, err := fs.Stat(Skeleton, path)
		if err == nil {
			t.Errorf("Skeleton NO debe contener %s", path)
			continue
		}
		if !errors.Is(err, fs.ErrNotExist) {
			t.Errorf("Stat(%s) devolvió un error inesperado: %v", path, err)
		}
	}
}
