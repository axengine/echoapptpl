package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func TestToYaml(t *testing.T) {
	cfg := Config{}

	yaml.NewEncoder(os.Stdout).Encode(&cfg)
}
