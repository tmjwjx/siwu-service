package controllers

import (
	"fmt"
	"forum/internal/AI/logics"
	requests "forum/internal/AI/requests/code"
	token2 "forum/pkg/AI/token"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func GetCodeExplain(c *gin.Context) {
	//接收前端的代码
	var codeRequest requests.CodeReq
	if err := c.ShouldBind(&codeRequest); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetCodeExplain() -> %v", err), nil))
		return
	}

	// 生成token
	codeToken := token2.CodeToken{
		GenerateTokenKey: "123456",
	}
	token, err := codeToken.GetToken()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("etcdToken.GetToken() -> %v", err), nil))
		return
	}

	// 调用逻辑层, 返回的stream为流式响应
	getCodeExplainLogic := logics.GetCodeExplainLogic{CodeReq: codeRequest}
	stream, err := getCodeExplainLogic.ExplainCodeLogic(token)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetCodeExplain() -> %v", err), nil))
		return
	}

	// 设置 SSE 响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")

	// 处理 gRPC 流并通过 SSE 发送
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			// 流结束
			break
		}
		if err != nil {
			fmt.Printf("sse从siwuai服务接收流失败: %v", err)
			response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("从siwuai接收流失败 -> %v", err), nil))
			return
		}

		// 发送 SSE 数据
		fmt.Println(resp.CodeExplain)

		fmt.Fprintf(c.Writer, "data: %s\n\n", resp.CodeExplain)
		c.Writer.Flush() // 立即推送给客户端

	}

	// 发送结束标志
	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")

	c.Writer.Flush()

}
