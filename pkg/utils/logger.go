package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

var Log = logrus.New()

type FileDateHook struct {
	file     *os.File
	logPath  string
	fileDate string // 判断日期切换目录
	appName  string
}

func (hook *FileDateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (hook *FileDateHook) Fire(entry *logrus.Entry) error {
	now := time.Now()
	timerDate := now.Format("2006-01-02")
	line, _ := entry.String()

	// 如果日期不同，关闭当前文件，创建新文件
	if hook.fileDate != timerDate {
		if hook.file != nil {
			hook.file.Close()
		}

		// 更新日期
		hook.fileDate = timerDate

		// //创建目录（按日期分目录）
		// dirPath := hook.logPath
		// if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
		//	return err
		// }

		// 创建新文件（按天分文件）
		filename := fmt.Sprintf("%s/%s-%s.log", hook.logPath, hook.appName, timerDate)
		var err error
		hook.file, err = os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
		if err != nil {
			return err
		}
	}

	// 写入日志
	if _, err := hook.file.Write([]byte(line)); err != nil {
		return err
	}
	return nil
}

// InitFile 初始化日志文件
func InitFile(logPath, appName string) {
	// now := time.Now()
	// fileDate := now.Format("2006-01-02")

	// 创建文件（按天分文件）
	// filename := fmt.Sprintf("%s/%s-%s.log", logPath, appName, fileDate)
	// file, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	// if err != nil {
	//	Log.Error(err)
	//	return
	// }

	fileHook := FileDateHook{
		// file:     file,
		logPath: logPath,
		// fileDate: fileDate,
		appName: appName,
	}
	Log.AddHook(&fileHook)

	// 包含调用者信息
	Log.SetReportCaller(true)
}
