package requests

type CodeRunnerReq struct {
	Id       string `json:"id"`       // 代码块id
	Language string `json:"language"` // 代码语言
	CodeArea string `json:"codeArea"` // 代码块
}
