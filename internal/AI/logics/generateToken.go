package logics

import (
	"context"
	"fmt"
	"forum/pkg/AI/token/proto"
	"forum/pkg/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type CodeToken struct {
	GenerateTokenKey string
}

func (c *CodeToken) GetEtcdToken() (string, error) {
	// 1. 获取服务地址
	srv := "/services/siwuai"
	addr, err := globals.EtcdClient.GetService(srv)
	if err != nil {
		log.Printf("获取服务地址失败: %v", err)
		return "", fmt.Errorf("服务不可用")
	}

	// 2. 建立 gRPC 连接
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		log.Printf("gRPC 连接失败: %v", err)
		return "", fmt.Errorf("连接服务失败")
	}
	defer conn.Close() // 确保连接关闭
	// 3. 创建正确的客户端
	client := proto.NewTokenServiceClient(conn)
	// 4. 构造请求参数
	req := &proto.TokenRequest{
		GenerateTokenKey: c.GenerateTokenKey,
	}

	// 5. 调用 RPC 方法
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GenerateToken(ctx, req)
	if err != nil {
		log.Printf("生成 Token 失败: %v", err)
		return "", fmt.Errorf("服务调用失败")
	}

	// 6. 返回结果
	return resp.Token, nil
}
