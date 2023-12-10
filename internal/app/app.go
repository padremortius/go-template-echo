package app

import (
	"context"
	"go-template-echo/internal/config"
	"go-template-echo/internal/controller/baserouting"
	v1 "go-template-echo/internal/controller/v1"
	"go-template-echo/internal/crontab"
	"go-template-echo/internal/httpserver"
	"go-template-echo/internal/svclogger"
	"os"
	"os/signal"
	"syscall"
)

func Run() {
	log := svclogger.New("")

	if err := config.NewConfig(); err != nil {
		log.Logger.Fatal().Msgf("Config error: %v", err)
	}

	shutdownTimeout := config.Cfg.HTTP.Timeouts.Shutdown

	ctxParent, cancel := context.WithCancel(context.Background())
	ctxLogger := context.WithValue(ctxParent, "log", log)
	ctxProfile := context.WithValue(ctxLogger, "profile", &config.Cfg.ProfileName)
	defer cancel()

	log.Logger.Info().Msgf("Start application. Version: %v", config.Cfg.Version.Version)

	log.ChangeLogLevel(config.Cfg.Log.Level)

	//Init crontab
	ctb := crontab.New(ctxProfile, log, &config.Cfg.Crontab)
	ctb.LoadTasks(ctxParent, &config.Cfg.Crontab)
	go ctb.StartCron()

	// HTTP Server
	log.Logger.Info().Msg("Start web-server on port " + config.Cfg.HTTP.Port)

	httpServer := httpserver.New(ctxProfile, log, &config.Cfg.HTTP)
	baserouting.InitBaseRouter(httpServer.Handler)
	v1.InitAppRouter(httpServer.Handler)
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Logger.Info().Msgf("app - Run - signal: %v", s.String())
	case err := <-httpServer.Notify():
		log.Logger.Error().Msgf("app - Run - httpServer.Notify: %v", err)
	}

	// Shutdown
	ctb.StopCron()
	if err := httpServer.Shutdown(shutdownTimeout); err != nil {
		log.Logger.Error().Msgf("app - Run - httpServer.Shutdown: %v", err)
	}
}
