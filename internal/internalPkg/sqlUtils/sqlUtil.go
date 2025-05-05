package sqlUtils

import (
	"fmt"
	"forum/pkg/globals"
	"forum/pkg/response"
	"gorm.io/gorm"
	"reflect"
)

// 封装了数据库增、删、改操作方法（只适用于操作一个表）

// InsertObject
// @Description: 插入一条数据。
// @Author lizhuang 2024-10-04 21:01:19
// @param        db *gorm.DB GORM的数据库实例。
// @param        model interface{} 指针类型，里面要包含要添加的结构体字段。
// @return       error
func InsertObject(db *gorm.DB, model interface{}) error {
	// 检查 model 是否为指针类型
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		// return fmt.Errorf("InsertObject() err: 数据模型必须是指针类型")
		globals.Log.Errorf(response.ErrModelNotPointer)
		return fmt.Errorf(response.ErrModelNotPointer)
	}

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 插入数据
		result := tx.Create(model)
		if result.Error != nil {
			// 如果插入时出错，事务会自动回滚
			// return fmt.Errorf("InsertObject() err: %v", result.Error)
			globals.Log.Errorf(result.Error.Error())
			return result.Error
		}
		if result.RowsAffected == 0 {
			// 如果没有插入任何数据，也进行回滚
			// return fmt.Errorf("InsertObject() err: 插入数据失败")
			globals.Log.Errorf(response.ErrInsertDataFail)
			return fmt.Errorf(response.ErrInsertDataFail)
		}

		// 如果成功，事务会自动提交
		return nil
	})
}

// InsertObjects
// @Description: 插入多条数据。
// @Author lizhuang 2024-10-04 21:24:33
// @param        db *gorm.DB GORM的数据库实例。
// @param        models interface{} 切片类型。切片的每一个元素要包含要添加的结构体字段。
// @return       error
// func InsertObjects(db *gorm.DB, models interface{}) error {
// 	// 确保传入的 models 是一个切片
// 	modelsValue := reflect.ValueOf(models)
// 	if modelsValue.Kind() != reflect.Slice {
// 		return fmt.Errorf("InsertObjects() err: 传入的参数必须是切片类型")
// 	}
//
// 	// 开启事务
// 	return db.Transaction(func(tx *gorm.DB) error {
// 		// 执行批量插入
// 		result := tx.Create(models)
// 		if result.Error != nil {
// 			return fmt.Errorf("InsertObjects() err: %v", result.Error)
// 		}
// 		if result.RowsAffected == 0 {
// 			return fmt.Errorf("InsertObjects() err: 插入数据失败")
// 		}
// 		return nil
// 	})
// }

// DeleteObjectsByModel
// @Description: 按照model模型，根据一个或多个条件删除一个或多个对象。
// @Author lizhuang 2024-10-04 21:00:45
// @param        db *gorm.DB GORM 的数据库实例。
// @param        modelType interface{} 要删除的模型类型的指针。
// @param        condition map[string]interface{} 删除条件的键值对。
// @return       int64 删除的记录数。
// @return       error 错误。
func DeleteObjectsByModel(db *gorm.DB, modelType interface{}, condition map[string]interface{}) (int64, error) {
	// 确保 modelType 是指针类型
	if reflect.TypeOf(modelType).Kind() != reflect.Ptr {
		// return 0, fmt.Errorf("DeleteObjectsByModel() err: 数据模型必须是指针类型")
		globals.Log.Errorf(response.ErrModelNotPointer)
		return 0, fmt.Errorf(response.ErrModelNotPointer)
	}
	var rowsAffected int64 // 用于保存删除的记录数

	// 开启事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 构建查询条件
		query := tx.Model(modelType)
		for key, value := range condition {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}

		// 执行删除操作
		result := query.Delete(modelType)
		if result.Error != nil {
			// return fmt.Errorf("DeleteObjectsByModel() err: %v", result.Error)
			globals.Log.Errorf(result.Error.Error())
			return result.Error
		}

		// 保存删除的记录数
		rowsAffected = result.RowsAffected
		return nil
	})
	return rowsAffected, err
}

// DeleteObjectsByTable
// @Description: 按照表名，根据一个或多个条件删除一个或多个对象。
// @Author lizhuang 2024-10-04 21:25:06
// @param        db *gorm.DB GORM的数据库实例。
// @param        tableName string 要删除的表名。
// @param        condition map[string]interface{} 删除条件的键值对。
// @return       int64 删除的记录数。
// @return       error 错误。
func DeleteObjectsByTable(db *gorm.DB, tableName string, condition map[string]interface{}) (int64, error) {
	var rowsAffected int64 // 用于保存删除的记录数

	// 开启事务
	err := db.Transaction(func(tx *gorm.DB) error {
		// 构建查询条件
		query := tx.Table(tableName)
		for key, value := range condition {
			query = query.Where(fmt.Sprintf("%s = ?", key), value)
		}

		// 执行删除操作。query.Delete(nil)：删除符合条件的记录，不需要指定具体的模型类型。
		result := query.Delete(nil)
		if result.Error != nil {
			// return fmt.Errorf("DeleteObjectsByTable() err: %v", result.Error)
			globals.Log.Errorf(result.Error.Error())
			return result.Error
		}

		// 保存删除的记录数
		rowsAffected = result.RowsAffected
		return nil
	})
	return rowsAffected, err
}

// UpdateObjects
// @Description: 根据一个或多个参数更新一个或多个对象。
// @Author lizhuang 2024-10-04 21:26:25
// @param        db *gorm.DB GORM的数据库实例。
// @param        model interface{} 要更新的模型类型的指针，模型中要包含查询的参数。
// @param        updates map[string]interface{} 更新值的键值对。
// @return       error
func UpdateObjects(db *gorm.DB, model interface{}, updates map[string]interface{}) error {
	// 检查 model 是否为指针类型
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		// return fmt.Errorf("UpdateObjects() err: 数据模型必须是指针类型")
		globals.Log.Errorf(response.ErrModelNotPointer)
		return fmt.Errorf(response.ErrModelNotPointer)
	}

	// 开启事务
	return db.Transaction(func(tx *gorm.DB) error {
		// 执行更新操作
		result := tx.Model(model).Where(model).Updates(updates)
		if result.Error != nil {
			globals.Log.Errorf(result.Error.Error())
			return result.Error
			// return fmt.Errorf("UpdateObjects() err: %v", result.Error)
		}
		return nil
	})
}

// Paginate
// @Description: 分页查询的复用逻辑。用于db.Scopes()。（如果参数有负数，最终分页查询得到的结果为空。）
// @Author lizhuang 2024-10-04 21:28:13
// @param        page int
// @param        limit int
// @return       func(db *gorm.DB) *gorm.DB
func Paginate(page, limit int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
