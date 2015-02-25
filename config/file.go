package config

import (
	"io/ioutil"
	"os"

	"github.com/naoina/toml"
)

var DefaultPaths []string = []string{
	"/etc/quorra/quorra.conf",
	"~/.quorra.conf",
	"./quorra.conf",
}

func ParseDefaultFiles() (*Config, error) {
	var c Config

	for _, p := range DefaultPaths {
		if err := c.ParseFile(p); err != nil {
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}
	}

	return &c, nil
}

func ParseFile(path string) (*Config, error) {
	var config Config

	if err := config.ParseFile(path); err != nil {
		return nil, err
	}

	return &config, nil
}

func (config *Config) ParseFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	if err := toml.Unmarshal(buf, config); err != nil {
		return err
	}

	return nil
}
