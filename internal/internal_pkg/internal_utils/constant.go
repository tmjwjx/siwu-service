package internal_utils

import "time"

// 常量

type Kind int

const (
	User Kind = iota + 1
	Article
)

// 验证码信息常量
const (
	VerifyCodeLen      = 6               // 验证码长度
	VerifyCodeDuration = 5 * time.Minute // 验证码有效时长
	SendInterval       = 1 * time.Second // 发送的间隔时间
)
