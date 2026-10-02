package template

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

// Empty es el valor canónico para un campo vacío: "_—_".
const Empty = "_\u2014_"

// RenderFile sustituye cada {{CLAVE}} de la plantilla por su valor y escribe el
// destino. Las claves no presentes se dejan intactas.
func RenderFile(templatePath, targetPath string, vars map[string]string) error {
	data, err := os.ReadFile(templatePath)
	if err != nil {
		return err
	}
	content := string(data)
	for key, value := range vars {
		content = strings.ReplaceAll(content, "{{"+key+"}}", value)
	}
	return os.WriteFile(targetPath, []byte(content), 0o644)
}

// ResetCurrent restablece current.md a su estado inicial.
func ResetCurrent(p *paths.Project) error {
	if err := os.MkdirAll(p.ProgressDir(), 0o755); err != nil {
		return err
	}
	return RenderFile(p.TemplatesDir()+"/current.md", p.CurrentFile(), map[string]string{
		"WORK_ITEM": "_ninguno_",
		"TYPE":      Empty,
		"STATE":     Empty,
		"START":     Empty,
		"AGENT":     Empty,
		"PLAN":      Empty,
		"LOG":       Empty,
		"NEXT":      Empty,
	})
}

// StartSession inicia la sesión activa para un Work Item nuevo.
func StartSession(p *paths.Project, id, typ string) error {
	if err := os.MkdirAll(p.ProgressDir(), 0o755); err != nil {
		return err
	}
	return RenderFile(p.TemplatesDir()+"/current.md", p.CurrentFile(), map[string]string{
		"WORK_ITEM": id,
		"TYPE":      typ,
		"STATE":     "pending",
		"START":     time.Now().Format(time.RFC3339),
		"AGENT":     "spec_author",
		"PLAN":      Empty,
		"LOG":       Empty,
		"NEXT":      "Planificación pendiente.",
	})
}

// ArchiveSession mueve el resumen de current.md a history.md y restablece
// current.md. Devuelve el id archivado, o "" si no había sesión activa.
func ArchiveSession(p *paths.Project) (string, error) {
	data, err := os.ReadFile(p.CurrentFile())
	if err != nil {
		return "", nil
	}
	content := string(data)
	workItem := field(content, "Work Item")
	if workItem == "" || workItem == "_ninguno_" {
		return "", nil
	}

	history, err := os.OpenFile(p.HistoryFile(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return "", err
	}
	defer history.Close()

	header := fmt.Sprintf("\n## %s \u2014 %s\n\n", time.Now().Format("2006-01-02 15:04"), workItem)
	if _, err := history.WriteString(header + content + "\n"); err != nil {
		return "", err
	}
	if err := ResetCurrent(p); err != nil {
		return "", err
	}
	return workItem, nil
}

func field(content, label string) string {
	prefix := "- **" + label + ":**"
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}
