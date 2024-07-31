package models

import (
	"time"
)

// User 用户简略信息
type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"` // 主键
	Nickname     string `json:"nickname"`             // 昵称
	Email        string `json:"email" gorm:"unique"`  // 邮箱，唯一
	Password     string `json:"password"`             // 密码
	PersonalHeat int    `json:"personal_heat"`        // 个人热度
	FansCount    uint   `json:"fans_count"`           // 粉丝数
	// many2many:user_follows_tags 中间表名称
	FollowsTags []Tag `gorm:"many2many:user_follows_tags;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:TagID"` // 关联中间表 user_follows_tags , USerID -> ID
}

// UserHeadImage 头像信息
type UserHeadImage struct {
	ID      uint   `json:"id" gorm:"primaryKey"`  // 用户头像图片
	UserID  uint   `json:"user_id" gorm:"unique"` // 用户ID，外键，唯一
	HeadImg []byte `json:"head_img"`              // 头像照片
	ImgName string `json:"img_name"`              // 图像名称
	ImgType string `json:"img_type"`              // 图像类型
	ImgSize uint   `json:"img_size"`              // 图像大小
	// references:ID 表示User表中的ID字段
	User User `gorm:"foreignKey:UserID;references:ID"` // UserHeadImage UserID -> User ID
}

// UserDetail 用户详细信息
type UserDetail struct {
	ID                uint       `json:"id" gorm:"primaryKey"`            // 主键
	UserID            uint       `json:"user_id" gorm:"unique"`           // 用户ID，外键，唯一索引
	CareerDirection   string     `json:"career_direction"`                // 职业方向
	UserHomePage      string     `json:"user_home_page"`                  // 个人主页
	UserSignature     string     `json:"user_signature"`                  // 个人签名
	UserTags          string     `json:"user_tags"`                       // 作者标签
	BlogLink          string     `json:"blog_link"`                       // 个人博客链接
	WeiboLink         string     `json:"weibo_link"`                      // 新浪微博链接
	GithubLink        string     `json:"github_link"`                     // Github链接
	LastLoginTime     *time.Time `json:"last_login_time"`                 // 上次登录时间
	CreateAccountTime *time.Time `json:"create_account_time"`             // 创建账号时间
	User              User       `gorm:"foreignKey:UserID;references:ID"` // UserDetail UserID -> User ID
}

// UserMessage 用户消息表(小铃铛专属)
type UserMessage struct {
	ID                uint `json:"id" gorm:"primaryKey"`            // 主键
	UserID            uint `json:"user_id" gorm:"uniqueIndex"`      // 用户ID，外键，唯一索引
	UnreadComments    int  `json:"unread_comments"`                 // 未读评论数量
	UnreadLikes       int  `json:"unread_likes"`                    // 未读点赞数量
	UnreadCollections int  `json:"unread_collections"`              // 未读收藏数量
	UnreadFans        int  `json:"unread_fans"`                     // 未读粉丝数量
	UnreadMessages    int  `json:"unread_messages"`                 // 未读私信数量
	UnreadSystem      int  `json:"unread_system"`                   // 未读系统消息数量
	User              User `gorm:"foreignKey:UserID;references:ID"` // UserMessage UserID -> User ID
}

// Follow 关注关系表
type Follow struct {
	ID         uint `json:"id" gorm:"primaryKey"`     // 主键
	FollowerID uint `json:"follower_id" gorm:"index"` // 关注者ID
	FollowedID uint `json:"followed_id" gorm:"index"` // 被关注者ID

	// FollowerID uint `json:"follower_id" gorm:"foreignKey:FollowerID;references:ID"`  // 关注者ID
	// FollowedID uint `json:"followed_id"  gorm:"foreignKey:FollowedID;references:ID"` // 被关注者ID
	User1 User `gorm:"foreignKey:FollowerID;references:ID"`
	User2 User `gorm:"foreignKey:FollowedID;references:ID"`
}
