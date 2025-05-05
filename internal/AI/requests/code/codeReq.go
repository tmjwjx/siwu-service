package code

type CodeReq struct {
	CodeQuestion string `form:"code_question"` // 用户提问的代码
	UserId       uint32 `form:"user_id"`       // 用户的id
	CodeType     string `form:"code_type"`     // 代码语言
}
