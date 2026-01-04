package testcase

import (
	"context"

	"github.com/padremortius/go-template-echo/pkgs/svclogger"
)

func RunTask(actx context.Context, alog *svclogger.Log) {
	alog.Logger.Info().Msgf("Start task 'Test usecase'")
	//
	alog.Logger.Info().Msgf("End task 'Test usecase'")
}
