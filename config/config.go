package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Rest RestConfig `koanf:"rest"`
}

type RestConfig struct {
	Server ServerConfig   `koanf:"server"`
	Db     DatabaseConfig `koanf:"db"`
}

type ServerConfig struct {
	Host string `koanf:"host"`
}

type DatabaseConfig struct {
	DataStore  string `koanf:"datastore"`
	NumberConn int    `koanf:"nconn"`
}

// LoadConfig loads configuration from a YAML file and merges it with environment
// variables that share the given prefix. Environment variables take priority over
// file values when both define the same key.
//
// Env vars are mapped by stripping the prefix, lowercasing, and replacing "_"
// with ".". For example, with prefix "APP", APP_REST_SERVER_HOST maps to
// rest.server.host.
func LoadConfig(configPath, envPrefix string) (Config, error) {
	k := koanf.New(".")

	// 1. Load the YAML file (lower priority).
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("loading config file %q: %w", configPath, err)
	}

	// 2. Overlay environment variables (higher priority).
	prefix := strings.ToUpper(envPrefix) + "_"
	err := k.Load(env.Provider(prefix, ".", func(s string) string {
		// Strip the prefix, lowercase, replace "_" with "." to form the koanf key.
		s = strings.TrimPrefix(s, prefix)
		return strings.ReplaceAll(strings.ToLower(s), "_", ".")
	}), nil)
	if err != nil {
		return Config{}, fmt.Errorf("loading env vars with prefix %q: %w", envPrefix, err)
	}

	// 3. Unmarshal into Config struct.
	cfg := Config{}
	if err := k.UnmarshalWithConf("", &cfg, koanf.UnmarshalConf{Tag: "koanf"}); err != nil {
		return Config{}, fmt.Errorf("unmarshaling config: %w", err)
	}

	return cfg, nil
}
