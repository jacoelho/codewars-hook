package config

import (
	"errors"
	"flag"
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

func (c *Config) Flags(fs *flag.FlagSet) {
	fs.StringVar(&c.Port, "port", getEnvOrDefault("PORT", "8080"), "listen port (env: PORT)")
	fs.StringVar(&c.SlackHook, "slack-webhook", getEnvOrDefault("SLACK_WEBHOOK", ""), "slack webhook (env: SLACK_WEBHOOK)")
	fs.StringVar(&c.Secret, "secret", getEnvOrDefault("SECRET", "secret"), "authentication secret (env: SECRET)")
}

func (c *Config) Validate() []error {
	var errs []error

	if c.Port == "" {
		errs = append(errs, errors.New("invalid port"))
	}

	if c.SlackHook == "" {
		errs = append(errs, errors.New("invalid slack webhook"))
	}

	if c.Secret == "" {
		errs = append(errs, errors.New("invalid secret"))
	}

	return errs
}
