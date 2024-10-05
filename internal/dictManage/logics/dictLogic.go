package logics

import (
	"fmt"
	"forum/internal/dictManage/repositories"
	"forum/internal/dictManage/requests"
	"forum/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DictReqContext 用于在处理请求时传递数据库连接和请求上下文信息
type DictReqContext struct {
	DB  *gorm.DB
	Ctx *gin.Context
}

func NewDictReqContext(db *gorm.DB, c *gin.Context) *DictReqContext {
	return &DictReqContext{
		DB:  db,
		Ctx: c,
	}
}

// AddType 新增字典类型
func (d *DictReqContext) AddType(req requests.AddTypeReq) (uint, error) {
	// 查询是否已经有 Code 了
	dictType := repositories.QueryDictTypeByCode(d.DB, req.Code)
	if dictType != nil {
		return 0, fmt.Errorf("DictReqContext.AddType() err: 已经存在Code为 %v 的数据了", req.Code)
	}

	// 查询是否软删除了这个code，如果软删除了就恢复、更新这个code
	dictType2 := repositories.QueryDictTypeSoftDeletedByCode(d.DB, req.Code)
	if dictType2 != nil { // 这个code被软删除过
		// 恢复软删除的数据
		if err := repositories.UpdateDictTypeSoftDelete(d.DB, req.Code); err != nil {
			return 0, fmt.Errorf("DictReqContext.AddType() -> %v", err)
		}
		// 更新数据
		if err := repositories.UpdateObjects(d.DB, &models.DictType{Code: req.Code}, map[string]interface{}{"name": req.Name, "status": req.Status, "description": req.Description}); err != nil {
			return 0, fmt.Errorf("DictReqContext.AddType() -> %v", err)
		}
	} else { // 这个code没有被软删除过
		// 插入数据
		if err := repositories.InsertObject(d.DB, &models.DictType{Name: req.Name, Code: req.Code, Status: req.Status, Description: req.Description}); err != nil {
			return 0, fmt.Errorf("DictReqContext.AddType() -> %v", err)
		}
	}

	// 查询数据库中code对应的id
	dictType = repositories.QueryDictTypeByCode(d.DB, req.Code)
	if dictType == nil {
		return 0, fmt.Errorf("DictReqContext.AddType() err: 没有查询到code为 %v 的字典类型", req.Code)
	}
	return dictType.ID, nil
}

// DeleteType 批量删除字典类型
func (d *DictReqContext) DeleteType(req requests.DeleteTypeReq) error {
	// 删除数据
	for _, v := range req.IdList {
		if _, err := repositories.DeleteObjectsByModel(d.DB, &models.DictType{}, map[string]interface{}{"id": v}); err != nil {
			return fmt.Errorf("DictReqContext.DeleteType() -> %v", err)
		}
	}
	return nil
}

// UpdateType 修改字典类型
func (d *DictReqContext) UpdateType(req requests.UpdateTypeReq) error {
	// 判断 id 是否存在
	if dictType := repositories.QueryDictTypeById(d.DB, req.Id); dictType == nil {
		return fmt.Errorf("DictReqContext.UpdateType() err: id为 %v 的字典类型不存在", req.Id)
	}

	m := map[string]interface{}{
		"name":        req.Name,
		"code":        req.Code,
		"status":      req.Status,
		"description": req.Description,
	}
	// 更新
	if err := repositories.UpdateObjects(d.DB, &models.DictType{Model: gorm.Model{ID: req.Id}}, m); err != nil {
		return fmt.Errorf("DictReqContext.UpdateType() -> %v", err)
	}
	return nil
}

// GetType 获取字典类型
func (d *DictReqContext) GetType(req requests.GetTypeReq) ([]*requests.GetTypeRes, int, error) {
	conditions := map[string]interface{}{}
	// 判断是否该添加某些查询条件（如果某些条件为空，那么就不查询这个条件）
	if req.Name != "" {
		conditions["name"] = req.Name
	}
	if req.Code != "" {
		conditions["code"] = req.Code
	}
	// status == 0 代表着全部
	if req.Status != 0 {
		conditions["status"] = req.Status
	}

	// 查询
	dictTypes, total, err := repositories.QueryDictTypeByPage(d.DB, conditions, req.Page, req.Limit, req.CreateAtBegin, req.CreateAtEnd)
	if err != nil {
		return nil, 0, fmt.Errorf("DictReqContext.GetType() err: %v", err)
	}

	// 返回数据
	var dictTypeList []*requests.GetTypeRes = make([]*requests.GetTypeRes, 0)
	for _, v := range dictTypes {
		dictTypeList = append(dictTypeList, &requests.GetTypeRes{
			Id:          v.ID,
			CreatedAt:   v.CreatedAt.Format("2006-01-02 15:04:05"),
			Name:        v.Name,
			Code:        v.Code,
			Status:      v.Status,
			Description: v.Description,
		})
	}

	return dictTypeList, total, nil
}

// AddItem 新增字典项
func (d *DictReqContext) AddItem(req requests.AddItemReq) (uint, error) {
	// 判断该 DictTypeCode 是否存在
	dictType := repositories.QueryDictTypeByCode(d.DB, req.DictTypeCode)
	if dictType == nil {
		return 0, fmt.Errorf("DictReqContext.AddItem() err: 不存在code为 %v 的DictType", req.DictTypeCode)
	}

	// 插入数据
	if err := repositories.InsertObject(d.DB, &models.DictItem{
		DictTypeCode: req.DictTypeCode,
		Label:        req.Label,
		Value:        req.Value,
		Sort:         req.Sort,
		Status:       req.Status,
		Description:  req.Description,
		ExtendValue:  req.ExtendValue,
	}); err != nil {
		return 0, fmt.Errorf("DictReqContext.AddItem() -> %v", err)
	}

	// 查询id
	dictItem := repositories.QueryDictItemByCodeAndLabel(d.DB, req.DictTypeCode, req.Label)
	if dictItem == nil {
		return 0, fmt.Errorf("DictReqContext.AddItem() err: 不存在DictItem和Code %v 关联的数据", dictType)
	}
	return dictItem.ID, nil
}

// DeleteItem 批量删除字典项
func (d *DictReqContext) DeleteItem(req requests.DeleteItemReq) error {
	// 判断该 DictTypeCode 是否存在
	dictType := repositories.QueryDictTypeByCode(d.DB, req.DictTypeCode)
	if dictType == nil {
		return fmt.Errorf("DictReqContext.AddItem() err: 不存在code为 %v 的DictType", req.DictTypeCode)
	}

	// 删除数据
	for _, v := range req.IdList {
		if _, err := repositories.DeleteObjectsByModel(d.DB, &models.DictItem{}, map[string]interface{}{"id": v, "dict_type_code": req.DictTypeCode}); err != nil {
			return fmt.Errorf("DictReqContext.DeleteItem() -> %v", err)
		}
	}
	return nil
}

// UpdateItem 修改字典项
func (d *DictReqContext) UpdateItem(req requests.UpdateItemReq) error {
	// 判断 id 是否存在，
	if dictItem := repositories.QueryDictItemById(d.DB, req.Id); dictItem == nil {
		return fmt.Errorf("DictReqContext.UpdateItem() err: id为 %v 的字典项不存在", req.Id)
	}

	m := map[string]interface{}{
		"label":        req.Label,
		"value":        req.Value,
		"sort":         req.Sort,
		"status":       req.Status,
		"description":  req.Description,
		"extend_value": req.ExtendValue,
	}
	// 更新
	if err := repositories.UpdateObjects(d.DB, &models.DictItem{Model: gorm.Model{ID: req.Id}}, m); err != nil {
		return fmt.Errorf("DictReqContext.UpdateItem() -> %v", err)
	}
	return nil
}

// GetItem 获取字典项
func (d *DictReqContext) GetItem(req requests.GetItemReq) ([]*requests.GetItemRes, int, error) {
	conditions := map[string]interface{}{}
	// 判断是否该添加某些查询条件（如果某些条件为空，那么就不查询这个条件）
	// 判断该 DictTypeCode 是否存在
	dictType := repositories.QueryDictTypeByCode(d.DB, req.DictTypeCode)
	if dictType == nil {
		return nil, 0, fmt.Errorf("DictReqContext.GetItem() err: 不存在code为 %v 的DictType", req.DictTypeCode)
	}
	conditions["dict_type_code"] = req.DictTypeCode
	if req.Label != "" {
		conditions["label"] = req.Label
	}
	// status == 0 代表着全部
	if req.Status != 0 {
		conditions["status"] = req.Status
	}

	// 查询
	dictItems, total, err := repositories.QueryDictItemByPage(d.DB, conditions, req.Page, req.Limit, req.CreateAtBegin, req.CreateAtEnd)
	if err != nil {
		return nil, 0, fmt.Errorf("DictReqContext.GetItem() err: %v", err)
	}

	// 返回数据
	var dictItemList []*requests.GetItemRes = make([]*requests.GetItemRes, 0)
	for _, v := range dictItems {
		dictItemList = append(dictItemList, &requests.GetItemRes{
			Id:           v.ID,
			CreateAt:     v.CreatedAt.Format("2006-01-02 15:04:05"),
			DictTypeCode: v.DictTypeCode,
			Label:        v.Label,
			Value:        v.Value,
			Sort:         v.Sort,
			Status:       v.Status,
			Description:  v.Description,
			ExtendValue:  v.ExtendValue,
		})
	}
	return dictItemList, total, nil
}
