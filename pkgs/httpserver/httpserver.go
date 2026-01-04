package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/padremortius/go-template-echo/pkgs/svclogger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	slogecho "github.com/sse-open/slog-echo"

	echoPrometheus "github.com/globocom/echo-prometheus"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HTTP struct {
	Cors struct {
		Headers []string `yaml:"headers" json:"headers" validate:"required"`
		Methods []string `yaml:"methods" json:"methods" validate:"required"`
		Origins []string `yaml:"origins" json:"origins" validate:"required"`
	} `yaml:"cors" json:"cors"`
	Port     string `yaml:"port" json:"port"`
	Timeouts struct {
		Read     time.Duration `yaml:"read" json:"read"`
		Write    time.Duration `yaml:"write" json:"write"`
		Idle     time.Duration `yaml:"idle" json:"idle"`
		Shutdown time.Duration `yaml:"shutdown" json:"shutdown"`
	} `yaml:"timeouts" json:"timeouts"`
}

type Server struct {
	ctx     context.Context
	server  *http.Server
	Handler *echo.Echo
	notify  chan error
}

// New -.
func New(c context.Context, log *svclogger.Log, opts *HTTP) *Server {
	handler := echo.New()
	//logger middleware
	handler.Use(slogecho.New(log.Logger).
		WithFilter(slogecho.IgnorePath("/health", "/env", "/info", "/favicon.ico", "/prometheus")).
		EchoMiddleware())

	// recovery middleware
	handler.Use(middleware.Recover())

	// CORS settings
	corsMW := middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     opts.Cors.Headers,
		AllowMethods:     opts.Cors.Methods,
		AllowOrigins:     opts.Cors.Origins,
		Skipper:          mySkipper,
	})
	handler.Use(corsMW)

	//metrics settings
	handler.Use(echoPrometheus.MetricsMiddleware()) // adds middleware to gather metrics
	handler.GET("/prometheus", echo.WrapHandler(promhttp.Handler()))

	s := &Server{
		server: &http.Server{
			Handler:      handler,
			IdleTimeout:  opts.Timeouts.Idle,
			ReadTimeout:  opts.Timeouts.Read,
			WriteTimeout: opts.Timeouts.Write,
			Addr:         fmt.Sprint(":", opts.Port),
		},
		notify:  make(chan error, 1),
		Handler: handler,
		ctx:     c,
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown(shutdownTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(s.ctx, shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
