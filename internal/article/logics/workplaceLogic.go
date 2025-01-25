package logics

import (
	"fmt"
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// GetWorkplaceDataLogic
// @Description: 获取工作台数据
// @param        db *gorm.DB
// @return       date
// @return       err
// @Author tianjiajie 2025-01-16 14:43:51
func GetWorkplaceDataLogic(db *gorm.DB) (data interface{}, err error) {
	// 所有文章的数量
	articleTotal, err := repositories.GetArticleCountRep(db)
	if err != nil {
		return nil, err
	}

	// 今天新增文章的数量
	newArticle, err := repositories.GetNewAddArticleRep(db)
	if err != nil {
		return nil, err
	}
	// 今天新增文章的比例
	newAdd := fmt.Sprintf("%.2f", float64(newArticle)*100/float64(articleTotal))

	// 今天所有文章的访问量之和
	todayViews, err := repositories.GetTodayViewsRep(db)
	if err != nil {
		return nil, err
	}

	// 今天发布评论数之和
	todayComments, err := repositories.GetTodayCommentsRep(db)
	if err != nil {
		return nil, err
	}

	// 封装数据
	totalData := requests.TotalData{
		ArticleTotal:  strconv.FormatInt(articleTotal, 10),
		NewAdd:        newAdd,
		TodayViews:    todayViews,
		TodayComments: todayComments,
	}
	data = gin.H{"total_data": totalData}
	return
}

// GetHotTagsLogic
// @Description: 查询热门标签
// @param        db *gorm.DB
// @return       date
// @return       err
// @Author tianjiajie 2025-01-16 14:34:42
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
