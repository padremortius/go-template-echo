package httpserver

import (
	"strconv"
	"time"

	"github.com/VictoriaMetrics/metrics"
	echo "github.com/labstack/echo/v5"
)

var (
	reqCounter  = metrics.GetOrCreateCounter(`http_requests_total{path="$path", method="$method", status="$status"}`)
	reqDuration = metrics.GetOrCreateHistogram(`http_request_duration_seconds{path="$path", method="$method"}`)
)

// MetricsMiddleware собирает базовые HTTP-метрики с использованием github.com/VictoriaMetrics/metrics
func MetricsMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if mySkipper(c) {
				return next(c)
			}

			start := time.Now()
			err := next(c)
			duration := time.Since(start).Seconds()

			status := ""
			if resp, err := echo.UnwrapResponse(c.Response()); err == nil {
				status = strconv.Itoa(resp.Status)
			}
			method := c.Request().Method

			// ВАЖНО: c.Path() возвращает зарегистрированный шаблон пути (например, "/users/:id"),
			// а не реальный URL запроса. Это предотвращает проблему высокой кардинальности метрик.
			path := c.Path()

			// 1. Метрика количества запросов (Counter)
			counterName := `http_requests_total{method="` + method + `",path="` + path + `",status="` + status + `"}`
			metrics.GetOrCreateCounter(counterName).Inc()

			// 2. Метрика времени выполнения запроса (Histogram)
			// VictoriaMetrics использует автоматические бакеты, что эффективнее и проще в настройке,
			// чем стандартные Prometheus-бакеты с лейблами "le".
			histogramName := `http_request_duration_seconds{method="` + method + `",path="` + path + `"}`
			metrics.GetOrCreateHistogram(histogramName).Update(duration)

			return err
		}
	}
}

// MetricsHandler отдает метрики в формате Prometheus
func MetricsHandler(c *echo.Context) error {
	c.Response().Header().Set("Content-Type", "text/plain; version=0.0.4")
	metrics.WritePrometheus(c.Response(), true)
	return nil
}
