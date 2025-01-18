package controllers

// BsLogout 后台登出
// func BsLogout(c *gin.Context) {
// 	// 获取token
// 	tokenString := c.GetHeader("Authorization")
//
// 	// 业务逻辑
// 	bsManageContext := logics.NewBsManageContext(globals.DB, c)
// 	if err := bsManageContext.BsLogout(tokenString[len("Bearer "):]); err != nil {
// 		response.Failed(c, http.StatusUnauthorized, response.NewAppErr(globals.StatusUnauthorized, fmt.Errorf("BsLogout() err = %v", err), nil))
// 		return
// 	}
//
// 	// 成功
// 	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
// }
