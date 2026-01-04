package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/padremortius/go-template-echo/internal/config"
	"github.com/padremortius/go-template-echo/internal/cron"
	v1 "github.com/padremortius/go-template-echo/internal/handlers/v1"
	"github.com/padremortius/go-template-echo/internal/storage"
	"github.com/padremortius/go-template-echo/pkgs/crontab"
	"github.com/padremortius/go-template-echo/pkgs/httpserver"
	"github.com/padremortius/go-template-echo/pkgs/svclogger"
)

func Run(aBuildNumber, aBuildTimeStamp, aGitBranch, aGitHash string) {
	log := svclogger.New("")
	appCfg, err := config.NewConfig(aBuildNumber, aBuildTimeStamp, aGitBranch, aGitHash)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Config error: %v", err))
		os.Exit(-1)
	}

	shutdownTimeout := appCfg.HTTP.Timeouts.Shutdown

	ctxParent, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Logger.Info(fmt.Sprintf("Start application. Version: %v", appCfg.Version.BuildVersion))

	log.ChangeLogLevel(ctxParent, appCfg.Log.Level)

	//init storage
	store, err := storage.New(ctxParent, appCfg.Storage.Path, log)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Storage error: %v", err))
		os.Exit(-1)
	}

	if err := store.InitDB(); err != nil {
		log.Logger.Error(fmt.Sprintf("Storage error: %v", err))
	}

	//Init crontab
	ctb := crontab.New(ctxParent, &appCfg.Crontab)
	cron.LoadTasks(ctxParent, ctb, &appCfg.Crontab, log)
	log.Logger.Info("Starting cron server")
	go ctb.StartCron()

	// HTTP Server
	log.Logger.Info(fmt.Sprintf("Start web-server on port %v", appCfg.HTTP.Port))

	httpServer := httpserver.New(ctxParent, log, &appCfg.HTTP)
	httpserver.InitBaseRouter(httpServer.Handler, *appCfg, appCfg.Version)
	v1.InitAppRouter(httpServer.Handler, *appCfg, *log, *store)
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Logger.Info(fmt.Sprintf("app - Run - signal: %v", s.String()))
	case err := <-httpServer.Notify():
		log.Logger.Error(fmt.Sprintf("app - Run - httpServer.Notify: %v", err))
	}

	// Shutdown
	log.Logger.Info("Waiting for stop crontab")
	ctb.StopCron()
	if err := httpServer.Shutdown(shutdownTimeout); err != nil {
		log.Logger.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %v", err))
	}
}
