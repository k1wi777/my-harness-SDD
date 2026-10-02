package cli

import (
	"strings"
	"testing"
)

func TestCommandHelpTextConocidos(t *testing.T) {
	for _, name := range []string{
		"check", "doctor", "new", "session", "items", "status",
		"commit", "validate", "review-diff", "test", "item", "version", "help",
	} {
		text, ok := commandHelpText(name)
		if !ok {
			t.Errorf("commandHelpText(%q) no encontrado", name)
			continue
		}
		if !strings.Contains(text, "Uso:") || !strings.Contains(text, name) {
			t.Errorf("ayuda de %q incompleta:\n%s", name, text)
		}
	}
}

func TestCommandHelpTextDesconocido(t *testing.T) {
	if _, ok := commandHelpText("no-existe"); ok {
		t.Fatalf("commandHelpText(no-existe) devolvió ok")
	}
}

func TestPrintCommandHelpDesconocido(t *testing.T) {
	if code := printCommandHelp("no-existe"); code != 2 {
		t.Fatalf("printCommandHelp(no-existe) = %d, esperaba 2", code)
	}
}

func TestHelpTestMencionaMakeTest(t *testing.T) {
	text, ok := commandHelpText("test")
	if !ok {
		t.Fatal("no hay ayuda para test")
	}
	if !strings.Contains(text, "make test") {
		t.Fatalf("la ayuda de test debe mencionar 'make test':\n%s", text)
	}
}

func TestHelpGeneralIncluyeNuevosComandos(t *testing.T) {
	text := helpText()
	for _, want := range []string{"test", "item"} {
		if !strings.Contains(text, want) {
			t.Errorf("la ayuda general no incluye %q:\n%s", want, text)
		}
	}
}
