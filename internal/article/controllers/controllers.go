package controllers

/*func SearchHandler(c *gin.Context) {

	db := globals.DB
	var req requests.SearchRequest // 创建一个 SearchRequest 类型的变量，用于存储请求参数

	// 绑定查询参数到 req 变量，如果绑定失败，返回错误信息
	if err := c.ShouldBindQuery(&req); err != nil {
		globals.Log.Errorf("err = %s", err)
		c.JSON(400, response.StatusBadRequestErr) // 返回 400 错误`
		return                                    // 结束函数执行
	}
	fmt.Println(req)
	articles, err := logics.Search(db, req)
	if err != nil {
		globals.Log.Errorf("搜索文章失败 err = %s", err)
		c.JSON(500, response.StatusInternalServerErr)
		return
	}
	data := response.AppData{
		Code: globals.StatusOK,
		Msg:  globals.CodeMsgMap[globals.StatusOK],
		Data: articles,
	}
	//c.JSON(200, data)
	response.Success(c, &data)
}
func Test1(c *gin.Context) {
	c.JSON(200, gin.H{"msg": "ok"})
}*/
