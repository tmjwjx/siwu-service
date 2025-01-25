package requests

import (
	"time"
)

// ClickAttentionReq 点击关注和点击取消关注请求
type ClickAttentionReq struct {
	FollowerId uint `json:"follower_id"`
	FollowedId uint `json:"followed_id"`
}

// UserRankReq 用户排行请求
type UserRankReq struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

// UserRankRes 用户排行响应
type UserRankRes struct {
	Id              uint   `json:"id"`
	Nickname        string `json:"nickname"`
	CareerDirection string `json:"career_direction"`
	AvatarPath      string `json:"avatar_path"` // 头像图片
	IsFollowed      int    `json:"is_followed"` // 用户是否关注了这个排行榜上的用户：未关注：0，已关注：1，这个用户是自己：2
}

// AttentionReq 搜索用户关注的人请求
type AttentionReq struct {
	UserId  uint   `json:"user_id"`
	Keyword string `json:"keyword"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
}

// AttentionRes 搜索用户关注的人响应
type AttentionRes struct {
	Ids []uint `json:"ids"`
}

// GetBasicInfoReq 通过ids获取到用户简略信息请求
type GetBasicInfoReq struct {
	Ids []uint `json:"ids"`
}

// GetBasicInfoRes 通过ids获取到用户简略信息响应
type GetBasicInfoRes struct {
	ID              uint      `json:"id"`               // id
	CreatedAt       time.Time `json:"created_at"`       // 创建时间
	UpdatedAt       time.Time `json:"updated_at"`       // 更新时间
	Nickname        string    `json:"nickname"`         // 昵称
	Email           string    `json:"email"`            // 邮箱，唯一
	Heat            int       `json:"heat"`             // 个人热度
	AttentionCount  uint      `json:"attention_count"`  // 关注了多少人数
	FansCount       int       `json:"fans_count"`       // 粉丝数
	PrivateSettings string    `json:"private_settings"` // 私信设置
	Status          int       `json:"status"`           // 用户状态：0全部 1正常 2封禁
	LastLoginTime   time.Time `json:"last_login_time"`  // 最后一次的登录时间。设置默认值为零值
	IsFollowed      int       `json:"is_followed"`      // 是否已关注该用户。未关注：0，已关注：1。
	AvatarPath      string    `json:"avatar_path"`      // 头像路径
	AuthorArticles  int       `json:"author_articles"`  // 拥有文章数
}

// // UserDataRequest 用户 文章 请求
// type UserDataRequest struct {
//	Id    uint `json:"id" form:"id"`
//	Page  int  `json:"page" form:"page"`
//	Limit int  `json:"limit" form:"limit"`
// }
