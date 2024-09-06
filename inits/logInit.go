package inits

import (
	"forum/pkg/globals"
	"forum/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func LogInit(logPath, appName string) {
	writeSyncer := logger.GetLogWriter(logPath, appName)
	encoder := logger.GetEncoder()

	// 新增部分：将日志输出到控制台
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.ErrorLevel)

	// 新增部分：将日志输出到文件
	fileCore := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// 修改部分：合并控制台输出和文件输出
	core := zapcore.NewTee(consoleCore, fileCore)

	logger := zap.New(core, zap.AddCaller())

	globals.Log = logger.Sugar()
}
