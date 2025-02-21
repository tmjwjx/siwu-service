package repositories

import (
	"forum/internal/administrator/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"gorm.io/gorm"
	"time"
)

// DeleteAdministratorRep
// @Description: 删除管理员
// @param        db *gorm.DB
// @param        deleteId int
// @return       error
// @Author tianjiajie 2025-02-21 15:06:32
func DeleteAdministratorRep(db *gorm.DB, req requests.DeleteAdministratorReq) (err error) {
	administrator := models.Administrator{}
	err = db.Where("id = ?", req.ID).Delete(&administrator).Error
	return err
}

// CheckAdministratorEmail
// @Description: 邮箱查重
// @param        db *gorm.DB
// @param        email string
// @return       bool
// @Author tianjiajie 2025-02-20 20:57:08
func CheckAdministratorEmail(db *gorm.DB, req requests.AddAdministratorReq) bool {
	administrator := models.Administrator{}
	db.Where("email = ?", req.Email).First(&administrator)
	// 如果ID不为0，说明存在
	return administrator.ID != 0
}

// AddAdministratorRep
// @Description: 添加管理员
// @param        db *gorm.DB
// @param        email string
// @return       data
// @return       err
// @Author tianjiajie 2025-02-20 20:51:39
func AddAdministratorRep(db *gorm.DB, req requests.AddAdministratorReq) (id uint, err error) {

	password, err := internalUtils.HashPassword(req.Email)
	if err != nil {
		return 0, err
	}

	admin := models.Administrator{}
	db.Unscoped().Where("email = ?", req.Email).First(&admin)

	if admin.ID != 0 { // 存在这个管理员，不过已经软删除了，回复即可

		admin.Name = req.Email
		admin.Email = req.Email
		admin.Password = password
		admin.DeletedAt = gorm.DeletedAt{
			Time:  time.Time{},
			Valid: false,
		}

		err = db.Save(&admin).Error
		if err != nil {
			return 0, err
		}

	} else {
		admin = models.Administrator{
			Name:     req.Email,
			Email:    req.Email,
			Password: password,
		}
		err = db.Save(&admin).Error
		if err != nil {
			return 0, err
		}

	}
	return admin.ID, nil
}
