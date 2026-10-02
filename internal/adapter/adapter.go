// Package adapter implementa el formato canónico de rol de REI Harness y su
// adaptación a runtimes nativos. En esta versión el único runtime implementado
// es OpenCode V2: el paquete extrae el frontmatter y la sección ## Contrato de
// .rei/agents/<rol>.md, los mapea al frontmatter nativo y genera
// .opencode/agents/<rol>.md.
//
// No tiene dependencias externas: el parser de frontmatter es el subconjunto
// mínimo descrito en design.md (clave: valor y tools como lista inline).
package adapter

import (
	"errors"
	"fmt"
	"strings"
)

// Roles es la lista canónica y ordenada de roles del arnés. Coincide con los
// archivos de rol requeridos por check.RequiredFiles.
var Roles = []string{"leader", "spec_author", "implementer", "reviewer", "initializer"}

// genericTools es el conjunto de herramientas genéricas admitidas en el
// frontmatter canónico.
var genericTools = map[string]bool{
	"read":     true,
	"write":    true,
	"edit":     true,
	"search":   true,
	"shell":    true,
	"subagent": true,
}

// Role es un rol canónico parseado desde .rei/agents/<rol>.md.
type Role struct {
	Name        string
	Description string
	Mode        string   // "primary" | "subagent"
	Tools       []string // genéricas: read, write, edit, search, shell, subagent
	Model       string   // opcional
	Contract    string   // texto de ## Contrato, sin los encabezados ## Contrato y ## Referencia
}

// ParseRole parsea un archivo de rol canónico. Devuelve error si falta el
// frontmatter, si `mode` no es `primary`/`subagent`, si hay una herramienta
// desconocida, o si no existen las secciones ## Contrato y ## Referencia en el
// orden correcto.
func ParseRole(data []byte) (Role, error) {
	fm, body, err := splitFrontmatter(data)
	if err != nil {
		return Role{}, err
	}
	scalars, tools, err := parseFrontmatter(fm)
	if err != nil {
		return Role{}, err
	}

	var role Role
	role.Name = scalars["name"]
	role.Description = scalars["description"]
	role.Mode = scalars["mode"]
	role.Model = scalars["model"]
	role.Tools = tools

	if role.Name == "" {
		return Role{}, errors.New("frontmatter: falta la clave 'name'")
	}
	if role.Description == "" {
		return Role{}, errors.New("frontmatter: falta la clave 'description'")
	}
	switch role.Mode {
	case "primary", "subagent":
	default:
		return Role{}, fmt.Errorf("frontmatter: 'mode' inválido: %q (se espera primary o subagent)", role.Mode)
	}
	if len(role.Tools) == 0 {
		return Role{}, errors.New("frontmatter: falta la clave 'tools' o está vacía")
	}
	for _, tool := range role.Tools {
		if !genericTools[tool] {
			return Role{}, fmt.Errorf("frontmatter: herramienta genérica desconocida: %q", tool)
		}
	}

	contract, err := ExtractContract(body)
	if err != nil {
		return Role{}, err
	}
	role.Contract = contract
	return role, nil
}

// ExtractContract devuelve el texto entre el encabezado ## Contrato y el
// encabezado ## Referencia, excluyendo ambos encabezados y recortando los
// límites. Devuelve error si falta alguna sección, si hay duplicados o si el
// orden es incorrecto.
func ExtractContract(body string) (string, error) {
	lines := strings.Split(body, "\n")
	contractIdx, referenceIdx := -1, -1
	for i, line := range lines {
		switch strings.TrimSpace(line) {
		case "## Contrato":
			if contractIdx >= 0 {
				return "", errors.New("sección ## Contrato duplicada")
			}
			contractIdx = i
		case "## Referencia":
			if referenceIdx >= 0 {
				return "", errors.New("sección ## Referencia duplicada")
			}
			referenceIdx = i
		}
	}
	if contractIdx < 0 {
		return "", errors.New("falta la sección ## Contrato")
	}
	if referenceIdx < 0 {
		return "", errors.New("falta la sección ## Referencia")
	}
	if referenceIdx < contractIdx {
		return "", errors.New("## Contrato debe preceder a ## Referencia")
	}
	return strings.TrimSpace(strings.Join(lines[contractIdx+1:referenceIdx], "\n")), nil
}

// splitFrontmatter separa el bloque YAML del cuerpo. La primera línea debe ser
// el delimitador de apertura --- y debe existir un delimitador de cierre.
func splitFrontmatter(data []byte) (frontmatter, body string, err error) {
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	lines := strings.Split(s, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return "", "", errors.New("frontmatter: la primera línea debe ser '---'")
	}
	end := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			end = i
			break
		}
	}
	if end < 0 {
		return "", "", errors.New("frontmatter: falta el delimitador de cierre '---'")
	}
	return strings.Join(lines[1:end], "\n"), strings.Join(lines[end+1:], "\n"), nil
}

// parseFrontmatter interpreta líneas `clave: valor` y la lista inline `tools`.
func parseFrontmatter(frontmatter string) (map[string]string, []string, error) {
	scalars := map[string]string{}
	var tools []string
	for _, raw := range strings.Split(frontmatter, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, nil, fmt.Errorf("frontmatter: línea inválida: %q", line)
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key == "tools" {
			list, err := parseInlineList(value)
			if err != nil {
				return nil, nil, err
			}
			tools = list
			continue
		}
		scalars[key] = value
	}
	return scalars, tools, nil
}

// parseInlineList interpreta `[read, write, edit]`.
func parseInlineList(value string) ([]string, error) {
	v := strings.TrimSpace(value)
	if !strings.HasPrefix(v, "[") || !strings.HasSuffix(v, "]") {
		return nil, fmt.Errorf("frontmatter: 'tools' debe ser una lista inline entre corchetes: %q", value)
	}
	inner := strings.TrimSpace(v[1 : len(v)-1])
	if inner == "" {
		return nil, nil
	}
	var out []string
	for _, part := range strings.Split(inner, ",") {
		if item := strings.TrimSpace(part); item != "" {
			out = append(out, item)
		}
	}
	return out, nil
}
