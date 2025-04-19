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

type GetArticleInfoFirstLogic struct {
	Token string
}

func (a *GetArticleInfoFirstLogic) GetArticleInfoFirstLogic(articleInfo *requests.GetArticleInfoFirstReq) (*requests.GetArticleInfoFirstRes, error) {
	// 1. 建立 gRPC 连接
	conn, err := utils.GenerateGrpcConn()
	if err != nil {
		return nil, err
	}

	defer conn.Close() // 确保连接关闭

	// 3. 创建正确的客户端
	client := article.NewArticleServiceClient(conn)

	// 4. 构造请求参数
	req := &article.GetArticleInfoFirstRequest{
		Content: articleInfo.Content,
		Tags:    articleInfo.Tags,
	}

	// 5. 构建token
	md := metadata.Pairs(
		"authorization", "Bearer "+a.Token, // 这里与服务端拦截器的 md["authorization"] 键名对应
	)
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)
	ctx, cancel := context.WithTimeout(ctxWithToken, 100*time.Second)
	defer cancel()

	// 6. 调用 gRPC 方法
	resGrpc, err := client.GetArticleInfoFirst(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("GetArticleInfoFirstLogic -> 服务调用失败 -> %v", err)
	}

	res := &requests.GetArticleInfoFirstRes{
		Key:      resGrpc.Key,
		Abstract: resGrpc.Abstract,
		Summary:  resGrpc.Summary,
		Tags:     resGrpc.Tags,
	}

	// 6. 返回结果
	return res, nil
}
