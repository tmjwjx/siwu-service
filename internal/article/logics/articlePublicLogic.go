package logics

import (
	"forum/internal/article/repositories"
	"forum/internal/article/requests"
	"gorm.io/gorm"
)

// ArticleEditLogic 编辑文章的界面所需的数据
func ArticleEditLogic(db *gorm.DB) (data interface{}, err error) {
	var t []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}
	var c []struct {
		Value string `json:"value"`
		Label string `json:"label"`
	}

	tags, err := repositories.QueryTag(db)
	if err != nil {
		return nil, err
	}
	for _, tag := range tags {
		t = append(t, struct {
			Value string `json:"value"`
			Label string `json:"label"`
		}{Value: tag.Name, Label: tag.Name})
	}

	categories, err := repositories.QueryCategory(db)
	if err != nil {
		return nil, err
	}
	for _, category := range categories {
		c = append(c, struct {
			Value string `json:"value"`
			Label string `json:"label"`
		}{Value: category.Name, Label: category.Name})
	}

	data = struct {
		Tags []struct {
			Value string `json:"value"`
			Label string `json:"label"`
		} `json:"tags"`
		Categories []struct {
			Value string `json:"value"`
			Label string `json:"label"`
		} `json:"categories"`
	}{
		Tags:       t,
		Categories: c,
	}

	return data, nil
}

// ArticlePublicLogic 发布文章
func ArticlePublicLogic(db *gorm.DB, req requests.ReqPublish) error {
	repositories.InsertArticlesRep(db, req)
	return nil
}