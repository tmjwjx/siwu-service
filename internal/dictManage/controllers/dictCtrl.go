package controllers

import (
	"fmt"
	"forum/internal/dictManage/logics"
	"forum/internal/dictManage/requests"
	"forum/pkg/globals"
	"forum/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AddType 新增字典类型
func AddType(c *gin.Context) {
	// 绑定数据
	var addTypeReq requests.AddTypeReq
	if err := c.ShouldBind(&addTypeReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("AddType() err: 绑定数据错误"), nil))
		return
	}

	// 检验数据
	if addTypeReq.Name == "" {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("AddType() err: name不能为空"), nil))
		return
	}
	if addTypeReq.Code == "" {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("AddType() err: code不能为空"), nil))
		return
	}

	// 业务逻辑
	dictReqContext := logics.NewDictReqContext(globals.DB, c)
	id, err := dictReqContext.AddType(addTypeReq)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("AddType() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"id": id}))
}

// DeleteType 批量删除字典类型
func DeleteType(c *gin.Context) {
	// 绑定数据
	var deleteIdReq requests.DeleteTypeReq
	if err := c.ShouldBind(&deleteIdReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("DeleteType() err: 绑定数据错误"), nil))
		return
	}

	// 业务逻辑
	dictReqContext := logics.NewDictReqContext(globals.DB, c)
	if err := dictReqContext.DeleteType(deleteIdReq); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("DeleteType() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// UpdateType 修改字典类型
func UpdateType(c *gin.Context) {
	// 绑定数据
	var updateTypeReq requests.UpdateTypeReq
	if err := c.ShouldBind(&updateTypeReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdateType() err: 绑定数据错误"), nil))
		return
	}

	// 检验数据
	if updateTypeReq.Name == "" {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdateType() err: name不能为空"), nil))
		return
	}
	if updateTypeReq.Code == "" {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdateType() err: code不能为空"), nil))
		return
	}

	// 业务逻辑
	dictReqContext := logics.NewDictReqContext(globals.DB, c)
	if err := dictReqContext.UpdateType(updateTypeReq); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UpdateType() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// GetType 获取字典类型
func GetType(c *gin.Context) {
	// 绑定数据
	var getTypeReq requests.GetTypeReq
	if err := c.ShouldBind(&getTypeReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetType() err: 绑定数据错误"), nil))
		return
	}

	// 检验数据
	if getTypeReq.Page <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetType() err: Page参数必须为正数"), nil))
		return
	}
	if getTypeReq.Limit <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetType() err: limit参数必须为正数"), nil))
		return
	}

	// 业务逻辑
	dictReqContext := logics.NewDictReqContext(globals.DB, c)
	dictTypeList, total, err := dictReqContext.GetType(getTypeReq)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetType() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"dict_type_list": dictTypeList, "total": total}))
}

// AddItem 新增字典项
func AddItem(c *gin.Context) {
	// 绑定数据
	var addItemReq requests.AddItemReq
	if err := c.ShouldBind(&addItemReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("AddItem() err: 绑定数据错误"), nil))
		return
	}

	// 检验数据
	if addItemReq.Label == "" {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("AddItem() err: label不能为空"), nil))
		return
	}

	// 业务逻辑
	dictReqContext := logics.NewDictReqContext(globals.DB, c)
	id, err := dictReqContext.AddItem(addItemReq)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("AddType() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"id": id}))
}

// DeleteItem 批量删除字典项
func DeleteItem(c *gin.Context) {
	// 绑定数据
	var deleteItemReq requests.DeleteItemReq
	if err := c.ShouldBind(&deleteItemReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("DeleteItem() err: 绑定数据错误"), nil))
		return
	}

	// 业务逻辑
	dictReqContext := logics.NewDictReqContext(globals.DB, c)
	if err := dictReqContext.DeleteItem(deleteItemReq); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("DeleteItem() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// UpdateItem 修改字典项
func UpdateItem(c *gin.Context) {
	// 绑定数据
	var updateItemReq requests.UpdateItemReq
	if err := c.ShouldBind(&updateItemReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdateItem() err: 绑定数据错误"), nil))
		return
	}

	// 检验数据
	if updateItemReq.Label == "" {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("UpdateItem() err: name不能为空"), nil))
		return
	}

	// 业务逻辑
	dictReqContext := logics.NewDictReqContext(globals.DB, c)
	if err := dictReqContext.UpdateItem(updateItemReq); err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("UpdateItem() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", nil))
}

// GetItem 获取字典项
func GetItem(c *gin.Context) {
	// 绑定数据
	var getItemReq requests.GetItemReq
	if err := c.ShouldBind(&getItemReq); err != nil {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetItem() err: 绑定数据错误"), nil))
		return
	}

	// 检验数据
	if getItemReq.Page <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetItem() err: Page参数必须为正数"), nil))
		return
	}
	if getItemReq.Limit <= 0 {
		response.Failed(c, http.StatusBadRequest, response.NewAppErr(globals.StatusBadRequest, fmt.Errorf("GetItem() err: limit参数必须为正数"), nil))
		return
	}

	// 业务逻辑
	dictReqContext := logics.NewDictReqContext(globals.DB, c)
	dictTypeList, total, err := dictReqContext.GetItem(getItemReq)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, response.NewAppErr(globals.StatusInternalServerError, fmt.Errorf("GetItem() -> %v", err), nil))
		return
	}
	response.Success(c, http.StatusOK, response.NewAppData(globals.StatusOK, "成功", gin.H{"dict_item_list": dictTypeList, "total": total}))
}
