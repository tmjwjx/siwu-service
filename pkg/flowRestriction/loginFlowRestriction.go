// // // 固定窗口记录失败次数
// //
// // package flowRestriction
// //
// // import (
// // 	"errors"
// // 	"fmt"
// // 	"github.com/gin-gonic/gin"
// // 	"github.com/go-redis/redis/v8"
// // 	"math"
// // 	"time"
// // )
// //
// // // 登录限流相关常量
// // const (
// // 	MaxAttemptsBeforeLock = 3  // 开始限制前的最大尝试次数
// // 	SingleFailureTTL      = 10 // 单次失败记录的过期时间（分钟）
// // 	BaseLockTime          = 1  // 基础锁定时间（分钟）
// // 	MaxLockTime           = 24 // 最大锁定时间（小时）
// //
// // 	// Redis键前缀
// // 	LoginFailCountKeyPrefix = "login_fail_count:%s" // 失败次数计数键
// // 	LoginLockKeyPrefix      = "login_lock:%s"       // 锁定状态键
// // 	LockValue               = "locked"              // 锁定键的值
// // )
// //
// // // CalculateLockTime 计算锁定时间（指数增长）
// // func CalculateLockTime(failCount int) time.Duration {
// // 	if failCount < MaxAttemptsBeforeLock {
// // 		return 0 // 未达到限制次数时不锁定
// // 	}
// // 	// 从第3次失败开始，锁定时间指数增长：第3次1分钟，第4次2分钟，第5次4分钟...
// // 	baseTime := time.Duration(BaseLockTime) * time.Minute
// // 	waitTime := baseTime * time.Duration(math.Pow(2, float64(failCount-MaxAttemptsBeforeLock)))
// // 	maxWait := time.Duration(MaxLockTime) * time.Hour
// // 	if waitTime > maxWait {
// // 		return maxWait
// // 	}
// // 	return waitTime
// // }
// //
// // // CheckRateLimit 检查登录限制并返回剩余次数
// // func CheckRateLimit(rdb *redis.Client, c *gin.Context) (bool, int, string) {
// // 	ip := c.ClientIP()
// // 	countKey := fmt.Sprintf(LoginFailCountKeyPrefix, ip)
// // 	lockKey := fmt.Sprintf(LoginLockKeyPrefix, ip)
// //
// // 	// 检查是否处于锁定状态
// // 	if locked, err := rdb.Exists(c, lockKey).Result(); err == nil && locked > 0 {
// // 		ttl, _ := rdb.TTL(c, lockKey).Result()
// // 		if ttl > 0 {
// // 			// 向上取整
// // 			remainingMinutes := int(math.Ceil(ttl.Minutes()))
// // 			return false, 0, fmt.Sprintf("登录尝试过于频繁，请在%d分钟后再试", remainingMinutes)
// // 			// return false, 0, fmt.Sprintf("登录尝试过于频繁，请在%d分钟后再试", int(ttl.Minutes()))
// // 		}
// // 	}
// //
// // 	// 获取当前失败次数
// // 	count, err := rdb.Get(c, countKey).Int()
// // 	if err != nil && !errors.Is(err, redis.Nil) {
// // 		return false, 0, "系统错误"
// // 	}
// //
// // 	// 计算剩余尝试次数（基于3次限制）
// // 	remaining := MaxAttemptsBeforeLock - count
// // 	if remaining < 0 {
// // 		remaining = 0
// // 	}
// //
// // 	return true, remaining, ""
// // }
// //
// // // RecordFailedAttempt 记录登录失败
// // func RecordFailedAttempt(rdb *redis.Client, c *gin.Context) {
// // 	ip := c.ClientIP()
// // 	countKey := fmt.Sprintf(LoginFailCountKeyPrefix, ip)
// // 	lockKey := fmt.Sprintf(LoginLockKeyPrefix, ip)
// //
// // 	// 增加失败计数并设置10分钟过期
// // 	count, err := rdb.Incr(c, countKey).Result()
// // 	if err != nil {
// // 		return
// // 	}
// // 	rdb.Expire(c, countKey, SingleFailureTTL*time.Minute)
// //
// // 	// 如果达到或超过限制次数，设置指数增长的锁定时间
// // 	if count >= MaxAttemptsBeforeLock {
// // 		waitTime := CalculateLockTime(int(count))
// // 		rdb.Set(c, lockKey, LockValue, waitTime)
// // 	}
// // }

