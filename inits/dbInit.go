package inits

import (
	"fmt"
	"forum/internal/models"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
		DisableForeignKeyConstraintWhenMigrating: true, // 取消外键约束
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "sw_", // 设置表前缀
		},
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

// TableInit
//
//	@Description: 初始化表
func TableInit() {
	// 用户模块

	err := globals.DB.AutoMigrate(
		&models.Administrator{},
		&models.Advertisement{},
		&models.Attachment{},
		&models.Resource{},

		&models.User{},
		&models.Category{},
		&models.Article{},
		&models.Tag{},

		&models.ArticleComment{},
		&models.ArticleCollection{},
		&models.ArticleLike{},
		&models.ArticleTag{},
		&models.CommentLike{},

		&models.UserDetail{},
		&models.UserFollow{},
		&models.UserMessage{},
		&models.UserVerifyCode{},
		&models.UserTag{},

		// 角色模块
		&models.Role{},
		&models.UserRole{},
	)

	if err != nil {
		globals.Log.Errorf("db.AutoMigrate err = %s", err)
		return
	}
}
