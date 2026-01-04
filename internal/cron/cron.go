package cron

import (
	"context"
	"fmt"

	"github.com/padremortius/go-template-echo/internal/usecase/testcase"
	"github.com/padremortius/go-template-echo/pkgs/crontab"
	"github.com/padremortius/go-template-echo/pkgs/svclogger"
)

func LoadTasks(aCtx context.Context, ct crontab.Crontab, opts *crontab.CronOpts, log *svclogger.Log) {
	ctx, cancel := context.WithCancel(aCtx)
	defer cancel()

	taskCount := 0
	for _, job := range opts.Jobs {
		if !job.Disable {
			taskCount++
		}
	}
	if taskCount > 0 {
		log.Logger.Debug(fmt.Sprintf("taskCount = %v", taskCount))
		ct.WGroup.Add(taskCount)
		if !opts.Jobs[0].Disable {
			log.Logger.Info(fmt.Sprintf("Add new task. { Name: %v, Schedule: %v }", opts.Jobs[0].Name, opts.Jobs[0].Schedule))
			_, _ = ct.CronServer.AddFunc(opts.Jobs[0].Schedule, func() {
				testcase.RunTask(ctx, log)
			})
		}
	}
}
