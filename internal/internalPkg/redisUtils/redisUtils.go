package redisUtils

import (
	"context"
	"errors"
	"fmt"
	"forum/pkg/globals"
	"github.com/go-redis/redis/v8"
	"time"
)

// Set 设置缓存
func Set(rdb *redis.Client, ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// 设置缓存值，使用指定的过期时间
	err := rdb.Set(ctx, key, value, expiration).Err()
	if err != nil {
		globals.Log.Errorf("设置 Redis 键 %s 时出错: %v", key, err)
		return err
	}
	return nil
}

// Get 获取缓存
func Get(rdb *redis.Client, ctx context.Context, key string) (string, error) {
	// 获取缓存值
	val, err := rdb.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 键不存在
		return "", nil
	}
	if err != nil {
		globals.Log.Errorf("获取 Redis 键 %s 时出错: %v", key, err)
		return "", err
	}
	return val, nil
}

// Delete 删除缓存
func Delete(rdb *redis.Client, ctx context.Context, key string) error {
	// 删除缓存值
	err := rdb.Del(ctx, key).Err()
	if err != nil {
		globals.Log.Errorf("删除 Redis 键 %s 时出错: %v", key, err)
		return err
	}
	return nil
}

// Exists 检查缓存是否存在
func Exists(rdb *redis.Client, ctx context.Context, key string) (bool, error) {
	// 检查键是否存在
	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		globals.Log.Errorf("检查 Redis 键 %s 是否存在时出错: %v", key, err)
		return false, err
	}
	return exists > 0, nil
}

// Increment 增加缓存值
func Increment(rdb *redis.Client, ctx context.Context, key string, increment int64) (int64, error) {
	// 增加缓存值
	newValue, err := rdb.IncrBy(ctx, key, increment).Result()
	if err != nil {
		globals.Log.Errorf("增加 Redis 键 %s 的值时出错: %v", key, err)
		return 0, err
	}
	return newValue, nil
}

// Decrement 递减缓存值
func Decrement(rdb *redis.Client, ctx context.Context, key string, decrement int64) (int64, error) {
	// 执行递减操作
	val, err := rdb.DecrBy(ctx, key, decrement).Result()
	if err != nil {
		globals.Log.Errorf("递减 Redis 键 %s 时出错: %v", key, err)
		return 0, err
	}
	return val, nil
}

// HSet 设置哈希表中的字段值
func HSet(rdb *redis.Client, ctx context.Context, hash, key string, value interface{}) error {
	// 向哈希表中设置字段值
	err := rdb.HSet(ctx, hash, key, value).Err()
	if err != nil {
		globals.Log.Errorf("设置 Redis 哈希表 %s 中的字段 %s 时出错: %v", hash, key, err)
		return err
	}
	return nil
}

// HGet 获取哈希表中的字段值
func HGet(rdb *redis.Client, ctx context.Context, hash, key string) (string, error) {
	// 获取哈希表中字段的值
	val, err := rdb.HGet(ctx, hash, key).Result()
	if errors.Is(err, redis.Nil) {
		// 字段不存在
		err = rdb.HSet(ctx, hash, key, 0).Err()
		if err != nil {
			globals.Log.Errorf("设置哈希表 %s 中字段 %s 的默认值时出错: %v", hash, key, err)
			return "", fmt.Errorf("内部错误: %v", err)
		}
	}
	if err != nil {
		globals.Log.Errorf("获取 Redis 哈希表 %s 中的字段 %s 时出错: %v", hash, key, err)
		return "", err
	}
	return val, nil
}

// HDel 删除哈希表中的字段
func HDel(rdb *redis.Client, ctx context.Context, hash, key string) error {
	// 删除哈希表中的字段
	err := rdb.HDel(ctx, hash, key).Err()
	if err != nil {
		globals.Log.Errorf("删除 Redis 哈希表 %s 中的字段 %s 时出错: %v", hash, key, err)
		return err
	}
	return nil
}

// IncrementHash 增加哈希表中字段的值
func IncrementHash(rdb *redis.Client, ctx context.Context, hash, field string, increment int64) (int64, error) {
	// 使用 HIncrBy 增加哈希表字段的值
	newValue, err := rdb.HIncrBy(ctx, hash, field, increment).Result()
	if err != nil {
		globals.Log.Errorf("增加 Redis 哈希表 %s 字段 %s 的值时出错: %v", hash, field, err)
		return 0, err
	}
	return newValue, nil
}

// DecrementHash 减少哈希表中字段的值
func DecrementHash(rdb *redis.Client, ctx context.Context, hash, field string, decrement int64) (int64, error) {
	// 使用 HIncrBy 减少哈希表字段的值（传递负数）
	newValue, err := rdb.HIncrBy(ctx, hash, field, -decrement).Result()
	if err != nil {
		globals.Log.Errorf("减少 Redis 哈希表 %s 字段 %s 的值时出错: %v", hash, field, err)
		return 0, err
	}
	return newValue, nil
}

// LPush 向列表左侧推送数据
func LPush(rdb *redis.Client, ctx context.Context, key string, values ...interface{}) error {
	// 向列表左侧推送元素
	err := rdb.LPush(ctx, key, values...).Err()
	if err != nil {
		globals.Log.Errorf("向 Redis 列表 %s 左侧推送数据时出错: %v", key, err)
		return err
	}
	return nil
}

// RPush 向列表右侧推送数据
func RPush(rdb *redis.Client, ctx context.Context, key string, values ...interface{}) error {
	// 向列表右侧推送元素
	err := rdb.RPush(ctx, key, values...).Err()
	if err != nil {
		globals.Log.Errorf("向 Redis 列表 %s 右侧推送数据时出错: %v", key, err)
		return err
	}
	return nil
}

// LPop 从列表左侧弹出数据
func LPop(rdb *redis.Client, ctx context.Context, key string) (string, error) {
	// 从列表左侧弹出元素
	val, err := rdb.LPop(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 列表为空
		return "", fmt.Errorf("Redis 列表 %s 中没有数据", key)
	}
	if err != nil {
		globals.Log.Errorf("从 Redis 列表 %s 左侧弹出数据时出错: %v", key, err)
		return "", err
	}
	return val, nil
}

// RPop 从列表右侧弹出数据
func RPop(rdb *redis.Client, ctx context.Context, key string) (string, error) {
	// 从列表右侧弹出元素
	val, err := rdb.RPop(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		// 列表为空
		return "", fmt.Errorf("Redis 列表 %s 中没有数据", key)
	}
	if err != nil {
		globals.Log.Errorf("从 Redis 列表 %s 右侧弹出数据时出错: %v", key, err)
		return "", err
	}
	return val, nil
}
