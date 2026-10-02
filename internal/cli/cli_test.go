package cli

import (
	"errors"
	"testing"

	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

func TestProjectMessage(t *testing.T) {
	const want = "REI no está inicializado aquí; ejecuta `rei init` para inicializarlo."
	if got := projectMessage(paths.ErrNotFound); got != want {
		t.Fatalf("projectMessage(ErrNotFound) = %q, want %q", got, want)
	}
	if got := projectMessage(errors.New("boom")); got != "boom" {
		t.Fatalf("projectMessage(other) = %q, want %q", got, "boom")
	}
}
