package inits

import (
	"forum/internal/internal_pkg/internal_utils"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"log"
)

// SendEmailCfgInit 初始化发送邮件配置
func SendEmailCfgInit() {
	if err := viper.UnmarshalKey("send_email", &globals.AppConfig.SendEmailCfg); err != nil {
		log.Fatalf("无法解码为结构: %s", err)
	}

	globals.SendEmailCfg = &globals.AppConfig.SendEmailCfg
	if !internal_utils.IsValidEmail(globals.SendEmailCfg.From) {
		log.Fatalln("SendEmailCfgInit() err: 发送者邮箱错误")
	}
	// fmt.Println(globals.SendEmailCfg) // &{3174285493@qq.com 验证码 6 100 1 smtp.qq.com 587 3174285493@qq.com mmureuzrdnmndfef}

}
