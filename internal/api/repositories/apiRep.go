package repositories

import (
	"fmt"
	"forum/internal/api/requests"
	"forum/internal/models"
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
	var groups []models.Group
	// 查询所有的分组id
	err := db.Model(models.Group{}).Find(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("GetAllApiRep -> 查询所有的分组id失败 -> %s", err)
	}

	for _, group := range groups {

		list := &requests.List{
			GroupID:   group.ID,
			GroupName: group.Name,
			Children:  make([]*requests.Children, 0),
		}
		var apiIDs []uint
		// 查询每个分组下面拥有的apiID
		err = db.Model(&models.ApiGroup{}).Select("api_id").Where("group_id = ?", group.ID).Find(&apiIDs).Error
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
//		err = db.Model(&models.ApiRequestMethod{}).Select("RequestMethodId").Where("api_id = ?", api.ID).First(&requestMethodId).Error
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
//		err = db.Model(models.ApiGroup{}).Select("group_id").Where("api_id = ?", api.ID).First(&groupId).Error
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

	var api models.Api
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

	return apiDetailsRes, nil

}

// GetGroupListRep 获取所有api分组列表
func GetGroupListRep(db *gorm.DB) (*requests.ApiGroupRes, error) {

	var apiGroupRes requests.ApiGroupRes
	var groups []models.Group
	// 从group表中查询 分组信息
	err := db.Find(&groups).Error
	if err != nil {
		return nil, fmt.Errorf("GetGroupListRep -> 从group表中查询 分组信息失败 -> %s", err)
	}

	for _, group := range groups {

		apiGroup := &requests.ApiGroup{
			Value: group.Name,
			Label: group.Name,
		}

		apiGroupRes.Groups = append(apiGroupRes.Groups, apiGroup)
	}

	return &apiGroupRes, nil

}

// GetRequestMethodRep 获取所有请求方法
func GetRequestMethodRep(db *gorm.DB) (*requests.ApiReqMethodRes, error) {

	var apiReqMethodRes requests.ApiReqMethodRes
	var reqMethods []models.RequestMethod
	// 从request_method表中查询 分组信息
	err := db.Find(&reqMethods).Error
	if err != nil {
		return nil, fmt.Errorf("GetRequestMethodRep -> 从request_method表中查询 分组信息失败 -> %s", err)
	}

	for _, reqMethod := range reqMethods {

		apiReqMethod := &requests.ApiReqMethod{
			Value: reqMethod.Name,
			Label: reqMethod.Name,
		}

		apiReqMethodRes.Methods = append(apiReqMethodRes.Methods, apiReqMethod)
	}

	return &apiReqMethodRes, nil
}

// DeleteApiRep 删除api
func DeleteApiRep(req *requests.DeleteApiReq, db *gorm.DB) error {

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("DeleteApiRep -> 开启事务失败 -> %s", tx.Error)
	}

	for _, id := range req.ID {

		// 删除api表中的信息
		result := tx.Model(&models.Api{}).Where("id = ?", id).Delete(nil)
		if result.Error != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("DeleteApiRep ->  删除api表中的信息失败 -> %s", result.Error)
		} else if result.RowsAffected == 0 {
			return fmt.Errorf("DeleteApiRep ->  该api不存在")
		}

		// 删除api_group表中的信息
		result2 := tx.Model(&models.ApiGroup{}).Where("api_id = ?", id).Delete(nil)
		if result2.Error != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("DeleteApiRep ->  删除api_group表中的信息失败 -> %s", result2.Error)
		}

		// 删除api_request_method表中的信息
		result3 := tx.Model(&models.ApiRequestMethod{}).Where("api_id = ?", id).Delete(nil)
		if result3.Error != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("DeleteApiRep ->  删除api_request_method表中的信息失败 -> %s", result3.Error)
		}
	}

	// 提交事务
	err := tx.Commit().Error
	if err != nil {
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

	// 查询 group_id
	var groupId uint
	err = tx.Model(&models.Group{}).Select("ID").First(&groupId).Error
	if err != nil {
		return fmt.Errorf("UpdateApiRep -> 查询 group_id失败 -> %s", err)
	}

	// 更新 api_group 表中的 group_id
	err = tx.Model(&models.ApiGroup{}).Where("api_id = ?", api.ID).Update("GroupId", groupId).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UpdateApiRep -> 更新 api_group 表中的 group_id失败 -> %s", err)
	}

	// 查询 request_method_id
	var reqMethodId uint
	err = tx.Model(&models.RequestMethod{}).Select("ID").First(&reqMethodId).Error
	if err != nil {
		return fmt.Errorf("UpdateApiRep -> 查询 group_id失败 -> %s", err)
	}

	// 更新 api_request_method 表中的 request_method_id
	err = tx.Model(&models.ApiRequestMethod{}).Where("api_id = ?", api.ID).Update("RequestMethodId", reqMethodId).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("UpdateApiRep -> 更新 api_request_method 表中的 request_method_id失败 -> %s", err)
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

	//var groupId uint
	//var reqMethodId uint

	// 开启事务
	tx := db.Begin()
	if tx.Error != nil {
		return fmt.Errorf("CreateApiRep -> 开启事务失败 -> %s", tx.Error)
	}

	api := models.Api{
		Path:              req.Path,
		BriefIntroduction: req.BriefIntroduction,
	}
	//查询该api是否存在
	err := tx.Model(models.Api{}).Where("path = ?", req.Path).First(&api).Error
	if err == nil {
		return fmt.Errorf("CreateApiRep -> 该api已经存在")
	}

	// 向api表中添加api相关信息
	err = tx.Model(models.Api{}).Create(&api).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateApiRep -> 向api表中添加api相关信息失败 -> %s", err)
	}

	//// 查询新增api的ID
	//err = tx.Model(models.Api{}).First(&api).Error
	//if err != nil {
	//	tx.Rollback() // 回滚事务
	//	return fmt.Errorf("CreateApiRep -> 查询新增api的ID失败 -> %s", err)
	//}

	// 查询分组id
	group := &models.Group{
		Name: req.Grouping,
	}
	err = tx.Model(&models.Group{}).Where("name = ?", req.Grouping).First(&group).Error
	if err != nil {
		// 没有查到，说明表中还没有该分组，将分组添加进表中
		err = tx.Model(&models.Group{}).Create(group).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("CreateApiRep -> 将分组添加进表中失败 -> %s", err)
		}
		//// 查询分组id
		//err = tx.Model(&models.Group{}).Select("ID").First(&groupId).Error
		//if err != nil {
		//	return fmt.Errorf("CreateApiRep -> 查询分组id失败 -> %s", err)
		//}
	}
	//else {
	//	groupId = group.ID
	//}

	// 向api_group表中添加api相关信息
	apiGroup := &models.ApiGroup{
		ApiId:   api.ID,
		GroupId: group.ID,
	}

	err = tx.Model(&models.ApiGroup{}).Create(apiGroup).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateApiRep -> 向api_group表中添加api相关信息失败 -> %s", err)
	}

	// 查询请求方法id
	reqMethod := &models.RequestMethod{
		Name: req.RequestMethod,
	}
	err = tx.Model(&models.RequestMethod{}).Where("name = ?", req.RequestMethod).First(reqMethod).Error
	if err != nil {
		// 没有查到，说明表中还没有该分组，将分组添加进表中
		err = tx.Model(&models.RequestMethod{}).Create(reqMethod).Error
		if err != nil {
			tx.Rollback() // 回滚事务
			return fmt.Errorf("CreateApiRep -> 将请求方法添加进表中失败 -> %s", err)
		}
		//// 查询请求方法id
		//err = tx.Model(&models.RequestMethod{}).First(&reqMethodId).Error
		//if err != nil {
		//	return fmt.Errorf("CreateApiRep -> 查询请求方法id失败 -> %s", err)
		//}
	}
	//else {
	//	reqMethodId = reqMethod.ID
	//}

	// 向api_request_method表中添加api相关信息
	apiReqMethod := &models.ApiRequestMethod{
		ApiId:           api.ID,
		RequestMethodId: reqMethod.ID,
	}

	err = tx.Model(&models.ApiRequestMethod{}).Create(apiReqMethod).Error
	if err != nil {
		tx.Rollback() // 回滚事务
		return fmt.Errorf("CreateApiRep -> 向api_request_method表中添加api相关信息失败 -> %s", err)
	}

	// 提交事务
	err = tx.Commit().Error
	if err != nil {
		return fmt.Errorf("CreateApiRep -> 提交事务失败 -> %s", err)
	}

	return nil

}

