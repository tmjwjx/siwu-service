package inits

import (
	"fmt"
	"forum/internal/internal_pkg/internal_utils"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"log"
)

// VerifyCodeInit 初始化验证码配置
func VerifyCodeInit() {
	if err := viper.UnmarshalKey("verifyCode", &globals.AppConfig.VerifyCodeConfig); err != nil {
		log.Fatalf("无法解码为结构: %s", err)
	}

	globals.VerifyCode = &globals.AppConfig.VerifyCodeConfig
	if !internal_utils.IsValidEmail(globals.VerifyCode.From) {
		log.Fatalln("VerifyCodeInit() err: 发送者邮箱错误")
	}
	fmt.Println(globals.VerifyCode) // &{3174285493@qq.com 验证码 6 100 1 smtp.qq.com 587 3174285493@qq.com mmureuzrdnmndfef}

}
