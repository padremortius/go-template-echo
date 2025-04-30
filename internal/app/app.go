package app

import (
	"context"
	"go-template-echo/internal/config"
	"go-template-echo/internal/crontab"
	"go-template-echo/internal/handlers/actuators"
	v1 "go-template-echo/internal/handlers/v1"
	"go-template-echo/internal/httpserver"
	"go-template-echo/internal/storage/sqlite"
	"go-template-echo/internal/svclogger"
	"os"
	"os/signal"
	"syscall"
)

func Run(aBuildNumber, aBuildTimeStamp, aGitBranch, aGitHash string) {
	log := svclogger.New("")
	appCfg, err := config.NewConfig()
	if err != nil {
		log.Logger.Fatal().Msgf("Config error: %v", err)
	}
	appCfg.Version = *config.InitVersion(aBuildNumber, aBuildTimeStamp, aGitBranch, aGitHash)
	shutdownTimeout := appCfg.HTTP.Timeouts.Shutdown

	ctxParent, cancel := context.WithCancel(context.Background())
	ctxLogger := context.WithValue(ctxParent, "log", log)
	defer cancel()

	log.Logger.Info().Msgf("Start application. Version: %v", appCfg.Version.Version)

	log.ChangeLogLevel(appCfg.Log.Level)

	//init storage
	storage, err := sqlite.New(ctxLogger, appCfg.Storage.Path, log)
	if err != nil {
		log.Logger.Fatal().Msgf("Storage error: %v", err)
	}

	if err := storage.InitDB(); err != nil {
		log.Logger.Fatal().Msgf("Storage error: %v", err)
	}

	ctxDb := context.WithValue(ctxLogger, "db", storage)

	//Init crontab
	ctb := crontab.New(ctxDb, log, &appCfg.Crontab)
	ctb.LoadTasks(ctxParent, &appCfg.Crontab)
	go ctb.StartCron()

	// HTTP Server
	log.Logger.Info().Msg("Start web-server on port " + appCfg.HTTP.Port)

	httpServer := httpserver.New(ctxLogger, log, &appCfg.HTTP)
	actuators.InitBaseRouter(httpServer.Handler)
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
