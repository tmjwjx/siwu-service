package server

import (
	"context"
	"errors"
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/sendEmail"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Run 启动路由
func Run() {
	// 运行结束时，刷新日志的缓冲区（缓存区的信息写入到文件中）
	defer globals.Log.Sync()

	// 开协程
	go sendEmail.StartEmailTaskConsumer(globals.RDB)

	// 启动处理函数
	SetupRouter()

	// 启动Http服务 + 平滑关闭
	Start()
}

// Start 启动Http服务+平滑关闭（软关闭）
func Start() {
	address := fmt.Sprintf("%s:%d", globals.AppConfig.App.Host, globals.AppConfig.App.Port)

	// 启动路由
	// err := globals.Router.Run(address)
	// if err != nil {
	// 	globals.Log.Errorf("路由启动错误")
	// 	return
	// }

	// 创建一个 HTTP 服务（处理器）
	server := &http.Server{
		Addr:    address,
		Handler: globals.Router,
	}

	// 创建一个 channel 用于接收系统信号
	signalChan := make(chan os.Signal, 1)

	// 监听 OS 中断信号（SIGINT，SIGTERM）
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			globals.Log.Fatalf("ListenAndServe() failed: %v", err)
		}
	}()

	// 等待接收终止信号
	sig := <-signalChan
	fmt.Printf("接收到信号: %s. 正在关闭中...\n", sig)

	// 创建一个具有超时的 context，用于优雅地关闭服务器
	// 如果正在进行连接的接口10s之内还没有结束就强制结束了
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 调用 Shutdown() 平滑关闭服务器
	if err := server.Shutdown(ctx); err != nil {
		globals.Log.Fatalf("服务器关闭失败: %v", err)
	}

	fmt.Println("服务器已正常停止")
}
