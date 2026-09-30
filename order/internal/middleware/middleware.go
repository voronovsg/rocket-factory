package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/voronovsg/rocket-factory/order/internal/metrics"
)

// responseWriter обертка для http.ResponseWriter для отслеживания кода состояния
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newResponseWriter создает новый responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

// WriteHeader переопределяет метод для сохранения кода состояния
func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// MetricsMiddleware middleware для сбора метрик
func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)
		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.statusCode)

		metrics.ObserveHttpDuration(
			r.Context(),
			r.Method,
			chi.RouteContext(r.Context()).RoutePattern(),
			status,
			duration,
		)
	})
}
