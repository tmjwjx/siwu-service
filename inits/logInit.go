package inits

import "forum/pkg/utils"

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

	fileHook := utils.FileDateHook{
		// file:     file,
		LogPath: logPath,
		// fileDate: fileDate,
		AppName: appName,
	}
	utils.Log.AddHook(&fileHook)

	// 包含调用者信息
	utils.Log.SetReportCaller(true)
}
