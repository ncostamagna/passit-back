package secrets

import (
	"context"
	"time"

	"github.com/ncostamagna/passit-back/internal/entities"
	"github.com/prometheus/client_golang/prometheus"
)

type (
	instrumenting struct {
		requestCount          *prometheus.CounterVec
		requestLatency        *prometheus.HistogramVec
		requestLatencySummary *prometheus.SummaryVec
		s                     Service
	}

	Instrumenting interface {
		Service
	}
)

func NewInstrumenting(requestCount *prometheus.CounterVec, requestLatencySummary *prometheus.SummaryVec, requestLatency *prometheus.HistogramVec, s Service) Instrumenting {
	return &instrumenting{
		requestCount:          requestCount,
		requestLatencySummary: requestLatencySummary,
		requestLatency:        requestLatency,
		s:                     s,
	}
}

func (i *instrumenting) Create(ctx context.Context, secret *entities.Secret) (string, error) {
	defer func(begin time.Time) {
		i.requestCount.WithLabelValues("Create").Inc()
		i.requestLatencySummary.WithLabelValues("Create").Observe(time.Since(begin).Seconds())
		i.requestLatency.WithLabelValues("Create").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.s.Create(ctx, secret)
}

func (i *instrumenting) Get(ctx context.Context, key string) (*entities.Secret, error) {
	defer func(begin time.Time) {
		i.requestCount.WithLabelValues("Get").Inc()
		i.requestLatencySummary.WithLabelValues("Get").Observe(time.Since(begin).Seconds())
		i.requestLatency.WithLabelValues("Get").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.s.Get(ctx, key)
}

func (i *instrumenting) Delete(ctx context.Context, key string) error {
	defer func(begin time.Time) {
		i.requestCount.WithLabelValues("Delete").Inc()
		i.requestLatencySummary.WithLabelValues("Delete").Observe(time.Since(begin).Seconds())
		i.requestLatency.WithLabelValues("Delete").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return i.s.Delete(ctx, key)
}
