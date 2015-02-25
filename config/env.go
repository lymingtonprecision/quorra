package config

import (
	"os"
	str "strings"
)

type VarSetter func(*Config, string)
type VarSetters map[string]VarSetter

const EnvPrefix string = "QUORRA_"

var ValidVars VarSetters = VarSetters{
	"USERNAME": func(c *Config, v string) { c.Username = v },
	"PASSWORD": func(c *Config, v string) { c.Password = v },

	"API_URL":    func(c *Config, v string) { c.ApiUrl = v },
	"DATACENTER": func(c *Config, v string) { c.Datacenter = v },
	"DATASTORE":  func(c *Config, v string) { c.Default.Datastore = v },
	"HOST":       func(c *Config, v string) { c.Default.Host = v },
	"OVA":        func(c *Config, v string) { c.Default.OVA = v },
}

func ParseEnv() (*Config, error) {
	return ParseEnvArray(os.Environ())
}

func ParseEnvArray(env []string) (*Config, error) {
	var config Config

	if err := config.ParseEnvArray(env); err != nil {
		return nil, err
	}

	return &config, nil
}

func (config *Config) ParseEnv() error {
	return config.ParseEnvArray(os.Environ())
}

func (config *Config) ParseEnvArray(env []string) error {
	m := envMap(env)

	for k, s := range ValidVars {
		if v, ok := m[k]; ok {
			s(config, v)
		}
	}

	return nil
}

func envMap(env []string) map[string]string {
	m := make(map[string]string)

	for _, e := range env {
		kv := str.SplitN(e, "=", 2)

		if len(kv) == 2 && str.HasPrefix(kv[0], EnvPrefix) {
			m[str.TrimPrefix(kv[0], EnvPrefix)] = kv[1]
		}
	}

	return m
}
