// // package sendEmailAsynchronous
// //
// // import (
// // 	"context"
// // 	"encoding/json"
// // 	"errors"
// // 	"fmt"
// // 	"forum/pkg/globals"
// // 	"github.com/go-redis/redis/v8"
// // 	"time"
// // )
// //
// // // EmailTask 存储邮件发送任务的消息
// // type EmailTask struct {
// // 	Email    string `json:"email"`
// // 	Subject  string `json:"subject"`
// // 	Body     string `json:"body"`
// // 	Retries  int    `json:"retries"`   // 记录重试次数
// // 	OriginID string `json:"origin_id"` // 原始消息 ID（用于追踪）
// // }
// //
// // // PushEmailTaskToStream 生产者，推送邮件任务到 Redis Stream
// // func PushEmailTaskToStream(rdb *redis.Client, task EmailTask) error {
// // 	ctx := context.Background()
// //
// // 	// 序列化任务为 JSON
// // 	taskJSON, err := json.Marshal(task)
// // 	if err != nil {
// // 		return fmt.Errorf("序列化邮件任务失败: %v", err)
// // 	}
// //
// // 	// 添加消息到 Stream
// // 	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
// // 		Stream: globals.EmailStreamKey,
// // 		Values: map[string]interface{}{"task": taskJSON},
// // 	}).Result()
// //
// // 	if err != nil {
// // 		return fmt.Errorf("redis XAdd 失败: %v", err)
// // 	}
// //
// // 	return nil
// // }
// //
// // // StartEmailTaskConsumer 消费者，启动邮件任务消费者
// // func StartEmailTaskConsumer(rdb *redis.Client) {
// // 	ctx := context.Background()
// //
// // 	for {
// // 		// 阻塞读取消息
// // 		messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
// // 			Group:    globals.EmailGroupKey,
// // 			Consumer: "server",
// // 			Streams:  []string{globals.EmailStreamKey, ">"},
// // 			Block:    0,
// // 			Count:    1,
// // 		}).Result()
// //
// // 		if err != nil {
// // 			if errors.Is(err, context.Canceled) {
// // 				globals.Log.Info("邮件消费者退出")
// // 				return
// // 			}
// // 			globals.Log.Errorf("读取 Redis Stream 失败: %v", err)
// // 			time.Sleep(1 * time.Second)
// // 			continue
// // 		}
// //
// // 		// 处理消息
// // 		for _, msg := range messages[0].Messages {
// // 			// 解析任务
// // 			var task EmailTask
// // 			taskJSON, ok := msg.Values["task"].(string)
// // 			if !ok {
// // 				globals.Log.Error("消息格式错误: task 字段缺失或类型错误")
// // 				continue
// // 			}
// //
// // 			if err = json.Unmarshal([]byte(taskJSON), &task); err != nil {
// // 				globals.Log.Errorf("反序列化邮件任务失败: %v", err)
// // 				continue
// // 			}
// //
// // 			// 发送邮件（含重试逻辑）
// // 			err = sendEmailWithRetry(rdb, task)
// // 			if err != nil {
// // 				globals.Log.Errorf("邮件任务最终失败: %v", err)
// // 				continue
// // 			}
// //
// // 			// 确认消息（可选：使用 XACK 确认消费）
// // 			rdb.XAck(ctx, globals.EmailStreamKey, globals.EmailGroupKey, msg.ID)
// // 		}
// // 	}
// // }
// //
// // // sendEmailWithRetry 带重试的邮件发送
// // func sendEmailWithRetry(rdb *redis.Client, task EmailTask) error {
// // 	for {
// // 		// 发送邮件
// // 		err := SendEmail(globals.SendEmailCfg, task.Email, task.Subject, task.Body)
// // 		if err == nil {
// // 			return nil // 成功
// // 		}
// //
// // 		// 记录重试次数
// // 		task.Retries++
// // 		fmt.Printf("邮件发送失败，进行第 %d 次重试: %v\n", task.Retries, err)
// // 		globals.Log.Infof("邮件发送失败，进行第 %d 次重试: %v", task.Retries, err)
// //
// // 		if task.Retries >= globals.MaxRetries {
// // 			// 转移到死信队列
// // 			if err := moveToDeadLetterQueue(rdb, task); err != nil {
// // 				return fmt.Errorf("重试次数耗尽且转移死信队列失败: %v", err)
// // 			}
// // 			return fmt.Errorf("重试次数耗尽，任务已进入死信队列")
// // 		}
// //
// // 		// 重新推送任务到 Stream（延迟重试）
// // 		if err := repushTaskWithDelay(rdb, task); err != nil {
// // 			return fmt.Errorf("重新推送任务失败: %v", err)
// // 		}
// //
// // 		return nil // 退出当前处理，等待下一次重试
// // 	}
// // }
// //
// // // repushTaskWithDelay 重新推送任务并添加延迟
// // func repushTaskWithDelay(rdb *redis.Client, task EmailTask) error {
// // 	fmt.Println("重新推送任务并添加延迟.........")
// //
// // 	ctx := context.Background()
// //
// // 	// 序列化更新后的任务
// // 	taskJSON, err := json.Marshal(task)
// // 	if err != nil {
// // 		return fmt.Errorf("序列化任务失败: %v", err)
// // 	}
// //
// // 	// 添加延迟（例如：10秒、30秒、60秒）
// // 	// delay := time.Duration(task.Retries) * 10 * time.Second
// // 	// time.Sleep(delay)
// // 	// time.Sleep(1*time.Second)
// //
// // 	// 推送到原 Stream
// // 	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
// // 		Stream: globals.EmailStreamKey,
// // 		Values: map[string]interface{}{"task": taskJSON},
// // 	}).Result()
// //
// // 	return err
// // }
// //
// // // moveToDeadLetterQueue 转移到死信队列
// // func moveToDeadLetterQueue(rdb *redis.Client, task EmailTask) error {
// // 	fmt.Println("转移到死信队列.........")
// //
// // 	ctx := context.Background()
// // 	// 序列化任务
// // 	taskJSON, err := json.Marshal(task)
// // 	if err != nil {
// // 		return err
// // 	}
// //
// // 	// 推送到死信队列
// // 	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
// // 		Stream: globals.DeadLetterStream,
// // 		Values: map[string]interface{}{
// // 			"task":      taskJSON,
// // 			"error":     "max retries exceeded",
// // 			"timestamp": time.Now().Unix(),
// // 		},
// // 	}).Result()
// //
// // 	return err
// // }
//
// package sendEmailAsynchronous
//
// import (
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"forum/pkg/globals"
// 	"github.com/go-redis/redis/v8"
// 	"time"
// )
//
// // EmailTask 存储邮件发送任务的消息
// type EmailTask struct {
// 	Email    string `json:"email"`
// 	Subject  string `json:"subject"`
// 	Body     string `json:"body"`
// 	Retries  int    `json:"retries"`   // 记录重试次数
// 	OriginID string `json:"origin_id"` // 原始消息 ID（用于追踪）
// }
//
// // PushEmailTaskToStream 生产者，推送邮件任务到 Redis Stream
// func PushEmailTaskToStream(rdb *redis.Client, task EmailTask) error {
// 	ctx := context.Background()
//
// 	// 序列化任务为 JSON
// 	taskJSON, err := json.Marshal(task)
// 	if err != nil {
// 		return fmt.Errorf("序列化邮件任务失败: %v", err)
// 	}
//
// 	// 添加消息到 Stream
// 	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
// 		Stream: globals.EmailStreamKey,
// 		Values: map[string]interface{}{"task": taskJSON},
// 	}).Result()
//
// 	if err != nil {
// 		return fmt.Errorf("redis XAdd 失败: %v", err)
// 	}
//
// 	return nil
// }
//
// // StartEmailTaskConsumer 消费者，启动邮件任务消费者
// func StartEmailTaskConsumer(rdb *redis.Client) {
// 	ctx := context.Background()
//
// 	// 确保消费者组存在
// 	err := rdb.XGroupCreateMkStream(ctx, globals.EmailStreamKey, globals.EmailGroupKey, "$").Err()
// 	if err != nil && !errors.Is(err, redis.Nil) {
// 		globals.Log.Warnf("创建消费者组失败（可能已存在）: %v", err)
// 	}
//
// 	// 这里无法读取到数据
//
// 	for {
// 		// 阻塞读取消息
// 		messages, err := rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
// 			Group:    globals.EmailGroupKey,
// 			Consumer: "server",
// 			Streams:  []string{globals.EmailStreamKey, ">"},
// 			Block:    0,
// 			Count:    1,
// 		}).Result()
//
// 		if err != nil {
// 			if errors.Is(err, context.Canceled) {
// 				globals.Log.Info("邮件消费者退出")
// 				return
// 			}
// 			globals.Log.Errorf("读取 Redis Stream 失败: %v", err)
// 			fmt.Printf("读取 Redis Stream 失败: %v", err)
// 			time.Sleep(1 * time.Second)
// 			continue
// 		}
//
// 		fmt.Println("==================")
//
// 		// 处理消息
// 		for _, msg := range messages[0].Messages {
// 			// 解析任务
// 			var task EmailTask
// 			taskJSON, ok := msg.Values["task"].(string)
// 			if !ok {
// 				globals.Log.Error("消息格式错误: task 字段缺失或类型错误")
// 				rdb.XAck(ctx, globals.EmailStreamKey, globals.EmailGroupKey, msg.ID)
// 				continue
// 			}
//
// 			if err = json.Unmarshal([]byte(taskJSON), &task); err != nil {
// 				globals.Log.Errorf("反序列化邮件任务失败: %v", err)
// 				rdb.XAck(ctx, globals.EmailStreamKey, globals.EmailGroupKey, msg.ID)
// 				continue
// 			}
//
// 			fmt.Println(task)
//
// 			// 发送邮件（含重试逻辑）
// 			err = sendEmailWithRetry(rdb, task)
// 			if err != nil {
// 				globals.Log.Errorf("邮件任务最终失败: %v", err)
// 			}
//
// 			// 确认消息
// 			rdb.XAck(ctx, globals.EmailStreamKey, globals.EmailGroupKey, msg.ID)
// 		}
// 	}
// }
//
// // sendEmailWithRetry 带重试的邮件发送
// func sendEmailWithRetry(rdb *redis.Client, task EmailTask) error {
// 	// 发送邮件
// 	err := SendEmail(globals.SendEmailCfg, task.Email, task.Subject, task.Body)
// 	if err == nil {
// 		globals.Log.Info("邮件发送成功")
// 		fmt.Println("邮件发送成功")
// 		return nil // 成功
// 	} else {
// 		fmt.Println("邮件发送失败")
// 	}
//
// 	// 记录重试次数
// 	task.Retries++
// 	globals.Log.Infof("邮件发送失败，进行第 %d 次重试: %v", task.Retries, err)
// 	fmt.Printf("邮件发送失败，进行第 %d 次重试: %v", task.Retries, err)
//
// 	if task.Retries >= globals.MaxRetries {
// 		// 转移到死信队列
// 		if err := moveToDeadLetterQueue(rdb, task); err != nil {
// 			return fmt.Errorf("重试次数耗尽且转移死信队列失败: %v", err)
// 		}
// 		return fmt.Errorf("重试次数耗尽，任务已进入死信队列")
// 	}
//
// 	// 异步重新推送任务到 Stream，确保这个任务不影响其他的任务
// 	delay := time.Duration(task.Retries) * 10 * time.Second
// 	repushTaskWithDelay(rdb, task, delay)
// 	return nil // 当前任务处理结束，等待消费者再次消费
// }
//
// // repushTaskWithDelay 异步重新推送任务并添加延迟
// func repushTaskWithDelay(rdb *redis.Client, task EmailTask, delay time.Duration) {
// 	go func() {
// 		fmt.Println("异步重新推送任务并添加延迟")
//
// 		// 添加延迟
// 		time.Sleep(delay)
//
// 		ctx := context.Background()
// 		// 序列化更新后的任务
// 		taskJSON, err := json.Marshal(task)
// 		if err != nil {
// 			globals.Log.Errorf("序列化任务失败: %v", err)
// 			return
// 		}
//
// 		// 推送到原 Stream
// 		_, err = rdb.XAdd(ctx, &redis.XAddArgs{
// 			Stream: globals.EmailStreamKey,
// 			Values: map[string]interface{}{"task": taskJSON},
// 		}).Result()
//
// 		if err != nil {
// 			globals.Log.Errorf("重新推送任务失败: %v", err)
// 		}
// 	}()
// }
//
// // moveToDeadLetterQueue 转移到死信队列
// func moveToDeadLetterQueue(rdb *redis.Client, task EmailTask) error {
// 	fmt.Println("转移到死信队列")
// 	ctx := context.Background()
// 	// 序列化任务
// 	taskJSON, err := json.Marshal(task)
// 	if err != nil {
// 		return err
// 	}
//
// 	// 推送到死信队列
// 	_, err = rdb.XAdd(ctx, &redis.XAddArgs{
// 		Stream: globals.DeadLetterStream,
// 		Values: map[string]interface{}{
// 			"task":      taskJSON,
// 			"error":     "max retries exceeded",
// 			"timestamp": time.Now().Unix(),
// 		},
// 	}).Result()
//
// 	return err
// }

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

