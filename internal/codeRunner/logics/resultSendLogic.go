package logics

import (
	"encoding/json"
	"forum/internal/codeRunner/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"strconv"
)

type ResultSendLogic struct {
	err error
	req *requests.CodeRunnerReq
}

func NewResultSendLogic(req *requests.CodeRunnerReq) *ResultSendLogic {
	return &ResultSendLogic{req: req}
}

func (r *ResultSendLogic) Send() error {
	data, err := json.Marshal(*r.req)
	if err != nil {
		globals.Log.Errorf("消息序列化失败 err = %s", err)
		return err
	}

	// 将消息push到sse
	internalUtils.MessagePush(string(data), strconv.Itoa(int(r.req.Uid)))

	return nil
}
