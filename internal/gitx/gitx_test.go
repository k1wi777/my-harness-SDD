package gitx

import (
	"reflect"
	"testing"
)

func TestFilterPaths(t *testing.T) {
	id := "2026-10-01_10-00__demo"
	in := []string{
		"src/a.go",
		".rei/progress/current.md",
		".rei/specs/" + id + "/tasks.md",
		".rei/specs/otro/tasks.md",
		".rei/docs/harness/specs.md",
		"README.md",
	}
	got := filterPaths(in, id)
	want := []string{"src/a.go", ".rei/specs/" + id + "/tasks.md", "README.md"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterPaths = %v, want %v", got, want)
	}
}

func TestFilterCode(t *testing.T) {
	in := []string{"src/a.go", ".rei/specs/x/tasks.md", "main.go"}
	got := filterCode(in)
	want := []string{"src/a.go", "main.go"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("filterCode = %v, want %v", got, want)
	}
}
