package testcase

import (
	"context"

	"github.com/padremortius/go-template-echo/pkgs/svclogger"
)

func RunTask(actx context.Context, alog *svclogger.Log) {
	alog.Info("Start task 'Test usecase'")
	//
	alog.Info("End task 'Test usecase'")
}
