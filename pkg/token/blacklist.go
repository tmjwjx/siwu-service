package token

import (
	"context"
	"errors"
	"fmt"
	"forum/pkg/globals"
	"github.com/go-redis/redis/v8"
	"time"
)

// BlacklistKey 定义黑名单的 Redis Sorted Set 键名
const BlacklistKey = "blacklist:tokens"

// AddTokenToBlacklist 将 token 加入 Redis 黑名单（Sorted Set）
func AddTokenToBlacklist(rdb *redis.Client, tokenString string, expiration time.Duration) error {
	ctx := context.Background()

	// 计算过期时间戳
	expireAt := time.Now().Add(expiration).Unix()

	// 将 token 添加到 Sorted Set，分数为过期时间戳
	err := rdb.ZAdd(ctx, BlacklistKey, &redis.Z{
		Score:  float64(expireAt),
		Member: tokenString,
	}).Err()
	if err != nil {
		return fmt.Errorf("添加 token 到黑名单失败: %v", err)
	}

	globals.Log.Infof("Token 已加入黑名单: %v, 过期时间: %v", tokenString, time.Unix(expireAt, 0))
	return nil
}

// IsTokenBlacklisted 检查 token 是否在黑名单中
func IsTokenBlacklisted(rdb *redis.Client, tokenString string) bool {
	ctx := context.Background()

	// 获取 token 的过期时间戳
	score, err := rdb.ZScore(ctx, BlacklistKey, tokenString).Result()
	if errors.Is(err, redis.Nil) {
		// token 不存在
		globals.Log.Infof("Token 不在黑名单: %v", tokenString)
		return false
	} else if err != nil {
		// 查询失败
		globals.Log.Errorf("Redis 查询失败: %v", err)
		return false
	}

	// 检查是否过期
	currentTime := time.Now().Unix()
	if currentTime > int64(score) {
		// 已过期，异步移除
		go rdb.ZRem(ctx, BlacklistKey, tokenString)
		globals.Log.Infof("Token 已过期并移除: %v", tokenString)
		return false
	}

	globals.Log.Infof("Token 在黑名单中: %v, 过期时间: %v", tokenString, time.Unix(int64(score), 0))
	return true
}

// RemoveTokenFromBlacklist 从黑名单中移除 token
func RemoveTokenFromBlacklist(rdb *redis.Client, tokenString string) error {
	ctx := context.Background()

	err := rdb.ZRem(ctx, BlacklistKey, tokenString).Err()
	if err != nil {
		return fmt.Errorf("从黑名单移除 token 失败: %v", err)
	}

	globals.Log.Infof("Token 已从黑名单移除: %v", tokenString)
	return nil
}

// CleanExpiredTokens 清理过期的 token
func CleanExpiredTokens(rdb *redis.Client) error {
	ctx := context.Background()

	currentTime := time.Now().Unix()
	// 移除分数小于当前时间戳的成员（即已过期的 token）
	_, err := rdb.ZRemRangeByScore(ctx, BlacklistKey, "0", fmt.Sprintf("%d", currentTime)).Result()
	if err != nil {
		return fmt.Errorf("清理过期 token 失败: %v", err)
	}

	globals.Log.Infof("已清理过期的黑名单 token")
	return nil
}

// GetAllBlacklistedTokens 获取所有未过期的黑名单 token
func GetAllBlacklistedTokens(rdb *redis.Client) ([]string, error) {
	ctx := context.Background()

	// 获取所有 token
	tokens, err := rdb.ZRange(ctx, BlacklistKey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("获取黑名单 token 失败: %v", err)
	}

	// 过滤未过期的 token
	var validTokens []string
	currentTime := time.Now().Unix()
	for _, token := range tokens {
		score, err := rdb.ZScore(ctx, BlacklistKey, token).Result()
		if err == nil && currentTime <= int64(score) {
			validTokens = append(validTokens, token)
		}
	}

	return validTokens, nil
}