// SearchApiListRep 检索api列表
func SearchApiListRep(db *gorm.DB, req *requests.SearchApiListReq) (*requests.SearchApiListRes, error) {

	if req.Limit == 0 {
		return nil, fmt.Errorf("SearchApiListRep -> Limit的值不能为0")
	}

	var searchApiRes []requests.SearchApiRes
	var res requests.SearchApiListRes

	//建立表关联

	query := db.Table("sw_apis").
		Debug().
		Select("sw_apis.id, sw_apis.path, sw_apis.brief_introduction, sw_groups.id as group_id, sw_groups.name as group_name, sw_request_methods.id as request_method_id, sw_request_methods.name as request_method_name").
		Joins("left join sw_api_groups on sw_api_groups.api_id = sw_apis.id").
		Joins("left join sw_groups on sw_api_groups.group_id = sw_groups.id").
		Joins("left join sw_api_request_methods on sw_api_request_methods.api_id = sw_apis.id").
		Joins("left join sw_request_methods on sw_api_request_methods.request_method_id = sw_request_methods.id")

	// 添加查询条件
	if req.Path != "" {
		query = query.Where("sw_apis.path = ?", req.Path)
	}
	if req.RequestMethod != "" {
		query = query.Where("sw_request_methods.name = ?", req.RequestMethod)
	}
	if req.Grouping != "" {
		query = query.Where("sw_groups.name = ?", req.Grouping)
	}
	if req.BriefIntroduction != "" {
		query = query.Where("sw_apis.brief_introduction = ?", req.BriefIntroduction)
	}

	query = query.Where("sw_apis.deleted_at IS NULL")

	err := query.Limit(req.Limit).Offset(req.Page).Scan(&searchApiRes).Error
	if err != nil {
		return nil, fmt.Errorf("SearchApiListRep -> 查询api异常 -> %s", err)
	}
	if len(searchApiRes) == 0 {
		//return nil, fmt.Errorf("SearchApiListRep -> 不存在该api")
		res2 := &requests.SearchApiListRes{
			Api:   make([]*requests.SearchApiRes, 0),
			Total: 0,
		}
		return res2, nil
	}

	for _, searchApi := range searchApiRes {
		res.Api = append(res.Api, &searchApi)
		res.Total += 1
	}

	return &res, nil
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
	err := db.Table("sw_api_request_methods").
		Select("sw_api_request_methods.request_method_id, sw_request_methods.name as request_method_name, sw_api_groups.group_id, sw_groups.name as group_name").
		Joins("join sw_request_methods on sw_api_request_methods.request_method_id = sw_request_methods.id").
		Joins("join sw_api_groups on sw_api_request_methods.api_id = sw_api_groups.api_id").
		Joins("join sw_groups on sw_api_groups.group_id = sw_groups.id").
		Where("sw_api_request_methods.api_id = ?", apiID).
		Scan(&result).Error
	if err != nil {
		return nil, fmt.Errorf("GetApiGrouGetApiGroupAndMethod -> 根据api的id获取api的 分组的id和name 请求方式的id和name失败 -> %s", err)
	}
	return &result, nil
}

