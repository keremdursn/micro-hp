package metrics

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)

	// Login başarısızlıklarını sayan özel metrik
	LoginFailCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "login_fail_total",
			Help: "Total number of failed login attempts",
		},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(LoginFailCounter)
}

// PrometheusHandler Fiber ile uyumlu /metrics endpointi için handler döndürür
func PrometheusHandler() func(*fiber.Ctx) error {
	return adaptor.HTTPHandler(promhttp.Handler())
}
