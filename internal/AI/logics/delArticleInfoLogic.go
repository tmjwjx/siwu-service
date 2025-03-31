package logics

import (
	"context"
	"fmt"
	"forum/internal/AI/proto/article"
	"forum/internal/AI/requests"
	"forum/pkg/AI/utils"
	"google.golang.org/grpc/metadata"
	"time"
)

type DelArticleInfoLogic struct {
	Token string
}

func (s *DelArticleInfoLogic) DelArticleInfoLogic(req *requests.DelArticleInfoReq) (string, error) {
	// 1. 建立 gRPC 连接
	conn, err := utils.GenerateGrpcConn()
	if err != nil {
		return "", err
	}

	defer conn.Close() // 确保连接关闭

	// 3. 创建正确的客户端
	client := article.NewArticleServiceClient(conn)

	// 4. 构造请求参数
	reqGrpc := &article.DelArticleInfoRequest{
		ArticleID: uint32(req.ArticleID),
	}

	// 5. 构建token
	md := metadata.Pairs(
		"authorization", "Bearer "+s.Token, // 这里与服务端拦截器的 md["authorization"] 键名对应
	)
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)
	ctx, cancel := context.WithTimeout(ctxWithToken, 5*time.Second)
	defer cancel()

	// 6. 调用 gRPC 方法
	resGrpc, err := client.DelArticleInfo(ctx, reqGrpc)
	if err != nil {
		return "", fmt.Errorf("DelArticleInfoLogic -> 服务调用失败 -> %v", err)
	}

	// 6. 返回结果
	return resGrpc.Inform, nil
}
