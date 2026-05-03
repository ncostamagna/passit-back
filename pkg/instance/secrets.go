package instance

import (
	"log/slog"

	"github.com/ncostamagna/passit-back/adapters/cache"
	"github.com/ncostamagna/passit-back/internal/secrets"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	metricNamespace = "api"
	metricSubsystem = "passit_service"
	metricMethod    = "method"
)

func NewSecretService(cache cache.Cache, logger *slog.Logger) secrets.Service {
	service := secrets.NewService(logger, cache)

	requestCount := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "request_count_total",
		Help:      "Number of requests received.",
	}, []string{metricMethod})

	requestLatencySummary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "request_latency_seconds",
		Help:      "Total duration of requests in seconds.",
	}, []string{metricMethod})

	requestLatency := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: metricNamespace,
		Subsystem: metricSubsystem,
		Name:      "request_latency_seconds_2",
		Help:      "Total duration of requests in seconds.",
	}, []string{metricMethod})

	prometheus.MustRegister(requestCount, requestLatencySummary, requestLatency)

	return secrets.NewInstrumenting(requestCount, requestLatencySummary, requestLatency, service)
}
