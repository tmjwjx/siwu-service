package repositories

import (
	"errors"
	"forum/internal/administrator/requests"
	"forum/internal/internalPkg/internalUtils"
	"forum/internal/models"
	"gorm.io/gorm"
	"time"
)

// UpdateAdministratorRep
// @Description: 编辑管理员信息
// @param        db *gorm.DB
// @param        req requests.UpdateAdministratorReq
// @return       err
// @Author tianjiajie 2025-02-21 21:58:11
func UpdateAdministratorRep(db *gorm.DB, req requests.UpdateAdministratorReq) (err error) {
	// 更新管理员信息
	tx := db.Begin()
	tx = tx.Model(models.Administrator{}).Where("id = ?", req.ID)

	if req.Name != "" {
		err = tx.Update("name", req.Name).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if req.Password != "" {
		// 密码格式验证
		b := internalUtils.IsValidPassword(req.Password)
		if !b {
			tx.Rollback()
			return errors.New("密码格式不正确")
		}

		// 密码加密
		password, err := internalUtils.HashPassword(req.Password)
		if err != nil {
			tx.Rollback()
			return err
		}

		// 更新密码
		err = tx.Update("password", password).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if req.Avatar != "" {
		str := []string{req.Avatar}
		url := internalUtils.UrlParam{
			UrlPath: str,
			Home:    "Avatar",
			HomeID:  uint(int(req.ID)),
			DB:      db,
		}
		err := internalUtils.StoreUrl(&url)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error

	return err
}

// ResetAdministratorPasswordRep
// @Description: 重置管理员密码
// @param        db *gorm.DB
// @param        req requests.AdministratorReq
// @return       err
// @Author tianjiajie 2025-02-21 21:34:02
func ResetAdministratorPasswordRep(db *gorm.DB, req requests.AdministratorReq) (err error) {
	// 查找用户
	admin := models.Administrator{}
	err = db.Where("id = ?", req.ID).First(&admin).Error

	// 重置密码
	password, err := internalUtils.HashPassword(admin.Email)
	if err != nil {
		return err
	}

	// 更新密码
	err = db.Model(&admin).Where("id = ?", req.ID).Update("password", password).Error
	return err
}

// GetAdministratorInfoRep
// @Description: 查询管理员详情
// @param        db *gorm.DB
// @param        id string
// @return       administrator
// @return       err
// @Author tianjiajie 2025-02-21 21:22:12
func GetAdministratorInfoRep(db *gorm.DB, id string) (administrator requests.GetAdministratorInfoRes, err error) {
	err = db.Table("sw_administrators").
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Find(&administrator).Error

	// 查询角色 todo

	// 格式化时间
	c, err := time.Parse("2006-01-02T15:04:05Z07:00", administrator.CreatedAt)
	if err != nil {
		return
	}
	l, err := time.Parse("2006-01-02T15:04:05Z07:00", administrator.LastLoginTime)
	if err != nil {
		return
	}
	administrator.CreatedAt = internalUtils.TimeFormatDetail(&c)
	administrator.LastLoginTime = internalUtils.TimeFormatDetail(&l)

	return administrator, err
}

// GetAdministratorListRep
// @Description: 查询管理员列表
// @param        db *gorm.DB
// @param        req requests.GetAdministratorListReq
// @return       res
// @return       err
// @Author tianjiajie 2025-02-21 20:35:37
func GetAdministratorListRep(db *gorm.DB, req requests.GetAdministratorListReq) (res []requests.GetAdministratorListRes, err error) {
	err = db.Table("sw_administrators").
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Find(&res).Error

	// 查询角色 todo

	// 格式化时间
	for i := range res {
		c, err := time.Parse("2006-01-02T15:04:05Z07:00", res[i].CreatedAt)
		if err != nil {
			return nil, err
		}
		l, err := time.Parse("2006-01-02T15:04:05Z07:00", res[i].LastLoginTime)
		if err != nil {
			return nil, err
		}
		res[i].CreatedAt = internalUtils.TimeFormatDetail(&c)
		res[i].LastLoginTime = internalUtils.TimeFormatDetail(&l)
	}

	return res, err
}

// DeleteAdministratorRep
// @Description: 删除管理员
// @param        db *gorm.DB
// @param        deleteId int
// @return       error
// @Author tianjiajie 2025-02-21 15:06:32
func DeleteAdministratorRep(db *gorm.DB, req requests.AdministratorReq) (err error) {
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
		admin.CreatedAt = time.Now()
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
