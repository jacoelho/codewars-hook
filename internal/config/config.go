package config

import (
	"flag"
	"fmt"
	"os"
)

func getEnvOrDefault(key, def string) string {
	v := os.Getenv(key)

	if v == "" {
		return def
	}

	return v
}

type Config struct {
	Port      string
	SlackHook string
	Secret    string
}

func New() *Config {
	return &Config{}
}

func (c *Config) Flags() {
	flag.StringVar(&c.Port, "port", getEnvOrDefault("PORT", "8080"), "listen port (env: PORT)")
	flag.StringVar(&c.SlackHook, "slack-webhook", getEnvOrDefault("SLACK_WEBHOOK", ""), "slack webhook (env: SLACK_WEBHOOK)")
	flag.StringVar(&c.Secret, "secret", getEnvOrDefault("SECRET", "secret"), "authentication secret (env: SECRET)")
}

func (c *Config) Validate() []error {
	errs := []error{}

	if c.Port == "" {
		errs = append(errs, fmt.Errorf("invalid port"))
	}

	if c.SlackHook == "" {
		errs = append(errs, fmt.Errorf("invalid slack webhook"))
	}

	if c.Secret == "" {
		errs = append(errs, fmt.Errorf("invalid secret"))
	}

	return errs
}
