package inits

import (
	"fmt"
	"forum/internal/models"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var DB *gorm.DB

// 初始化mysql
func dbInit() {

	username := viper.GetString("db.username")
	password := viper.GetString("db.password")
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")
	dbname := viper.GetString("db.dbname")
	timeout := viper.GetString("db.timeout")
	// 打印所有的配置供调试
	fmt.Println("All configurations:", viper.AllSettings())

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, dbname, timeout)
	// 连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 跳过默认事务，能获得 60% 的性能提升
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "t_", // 表名前缀
		},
	})
	if err != nil {
		log.Fatalln("dbInit err = ", err)
	}

	DB = db

	// 用户模块
	err = db.AutoMigrate(&models.User{}, &models.UserDetail{}, &models.UserMessage{}, &models.UserHeadImage{}, &models.Follow{})
	if err != nil {
		fmt.Println("db.AutoMigrate err = ", err)
		return
	}

	// 文章模块
	err = db.AutoMigrate(&models.Comment{}, &models.ArticleLike{}, &models.ArticleCollection{}, &models.Attachment{}, &models.Category{}, &models.Article{})
	if err != nil {
		fmt.Println("db.AutoMigrate err = ", err)
		return
	}

	// 标签模块
	err = db.AutoMigrate(&models.Tag{}, &models.Resource{}, &models.Administrator{}, &models.ArticleTag{}, &models.UserFollowsTag{})
	if err != nil {
		return
	}

}
