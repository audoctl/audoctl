package configs

import (
	"context"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
)

const (
	DefaultContextDeadline = 15 * time.Second
)

var Version string = "undefined"

var cfg *Config = &Config{}

type Config struct {
	Database    `yaml:"database" env:",prefix=DB_"`
	Application `yaml:"application"`
	Log         `yaml:"log"`
	HTTPServer  `yaml:"http_server" env:",prefix=HTTP_"`
	CORS        `yaml:"cors" env:",prefix=CORS_"`
	TLS         `yaml:"tls" env:",prefix=TLS_"`
	Security    `yaml:"security" env:",prefix=SECURITY_"`
}

type Bootstrap struct {
	ConfigPath string `env:"AUDOCTL_CONFIG,default=config.yaml"`
}

func LoadConfig(ctx context.Context) (*Config, error) {
	var b Bootstrap

	// 1. first get config file path from env
	if err := envconfig.Process(ctx, &b); err != nil {
		return nil, err
	}

	// 2. Read YAML file if it exists
	f, err := os.Open(b.ConfigPath)
	if err == nil {
		defer f.Close()
		if err := yaml.NewDecoder(f).Decode(cfg); err != nil {
			return nil, err
		}
	}
	// 3. Run envconfig again (overwrite YAML values with ENV values)
	// sethvargo/go-envconfig keeps existing struct values, only updates with env vars.
	if err := envconfig.Process(ctx, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func GetEnv() string {
	return cfg.Log.Env
}

func GetAppName() string {
	return cfg.Application.Name
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}
