package utils

import (
	"fmt"
	"forum/pkg/globals"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func GenerateGrpcConn() (*grpc.ClientConn, error) {
	// 1. 获取服务地址
	srv := "/services/siwuai"
	addr, err := globals.EtcdClient.GetService(srv)
	if err != nil {
		err = fmt.Errorf("globals.EtcdClient.GetService() %v", err)
		return nil, err
	}

	// 2. 建立 gRPC 连接
	conn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		err = fmt.Errorf("grpc.Dial() %v", err)
		return nil, err
	}
	return conn, nil
}
