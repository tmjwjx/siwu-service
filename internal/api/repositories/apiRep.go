package repositories

import (
	"errors"
	"fmt"
	"forum/internal/api/requests"
	"forum/internal/models"
	"forum/pkg/casbin"
	casbin2 "github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// GetAllApiRep
// @Description: 获取所有api列表
// @Author wangyulong 2024-10-15 08:50:14
// @param        db *gorm.DB
// @return       *requests.GetAllApiRes
// @return       error
func GetAllApiRep(db *gorm.DB) (*requests.GetAllApiRes, error) {

	var res requests.GetAllApiRes
	var groups []requests.MiddleResultG

	// 查询所有的分组id和label
	err := db.Table("sw_dict_types").
		Select("sw_dict_items.id, sw_dict_items.label").
		Joins("join sw_dict_items on sw_dict_items.dict_type_code = sw_dict_types.code").
		Where("sw_dict_types.name = ?", "API分组").
		Scan(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("GetAllApiRep -> 查询所有的分组id和label失败 -> %s", err)
	}

	for _, group := range groups {

		list := &requests.List{
			GroupID:   group.ID,
			GroupName: group.Label,
			Children:  make([]*requests.Children, 0),
		}
		var apiIDs []uint
		// 查询每个分组下面拥有的apiID
		err = db.Model(&models.ApiDictItemGroup{}).Where("group_id = ?", group.ID).Pluck("api_id", &apiIDs).Error
		if err != nil {
			return nil, fmt.Errorf("GetAllApiRep -> 查询每个分组下面拥有的apiID失败 -> %s", err)
		}

		for _, apiId := range apiIDs {

			var child requests.Children
			// 查询每个api的id和name
			err = db.Model(&models.Api{}).Select("id As api_id, path As api_name").Where("id = ?", apiId).Scan(&child).Error
			if err != nil {
				return nil, fmt.Errorf("GetAllApiRep -> 查询每个api的id和name失败 -> %s", err)
			}
			list.Children = append(list.Children, &child)
		}
		res.ApiList = append(res.ApiList, list)
	}
	return &res, nil
}

//// ApiInitRep 获取所有api列表
//func ApiInitRep(db *gorm.DB) (*[]*requests.ApiInitRes, error) {
//
//	var res []*requests.ApiInitRes
//	var apies []models.Api
//
//	// 查询api表
//	err := db.Find(&apies).Error
//	if err != nil {
//		return nil, fmt.Errorf("ApiInitRep -> 获取所有api列表失败 -> %s", err)
//	}
//
//	for _, api := range apies {
//
//		/*var requestMethodId uint
//		var requestMethodName string
//		var groupId uint
//		var groupName string
//
//		// 从api_request_method表中查询 request_method_id
//		err = db.Model(&models.ApiDictItemRequestMethod{}).Select("RequestMethodId").Where("api_id = ?", api.ID).First(&requestMethodId).Error
//		if err != nil {
//			return nil, fmt.Errorf("ApiInitRep -> 获取api请求方法id失败 -> %s", err)
//		}
//
//		// 从requestMethod表中查询 请求方法的name
//		err = db.Model(&models.RequestMethod{}).Select("Name").Where("id = ?", requestMethodId).First(&requestMethodName).Error
//		if err != nil {
//			return nil, fmt.Errorf("ApiInitRep -> 获取api请求方法name失败 -> %s", err)
//		}
//
//		// 从api_group表中查询 group_id
//		err = db.Model(models.ApiDictItemGroup{}).Select("group_id").Where("api_id = ?", api.ID).First(&groupId).Error
//		if err != nil {
//			return nil, fmt.Errorf("ApiInitRep -> 获取api分组id失败 -> %s", err)
//		}
//
//		// 从group表中查询 api分组的name
//		err = db.Model(models.Group{}).Select("name").Where("id = ?", groupId).First(&groupName).Error
//		if err != nil {
//			return nil, fmt.Errorf("ApiInitRep -> 获取api分组name失败 -> %s", err)
//		}*/
//
//		// 根据api的id获取api的请求方式的name和分组的name
//		groupId, groupName, requestMethodId, requestMethodName, err := GetApiGroupAndMethod(db, api.ID)
//		if err != nil {
//			return nil, fmt.Errorf("ApiInitRep -> %s", err)
//		}
//
//		apiInitRes := &requests.ApiInitRes{
//			ID:                api.ID,
//			Path:              api.Path,
//			GroupingId:        groupId,
//			Grouping:          groupName,
//			BriefIntroduction: api.BriefIntroduction,
//			RequestMethodId:   requestMethodId,
//			RequestMethod:     requestMethodName,
//		}
//
//		res = append(res, apiInitRes)
//	}
//	return &res, nil
//}

// GetApiDetailsRep 获取当前api详情
func GetApiDetailsRep(db *gorm.DB, id uint) (*requests.ApiDetailsRes, error) {

	var res requests.ApiDetailsRes
	query := db.Table("sw_apis").
		Select("sw_apis.id, sw_apis.path, sw_apis.brief_introduction, dg.id As grouping_id, dg.label As `grouping`, dm.id As request_method_id, dm.label As request_method").
		Joins("left join sw_api_dict_item_groups on sw_api_dict_item_groups.api_id = sw_apis.id").
		Joins("left join sw_api_dict_item_request_methods on sw_api_dict_item_request_methods.api_id = sw_apis.id").
		Joins("left join sw_dict_items As dg on sw_api_dict_item_groups.group_id = dg.id").
		Joins("left join sw_dict_items As dm on sw_api_dict_item_request_methods.request_method_id = dm.id")

	query = query.Where("sw_apis.deleted_at IS NULL AND sw_apis.id = ?", id).Order("sw_apis.created_at DESC")

	err := query.Scan(&res).Error
	if err != nil {
		return nil, fmt.Errorf("SearchApiListRep -> 查询api异常 -> %s", err)
	}

	return &res, nil

	/*var api models.Api
	// 从api表中查询 api
	err := db.Model(models.Api{}).Where("id = ?", id).Scan(&api).Error
	if err != nil {
		return nil, fmt.Errorf("GetApiDetailsRep -> 从api表中查询 api失败 -> %s", err)
	}
	// 根据api的id获取api的请求方式的name和分组的name
	result, err := GetApiGroupAndMethod(db, api.ID)
	if err != nil {
		return nil, fmt.Errorf("GetApiDetailsRep -> %s", err)
	}

	apiDetailsRes := &requests.ApiDetailsRes{
		ID:                api.ID,
		Path:              api.Path,
		GroupingId:        result.GroupId,
		Grouping:          result.GroupName,
		BriefIntroduction: api.BriefIntroduction,
		RequestMethodId:   result.RequestMethodId,
		RequestMethod:     result.RequestMethodName,
	}

	return apiDetailsRes, nil*/

}

// GetGroupListRep 获取所有api分组列表
func GetGroupListRep(db *gorm.DB) (*requests.ApiGroupRes, error) {

	var apiGroupRes requests.ApiGroupRes
	var groups []string

	// 查询分组信息
	query := db.Table("sw_dict_types").Joins("left join sw_dict_items on sw_dict_items.dict_type_code = sw_dict_types.code")

	err := query.Where("sw_dict_types.name = ?", "API分组").Pluck("sw_dict_items.label", &groups).Error
	if err != nil {
		return nil, fmt.Errorf("GetRequestMethodRep -> 查询分组信息失败 -> %s", err)
	}

	for _, group := range groups {

		apiGroup := &requests.ApiGroup{
			Value: group,
			Label: group,
		}

		apiGroupRes.Groups = append(apiGroupRes.Groups, apiGroup)
	}

	return &apiGroupRes, nil

}

// GetRequestMethodRep 获取所有请求方法
func GetRequestMethodRep(db *gorm.DB) (*requests.ApiReqMethodRes, error) {

	var apiReqMethodRes requests.ApiReqMethodRes
	var reqMethods []string

	// 查询请求方法
	query := db.Table("sw_dict_types").Joins("left join sw_dict_items on sw_dict_items.dict_type_code = sw_dict_types.code")

	err := query.Where("sw_dict_types.name = ?  AND sw_dict_items.deleted_at IS NULL", "API请求方法").Pluck("sw_dict_items.label", &reqMethods).Error
	if err != nil {
		return nil, fmt.Errorf("GetRequestMethodRep -> 查询请求方法失败 -> %s", err)
	}

	for _, reqMethod := range reqMethods {

		apiReqMethod := &requests.ApiReqMethod{
			Value: reqMethod,
			Label: reqMethod,
		}

		apiReqMethodRes.Methods = append(apiReqMethodRes.Methods, apiReqMethod)
	}

	return &apiReqMethodRes, nil
}

// DeleteApiRep 删除api
func DeleteApiRep(e *casbin2.Enforcer, req *requests.DeleteApiReq, db *gorm.DB) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("DeleteApiRep -> 开启事务失败 -> %s", tx.Error)
	}

	for _, id := range req.ID {

		// 删除api表中的信息
		result := tx.Where("id = ?", id).Delete(&models.Api{})
		if result.Error != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("DeleteApiRep ->  删除api表中的信息失败 -> %s", result.Error)
		}

		// 只有要删除的api存在，才有必要删除其他两张表中的信息，
		// 如果不存在，就变相等于删除成功了
		if result.RowsAffected != 0 {
			// 删除api_group表中的信息
			result2 := tx.Model(&models.ApiDictItemGroup{}).Where("api_id = ?", id).Delete(nil)
			if result2.Error != nil {
				tx.Rollback() // 回滚事务
				return fmt.Errorf("DeleteApiRep ->  删除api_group表中的信息失败 -> %s", result2.Error)
			}

			// 删除api_request_method表中的信息
			result3 := tx.Model(&models.ApiDictItemRequestMethod{}).Where("api_id = ?", id).Delete(nil)
			if result3.Error != nil {
				tx.Rollback() // 回滚事务
				return fmt.Errorf("DeleteApiRep ->  删除api_request_method表中的信息失败 -> %s", result3.Error)
			}
		}
	}

	// 开始casbin相关事务
	Ctx := e.GetAdapter().(*gormadapter.Adapter).GetDb().Begin()
	if Ctx.Error != nil {
		// 回滚普通事务
		tx.Rollback()
		// 回滚casbin相关事务
		Ctx.Rollback()
		return fmt.Errorf("DeleteApiRep ->  开始casbin相关事务失败 -> %s", Ctx.Error)
	}

	casbinServer := &casbin.CasbinService{
		Enforcer: e,
	}

	// 确保最新的策略数据
	err := casbinServer.Enforcer.LoadPolicy()
	if err != nil {
		// 回滚普通事务
		tx.Rollback()
		// 回滚casbin相关事务
		Ctx.Rollback()
		return fmt.Errorf("DeleteApiRep -> 策略加载失败， 已有的会忽略 -> %s", err)
	}

	// 删除casbin_rule表中的信息
	for _, id := range req.ID {

		err = casbinServer.DeletePermForAdminOrUser(fmt.Sprintf("%v", id))
		if err != nil {
			// 回滚普通事务
			tx.Rollback()
			// 回滚casbin相关事务
			Ctx.Rollback()
			return fmt.Errorf("DeleteApiRep ->  删除casbin_rule表中的信息失败 -> %s", err)
		}
	}

	// 如果需要持久化到数据库
	if err = casbinServer.Enforcer.SavePolicy(); err != nil {
		// 回滚普通事务
		tx.Rollback()
		// 回滚casbin相关事务
		Ctx.Rollback()
		return fmt.Errorf("DeleteApiRep -> 保存策略失败: %s", err)
	}

	// 提交casbin相关事务
	err = Ctx.Commit().Error
	if err != nil {
		// 回滚普通事务
		tx.Rollback()
		// 回滚casbin相关事务
		Ctx.Rollback()
		return fmt.Errorf("DeleteApiRep -> 提交casbin相关事务失败 -> %s", err)
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		// 回滚普通事务
		tx.Rollback()
		return fmt.Errorf("DeleteApiRep -> 提交事务失败 -> %s", err)
	}

	return nil

}

