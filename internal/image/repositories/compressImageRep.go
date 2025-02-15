package repositories

import (
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	cacheDuration = 24 * time.Hour
)

// GetFromCache 查询redis缓存中，是否存在该压缩过的图片
func GetFromCache(rdb *redis.Client, key string) ([]byte, error) {
	return rdb.Get(rdb.Context(), key).Bytes()
}

// SetToCache 将压缩过的图片存储到redis缓存中
func SetToCache(rdb *redis.Client, key string, data []byte) error {
	return rdb.Set(rdb.Context(), key, data, cacheDuration).Err()
}