// 滑动窗口记录失败次数

package flowRestriction

import (
	"errors"
	"fmt"
	"math"
	"time"

	"forum/pkg/globals"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// CalculateLockTime 计算锁定时间（指数增长）
func CalculateLockTime(failCount int) time.Duration {
	if failCount < globals.MaxAttemptsBeforeLock {
		return 0 // 未达到限制次数时不锁定
	}
	// 从第3次失败开始，锁定时间指数增长：第3次60秒，第4次120秒，第5次240秒...
	baseTime := time.Duration(globals.BaseLockTime) * time.Second
	waitTime := baseTime * time.Duration(math.Pow(2, float64(failCount-globals.MaxAttemptsBeforeLock)))
	maxWait := time.Duration(globals.MaxLockTime) * time.Second
	if waitTime > maxWait {
		return maxWait
	}
	return waitTime
}

// CheckRateLimit
// @Description: 检查登录限制并返回剩余次数
// @Author lizhuang 2025-02-21 20:53:06
// @param        rdb *redis.Client
// @param        c *gin.Context
// @return       bool 是否允许登录
// @return       int 剩余的登录尝试次数
// @return       int 锁定时间
// @return       string 错误信息
func CheckRateLimit(rdb *redis.Client, c *gin.Context) (bool, int, int, string) {
	ip := c.ClientIP()
	listKey := fmt.Sprintf(globals.LoginFailListKeyPrefix, ip)
	lockKey := fmt.Sprintf(globals.LoginLockKeyPrefix, ip)

	// 检查是否处于锁定状态
	if locked, err := rdb.Exists(c, lockKey).Result(); err == nil && locked > 0 {
		// 该ip地址处于锁定状态，检查一下过期时间
		ttl, _ := rdb.TTL(c, lockKey).Result()
		if ttl > 0 {
			lockTime := int(math.Ceil(ttl.Seconds())) // 改为秒，向上取整
			return false, 0, lockTime, fmt.Sprintf("登录尝试过于频繁，请在%d秒后再试", lockTime)
		}
	}

	// 获取滑动窗口内的失败次数
	now := time.Now().Unix() // 时间戳
	// 下面的代码意思是：计算最近 WindowSize（600秒）内的登录失败次数。范围是：[当前时间-窗口大小，当前时间]。使用 "+inf" 简化了代码，避免了动态计算最大时间戳的麻烦
	windowStart := now - int64(globals.WindowSize) // 窗口开始时间（秒）
	count, err := rdb.ZCount(c, listKey, fmt.Sprintf("%d", windowStart), "+inf").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return false, 0, 0, "系统内部错误"
	}

	// 计算剩余尝试次数
	remaining := globals.MaxAttemptsBeforeLock - int(count)
	if remaining < 0 {
		remaining = 0
	}

	return true, remaining, 0, ""
}

// RecordFailedAttempt 记录登录失败
func RecordFailedAttempt(rdb *redis.Client, c *gin.Context) {
	ip := c.ClientIP()
	listKey := fmt.Sprintf(globals.LoginFailListKeyPrefix, ip)
	lockKey := fmt.Sprintf(globals.LoginLockKeyPrefix, ip)

	// 记录当前失败的时间戳
	now := time.Now().Unix()
	err := rdb.ZAdd(c, listKey, &redis.Z{
		Score:  float64(now),
		Member: now, // 使用时间戳作为成员
	}).Err()
	if err != nil {
		globals.Log.Errorf(err.Error())
		return
	}

	// 设置滑动窗口的过期时间
	rdb.Expire(c, listKey, globals.WindowSize*time.Second)

	// 获取当前窗口内的失败次数
	windowStart := now - int64(globals.WindowSize)
	count, err := rdb.ZCount(c, listKey, fmt.Sprintf("%d", windowStart), "+inf").Result()
	if err != nil {
		globals.Log.Errorf(err.Error())
		return
	}

	// 如果达到或超过限制次数，设置指数增长的锁定时间
	if count >= globals.MaxAttemptsBeforeLock {
		waitTime := CalculateLockTime(int(count))
		rdb.Set(c, lockKey, globals.LockValue, waitTime)
	}

	// 清理窗口外的旧记录
	rdb.ZRemRangeByScore(c, listKey, "-inf", fmt.Sprintf("%d", windowStart))
}
