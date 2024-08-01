package models

import (
	"gorm.io/gorm"
)

// User 用户简略信息
type User struct {
	gorm.Model         //ID CreatedAt UpdatedAt DeletedAt
	Nickname    string `json:"nickname"`            // 昵称
	Email       string `json:"email" gorm:"unique"` // 邮箱，唯一
	Password    string `json:"password"`            // 密码
	Heat        int    `json:"heat"`                // 个人热度
	FansCount   uint   `json:"fans_count"`          // 粉丝数
	UserDetail  UserDetail
	UserMessage UserMessage
	//用户和标签之间有2个关系表(暂时只写一个)
	Tags               []Tag     `gorm:"many2many:user_tags"`
	Articles           []Article // 写的文章
	ArticleLikes       []Article `gorm:"many2many:article_likes"`                                                    //点赞的人
	ArticleCollections []Article `gorm:"many2many:article_collections"`                                              // 收藏的文章
	ArticleComments    []Article `gorm:"many2many:article_comments"`                                                 // 文章的评论
	Users              []User    `gorm:"many2many:user_follows;joinForeignKey:FollowedID;JoinReferences:FollowerID"` // 关注的用户列表
}
