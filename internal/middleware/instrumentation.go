package middleware

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// total number of requests
	requestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total", 
			Help: "Total number of HTTP requests",
		}, 
		[]string{"method", "path", "status"},
	)

	requestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request latency in seconds",
		},
		[]string{"method", "path"},
	)
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggerAndMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		statusStr := strconv.Itoa(rw.statusCode)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.statusCode,
			"duration", duration,
			"user_agent", r.UserAgent(),
		)

		requestCount.WithLabelValues(r.Method, r.URL.Path, statusStr).Inc()
		requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration.Seconds())
	})
}