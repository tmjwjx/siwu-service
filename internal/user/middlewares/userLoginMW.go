package middlewares

// // RegisterMW 注册中间件
// func RegisterMW(c *gin.Context) {
// 	// 绑定数据
// 	var registerReq requests.RegisterReq
// 	if err := c.ShouldBind(&registerReq); err != nil {
// 		globals.Log.Error(response.ErrBindDataIsWrong)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrBindDataIsWrong), nil))
// 		c.Abort() // 终止后续处理
// 		return
// 	}
//
// 	// 检验邮箱是否合法
// 	if !internalUtils.IsValidEmail(registerReq.Email) {
// 		globals.Log.Error(response.ErrEmailIsInvalid)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrEmailIsInvalid), nil))
// 		c.Abort()
// 		return
// 	}
//
// 	// 核对两次输入的密码
// 	if registerReq.Password != registerReq.RePassword {
// 		globals.Log.Error(response.ErrPasswordTwiceIsWrong)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrPasswordTwiceIsWrong), nil))
// 		c.Abort()
// 		return
// 	}
//
// 	// 检验密码是否合法
// 	if !internalUtils.IsValidPassword(registerReq.Password) {
// 		globals.Log.Error(response.ErrPasswordIsInvalid)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrPasswordIsInvalid), nil))
// 		c.Abort()
// 		return
// 	}
//
// 	// 将验证通过的数据存储在上下文中
// 	c.Set("req", registerReq)
//
// 	// 如果验证通过，继续处理
// 	c.Next()
// }
//
// // ReqVerifyCodeMW 用户请求验证码中间件
// func ReqVerifyCodeMW(c *gin.Context) {
// 	// 绑定数据
// 	email := c.Query("email")
//
// 	// 检验邮箱是否合法
// 	if !internalUtils.IsValidEmail(email) {
// 		globals.Log.Error(response.ErrEmailIsInvalid)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrEmailIsInvalid), nil))
// 		c.Abort()
// 		return
// 	}
//
// 	c.Set("req", email)
//
// 	c.Next()
// }
//
// // LoginMW 登录中间件
// func LoginMW(c *gin.Context) {
// 	// 绑定数据
// 	var logicReq requests.LogicReq
// 	if err := c.ShouldBind(&logicReq); err != nil {
// 		globals.Log.Errorf(response.ErrBindDataIsWrong)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrBindDataIsWrong), nil))
// 		c.Abort() // 终止后续处理
// 		return
// 	}
//
// 	// 检验邮箱是否合法
// 	if !internalUtils.IsValidEmail(logicReq.Email) {
// 		globals.Log.Errorf(response.ErrEmailIsInvalid)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrEmailIsInvalid), nil))
// 		c.Abort() // 终止后续处理
// 		return
// 	}
//
// 	// 将验证通过的数据存储在上下文中
// 	c.Set("req", logicReq)
//
// 	// 如果验证通过，继续处理
// 	c.Next()
// }
//
// // ForgotPasswordMW 忘记密码中间件
// func ForgotPasswordMW(c *gin.Context) {
// 	// 绑定数据
// 	var forgotPasswordReq requests.ForgotPasswordReq
// 	if err := c.ShouldBind(&forgotPasswordReq); err != nil {
// 		globals.Log.Error(response.ErrBindDataIsWrong)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrBindDataIsWrong), nil))
// 		c.Abort()
// 		return
// 	}
//
// 	// 检验邮箱是否合法
// 	if !internalUtils.IsValidEmail(forgotPasswordReq.Email) {
// 		globals.Log.Error(response.ErrEmailIsInvalid)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrEmailIsInvalid), nil))
// 		c.Abort()
// 		return
// 	}
//
// 	// 核对两次输入的密码
// 	if forgotPasswordReq.Password != forgotPasswordReq.RePassword {
// 		globals.Log.Error(response.ErrPasswordTwiceIsWrong)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrPasswordTwiceIsWrong), nil))
// 		c.Abort()
// 		return
// 	}
//
// 	// 检验密码是否合法
// 	if !internalUtils.IsValidPassword(forgotPasswordReq.Password) {
// 		globals.Log.Error(response.ErrPasswordIsInvalid)
// 		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf(response.ErrPasswordIsInvalid), nil))
// 		c.Abort()
// 		return
// 	}
//
// 	c.Set("req", forgotPasswordReq)
//
// 	c.Next()
// }
