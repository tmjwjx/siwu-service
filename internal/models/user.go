package models

import (
	"time"
)

// User 用户简略信息
type User struct {
	ID            uint          `json:"id" gorm:"primaryKey"`                // 主键
	Nickname      string        `json:"nickname"`                            // 昵称
	Email         string        `json:"email" gorm:"unique"`                 // 邮箱，唯一
	Password      string        `json:"password"`                            // 密码
	PersonalHeat  int           `json:"personal_heat"`                       // 个人热度
	FansCount     uint          `json:"fans_count"`                          // 粉丝数
	UserHeadImage UserHeadImage `gorm:"foreignKey:UserID;references:ID"`     // 关联 UserHeadImage 表，UserID -> ID
	UserDetail    UserDetail    `gorm:"foreignKey:UserID;references:ID"`     // 关联 UserDetail 表，UserID -> ID
	UserMessage   UserMessage   `gorm:"foreignKey:UserID;references:ID"`     // 关联 UserMessage 表，UserID -> ID
	Follow1       Follow        `gorm:"foreignKey:FollowerID;references:ID"` // 关联 Follow 表，FollowerID -> ID
	Follow2       Follow        `gorm:"foreignKey:FollowedID;references:ID"` // 关联 Follow 表，FollowedID -> ID
}

// UserHeadImage 头像信息
type UserHeadImage struct {
	ID      uint   `json:"id" gorm:"primaryKey"`  // 用户头像图片
	UserID  uint   `json:"user_id" gorm:"unique"` // 用户ID，外键，唯一
	HeadImg []byte `json:"head_img"`              // 头像照片
	ImgName string `json:"img_name"`              // 图像名称
	ImgType string `json:"img_type"`              // 图像类型
	ImgSize uint   `json:"img_size"`              // 图像大小
}

// UserDetail 用户详细信息
type UserDetail struct {
	ID                uint       `json:"id" gorm:"primaryKey"`  // 主键
	UserID            uint       `json:"user_id" gorm:"unique"` // 用户ID，外键，唯一索引
	CareerDirection   string     `json:"career_direction"`      // 职业方向
	UserHomePage      string     `json:"user_home_page"`        // 个人主页
	UserSignature     string     `json:"user_signature"`        // 个人签名
	UserTags          string     `json:"user_tags"`             // 作者标签
	BlogLink          string     `json:"blog_link"`             // 个人博客链接
	WeiboLink         string     `json:"weibo_link"`            // 新浪微博链接
	GithubLink        string     `json:"github_link"`           // Github链接
	LastLoginTime     *time.Time `json:"last_login_time"`       // 上次登录时间
	CreateAccountTime *time.Time `json:"create_account_time"`   // 创建账号时间
}

// UserMessage 用户消息表(小铃铛专属)
type UserMessage struct {
	ID                uint `json:"id" gorm:"primaryKey"`       // 主键
	UserID            uint `json:"user_id" gorm:"uniqueIndex"` // 用户ID，外键，唯一索引
	UnreadComments    int  `json:"unread_comments"`            // 未读评论数量
	UnreadLikes       int  `json:"unread_likes"`               // 未读点赞数量
	UnreadCollections int  `json:"unread_collections"`         // 未读收藏数量
	UnreadFans        int  `json:"unread_fans"`                // 未读粉丝数量
	UnreadMessages    int  `json:"unread_messages"`            // 未读私信数量
	UnreadSystem      int  `json:"unread_system"`              // 未读系统消息数量
}

// Follow 关注关系表
type Follow struct {
	ID         uint `json:"id" gorm:"primaryKey"`     // 主键
	FollowerID uint `json:"follower_id" gorm:"index"` // 关注者ID
	FollowedID uint `json:"followed_id" gorm:"index"` // 被关注者ID
}
