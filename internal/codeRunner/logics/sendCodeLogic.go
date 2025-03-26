package logics

import (
	"context"
	"fmt"
	"forum/internal/codeRunner/requests"
	"forum/pkg/globals"
	request "github.com/ningzhaoxing/codeRunnerProto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type SendCodeLogin struct {
	requests.CodeRunnerReq
}

func (s *SendCodeLogin) SendCode(token string) (string, error) {
	// 1. 获取服务地址
	srv := "/services/code-runner"
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
	client := request.NewCodeRunnerClient(conn)

	// 4. 构造请求参数
	req := &request.ExecuteRequest{
		Id:          s.Id,
		Uid:         s.Uid,
		CallBackUrl: "http://192.168.23.49:8081/codeRunner/getResult",
		CodeBlock:   s.CodeArea,
		Language:    s.Language,
	}
	md := metadata.Pairs(
		"token", token, // 这里与服务端拦截器的 md["token"] 键名对应
	)

	// 3. 将 metadata 附加到上下文
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)
	// 5. 调用 RPC 方法
	ctx, cancel := context.WithTimeout(ctxWithToken, 5*time.Second)
	defer cancel()

	_, err = client.Execute(ctx, req)
	if err != nil {
		log.Printf("发送代码: %v", err)
		return "", fmt.Errorf("服务调用失败")
	}

	// 6. 返回结果
	return "发送成功", nil
}
