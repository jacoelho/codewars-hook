package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"
	"go.uber.org/zap"

	"github.com/jacoelho/codewars/internal/log"
	"github.com/jacoelho/codewars/internal/notifier/slack"
	"github.com/jacoelho/codewars/internal/usecase"
	"github.com/jacoelho/codewars/internal/user/api"
	"github.com/jacoelho/codewars/internal/web"
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

func main() {
	logger := log.New()

	cfg := &Config{}
	cfg.Flags()

	flag.Parse()

	if validateErrs := cfg.Validate(); len(validateErrs) > 0 {
		logger.Fatal("invalid configuration", zap.Errors("error", validateErrs))
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          5,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	svc := usecase.UserHonorUpdatedCase(
		&slack.Webhook{
			Client:   httpClient,
			Endpoint: cfg.SlackHook,
		},
		api.New(httpClient),
	)

	r := web.Routes(svc)

	chain := alice.New(
		web.LoggingHandler(logger),
		web.LimitBodySize(1<<20),
		web.AllowedWebhookSecret(cfg.Secret),
		web.AllowedWebhookEventType("user"),
	).Then(r)

	web.HTTPServerRunWith(logger, ":"+cfg.Port, chain)
}
