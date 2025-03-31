package utils

import (
	"fmt"
	"forum/pkg/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func GenerateGrpcConn() (*grpc.ClientConn, error) {
	// 1. 获取服务地址
	srv := "/services/siwuai"
	addr, err := globals.EtcdClient.GetService(srv)
	if err != nil {
		log.Printf("获取服务地址失败: %v", err)
		return nil, fmt.Errorf("服务不可用")
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
		return nil, fmt.Errorf("连接服务失败")
	}
	return conn, nil
}
