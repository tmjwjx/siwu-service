package inits

import (
	"fmt"
	"forum/configs"
	"forum/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

var DB *gorm.DB

// 初始化mysql
func mysqlInit() {
	configMap, _ := configs.LoadMysqlConfig("mysqlConfig.json")
	dbConfig := configMap["db"].(map[string]interface{})
	username := dbConfig["username"].(string)
	password := dbConfig["password"].(string)
	host := dbConfig["host"].(string)
	port := int(dbConfig["port"].(float64))
	Dbname := dbConfig["dbname"].(string)
	timeout := dbConfig["timeout"].(string)
	// username := "root"  // 账号
	// password := ""      // 密码
	// host := "127.0.0.1" // 数据库地址，可以是Ip或者域名
	// port := 3306        // 数据库端口
	// Dbname := "forum"   // 数据库名
	// timeout := "10s"    // 连接超时，10秒

	// root:root@tcp(127.0.0.1:3306)/gorm?
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	// 连接MYSQL, 获得DB类型实例，用于后面的数据库读写操作。
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 跳过默认事务，能获得 60% 的性能提升
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "t_", // 表名前缀
		},
	})
	if err != nil {
		log.Fatalln("mysqlInit err = ", err)
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
