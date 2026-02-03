package instance

import (
	"log/slog"

	"github.com/ncostamagna/passit-back/adapters/cache"
	"github.com/ncostamagna/passit-back/internal/secrets"
	"github.com/prometheus/client_golang/prometheus"
)

func NewSecretService(cache cache.Cache, logger *slog.Logger) secrets.Service {
	service := secrets.NewService(logger, cache)

	requestCount := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "api",
		Subsystem: "passit_service",
		Name:      "request_count_total",
		Help:      "Number of requests received.",
	}, []string{"method"})

	requestLatencySummary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "api",
		Subsystem: "passit_service",
		Name:      "request_latency_seconds",
		Help:      "Total duration of requests in seconds.",
	}, []string{"method"})

	requestLatency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "api",
		Subsystem: "passit_service",
		Name:      "request_latency_seconds",
		Help:      "Total duration of requests in seconds.",
	}, []string{"method"})

	prometheus.MustRegister(requestCount, requestLatencySummary, requestLatency)

	return secrets.NewInstrumenting(requestCount, requestLatencySummary, requestLatency, service)
}
