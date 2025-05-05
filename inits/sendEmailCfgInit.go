package inits

import (
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"github.com/spf13/viper"
)

// SendEmailCfgInit
// @Description: 初始化发送邮件配置
// @Author lizhuang 2024-10-10 21:58:06
func SendEmailCfgInit() {
	if err := viper.UnmarshalKey("send_email", &globals.AppConfig.SendEmailCfg); err != nil {
		globals.Log.Panicf("无法解码为结构: %s", err)
	}

	globals.SendEmailCfg = &globals.AppConfig.SendEmailCfg
	if !internalUtils.IsValidEmail(globals.SendEmailCfg.From) {
		globals.Log.Panicf("SendEmailCfgInit() err: 发送者邮箱错误")
	}
}
