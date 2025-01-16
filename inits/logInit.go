package inits

import (
	"forum/pkg/globals"
	"forum/pkg/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func LogInit() {

	if err := viper.UnmarshalKey("log", &globals.AppConfig.Log); err != nil {
		globals.Log.Panicf("无法解码为结构: %s", err)
	}

	//level := globals.AppConfig.Log.Level
	logPath := globals.AppConfig.Log.LogPath
	appName := globals.AppConfig.Log.AppName

	writeSyncer := logger.GetLogWriter(logPath, appName)
	encoder := logger.GetEncoder()

	// 新增部分：将日志输出到控制台
	consoleCore := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)

	// 新增部分：将日志输出到文件
	fileCore := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	// 修改部分：合并控制台输出和文件输出
	core := zapcore.NewTee(consoleCore, fileCore)

	log := zap.New(core, zap.AddCaller())

	globals.Log = log.Sugar()
}
