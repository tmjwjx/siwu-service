package repositories

import (
	"fmt"
	"forum/internal/models"
	"gorm.io/gorm"
	"reflect"
	"time"
)

// InsertObject 插入一条数据。
//
// db: GORM 的数据库实例。
// model: 指针类型。
//
// Returns an error
func InsertObject(db *gorm.DB, model interface{}) error {
	// 检查 data 是否为指针类型
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("InsertObject() err: 数据模型必须是指针类型")
	}

	// 插入数据
	result := db.Create(model)
	if result.Error != nil {
		return fmt.Errorf("InsertObject() err: %v\t执行的查询语句为: %v", result.Error, result.Statement.SQL.String())
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("InsertObject() err: 插入数据失败")
	}

	return nil
}

// InsertObjects 插入多条数据。
//
// db: GORM 的数据库实例。
// models: 切片类型。
//
// Returns an error
func InsertObjects(db *gorm.DB, models interface{}) error {
	// 确保传入的 models 是一个切片
	modelsValue := reflect.ValueOf(models)
	if modelsValue.Kind() != reflect.Slice {
		return fmt.Errorf("InsertObjects() err: 传入的参数必须是切片类型")
	}

	// 执行批量插入
	result := db.Create(models)
	if result.Error != nil {
		return fmt.Errorf("InsertObjects() err: %v\t执行的查询语句为: %v", result.Error, result.Statement.SQL.String())
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("InsertObjects() err: 插入数据失败")
	}
	return nil
}

// DeleteObjectsByModel 按照model模型，根据一个或多个条件删除一个或多个对象。
//
// db: GORM 的数据库实例。
// modelType: 要删除的模型类型的指针。
// condition: 删除条件的键值对。
//
// Returns int64:返回删除的记录数。
func DeleteObjectsByModel(db *gorm.DB, modelType interface{}, condition map[string]interface{}) (int64, error) {
	// 确保 modelType 是指针类型
	if reflect.TypeOf(modelType).Kind() != reflect.Ptr {
		return 0, fmt.Errorf("DeleteObjectsByModel() err: 数据模型必须是指针类型")
	}

	// 构建查询条件
	query := db.Model(modelType)
	for key, value := range condition {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}

	// 执行删除操作
	result := query.Delete(modelType)
	if result.Error != nil {
		return 0, fmt.Errorf("DeleteObjectsByModel() err: %v\t执行的查询语句为: %v", result.Error, result.Statement.SQL.String())
	}

	// 返回删除的记录数
	return result.RowsAffected, nil
}

// UpdateObjects 根据一个或多个参数更新一个或多个对象。
//
// db: GORM 的数据库实例。
// model: 要更新的模型类型的指针，模型中要包含查询的参数。
// updates: 更新值的键值对。
//
// Returns error: 错误
func UpdateObjects(db *gorm.DB, model interface{}, updates map[string]interface{}) error {
	// 检查 model 是否为指针类型
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		return fmt.Errorf("UpdateObjects() err: 数据模型必须是指针类型")
	}

	// 执行更新操作
	result := db.Model(model).Where(model).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("UpdateObjects() err: %v", result.Error)
	}

	return nil
}

// QueryDictTypeByCode 通过Code查找字典类型
func QueryDictTypeByCode(db *gorm.DB, code string) *models.DictType {
	var dictType models.DictType
	d := db.Model(&models.DictType{}).Where("code = ?", code).Select("*").Scan(&dictType)
	// 没有找到字典
	if d.RowsAffected <= 0 {
		return nil
	}

	return &dictType
}

// QueryDictTypeById 通过id查找字典类型
func QueryDictTypeById(db *gorm.DB, id uint) *models.DictType {
	var dictType models.DictType
	d := db.Model(&models.DictType{}).Where("id = ?", id).Select("*").Scan(&dictType)
	// 没有找到字典
	if d.RowsAffected <= 0 {
		return nil
	}
	return &dictType
}

// QueryDictItemById 通过id查找字典项
func QueryDictItemById(db *gorm.DB, id uint) *models.DictItem {
	var dictItem models.DictItem
	d := db.Model(&models.DictItem{}).Where("id = ?", id).Select("*").Scan(&dictItem)
	// 没有找到字典
	if d.RowsAffected <= 0 {
		return nil
	}
	return &dictItem
}

