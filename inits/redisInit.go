package inits

import (
	"context"
	"fmt"
	"forum/pkg/globals"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// RedisInit 初始化Redis
func RedisInit() {
	if err := viper.UnmarshalKey("redis", &globals.AppConfig.Redis); err != nil {
		globals.Log.Panicf("无法解码为结构: %s", err)
	}

	globals.RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", globals.AppConfig.Redis.Host, globals.AppConfig.Redis.Port),
		Password:     globals.AppConfig.Redis.Password,
		DB:           globals.AppConfig.Redis.DB,
		PoolSize:     globals.AppConfig.Redis.PoolSize,
		MinIdleConns: globals.AppConfig.Redis.MinIdleConns,
		IdleTimeout:  globals.AppConfig.Redis.IdleTimeout,
		DialTimeout:  globals.AppConfig.Redis.DialTimeout,
		ReadTimeout:  globals.AppConfig.Redis.ReadTimeout,
		WriteTimeout: globals.AppConfig.Redis.WriteTimeout,
		MaxRetries:   globals.AppConfig.Redis.MaxRetries,
	})
	ctx := context.Background()
	_, err := globals.RDB.Ping(ctx).Result()
	if err != nil {
		globals.Log.Panicf("Redis连接失败: %v", err)
	} else {
		globals.Log.Infof("Redis连接成功")
	}
}
