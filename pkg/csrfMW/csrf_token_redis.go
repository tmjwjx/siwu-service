package csrfMW

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

// SetCSRFTokenToRedis 设置 CSRF token 到 Redis 并设置过期时间
func SetCSRFTokenToRedis(rdb *redis.Client, token string, exp time.Duration) error {
	return rdb.Set(context.Background(), "csrf:"+token, "unused", exp).Err()
}

// 获取 CSRF token 在 Redis 中的状态
func checkCSRFTokenStatus(rdb *redis.Client, token string) (string, error) {
	ctx := context.Background()
	key := "csrf:" + token

	ttl, err := rdb.TTL(ctx, key).Result()
	if err != nil {
		return "", err
	}
	if ttl <= 0 {
		return "", redis.Nil
	}

	status, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return status, nil
}

// 标记 CSRF token 为已使用
func markCSRFTokenAsUsed(rdb *redis.Client, token string) error {
	return rdb.Set(context.Background(), "csrf:"+token, "used", 0).Err()
}
