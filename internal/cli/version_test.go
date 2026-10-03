package cli

import "testing"

// TestVersionDefaultYInyectable comprueba que el valor por defecto de `version`
// es "dev" y que se trata de una variable mutable, de modo que el linker pueda
// sobrescribirla en build time con -ldflags -X.
func TestVersionDefaultYInyectable(t *testing.T) {
	if version != "dev" {
		t.Fatalf("version por defecto = %q, want %q", version, "dev")
	}

	original := version
	defer func() { version = original }()

	version = "v0.1.0-test"
	if version != "v0.1.0-test" {
		t.Fatalf("version reasignada = %q, want %q", version, "v0.1.0-test")
	}
}
