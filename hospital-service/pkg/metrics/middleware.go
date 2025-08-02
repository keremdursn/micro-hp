package metrics

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		method := c.Method()
		timer := prometheus.NewTimer(HttpRequestDuration.WithLabelValues(path, method))
		err := c.Next()
		status := fmt.Sprintf("%d", c.Response().StatusCode())
		HttpRequestsTotal.WithLabelValues(path, method, status).Inc()
		timer.ObserveDuration()
		return err
	}
}

