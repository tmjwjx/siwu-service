package inits

import (
	"context"
	"fmt"
	"forum/pkg/globals"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"log"
)

// RedisInit 初始化redis
func RedisInit() {

	if err := viper.UnmarshalKey("redis", &globals.AppConfig.Redis); err != nil {
		log.Fatalf("无法解码为结构: %s", err)
	}

	globals.RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", globals.AppConfig.Redis.Host, globals.AppConfig.Redis.Port),
		Password: globals.AppConfig.Redis.Password,
		DB:       globals.AppConfig.Redis.DB,
	})

	ctx := context.Background()
	_, err := globals.RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis连接失败: %v", err)
	}
}