// UpdateApiRep 编辑api
func UpdateApiRep(db *gorm.DB, req *requests.UpdateApiReq) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("UpdateApiRep -> 开启事务失败 -> %s", tx.Error)
	}

	// 查询该api是否存在
	var api models.Api
	err := tx.Where("id = ?", req.ID).First(&api).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UpdateApiRep -> 该api不存在 -> %s", err)
	}

	// 更新 Api 表中的 path , brief_introduction
	err = tx.Model(&api).Updates(map[string]interface{}{
		"path":               req.Path,
		"brief_introduction": req.BriefIntroduction,
	}).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UpdateApiRep -> 更新 Api 表中的 path , brief_introduction失败 -> %s", err)
	}

	// 查询新的分组的id
	var groupId uint
	resg := tx.Model(&models.DictItem{}).Select("ID").Where("label = ?", req.Grouping).First(&groupId)
	if resg.Error != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UpdateApiRep -> 查询 group_id 失败 -> %s", resg.Error)
	}

	// 查询 ApiDictItemGroup 表中是否存在该api相关的记录
	var count int64
	err = tx.Model(&models.ApiDictItemGroup{}).Where("api_id = ?", api.ID).Count(&count).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UpdateApiRep -> 查询 ApiDictItemGroup 表中是否存在该api相关的记录 失败 -> %s", err)
	}

	if count > 0 {
		// 更新 api_group 表中的 group_id
		resag := tx.Model(&models.ApiDictItemGroup{}).Where("api_id = ?", api.ID).Update("group_id", groupId)
		if resag.Error != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UpdateApiRep -> 更新 api_group 表中的 group_id失败 -> %s", resag.Error)
		}
	} else {
		// 如果表中不存在该api相关的记录，需要将更新数据插入到表中
		adig := models.ApiDictItemGroup{
			ApiId:   api.ID,
			GroupId: groupId,
		}
		err = tx.Model(&models.ApiDictItemGroup{}).Create(&adig).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UpdateApiRep -> 向 api_group 表中的 插入数据 失败 -> %s", err)
		}
	}

	// 查询 request_method_id
	var reqMethodId uint
	resm := tx.Model(&models.DictItem{}).Select("ID").Where("label = ?", req.RequestMethod).First(&reqMethodId)
	if resm.Error != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UpdateApiRep -> 查询 request_method_id 失败 -> %s", resm.Error)
	}

	// 查询 ApiDictItemRequestMethod 表中是否存在该api相关的记录
	var count2 int64
	err = tx.Model(&models.ApiDictItemRequestMethod{}).Where("api_id = ?", api.ID).Count(&count2).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UpdateApiRep -> 查询 ApiDictItemRequestMethod 表中是否存在该api相关的记录 失败 -> %s", err)
	}

	if count2 > 0 {
		// 更新 api_request_method 表中的 request_method_id
		resam := tx.Model(&models.ApiDictItemRequestMethod{}).Where("api_id = ?", api.ID).Update("request_method_id", reqMethodId)
		if resam.Error != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UpdateApiRep -> 更新 api_request_method 表中的 request_method_id失败 -> %s", resam.Error)
		}
	} else {
		// 如果表中不存在该api相关的记录，需要将更新数据插入到表中
		adim := models.ApiDictItemRequestMethod{
			ApiId:           api.ID,
			RequestMethodId: reqMethodId,
		}
		err = tx.Model(&models.ApiDictItemRequestMethod{}).Create(&adim).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("UpdateApiRep -> 向 api_request_method 表中的 插入数据 失败 -> %s", err)
		}
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("UserAccountRequest -> 提交事务失败 -> %s", err)
	}

	return nil

}

