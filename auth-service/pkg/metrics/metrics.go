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
	LoginSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "login_success_total",
			Help: "Total number of successful login attempts",
		},
	)

	// Register endpointi için sayaçlar
	RegisterSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "register_success_total",
			Help: "Total number of successful register attempts",
		},
	)
	RegisterFailCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "register_fail_total",
			Help: "Total number of failed register attempts",
		},
	)

	// ForgotPassword endpointi için sayaçlar
	ForgotPasswordSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "forgot_password_success_total",
			Help: "Total number of successful forgot password attempts",
		},
	)
	ForgotPasswordFailCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "forgot_password_fail_total",
			Help: "Total number of failed forgot password attempts",
		},
	)

	// ResetPassword endpointi için sayaçlar
	ResetPasswordSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "reset_password_success_total",
			Help: "Total number of successful reset password attempts",
		},
	)
	ResetPasswordFailCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "reset_password_fail_total",
			Help: "Total number of failed reset password attempts",
		},
	)

	// RefreshToken endpointi için sayaçlar
	RefreshTokenSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "refresh_token_success_total",
			Help: "Total number of successful refresh token attempts",
		},
	)
	RefreshTokenFailCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "refresh_token_fail_total",
			Help: "Total number of failed refresh token attempts",
		},
	)

	// Rate limiting metrikleri
	RateLimitExceededCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "rate_limit_exceeded_total",
			Help: "Total number of rate limit exceeded requests",
		},
		[]string{"endpoint", "ip"},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(LoginFailCounter)
	prometheus.MustRegister(LoginSuccessCounter)
	prometheus.MustRegister(RegisterSuccessCounter)
	prometheus.MustRegister(RegisterFailCounter)
	prometheus.MustRegister(ForgotPasswordSuccessCounter)
	prometheus.MustRegister(ForgotPasswordFailCounter)
	prometheus.MustRegister(ResetPasswordSuccessCounter)
	prometheus.MustRegister(ResetPasswordFailCounter)
	prometheus.MustRegister(RefreshTokenSuccessCounter)
	prometheus.MustRegister(RefreshTokenFailCounter)
	prometheus.MustRegister(RateLimitExceededCounter)
}

// PrometheusHandler Fiber ile uyumlu /metrics endpointi için handler döndürür
func PrometheusHandler() func(*fiber.Ctx) error {
	return adaptor.HTTPHandler(promhttp.Handler())
}
