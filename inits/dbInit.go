package inits

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// DBInit 初始化mysql
func DBInit() {

	if err := viper.UnmarshalKey("database", &globals.AppConfig.Database); err != nil {
		log.Fatalf("无法解码为结构: %s", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		globals.AppConfig.Database.User,
		globals.AppConfig.Database.Password,
		globals.AppConfig.Database.Host,
		globals.AppConfig.Database.Port,
		globals.AppConfig.Database.Name,
	)

	var err error
	globals.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// NamingStrategy: schema.NamingStrategy{
		//	TablePrefix: "t_", // 设置表前缀
		// },
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

// TableInit 初始化表
func TableInit() {
	// 用户模块
	err := globals.DB.AutoMigrate(&models.Category{},
		&models.UserDetail{}, &models.UserMessage{},
		&models.Tag{}, &models.Resource{},
		&models.Administrator{}, &models.ArticleLike{},
		&models.ArticleCollection{}, &models.Attachment{},
		&models.User{}, &models.Article{},
		&models.ArticleComment{}, &models.Advertisement{}, &models.UserVerifyCode{})
	if err != nil {
		globals.Log.Errorf("db.AutoMigrate err = %s", err)
		return
	}

}