// CreateApiRep 添加api
func CreateApiRep(db *gorm.DB, req *requests.CreateApiReq) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		tx.Rollback()
		return fmt.Errorf("CreateApiRep -> 开启事务失败 -> %s", tx.Error)
	}

	api := models.Api{
		Path:              req.Path,
		BriefIntroduction: req.BriefIntroduction,
	}

	// 查询该api是否存在
	resa := tx.Model(models.Api{}).Where("path = ?", req.Path).First(&api)
	if resa.Error != nil {
		if errors.Is(resa.Error, gorm.ErrRecordNotFound) {
			// 如果没有找到，说明该api在数据库表中不存在，可以进行添加
		} else {
			tx.Rollback()
			return fmt.Errorf("CreateApiRep -> 查询该api是否存在 失败 -> %v", resa.Error)
		}
	} else if resa.RowsAffected > 0 {
		tx.Rollback()
		return fmt.Errorf("该api已经存在")
	}

	// 向api表中添加api相关信息
	err := tx.Model(models.Api{}).Create(&api).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateApiRep -> 向api表中添加api相关信息失败 -> %s", err)
	}

	// 查询分组id
	group := models.DictItem{
		Label: req.Grouping,
	}

	// 查询dict_item表中是否已经存在该分组
	resg := tx.Model(&models.DictItem{}).Where("label = ?", req.Grouping).First(&group)
	if resg.Error != nil {
		if errors.Is(resa.Error, gorm.ErrRecordNotFound) {
			// 如果没有找到，说明该分组在数据库表中不存在，不可以进行添加
			tx.Rollback()
			return fmt.Errorf("该分组不存在")
		} else {
			tx.Rollback()
			return fmt.Errorf("CreateApiRep -> 查询dict_item表中是否已经存在该分组失败 -> %s", resg.Error)
		}
		/*if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有查到，说明表中还没有该分组，将分组添加进dict_item表中
			err = tx.Model(&models.DictItem{}).Create(&group).Error
			if err != nil {
				tx.Rollback() // 回滚事务
				return fmt.Errorf("CreateApiRep -> 将分组添加进表中失败 -> %s", err)
			}
		} else {
			return fmt.Errorf("CreateApiRep -> 查询dict_item表中是否已经存在该分组失败 -> %s", err)
		}*/
	}

	// 向api_group表中添加api相关信息
	apiGroup := &models.ApiDictItemGroup{
		ApiId:   api.ID,
		GroupId: group.ID,
	}

	err = tx.Model(&models.ApiDictItemGroup{}).Create(apiGroup).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateApiRep -> 向api_group表中添加api相关信息失败 -> %s", err)
	}

	// 查询请求方法id
	reqMethod := &models.DictItem{
		Label: req.RequestMethod,
	}

	// 查询dict_item表中是否已经存在该请求方法
	resm := tx.Model(&models.DictItem{}).Where("label = ?", req.RequestMethod).First(reqMethod)
	if resm.Error != nil {
		if errors.Is(resa.Error, gorm.ErrRecordNotFound) {
			// 如果没有找到，说明该请求方法在数据库表中不存在，不可以进行添加
			tx.Rollback()
			return fmt.Errorf("该请求方法不存在")
		} else {
			tx.Rollback()
			return fmt.Errorf("CreateApiRep -> 查询dict_item表中是否已经存在该请求方法失败 -> %s", resm.Error)
		}
		/*if errors.Is(err, gorm.ErrRecordNotFound) {
			// 没有查到，说明表中还没有该请求方法，将请求方法添加进dict_item表中
			err = tx.Model(&models.DictItem{}).Create(reqMethod).Error
			if err != nil {
				tx.Rollback() // 回滚事务
				return fmt.Errorf("CreateApiRep -> 将请求方法添加进表中失败 -> %s", err)
			}
		} else {
			return fmt.Errorf("CreateApiRep -> 查询dict_item表中是否已经存在该请求方法失败 -> %s", err)
		}*/

	}

	// 向api_request_method表中添加api相关信息
	apiReqMethod := &models.ApiDictItemRequestMethod{
		ApiId:           api.ID,
		RequestMethodId: reqMethod.ID,
	}

	err = tx.Model(&models.ApiDictItemRequestMethod{}).Create(apiReqMethod).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateApiRep -> 向api_request_method表中添加api相关信息失败 -> %s", err)
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("CreateApiRep -> 提交事务失败 -> %s", err)
	}

	return nil

}

