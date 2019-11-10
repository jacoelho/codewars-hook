package web

import (
	"encoding/json"
	"net/http"

	"github.com/jacoelho/codewars/internal/usecase"
	"github.com/jacoelho/codewars/internal/user"
)

func EventWebHook(use usecase.UserHonorUpdated) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var data user.Event

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "failed to decode: "+err.Error(), http.StatusBadRequest)
			return
		}

		if err := use(r.Context(), data); err != nil {
			http.Error(w, "failed to process honor updated: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		//nolint:errcheck
		w.Write([]byte(`{"status":"ok"}`))
	}
}

func Routes(honorUpdated usecase.UserHonorUpdated) *http.ServeMux {
	r := http.NewServeMux()

	r.Handle("/webhook", EventWebHook(honorUpdated))

	return r
}
