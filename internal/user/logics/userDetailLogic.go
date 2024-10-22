package logics

import (
	"fmt"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/user/repositories"
	"forum/internal/user/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PersonalDataLogic 将用户信息存入数据库
func PersonalDataLogic(userDataReq *requests.UserDataReq, c *gin.Context, db *gorm.DB) (error, int) {
	// 将除图片文件外的数据存入数据库
	err := repositories.UserDataRequest(userDataReq, db)
	if err != nil {
		return err, 500
	}

	// 将图片文件的路径相关信息存入数据库和文件系统
	u := &internalUtils.UrlParam{
		UrlPath: userDataReq.Path,
		Home:    "用户",
		HomeID:  userDataReq.ID,
		DB:      db,
	}
	err = internalUtils.StoreUrl(u)
	if err != nil {
		return err, 500
	}
	return nil, 200
}

// ResponsePersonDateLogic 将用户个人资料响应给前端
func ResponsePersonDateLogic(userID uint, db *gorm.DB) (*requests.UserDataRes, error) {
	userDataRes, err := repositories.UserDataResponse(userID, db)
	return userDataRes, err
}

// UserAccountRequestLogic 更新用户账号设置
func UserAccountRequestLogic(userAccountReq *requests.UserAccountReq, db *gorm.DB) (error, int) {
	// 验证Email是否唯一
	err := repositories.QueryPersonEmail(userAccountReq, db)
	if err == nil {
		// 查询到数据库已经存在了该 Email，返回错误
		return fmt.Errorf("UserAccountRequestLogic -> %s", "该 Email 已经存在"), 400
	}

	// 运行到这一步，已经证明数据库中还没有该 Email ，接下来将其信息存入数据库中
	err = repositories.UserAccountRequest(userAccountReq, db)
	if err != nil {
		return fmt.Errorf("UserAccountRequestLogic -> %s", err), 500
	}
	return nil, 200
}

// UserAccountResponseLogic 返回用户账号设置信息给前端
func UserAccountResponseLogic(userID uint, db *gorm.DB) (*requests.UserAccountRes, error) {
	userAccountRes, err := repositories.UserAccountResponse(userID, db)
	return userAccountRes, err
}

// UserPrivateSetRequestLogic 更新用户私信设置
func UserPrivateSetRequestLogic(userPrivateSetReq *requests.UserPrivateSettingsReq, db *gorm.DB) error {
	err := repositories.UserPrivateSetRequest(userPrivateSetReq, db)
	if err != nil {
		return err
	}
	return nil
}

// UserPrivateSetResponseLogic 响应用户私信设置
func UserPrivateSetResponseLogic(userID uint, db *gorm.DB) (*requests.UserPrivateSettingsRes, error) {
	UserPrivateSetRes, err := repositories.UserPrivateSetResponse(userID, db)
	return UserPrivateSetRes, err
}
