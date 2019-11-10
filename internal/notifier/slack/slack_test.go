package slack_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacoelho/codewars/internal/notifier/slack"
)

func TestWebhookNotifySuccess(t *testing.T) {
	expected := `{
  "attachments":[
    {
       "title":"user",
       "pretext":"Much programing, such honor",
       "text":"message"
    }
  ]
}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)

		stringBody := string(body)

		if stringBody != expected {
			t.Fatal(stringBody)
		}

		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write(body)
	}))
	defer server.Close()

	notifier := slack.Webhook{
		Client:   http.DefaultClient,
		Endpoint: server.URL,
	}

	err := notifier.Notify(context.Background(), "user", "message")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebhookNotifyFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	notifier := slack.Webhook{
		Client:   http.DefaultClient,
		Endpoint: server.URL,
	}

	err := notifier.Notify(context.Background(), "user", "message")
	if err == nil {
		t.Fatal("expected error")
	}
}
