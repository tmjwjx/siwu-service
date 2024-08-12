package internal_utils

import (
	"time"
)

// 验证码信息常量
const (
	VerifyCodeSubject           = "验证码"            // 验证码主题
	VerifyCodeLen               = 6                // 验证码长度
	VerifyCodeEffectiveDuration = 10 * time.Minute // 验证码有效时长
	VerifyCodeCoolTime          = 1 * time.Minute  // 发送验证码的冷却时间
)

// 用户默认常量
const (
	UserNameLen int = 8 // 用户默认姓名长度（随机姓名）
)

// 默认图片路径
const (
	// UserDefaultImage 默认用户头像路径
	UserDefaultImage = "/images/user_default_head_image"
	// ArticleDefaultImage 默认文章图片路径
	ArticleDefaultImage = "/images/article_default_image"
	// TagDefaultImage 默认标签图片
	TagDefaultImage = "/images/tag_default_image"
	// AdvertisementDefaultImage 默认广告图片
	AdvertisementDefaultImage = "/images/advertisement_default_image"
)

// 用于判断广告是否被使用
const (
	Use = iota
	NoUse
)