// SearchApiListRep 检索api列表
func SearchApiListRep(db *gorm.DB, req *requests.SearchApiListReq) (*requests.SearchApiListRes, error) {

	if req.Page <= 0 || req.Limit <= 0 {
		return nil, fmt.Errorf("SearchApiListRep -> Page或Limit的值不能小于或等于0")
	}

	var searchApiRes []requests.SearchApiRes
	var res requests.SearchApiListRes

	//建立表关联

	query := db.Table("sw_apis").
		Select("sw_apis.id, sw_apis.path, sw_apis.brief_introduction, dg.id as group_id, dg.label as group_name, dm.id as request_method_id, dm.label as request_method_name").
		Joins("left join sw_api_dict_item_groups on sw_api_dict_item_groups.api_id = sw_apis.id").
		Joins("left join sw_api_dict_item_request_methods on sw_api_dict_item_request_methods.api_id = sw_apis.id").
		Joins("left join sw_dict_items As dg on sw_api_dict_item_groups.group_id = dg.id").
		Joins("left join sw_dict_items As dm on sw_api_dict_item_request_methods.request_method_id = dm.id")

	// 添加查询条件
	if req.Path != "" {
		query = query.Where("sw_apis.path LIKE ?", "%"+req.Path+"%")
	}
	if req.RequestMethod != "" {
		query = query.Where("dm.label = ?", req.RequestMethod)
	}
	if req.Grouping != "" {
		query = query.Where("dg.label = ?", req.Grouping)
	}
	if req.BriefIntroduction != "" {
		query = query.Where("sw_apis.brief_introduction = ?", req.BriefIntroduction)
	}

	query = query.Where("sw_apis.deleted_at IS NULL").Order("sw_apis.created_at DESC")

	err := query.Scan(&searchApiRes).Error
	if err != nil {
		return nil, fmt.Errorf("SearchApiListRep -> 查询api异常 -> %s", err)
	}
	length := len(searchApiRes)
	if length == 0 {
		//return nil, fmt.Errorf("SearchApiListRep -> 不存在该api")
		res2 := &requests.SearchApiListRes{
			Api:   make([]requests.SearchApiRes, 0),
			Total: 0,
		}
		return res2, nil
	}

	// 分页返回数据
	offset := (req.Page - 1) * req.Limit
	if length > offset {
		end := offset + req.Limit
		if end > length {
			end = length
		}
		res.Api = searchApiRes[offset:end]
		res.Total = length
		return &res, nil
	}

	res2 := &requests.SearchApiListRes{
		Api:   make([]requests.SearchApiRes, 0),
		Total: 0,
	}
	return res2, nil
}

