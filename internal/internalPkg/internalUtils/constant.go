package internalUtils

import (
	"time"
)

// 验证码信息常量
const (
	VerifyCodeSubject           = "验证码"                                                            // 验证码主题
	VerifyCodeLen               = 6                                                                // 验证码长度
	CustomAlphabetVerifyCode    = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ" // 自定义字母表
	VerifyCodeEffectiveDuration = 10 * time.Minute                                                 // 验证码有效时长
	VerifyCodeCoolTime          = 1 * time.Minute                                                  // 发送验证码的冷却时间
)

// 用户默认常量
const (
	UserNameLen            int = 8            // 用户默认姓名长度（随机姓名）
	CustomAlphabetNickname     = "1234567890" // 自定义字母表
)

// 默认图片路径
// const (
//	// UserDefaultImage 默认用户头像路径
//	UserDefaultImage = "user_default_head_image.png"
//	// ArticleDefaultImage 默认文章图片路径
//	ArticleDefaultImage = "/images/article_default_image.png"
//	// TagDefaultImage 默认标签图片
//	TagDefaultImage = "/images/tag_default_image.png"
//	// AdvertisementDefaultImage 默认广告图片
//	AdvertisementDefaultImage = "/images/advertisement_default_image.png"
//	CommentDefaultImage       = ""
//	CategoryDefaultImage      = ""
// )

// 用于判断广告是否被使用
const (
	Use = iota
	NoUse
)

// 默认图片路径
var (
	// UserDefaultImage 默认用户头像路径
	UserDefaultImage = ""

	// ArticleDefaultImage 默认文章图片路径
	ArticleDefaultImage = ""

	// TagDefaultImage 默认标签图片
	TagDefaultImage = ""

	// AdvertisementDefaultImage 默认广告图片
	AdvertisementDefaultImage = ""

	// CommentDefaultImage 默认评论图片
	CommentDefaultImage = ""

	// CategoryDefaultImage 默认类目图片
	CategoryDefaultImage = ""
)
