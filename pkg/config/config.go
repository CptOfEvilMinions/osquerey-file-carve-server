package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Config
type Config struct {
	Cleanup struct {
		TimeInvterval  int `yaml:"time_invterval"`
		ExpireInterval int `yaml:"expire_interval"`
	} `yaml:"Cleanup"`
	Webserver struct {
		Port int `yaml:"port"`
	} `yaml:"Webserver"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}
