package config

import (
	"encoding/json"
	"os"
)

// Check es una verificación del proyecto (checkpoint rápido).
type Check struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Command     []string `json:"command"`
}

// Config es el contenido de .rei/config.json.
type Config struct {
	Checks []Check `json:"checks"`
}

// Load lee .rei/config.json. Si no existe, devuelve una configuración vacía.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{}, nil
	}
	if err != nil {
		return nil, err
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
