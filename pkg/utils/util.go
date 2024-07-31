package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// 和配置数据库相关的结构体

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
}

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
