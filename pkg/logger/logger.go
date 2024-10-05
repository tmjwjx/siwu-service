package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

// GetEncoder
// @Description: 获取编码器
// @return       zapcore.Encoder
// @Author tianjiajie 2024-10-05 16:02:57
func GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

// GetLogWriter
// @Description: 获取日志写入器
// @param        logPath string
// @param        appName string
// @return       zapcore.WriteSyncer
// @Author tianjiajie 2024-10-05 16:03:01
func GetLogWriter(logPath, appName string) zapcore.WriteSyncer {
	currentDate := time.Now().Format("2006-01-02")
	fileName := fmt.Sprintf("./%s/%s-%s.log", logPath, appName, currentDate)
	file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	return zapcore.AddSync(file)
}