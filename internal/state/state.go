package state

import (
	"os"
	"sort"
	"strings"

	"github.com/k1wi777/my-harness-SDD/internal/meta"
	"github.com/k1wi777/my-harness-SDD/internal/paths"
)

const emptyWorkItem = "_ninguno_"

// Session es el estado vivo reflejado en .rei/progress/current.md.
type Session struct {
	WorkItem string
	Type     string
	State    string
	Agent    string
}

// Active indica si hay un Work Item en curso.
func (s Session) Active() bool {
	return s.WorkItem != "" && s.WorkItem != emptyWorkItem
}

// ReadSession parsea .rei/progress/current.md.
func ReadSession(p *paths.Project) (Session, error) {
	data, err := os.ReadFile(p.CurrentFile())
	if err != nil {
		return Session{}, err
	}
	content := string(data)
	return Session{
		WorkItem: markdownField(content, "Work Item"),
		Type:     markdownField(content, "Tipo"),
		State:    markdownField(content, "Estado"),
		Agent:    markdownField(content, "Agente activo"),
	}, nil
}

// ItemStatus resume el estado de un Work Item.
type ItemStatus struct {
	ID     string
	Status string
}

// ListWorkItems lee todos los meta.json y devuelve su estado (ordenado por id).
func ListWorkItems(p *paths.Project) ([]ItemStatus, error) {
	files, err := p.MetaFiles()
	if err != nil {
		return nil, err
	}
	items := make([]ItemStatus, 0, len(files))
	for _, f := range files {
		id := idFromMetaPath(f)
		status := "desconocido"
		if m, err := meta.Load(f); err == nil {
			if m.ID != "" {
				id = m.ID
			}
			if m.Status != "" {
				status = m.Status
			}
		}
		items = append(items, ItemStatus{ID: id, Status: status})
	}
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, nil
}

func markdownField(content, label string) string {
	prefix := "- **" + label + ":**"
	for _, line := range strings.Split(content, "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.TrimSpace(strings.TrimPrefix(line, prefix))
		}
	}
	return ""
}

func idFromMetaPath(path string) string {
	path = strings.TrimSuffix(path, "/meta.json")
	if idx := strings.LastIndex(path, "/"); idx >= 0 {
		return path[idx+1:]
	}
	return path
}
