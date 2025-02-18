package csrfMW

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// 设置 CSRF token 到 Redis 并设置过期时间
func setCSRFTokenToRedis(rdb *redis.Client, token string, tokenExpires time.Duration) error {
	return rdb.Set(context.Background(), token, "unused", tokenExpires).Err()
}

// 获取 CSRF token 在 Redis 中的状态
func checkCSRFTokenStatus(rdb *redis.Client, token string) (string, error) {
	return rdb.Get(context.Background(), token).Result()
}

// 标记 CSRF token 为已使用
func markCSRFTokenAsUsed(rdb *redis.Client, token string) error {
	return rdb.Set(context.Background(), token, "used", 0).Err()
}
