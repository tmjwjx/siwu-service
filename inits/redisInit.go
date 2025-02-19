package inits

import (
	"context"
	"errors"
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

// RedisEmailTaskStreamInit 创建 Redis 发送邮件 Stream 消费者组
func RedisEmailTaskStreamInit() {
	ctx := context.Background()
	// 创建消费者组，如果已经存在则忽略 BUSYGROUP 错误
	err := globals.RDB.XGroupCreateMkStream(ctx, globals.EmailStreamKey, globals.EmailGroupKey, "$").Err()
	if err != nil && !errors.Is(err, redis.Nil) && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		globals.Log.Panicf("Failed to create consumer group: %v", err)
	} else if err == nil {
		globals.Log.Infof("Consumer group '%s' created successfully.", globals.EmailGroupKey)
	} else {
		globals.Log.Infof("Consumer group '%s' already exists, continuing...", globals.EmailGroupKey)
	}
}
