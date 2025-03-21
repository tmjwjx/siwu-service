package logics

import (
	"context"
	"fmt"
	requests "forum/internal/AI/requests/code"
	"forum/pkg/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type GetCodeExplainLogic struct {
	requests.CodeReq
}

func (s *GetCodeExplainLogic) ExplainCodeLogic(token string) (grpc.ServerStreamingClient[requests.CodeResponse], error) {
	// 1. 获取服务地址
	srv := "/services/siwuai"
	addr, err := globals.EtcdClient.GetService(srv)
	if err != nil {
		log.Printf("获取服务地址失败: %v", err)
		return nil, fmt.Errorf("etcd服务发现失败")
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
		return nil, fmt.Errorf("gRPC 连接失败")
	}
	//defer conn.Close() // 确保连接关闭 //todo 连接关闭处理

	// 3. 创建正确的客户端
	client := requests.NewCodeServiceClient(conn)

	// 4. 构造请求参数
	req := &requests.CodeRequest{
		CodeQuestion: s.CodeQuestion,
		UserId:       s.UserId,
		CodeType:     s.CodeType,
	}

	// 5. 构建token
	md := metadata.Pairs(
		"token", token, // 这里与服务端拦截器的 md["token"] 键名对应
	)
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)
	ctx, _ := context.WithTimeout(ctxWithToken, 100*time.Second)

	// 6. 调用 gRPC 方法
	stream, err := client.ExplainCode(ctx, req)
	if err != nil {
		fmt.Println("logics.ExplainCode() client.ExplainCode(): ", err)
		return nil, fmt.Errorf("服务调用失败")
	}

	// 6. 返回结果
	return stream, nil
}