// MethodAndGroup
// @Description: 存储根据api的id获取api的 分组的id和name 请求方式的id和name
// @Author wangyulong 2024-10-16 20:15:36
type MethodAndGroup struct {
	RequestMethodId   uint   `json:"request_method_id"`
	RequestMethodName string `json:"request_method_name"`
	GroupId           uint   `json:"group_id"`
	GroupName         string `json:"group_name"`
}

type MiddleResult struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// GetApiGroupAndMethod
// @Description: 根据api的id获取api的 分组的id和name 请求方式的id和name
// @Author wangyulong 2024-10-16 20:37:20
// @param        db *gorm.DB
// @param        apiID uint
// @return       *MethodAndGroup
// @return       error
func GetApiGroupAndMethod(db *gorm.DB, apiID uint) (*MethodAndGroup, error) {

	var result MethodAndGroup

	// 使用GORM进行多表联查

	err := db.Table("sw_api_dict_item_groups").
		Select("dg.id As group_id, dg.label As group_name, dm.id As request_method_id, dm.label As request_method_name").
		Joins("join sw_dict_items As dg on dg.id = sw_api_dict_item_groups.group_id").
		Joins("join sw_dict_items As dm on dm.id = sw_api_dict_item_request_methods.request_method_id").
		Joins("join sw_api_dict_item_request_methods on sw_api_dict_item_request_methods.api_id = sw_api_dict_item_groups.api_id").
		Where("sw_api_dict_item_groups.api_id = ?", apiID).Scan(&result)
	if err != nil {
		return nil, fmt.Errorf("GetApiGrouGetApiGroupAndMethod -> 根据api的id获取api的 分组的id和name 请求方式的id和name失败 -> %s", err)
	}

	return &result, nil
}
