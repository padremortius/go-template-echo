package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/padremortius/go-template-echo/internal/config"
	"github.com/padremortius/go-template-echo/internal/crontab"
	v1 "github.com/padremortius/go-template-echo/internal/handlers/v1"
	"github.com/padremortius/go-template-echo/internal/storage"
	"github.com/padremortius/go-template-echo/pkgs/httpserver"
	"github.com/padremortius/go-template-echo/pkgs/svclogger"
)

func Run(ver config.Version) {
	log := svclogger.New("")
	appCfg, err := config.NewConfig()
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Config error: %v", err))
	}
	appCfg.Version = ver
	shutdownTimeout := appCfg.HTTP.Timeouts.Shutdown

	ctxParent, cancel := context.WithCancel(context.Background())
	defer cancel()

	log.Logger.Info(fmt.Sprintf("Start application. Version: %v", appCfg.Version.Version))

	log.ChangeLogLevel(ctxParent, appCfg.Log.Level)

	//init storage
	store, err := storage.New(ctxParent, appCfg.Storage.Path, log)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("Storage error: %v", err))
	}

	if err := store.InitDB(); err != nil {
		log.Logger.Error(fmt.Sprintf("Storage error: %v", err))
	}

	//Init crontab
	ctb := crontab.New(ctxParent, log, &appCfg.Crontab)
	ctb.LoadTasks(ctxParent, &appCfg.Crontab)
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
	ctb.StopCron()
	if err := httpServer.Shutdown(shutdownTimeout); err != nil {
		log.Logger.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %v", err))
	}
}
