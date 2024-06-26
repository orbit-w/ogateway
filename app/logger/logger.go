package logger

import (
	"github.com/orbit-w/golib/bases/logger"
	"go.uber.org/zap"
)

/*
   @Author: orbit-w
   @File: logger
   @2024 4月 周日 14:34
*/

var gLogger *zap.Logger

func InitLogger() {
	gLogger = logger.New("logs/o_gateway.log", zap.InfoLevel)
}

func StopLogger() {
	logger.Stop(gLogger)
}

func ZLogger() *zap.Logger {
	return gLogger
}
