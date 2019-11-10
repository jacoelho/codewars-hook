package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacoelho/codewars/internal/web"
)

func TestSecretHandlerSuccess(t *testing.T) {
	simple := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	h := web.AllowedWebhookSecret("abc")(http.HandlerFunc(simple))

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add(web.WebhookSecretHeader, "abc")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSecretHandlerUnauthorized(t *testing.T) {
	simple := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	h := web.AllowedWebhookSecret("abc")(http.HandlerFunc(simple))

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add(web.WebhookSecretHeader, "abcccc")

	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}
