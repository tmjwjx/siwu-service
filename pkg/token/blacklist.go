package token

import (
	"context"
	"errors"
	"forum/pkg/globals"
	"github.com/go-redis/redis/v8"
	"time"
)

// AddTokenToBlacklist 把 token 加入 Redis 黑名单
func AddTokenToBlacklist(rdb *redis.Client, tokenString string, expiration time.Duration) error {
	ctx := context.Background()
	return rdb.Set(ctx, tokenString, "blacklisted", expiration).Err()
}

// IsTokenBlacklisted 检查 Token 是否在黑名单
func IsTokenBlacklisted(rdb *redis.Client, tokenString string) bool {
	ctx := context.Background()
	s, err := rdb.Get(ctx, tokenString).Result()

	if errors.Is(err, redis.Nil) {
		globals.Log.Infof("Token 不在黑名单:%v", tokenString)
		return false
	} else if err != nil {
		globals.Log.Infof("Redis 查询失败:%v", err)
		return false
	}

	globals.Log.Infof("Token 在黑名单:%v", s)
	return true
}
