package inits

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/utils"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// DBInit 初始化mysql
func DBInit() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.AppConfig.Database.User,
		utils.AppConfig.Database.Password,
		utils.AppConfig.Database.Host,
		utils.AppConfig.Database.Port,
		utils.AppConfig.Database.Name,
	)

	var err error
	utils.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

// TableInit 初始化表
func TableInit() {
	// 用户模块
	err := utils.DB.AutoMigrate(&models.User{}, &models.UserDetail{}, &models.UserMessage{}, &models.UserHeadImage{}, &models.Follow{})
	if err != nil {
		fmt.Println("db.AutoMigrate err = ", err)
		return
	}

	// 文章模块
	err = utils.DB.AutoMigrate(&models.Comment{}, &models.ArticleLike{}, &models.ArticleCollection{}, &models.Attachment{}, &models.Category{}, &models.Article{})
	if err != nil {
		fmt.Println("db.AutoMigrate err = ", err)
		return
	}

	// 标签模块
	err = utils.DB.AutoMigrate(&models.Tag{}, &models.Resource{}, &models.Administrator{}, &models.ArticleTag{}, &models.UserFollowsTag{})
	if err != nil {
		return
	}

}
