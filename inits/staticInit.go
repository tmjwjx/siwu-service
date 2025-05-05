package inits

import (
	"forum/internal/internalPkg/internalUtils"
	"forum/pkg/globals"
	"github.com/spf13/viper"
	"strconv"
)

// StaticInit
// @Description: // 配置静态文件目录
// @Author wangyulong 2025-01-17 14:43:31
func StaticInit() {
	// 配置静态文件目录
	// 将文件系统中的目录映射到 URL 路径
	if err := viper.UnmarshalKey("static", &globals.SConfig); err != nil {
		globals.Log.Fatalf("Run -> 无法解码为结构: %s", err)
	}

	//globals.Router.Static(globals.SConfig.Prefix, globals.SConfig.Path)

	// 创建存储静态文件的目录路径文件夹
	err := internalUtils.CreateFolder(globals.SConfig.Path)
	if err != nil {
		globals.Log.Errorf("创建存储静态文件的目录路径文件夹")
	}
}

// PartPathPrefixInit
// @Description: 生产图片的部分路径前缀
// @Author wangyulong 2025-01-17 14:40:04
func PartPathPrefixInit() {
	path := "http://" + globals.AppConfig.App.Domain + ":" + strconv.Itoa(globals.AppConfig.App.Port) + globals.SConfig.Prefix + "/"

	// UserDefaultImage 默认用户头像路径
	internalUtils.UserDefaultImage = path + globals.SConfig.UserDefaultImage

	// ArticleDefaultImage 默认文章图片路径
	internalUtils.ArticleDefaultImage = path + globals.SConfig.ArticleDefaultImage

	// TagDefaultImage 默认标签图片
	internalUtils.TagDefaultImage = path + globals.SConfig.TagDefaultImage

	// AdvertisementDefaultImage 默认广告图片
	internalUtils.AdvertisementDefaultImage = path + globals.SConfig.AdvertisementDefaultImage

	// CommentDefaultImage 默认评论图片
	internalUtils.CommentDefaultImage = path + globals.SConfig.CommentDefaultImage

	// CategoryDefaultImage 默认类目图片
	internalUtils.CategoryDefaultImage = path + globals.SConfig.CategoryDefaultImage

}