func sendEmailWithRetry(rdb *redis.Client, task EmailTask) error {
	err := SendEmail(globals.SendEmailCfg, task.Email, task.Subject, task.Body)
	if err == nil {
		globals.Log.Info("邮件发送成功")
		fmt.Println("邮件发送成功")
		return nil
	}

	fmt.Println("邮件发送失败")
	task.Retries++
	globals.Log.Infof("邮件发送失败，进行第 %d 次重试: %v", task.Retries, err)
	fmt.Printf("邮件发送失败，进行第 %d 次重试: %v\n", task.Retries, err)

	if task.Retries >= globals.MaxRetries {
		if err := moveToDeadLetterQueue(rdb, task); err != nil {
			return fmt.Errorf("重试次数耗尽且转移死信队列失败: %v", err)
		}
		fmt.Println("任务进入死信队列")
		return fmt.Errorf("重试次数耗尽，任务已进入死信队列")
	}

	delay := time.Duration(task.Retries) * 10 * time.Second
	fmt.Println("准备异步重试，延迟:", delay)
	repushTaskWithDelay(rdb, task, delay)
	return nil
}

func repushTaskWithDelay(rdb *redis.Client, task EmailTask, delay time.Duration) {
	go func() {
		fmt.Println("异步重新推送任务，延迟:", delay)
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
			fmt.Println("重新推送任务成功，Retries:", task.Retries)
		}
	}()
}

func moveToDeadLetterQueue(rdb *redis.Client, task EmailTask) error {
	fmt.Println("转移到死信队列")
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
