package sendEmailAsynchronous

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"forum/pkg/globals"
	"github.com/go-redis/redis/v8"
	"strings"
	"time"
)

// EmailTask 存储邮件发送任务的消息
type EmailTask struct {
	Email    string `json:"email"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Retries  int    `json:"retries"`
	OriginID string `json:"origin_id"`
}

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
	err := rdb.XGroupCreateMkStream(ctx, globals.EmailStreamKey, globals.EmailGroupKey, "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		globals.Log.Warnf("Failed to create consumer group (may already exist): %v", err)
	}

	consumerName := fmt.Sprintf("server-%d", time.Now().UnixNano())
	globals.Log.Infof("Starting consumer with name: %s", consumerName)

	for {
		// Handle pending messages first
		pending, err := rdb.XPending(ctx, globals.EmailStreamKey, globals.EmailGroupKey).Result()
		if err == nil && pending.Count > 0 {
			pendingMessages, err := rdb.XRange(ctx, globals.EmailStreamKey, pending.Lower, pending.Higher).Result()
			if err == nil {
				for _, msg := range pendingMessages {
					if processMessage(rdb, ctx, consumerName, msg) {
						err := rdb.XAck(ctx, globals.EmailStreamKey, globals.EmailGroupKey, msg.ID).Err()
						if err != nil {
							globals.Log.Errorf("Failed to acknowledge pending message %s: %v", msg.ID, err)
						} else {
							globals.Log.Infof("Acknowledged pending message ID: %s", msg.ID)
						}
					}
				}
			}
		}

		// Read new messages
		messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    globals.EmailGroupKey,
			Consumer: consumerName,
			Streams:  []string{globals.EmailStreamKey, ">"},
			Block:    0,
			Count:    10,
		}).Result()

		if err != nil {
			if errors.Is(err, context.Canceled) {
				globals.Log.Info("Email consumer exiting")
				return
			}
			globals.Log.Errorf("Failed to read from Redis stream: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		globals.Log.Infof("Read time: %v, number of messages: %d", time.Now(), len(messages))

		for _, stream := range messages {
			for _, msg := range stream.Messages {
				if processMessage(rdb, ctx, consumerName, msg) {
					err := rdb.XAck(ctx, globals.EmailStreamKey, globals.EmailGroupKey, msg.ID).Err()
					if err != nil {
						globals.Log.Errorf("Failed to acknowledge message %s: %v", msg.ID, err)
					} else {
						globals.Log.Infof("Acknowledged message ID: %s", msg.ID)
					}
				}
			}
		}
	}
}

// processMessage 处理从 Redis Stream 中读取到的消息
func processMessage(rdb *redis.Client, ctx context.Context, consumerName string, msg redis.XMessage) bool {
	var task EmailTask
	taskJSON, ok := msg.Values["task"].(string)
	if !ok {
		globals.Log.Error("Message format error: task field missing or wrong type")
		return true
	}

	if err := json.Unmarshal([]byte(taskJSON), &task); err != nil {
		globals.Log.Errorf("Failed to unmarshal email task: %v", err)
		return true
	}

	globals.Log.Infof("Processing message ID: %s, task: %v", msg.ID, task)
	err := sendEmailWithRetry(rdb, task)
	if err != nil {
		globals.Log.Errorf("Email task finally failed: %v", err)
	}
	return true
}

// sendEmailWithRetry 带重试的邮件发送
func sendEmailWithRetry(rdb *redis.Client, task EmailTask) error {
	err := SendEmail(globals.SendEmailCfg, task.Email, task.Subject, task.Body)
	if err == nil {
		globals.Log.Info("邮件发送成功")
		return nil
	}

	task.Retries++
	globals.Log.Infof("邮件发送失败，进行第 %d 次重试: %v", task.Retries, err)

	if task.Retries >= globals.MaxRetries {
		if err := moveToDeadLetterQueue(rdb, task); err != nil {
			return fmt.Errorf("重试次数耗尽且转移死信队列失败: %v", err)
		}
		return fmt.Errorf("重试次数耗尽，任务已进入死信队列")
	}

	delay := time.Duration(task.Retries) * 10 * time.Second
	repushTaskWithDelay(rdb, task, delay)
	return nil
}

// repushTaskWithDelay 延迟后重新推送任务到 Redis Stream
func repushTaskWithDelay(rdb *redis.Client, task EmailTask, delay time.Duration) {
	go func() {
		time.Sleep(delay)

		ctx := context.Background()
		taskJSON, err := json.Marshal(task)
		if err != nil {
			globals.Log.Errorf("序列化任务失败: %v", err)
			return
		}

		_, err = rdb.XAdd(ctx, &redis.XAddArgs{
			Stream: globals.EmailStreamKey,
			Values: map[string]interface{}{"task": taskJSON},
		}).Result()
		if err != nil {
			globals.Log.Errorf("重新推送任务失败: %v", err)
		} else {
			globals.Log.Info("重新推送任务成功，Retries:", task.Retries)
		}
	}()
}

// moveToDeadLetterQueue 转移到死信队列
func moveToDeadLetterQueue(rdb *redis.Client, task EmailTask) error {
	ctx := context.Background()
	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}
	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: globals.DeadLetterStream,
		Values: map[string]interface{}{
			"task":      taskJSON,
			"error":     "max retries exceeded",
			"timestamp": time.Now().Unix(),
		},
	}).Result()
	return err
}