//// GetApiGroupAndMethod2 根据api的id获取api的 分组的id和name 请求方式的id和name
//func GetApiGroupAndMethod2(db *gorm.DB, id uint) (uint, string, uint, string, error) {
//	var requestMethodId uint
//	var requestMethodName string
//	var groupId uint
//	var groupName string
//
//	// 从api_request_method表中查询 request_method_id
//	err := db.Model(&models.ApiRequestMethod{}).Select("RequestMethodId").Where("api_id = ?", id).First(&requestMethodId).Error
//	if err != nil {
//		return 0, "", 0, "", fmt.Errorf("GetApiGrouGetApiGroupAndMethod -> 获取api请求方法id失败 -> %s", err)
//	}
//
//	// 从requestMethod表中查询 请求方法的name
//	err = db.Model(&models.RequestMethod{}).Select("Name").Where("id = ?", requestMethodId).First(&requestMethodName).Error
//	if err != nil {
//		return 0, "", 0, "", fmt.Errorf("GetApiGroupAndMethod -> 获取api请求方法name失败 -> %s", err)
//	}
//
//	// 从api_group表中查询 group_id
//	err = db.Model(models.ApiGroup{}).Select("group_id").Where("api_id = ?", id).First(&groupId).Error
//	if err != nil {
//		return 0, "", 0, "", fmt.Errorf("GetApiGroupAndMethod -> 获取api分组id失败 -> %s", err)
//	}
//
//	// 从group表中查询 api分组的name
//	err = db.Model(models.Group{}).Select("name").Where("id = ?", groupId).First(&groupName).Error
//	if err != nil {
//		return 0, "", 0, "", fmt.Errorf("GetApiGroupAndMethod -> 获取api分组name失败 -> %s", err)
//	}
//
//	return groupId, groupName, requestMethodId, requestMethodName, nil
//}
