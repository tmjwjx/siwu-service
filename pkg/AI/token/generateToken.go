package token

import (
	"context"
	"fmt"
	"forum/pkg/AI/token/proto"
	"forum/pkg/AI/utils"
	"time"
)

type CodeToken struct {
	GenerateTokenKey string
}

func (c *CodeToken) GetToken() (string, error) {
	// 1. 建立 gRPC 连接
	conn, err := utils.GenerateGrpcConn()
	if err != nil {
		err = fmt.Errorf("utils.GenerateGrpcConn() %v", err)
		return "", err
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
		err = fmt.Errorf("client.GenerateToken() %v", err)
		return "", err
	}

	// 6. 返回结果
	return resp.Token, nil
}
