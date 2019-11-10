package main

import (
	"flag"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/justinas/alice"

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

func main() {
	logger := log.New()

	port := flag.String("port", getEnvOrDefault("PORT", "8080"), "listen port")
	slackWebHook := flag.String("slack-webhook", getEnvOrDefault("SLACK_WEBHOOK", ""), "slack webhook")
	secret := flag.String("secret", getEnvOrDefault("SECRET", "secret"), "authentication secret")

	flag.Parse()

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
			Endpoint: *slackWebHook,
		},
		api.New(httpClient),
	)

	r := web.Routes(svc)

	chain := alice.New(
		web.LoggingHandler(logger),
		web.LimitBodySize(1<<20),
		web.AllowedWebhookSecret(*secret),
		web.AllowedWebhookEventType("user"),
	).Then(r)

	web.HTTPServerRunWith(logger, ":"+*port, chain)
}
