package inits

import (
	"forum/pkg/globals"
	"forum/pkg/logger"
	"github.com/sirupsen/logrus"
)

func LogInit(logPath, appName string) {

	globals.Log = logrus.New()

	fileHook := logger.FileDateHook{
		// file:     file,
		LogPath: logPath,
		// fileDate: fileDate,
		AppName: appName,
	}
	globals.Log.AddHook(&fileHook)

	// 包含调用者信息
	globals.Log.SetReportCaller(true)
}
