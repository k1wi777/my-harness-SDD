package gitx

import (
	"reflect"
	"testing"
)

func TestFilterPaths(t *testing.T) {
	in := []string{
		"src/a.go",
		"AGENTS.md",
		".rei/progress/current.md",
		".rei/progress/work-items/x/impl.md",
		".rei/specs/2026-10-01_10-00__demo/tasks.md",
		".rei/specs/otro/tasks.md",
		".rei/docs/harness/specs.md",
		".rei/agents/leader.md",
		".rei/adapters/opencode/map.json",
		".rei/templates/current.md",
		".rei/config.json",
		"README.md",
	}
	got := filterPaths(in)
	want := []string{
		"src/a.go",
		"AGENTS.md",
		".rei/specs/2026-10-01_10-00__demo/tasks.md",
		".rei/specs/otro/tasks.md",
		".rei/docs/harness/specs.md",
		".rei/agents/leader.md",
		".rei/adapters/opencode/map.json",
		".rei/templates/current.md",
		".rei/config.json",
		"README.md",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterPaths = %v, want %v", got, want)
	}
}
