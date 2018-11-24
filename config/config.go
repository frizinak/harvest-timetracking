package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

type ConfigLoader struct {
	path  string
	value Config
}

type Config interface {
	Validate() error
}

func New(path string, defaultValue Config) *ConfigLoader {
	return &ConfigLoader{path, defaultValue}
}

func DotFile(name string, defaultValue Config) (*ConfigLoader, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &ConfigLoader{filepath.Join(u.HomeDir, name), defaultValue}, nil
}

func (c *ConfigLoader) Path() string {
	return c.path
}

func (c *ConfigLoader) Read(v Config) error {
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	d := json.NewDecoder(file)
	if err := d.Decode(v); err != nil {
		return err
	}

	return v.Validate()
}

func (c *ConfigLoader) Create(v Config) error {
	file, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	e := json.NewEncoder(file)
	e.SetIndent("", "    ")
	return e.Encode(v)
}

func (c *ConfigLoader) CreateDefault() error {
	return c.Create(c.value)
}
