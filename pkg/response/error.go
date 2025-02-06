package response

const (
	// 2000
	DataSuccess = "成功"

	// 4000
	ErrPasswordIsInvalid    = "密码不合法"
	ErrPasswordIsWrong      = "密码错误"
	ErrPasswordTwiceIsWrong = "两次输入的密码不一致"
	ErrEmailIsInvalid       = "邮箱不合法"
	ErrEmailIsUse           = "邮箱已被使用"
	// ErrEmailIsNotExist         = "不存在该邮箱"
	ErrVerifyCodeIsWrong       = "验证码错误"
	ErrVerifyCodeIsExpired     = "验证码已失效"
	ErrReqVerifyCodeIsCooling  = "请求验证码正在冷却时间中"
	ErrUserIdNotGetFromContext = "无法从上下文中获取用户id"
	ErrPageRangeIsWrong        = "page 范围错误"
	ErrLimitRangeIsWrong       = "limit 范围错误"

	// 5000
	ErrUserIdNotExist       = "用户id不存在"
	ErrEmailNotExist        = "用户邮箱不存在"
	ErrBindDataIsWrong      = "绑定数据错误"
	ErrTypeAssertionFail    = "类型断言失败"
	ErrGetReqIsWrong        = "get req 错误"
	ErrUnableFindUserAvatar = "无法找到用户头像"
	ErrModelNotPointer      = "数据模型不是指针类型"
	ErrInsertDataFail       = "插入数据失败"
)
