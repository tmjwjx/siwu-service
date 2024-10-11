package globals

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (

	// AppConfig 项目总配置
	AppConfig Config

	// SConfig 静态文件配置
	SConfig StaticConfig

	// DB mysql链接
	DB *gorm.DB

	// RDB redis链接
	RDB *redis.Client

	// Log 日志记录
	Log *zap.SugaredLogger

	// Router 总路由
	Router *gin.Engine

	// Env 环境配置文件
	Env string

	// SendEmailCfg 发送邮件配置
	SendEmailCfg *SendEmailConfig

	// SubscriberChannels 全局 map，用于存储各文章 ID 的 SSE 订阅者通道
	SubscriberChannels = make(map[string]chan string)

	// // SystemMsgSubject 系统消息主题
	// SystemMsgSubject = observer.NewSystemMsgSubject()
)
