package internal_utils

import (
	"time"
)

// 验证码信息常量
const (
	Form          = "3174285493@qq.com" // 发送人
	Subject       = "验证码"               // 主题
	AuthorizeCode = "mmureuzrdnmndfef"  // 授权码

	VerifyCodeLen               = 6                // 验证码长度
	VerifyCodeEffectiveDuration = 10 * time.Minute // 验证码有效时长
	VerifyCodeCoolTime          = 1 * time.Minute  // 发送验证码的冷却时间
)
