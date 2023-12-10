package testcase

import (
	"context"
	"go-template-echo/internal/svclogger"
)

func RunTask(actx context.Context, alog *svclogger.Log) {
	alog.Logger.Info().Msgf("Start task 'Test usecase'")
	//
	alog.Logger.Info().Msgf("End task 'Test usecase'")
}
