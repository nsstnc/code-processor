package http

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Время выполнения запроса
	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"endpoint"},
	)

	// Использование фильтров и трансляторов
	FilterUsage = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "filter_usage_total",
			Help: "Count of each filter usage",
		},
		[]string{"filter"},
	)

	TranslatorUsage = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "translator_usage_total",
			Help: "Count of each translator usage",
		},
		[]string{"translator"},
	)
)

func init() {
	prometheus.MustRegister(RequestDuration, FilterUsage, TranslatorUsage)
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}

func useTranslator(translatorName string) {
	TranslatorUsage.WithLabelValues(translatorName).Inc()
}
