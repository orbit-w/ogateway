package logger

import (
	mlog "github.com/orbit-w/meteor/modules/mlog"
)

/*
   @Author: orbit-w
   @File: logger
   @2024 4月 周日 14:34
*/

var logger = mlog.NewFileLogger(mlog.WithLevel("info"),
	mlog.WithFormat("console"),
	mlog.WithRotation(500, 7, 3, false),
	mlog.WithInitialFields(map[string]any{"app": "content-moderation"}),
	mlog.WithOutputPaths("logs/gateway.log"))

func SetLogger(log *mlog.Logger) {
	logger = log
}

func ZLogger() *mlog.Logger {
	return logger
}

func StopLogger() {
	_ = logger.Sync()
}
