package internalhttp

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pyltsin/otusGolang/hw12_13_14_15_calendar/internal/logger"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK} //nolint:exhaustivestruct
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)
		diff := time.Since(start)
		logger.Log.Info(fmt.Sprintf("ip: %s, method: %s, path: %s, version: %s, status_code: %v, latency: %s, agent: %s",
			strings.Split(r.RemoteAddr, ":")[0],
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.statusCode,
			diff.String(),
			r.UserAgent()))
	})
}