// QueryDictTypeByPage 获取字典类型，支持分页
// conditions: 查询条件。
// page: 第几页。
// limit: 每页数据条数。
// 例子：如果 page = 2，limit = 10，那么会跳过前 10 条记录，返回第 11-20 条记录。
// 返回的int表示一共有多少条符合条件的数据
func QueryDictTypeByPage(db *gorm.DB, conditions map[string]interface{}, page int, limit int, createAtBegin, createAtEnd string) ([]*models.DictType, int, error) {
	var dictTypes []*models.DictType

	// 使用条件查询
	query := db.Model(&models.DictType{})
	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	// 解析 createAtBegin 和 createAtEnd 字符串为 time.Time 类型
	if createAtBegin != "" {
		parsedCreateAtBegin, err := time.Parse("2006-01-02", createAtBegin)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryDictTypeByPage() err: createAtBegin 解析错误: %v", err)
		}
		query = query.Where("created_at >= ?", parsedCreateAtBegin)
	}
	if createAtEnd != "" {
		parsedCreateAtEnd, err := time.Parse("2006-01-02", createAtEnd)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryDictTypeByPage() err: createAtEnd 解析错误: %v", err)
		}
		query = query.Where("created_at <= ?", parsedCreateAtEnd)
	}

	// 查看符合条件的数据一共有多少条
	query.Find(&dictTypes)
	total := len(dictTypes)

	// 添加分页逻辑
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit // 计算偏移量。Offset 是从第几条数据开始取。
		query = query.Offset(offset).Limit(limit)
	}

	// 执行查询
	if err := query.Find(&dictTypes).Error; err != nil {
		return nil, 0, err
	}

	return dictTypes, total, nil
}

// QueryDictItemByCodeAndLabel 通过关联的Code查找字典项
func QueryDictItemByCodeAndLabel(db *gorm.DB, code, label string) *models.DictItem {
	var dictItem models.DictItem
	d := db.Model(&models.DictItem{}).Where("dict_type_code = ? and label = ?", code, label).Select("*").Scan(&dictItem)
	// 没有找到字典
	if d.RowsAffected <= 0 {
		return nil
	}
	return &dictItem
}

// QueryDictItemByPage 获取字典项，支持分页
// conditions: 查询条件。
// page: 第几页。
// limit: 每页数据条数。
// 例子：如果 page = 2，limit = 10，那么会跳过前 10 条记录，返回第 11-20 条记录。
// 返回的int表示一共有多少条符合条件的数据
func QueryDictItemByPage(db *gorm.DB, conditions map[string]interface{}, page int, limit int, createAtBegin, createAtEnd string) ([]*models.DictItem, int, error) {
	var dictItems []*models.DictItem

	// 使用条件查询
	query := db.Model(&models.DictItem{})
	for key, value := range conditions {
		query = query.Where(fmt.Sprintf("%s = ?", key), value)
	}
	// 解析 createAtBegin 和 createAtEnd 字符串为 time.Time 类型
	if createAtBegin != "" {
		parsedCreateAtBegin, err := time.Parse("2006-01-02", createAtBegin)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryDictTypeByPage() err: createAtBegin 解析错误: %v", err)
		}
		query = query.Where("created_at >= ?", parsedCreateAtBegin)
	}
	if createAtEnd != "" {
		parsedCreateAtEnd, err := time.Parse("2006-01-02", createAtEnd)
		if err != nil {
			return nil, 0, fmt.Errorf("QueryDictTypeByPage() err: createAtEnd 解析错误: %v", err)
		}
		query = query.Where("created_at <= ?", parsedCreateAtEnd)
	}

	// 查看符合条件的数据一共有多少条
	query.Find(&dictItems)
	total := len(dictItems)

	// 添加分页逻辑
	if page > 0 && limit > 0 {
		offset := (page - 1) * limit // 计算偏移量。Offset 是从第几条数据开始取。
		query = query.Offset(offset).Limit(limit)
	}

	// 执行查询
	if err := query.Find(&dictItems).Error; err != nil {
		return nil, 0, err
	}

	return dictItems, total, nil
}

// QueryDictTypeSoftDeletedByCode 查询软删除的元素
func QueryDictTypeSoftDeletedByCode(db *gorm.DB, code string) *models.DictType {
	var dictType models.DictType
	// 使用 Unscoped 查询软删除的记录
	err := db.Unscoped().Where("code = ? AND deleted_at IS NOT NULL", code).First(&dictType).Error
	if err != nil {
		return nil // 发生错误或未找到记录，返回 nil
	}
	return &dictType // 找到记录，返回指向该记录的指针
}

// UpdateDictTypeSoftDelete 恢复字典类型已删除的记录
func UpdateDictTypeSoftDelete(db *gorm.DB, code string) error {
	var dictType models.DictType
	err := db.Unscoped().Where("code = ?", code).First(&dictType).Error
	if err != nil {
		return fmt.Errorf("UpdateDictTypeSoftDelete() err: %v", err)
	}
	// 恢复已删除的记录，将 DeletedAt 字段设置为零值。
	dictType.DeletedAt = gorm.DeletedAt{}

	// Save：保存（或更新）传入的模型实例。
	return db.Save(&dictType).Error
}
