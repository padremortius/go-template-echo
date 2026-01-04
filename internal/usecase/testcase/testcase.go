package testcase

import (
	"context"

	"github.com/padremortius/go-template-echo/pkgs/svclogger"
)

func RunTask(actx context.Context, alog *svclogger.Log) {
	alog.Logger.Info("Start task 'Test usecase'")
	//
	alog.Logger.Info("End task 'Test usecase'")
}
