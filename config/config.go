package config

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"
)

type Config struct {
	path  string
	value interface{}
}

func New(path string, defaultValue interface{}) *Config {
	return &Config{path, defaultValue}
}

func DotFile(name string, defaultValue interface{}) (*Config, error) {
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	return &Config{filepath.Join(u.HomeDir, name), defaultValue}, nil
}

func (c *Config) Path() string {
	return c.path
}

func (c *Config) Read(v interface{}) error {
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	d := json.NewDecoder(file)
	return d.Decode(v)
}

func (c *Config) Create(v interface{}) error {
	file, err := os.Create(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	e := json.NewEncoder(file)
	e.SetIndent("", "    ")
	return e.Encode(v)
}

func (c *Config) CreateDefault() error {
	return c.Create(c.value)
}
