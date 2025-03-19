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

type GetArticleInfoLogic struct {
	Token string
}

func (s *GetArticleInfoLogic) GetArticleInfoLogic(req *requests.GetArticleInfoReq) (*requests.GetArticleInfoRes, error) {
	// 1. 建立 gRPC 连接
	conn, err := utils.GenerateGrpcConn()
	if err != nil {
		return nil, err
	}

	defer conn.Close() // 确保连接关闭

	// 3. 创建正确的客户端
	client := article.NewArticleServiceClient(conn)

	// 4. 构造请求参数
	reqGrpc := &article.GetArticleInfoRequest{
		ArticleID: uint32(req.ArticleID),
	}

	// 5. 构建token
	md := metadata.Pairs(
		"token", s.Token, // 这里与服务端拦截器的 md["token"] 键名对应
	)
	ctxWithToken := metadata.NewOutgoingContext(context.Background(), md)
	ctx, cancel := context.WithTimeout(ctxWithToken, 5*time.Second)
	defer cancel()

	// 6. 调用 gRPC 方法
	resGrpc, err := client.GetArticleInfo(ctx, reqGrpc)
	if err != nil {
		return nil, fmt.Errorf("GetArticleInfoLogic -> 服务调用失败 -> %v", err)
	}

	res := &requests.GetArticleInfoRes{
		Abstract: resGrpc.Abstract,
		Summary:  resGrpc.Summary,
	}

	// 6. 返回结果
	return res, nil
}
