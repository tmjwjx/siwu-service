package logics

import (
	"encoding/json"
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	requests "github.com/ningzhaoxing/codeRunnerProto"
	"strconv"
)

type ResultSendLogic struct {
	err error
	req *requests.ExecuteResponse
}

func NewResultSendLogic(req *requests.ExecuteResponse) *ResultSendLogic {
	return &ResultSendLogic{req: req}
}

func (r *ResultSendLogic) Send() error {
	data, err := json.Marshal(*r.req)

	if err != nil {
		globals.Log.Errorf("消息序列化失败 err = %s", err)
		return err
	}

	// 将消息push到sse
	internalUtils.MessagePush2(string(data), strconv.Itoa(int(r.req.Uid)), globals.CodeRunnerType)
	fmt.Println(string(data))
	return nil
}
