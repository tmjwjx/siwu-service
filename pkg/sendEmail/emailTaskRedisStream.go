package sendEmail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"forum/pkg/globals"
	"github.com/go-redis/redis/v8"
	"time"
)

// EmailTask 存储邮件发送任务的消息
type EmailTask struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// const (
// 	EmailStreamKey = "email_tasks"          // 消费者组的流的名称
// 	EmailGroupKey  = "email_consumer_group" // 消费者组的名称
// )

// PushEmailTaskToStream 生产者，推送邮件任务到 Redis Stream
func PushEmailTaskToStream(rdb *redis.Client, task EmailTask) error {
	ctx := context.Background()

	// 序列化任务为 JSON
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("序列化邮件任务失败: %v", err)
	}

	// 添加消息到 Stream
	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: globals.EmailStreamKey,
		Values: map[string]interface{}{"task": taskJSON},
	}).Result()

	if err != nil {
		return fmt.Errorf("redis XAdd 失败: %v", err)
	}

	return nil
}

// StartEmailTaskConsumer 消费者，启动邮件任务消费者
func StartEmailTaskConsumer(rdb *redis.Client) {
	ctx := context.Background()

	for {
		// 阻塞读取消息
		messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    globals.EmailGroupKey,
			Consumer: "server",
			Streams:  []string{globals.EmailStreamKey, ">"},
			Block:    0,
			Count:    1,
		}).Result()

		if err != nil {
			if errors.Is(err, context.Canceled) {
				globals.Log.Info("邮件消费者退出")
				return
			}
			globals.Log.Errorf("读取 Redis Stream 失败: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		// 处理消息
		for _, msg := range messages[0].Messages {
			// 更新最后处理的消息 ID

			// 解析任务
			var task EmailTask
			taskJSON, ok := msg.Values["task"].(string)
			if !ok {
				globals.Log.Error("消息格式错误: task 字段缺失或类型错误")
				continue
			}

			if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
				globals.Log.Errorf("反序列化邮件任务失败: %v", err)
				continue
			}

			// 发送邮件
			if err := SendEmail(globals.SendEmailCfg, task.Email, task.Subject, task.Body); err != nil {
				globals.Log.Errorf("发送邮件失败: %v", err)
				// 可以添加重试逻辑
				continue
			}

			// 确认消息（可选：使用 XACK 确认消费）
			rdb.XAck(ctx, globals.EmailStreamKey, globals.EmailGroupKey, msg.ID)
		}
	}
}
