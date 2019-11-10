package web

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type statusWriter struct {
	http.ResponseWriter
	Flusher http.Flusher
	status  int
	length  int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func (w *statusWriter) Flush() {
	if w.Flusher != nil {
		w.Flusher.Flush()
	}
}

func newStatusWriter(w http.ResponseWriter) *statusWriter {
	sw := statusWriter{ResponseWriter: w}

	if f, ok := w.(http.Flusher); ok {
		sw.Flusher = f
	}

	return &sw
}

// LoggingHandler logs http requests
func LoggingHandler(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			sw := newStatusWriter(w)

			h.ServeHTTP(sw, r)

			logger.Info("client request",
				zap.String("host", r.Host),
				zap.String("remote", r.RemoteAddr),
				zap.String("method", r.Method),
				zap.String("url", r.RequestURI),
				zap.Int("status", sw.status),
				zap.Int("bytes", sw.length),
				zap.String("user_agent", r.Header.Get("User-Agent")),
				zap.String("duration", time.Since(start).String()))
		})
	}
}
