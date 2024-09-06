package models

import (
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

// User 用户简略信息
type User struct {
	gorm.Model                          // ID CreatedAt UpdatedAt DeletedAt
	Nickname           string           `json:"nickname" gorm:"size:16;not null"`  // 昵称
	Email              string           `json:"email" gorm:"not null;unique"`      // 邮箱，唯一
	Password           string           `json:"password" gorm:"size:100;not null"` // 密码
	Heat               int              `json:"heat" gorm:"default:0"`             // 个人热度
	AttentionCount     uint             `json:"attention_count" gorm:"default:0"`  // 关注了多少人数
	FansCount          uint             `json:"fans_count" gorm:"default:0"`       // 粉丝数
	PrivateSettings    string           `json:"private_settings"`
	UserDetail         UserDetail       // 用户详情
	UserMessage        UserMessage      // 通知用户信息
	Tags               []Tag            `gorm:"many2many:user_tags"`
	Articles           []Article        // 写的文章
	ArticleLikes       []Article        `gorm:"many2many:article_likes"`                                                    // 点赞的人
	ArticleCollections []Article        `gorm:"many2many:article_collections"`                                              // 收藏的文章
	ArticleComments    []ArticleComment `gorm:"many2many:article_comments"`                                                 // 文章的评论                                  // 用户对评论的点赞
	Users              []User           `gorm:"many2many:user_follows;joinForeignKey:FollowedID;JoinReferences:FollowerID"` // 关注的用户列表
	jwt.StandardClaims                  // StandardClaims 是jwt-go提供的标准声明结构体，包含了 exp（过期时间）、iss（发行者）等常见字段。
}