package paths

import (
	"errors"
	"os"
	"path/filepath"
)

// Project representa un proyecto REI, enraizado en el directorio que contiene .rei/.
type Project struct {
	Root string
}

// ErrNotFound indica que no se encontró un proyecto REI subiendo desde el cwd.
var ErrNotFound = errors.New("no se encontró un proyecto REI (.rei/)")

// Find localiza el proyecto REI subiendo desde el directorio actual.
func Find() (*Project, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return FindFrom(cwd)
}

// FindFrom localiza el proyecto REI subiendo desde un directorio dado.
func FindFrom(start string) (*Project, error) {
	dir := start
	for {
		rei := filepath.Join(dir, ".rei")
		if fi, err := os.Stat(rei); err == nil && fi.IsDir() {
			return &Project{Root: dir}, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, ErrNotFound
		}
		dir = parent
	}
}

func (p *Project) ReiDir() string       { return filepath.Join(p.Root, ".rei") }
func (p *Project) SpecsDir() string     { return filepath.Join(p.ReiDir(), "specs") }
func (p *Project) ProgressDir() string  { return filepath.Join(p.ReiDir(), "progress") }
func (p *Project) WorkItemsDir() string { return filepath.Join(p.ProgressDir(), "work-items") }
func (p *Project) TemplatesDir() string { return filepath.Join(p.ReiDir(), "templates") }
func (p *Project) CurrentFile() string  { return filepath.Join(p.ProgressDir(), "current.md") }
func (p *Project) HistoryFile() string  { return filepath.Join(p.ProgressDir(), "history.md") }

func (p *Project) SpecDir(id string) string     { return filepath.Join(p.SpecsDir(), id) }
func (p *Project) MetaFile(id string) string    { return filepath.Join(p.SpecDir(id), "meta.json") }
func (p *Project) WorkItemDir(id string) string { return filepath.Join(p.WorkItemsDir(), id) }

// MetaFiles devuelve las rutas de todos los meta.json bajo .rei/specs/.
func (p *Project) MetaFiles() ([]string, error) {
	return filepath.Glob(filepath.Join(p.SpecsDir(), "*", "meta.json"))
}
