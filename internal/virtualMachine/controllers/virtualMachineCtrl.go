package controllers

import (
	"context"
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	pvm "github.com/world-fish/proto/vm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

func CreateVM(c *gin.Context) {
	// 连接grpc服务 默认禁用安全连接 没有加密和认证
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Printf("连接失败: %v", err)
	}
	defer conn.Close()

	// 获取token
	token := c.GetHeader("Authorization")

	// 创建一个包含元数据的context
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("Authorization", token))

	// 建立连接
	client := pvm.NewVMManagerClient(conn)

	// 调用服务
	resp, err := client.CreateVM(ctx, &pvm.CreateVMReq{})
	if err != nil {
		log.Printf("调用失败: %v", err)
	} else {
		log.Printf("调用成功: %v", resp)
	}
}

// DestroyVM 销毁虚拟机
func DestroyVM(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		e := response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("DestroyVM -> 从token中获取用户email失败"), nil)
		response.Failed(c, 500, e)
		return
	}

	// 连接grpc服务 默认禁用安全连接 没有加密和认证
	conn, err := grpc.Dial("localhost:8888", grpc.WithInsecure())
	if err != nil {
		log.Printf("连接失败: %v", err)
	}
	defer conn.Close()

	// 获取token
	token := c.GetHeader("Authorization")

	// 创建一个包含元数据的context
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("Authorization", token))

	// 建立连接
	client := pvm.NewVMManagerClient(conn)

	// 调用服务（一个邮箱只允许分配一个vm虚拟机）
	_, err = client.DestroyVM(ctx, &pvm.DestroyVMReq{Email: email.(string)})
	if err != nil {
		log.Printf("调用失败: %v", err)
	}
}
