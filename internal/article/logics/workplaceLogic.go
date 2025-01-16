package logics

import (
	"forum/internal/article/repositories"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

func GetHotTagsLogic(db *gorm.DB) (date interface{}, err error) {

	// 查询热门标签
	tags, err := repositories.GetHotTagsRep(db, 5)
	if err != nil {
		return nil, err
	}

	date = gin.H{"pie_list": tags}

	return
}

// GetHotArticleLogic
// @Description: 查询前五篇热门文章数据
// @param        db *gorm.DB
// @return       articleList
// @return       err
// @Author tianjiajie 2025-01-15 14:53:04
func GetHotArticleLogic(db *gorm.DB) (date interface{}, err error) {

	// 查询前五篇热门文章数据
	articleList, err := repositories.GetHotArticleRep(db, 5)
	if err != nil {
		return nil, err
	}

	date = gin.H{"article_list": articleList}

	return

}

// GetTwoWeeksArticleSumLogic
// @Description: 查询近两周文章发布数量
// @param        db *gorm.DB
// @return       articleSum
// @return       err
// @Author tianjiajie 2025-01-15 09:35:16
func GetTwoWeeksArticleSumLogic(db *gorm.DB) (date interface{}, err error) {

	// 获取两周前时间
	ago := time.Now().AddDate(0, 0, -13)

	// 初始化一个数组，存放14天的文章数量
	articleSum := make([]int, 14)

	for i := 0; i < 14; i++ {
		// 计算某一天的日期
		day := ago.AddDate(0, 0, i)
		// 查询该天发布的文章数量
		count, err := repositories.GetArticleNumRep(db, day)
		if err != nil {
			return nil, err
		}
		// 将查询结果存入 articleSum 数组
		articleSum[i] = int(count)
	}

	date = gin.H{"article_sum": articleSum}

	return
}
