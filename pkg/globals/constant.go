package globals

// 全局常量

// 发送邮件异步处理
const (
	EmailStreamKey   = "email_tasks"          // 消费者组的流的名称
	EmailGroupKey    = "email_consumer_group" // 消费者组的名称
	MaxRetries       = 3                      // 发送消息的最大重试次数
	DeadLetterStream = "email_dead_letter"    // 死信队列
)

// 登录限流相关常量
const (
	MaxAttemptsBeforeLock = 3     // 开始限制前的最大尝试次数
	WindowSize            = 600   // 滑动窗口大小（秒，即10分钟）
	BaseLockTime          = 60    // 基础锁定时间（秒，即1分钟）
	MaxLockTime           = 86400 // 最大锁定时间（秒，即24小时）

	LoginFailListKeyPrefix = "login_fail_list:%s" // 失败时间戳列表键
	LoginLockKeyPrefix     = "login_lock:%s"      // 锁定状态键
	LockValue              = "locked"             // 锁定键的值
)
