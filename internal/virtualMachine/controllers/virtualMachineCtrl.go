package controllers

import (
	"context"
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
	
	// 创建一个包含元数据的context
	ctx := metadata.NewOutgoingContext(context.Background(), metadata.Pairs("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MzIsImVtYWlsIjoiMjkzNzY5MzkxOUBxcS5jb20iLCJpc3MiOiJzaXd1LXdlYi1zZXJ2aWNlIiwiZXhwIjoxNzQyMjU3OTU2LCJpYXQiOjE3NDE5NTU1NTZ9.tfX6T_Wa-dhYmK4ZvP5PGr8H8aUH4fMsTQD04VSJYVs"))
	
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
