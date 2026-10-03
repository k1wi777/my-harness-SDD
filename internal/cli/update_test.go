package cli

import (
	"strings"
	"testing"
)

func TestCmdUpdateUsoIncorrecto(t *testing.T) {
	cases := [][]string{
		{"extra"},
		{"--bogus"},
		{"--check", "extra"},
	}
	for _, args := range cases {
		if code := cmdUpdate(args); code != 2 {
			t.Errorf("cmdUpdate(%v) = %d, esperaba 2", args, code)
		}
	}
}

func TestHelpUpdateDocumentaCheck(t *testing.T) {
	text, ok := commandHelpText("update")
	if !ok {
		t.Fatal("no hay ayuda para update")
	}
	for _, want := range []string{"update [--check]", "--check", "checksums.txt", "Windows", "dev"} {
		if !strings.Contains(text, want) {
			t.Errorf("la ayuda de update debe mencionar %q:\n%s", want, text)
		}
	}
}

func TestHelpGeneralIncluyeUpdate(t *testing.T) {
	if !strings.Contains(helpText(), "update [--check]") {
		t.Fatal("la ayuda general debe incluir 'update [--check]'")
	}
}

func TestRunUpdateHelp(t *testing.T) {
	if code := Run([]string{"update", "--help"}); code != 0 {
		t.Fatalf("rei update --help = %d, esperaba 0", code)
	}
}
