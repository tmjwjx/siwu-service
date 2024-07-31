package inits

import (
	"context"
	"fmt"
	"forum/pkg/utils"
	"github.com/go-redis/redis/v8"
	"log"
)

// RedisInit 初始化redis
func RedisInit() {
	utils.RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", utils.AppConfig.Redis.Host, utils.AppConfig.Redis.Port),
		Password: utils.AppConfig.Redis.Password,
		DB:       utils.AppConfig.Redis.DB,
	})

	ctx := context.Background()
	_, err := utils.RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
