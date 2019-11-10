package main

import (
	"flag"

	"github.com/justinas/alice"
	"go.uber.org/zap"

	"github.com/jacoelho/codewars/internal/config"
	"github.com/jacoelho/codewars/internal/httpclient"
	"github.com/jacoelho/codewars/internal/log"
	"github.com/jacoelho/codewars/internal/notifier/slack"
	"github.com/jacoelho/codewars/internal/usecase"
	"github.com/jacoelho/codewars/internal/user/api"
	"github.com/jacoelho/codewars/internal/web"
)

func main() {
	logger := log.New()

	cfg := config.New()

	cfg.Flags()

	flag.Parse()

	if validateErrs := cfg.Validate(); len(validateErrs) > 0 {
		logger.Fatal("invalid configuration", zap.Errors("error", validateErrs))
	}

	httpClient := httpclient.WithDefaults()

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
