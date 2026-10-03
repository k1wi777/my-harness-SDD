package cli

import (
	"strings"
	"testing"
)

func TestCommandHelpTextConocidos(t *testing.T) {
	for _, name := range []string{
		"check", "init", "doctor", "new", "session", "items", "status",
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
	for _, want := range []string{"test", "item", "init [status|opencode|claude [--check]]"} {
		if !strings.Contains(text, want) {
			t.Errorf("la ayuda general no incluye %q:\n%s", want, text)
		}
	}
}

func TestHelpInitDocumentaStatusYInitializer(t *testing.T) {
	text, ok := commandHelpText("init")
	if !ok {
		t.Fatal("no hay ayuda para init")
	}
	for _, want := range []string{"status", "initializer"} {
		if !strings.Contains(text, want) {
			t.Errorf("la ayuda de init debe mencionar %q:\n%s", want, text)
		}
	}
}

func TestHelpInitDocumentaInstalador(t *testing.T) {
	text, ok := commandHelpText("init")
	if !ok {
		t.Fatal("no hay ayuda para init")
	}
	for _, want := range []string{"esqueleto", "git", "Códigos de salida", ".rei/config.json"} {
		if !strings.Contains(text, want) {
			t.Errorf("la ayuda de init debe mencionar %q:\n%s", want, text)
		}
	}
}

func TestHelpInitDocumentaOpencodeYCheck(t *testing.T) {
	text, ok := commandHelpText("init")
	if !ok {
		t.Fatal("no hay ayuda para init")
	}
	for _, want := range []string{"opencode", "--check", ".opencode/agents/", "GENERATED"} {
		if !strings.Contains(text, want) {
			t.Errorf("la ayuda de init debe mencionar %q:\n%s", want, text)
		}
	}
}

func TestHelpInitDocumentaClaudeYFallback(t *testing.T) {
	text, ok := commandHelpText("init")
	if !ok {
		t.Fatal("no hay ayuda para init")
	}
	for _, want := range []string{"claude", ".claude/agents/", "Cursor", "Codex", "fallback", "## Contrato"} {
		if !strings.Contains(text, want) {
			t.Errorf("la ayuda de init debe mencionar %q:\n%s", want, text)
		}
	}
}

func TestPrintCommandHelpInit(t *testing.T) {
	if code := printCommandHelp("init"); code != 0 {
		t.Fatalf("printCommandHelp(init) = %d, esperaba 0", code)
	}
}

func TestCmdInitArgumentoInvalido(t *testing.T) {
	if code := cmdInit([]string{"bogus"}); code != 2 {
		t.Fatalf("cmdInit(bogus) = %d, esperaba 2", code)
	}
	if code := cmdInit([]string{"status", "extra"}); code != 2 {
		t.Fatalf("cmdInit(status extra) = %d, esperaba 2", code)
	}
}

func TestCmdInitClaudeUsoIncorrecto(t *testing.T) {
	if code := cmdInit([]string{"claude", "bogus"}); code != 2 {
		t.Fatalf("cmdInit(claude bogus) = %d, esperaba 2", code)
	}
	if code := cmdInit([]string{"claude", "--check", "extra"}); code != 2 {
		t.Fatalf("cmdInit(claude --check extra) = %d, esperaba 2", code)
	}
}
