package web

import (
	"net/http"
)

// WebhookSecretHeader allows to restrict inbound requests
//nolint:gosec
const WebhookSecretHeader = "X-Webhook-Secret"

// WebhookEventHeader indicates webhook event type
const WebhookEventHeader = "X-Webhook-Event"

func LimitBodySize(size int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
			h.ServeHTTP(w, r)
		})
	}
}

// AllowedWebhookSecret will check if X-Webhook-Secret value matches expected secret
// returns 401 otherwise
func AllowedWebhookSecret(secret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(WebhookSecretHeader) != secret {
				http.Error(w, "invalid secret", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

// AllowedWebhookEventType will check if X-Webhook-Event value matches expected event type
// requests are ignored otherwise
func AllowedWebhookEventType(eventType string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(WebhookEventHeader) != eventType {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
