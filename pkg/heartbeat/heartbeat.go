package heartbeat

import (
	"context"
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"time"
)

func StartHeartbeat() {

	// 创建心跳定时器（每15秒发送一次）
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// 创建一个上下文（context），用于控制协程的生命周期
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动心跳协程
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				// 发送心跳信号
				fmt.Println("Sending heartbeat signal...")
				for userId, _ := range globals.SubscriberChannels {
					internalUtils.MessagePush2("", userId, globals.HeatBeatType, "heartbeat")
				}
			}
		}
	}()
}
