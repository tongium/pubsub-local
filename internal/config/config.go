package config

import (
	"os"

	"github.com/tongium/pubsub-local/pkg/models"
	"gopkg.in/yaml.v3"
)

func LoadSettings(path string) (*models.Settings, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var settings models.Settings
	if err := yaml.Unmarshal(data, &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}
