package globals

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// AppConfig 项目的mysql，redis配置
var AppConfig Config

// DB mysql链接
var DB *gorm.DB

// RDB redis链接
var RDB *redis.Client

// Log 日志记录
var Log = logrus.New()

// Router 总路由
var Router *gin.Engine
