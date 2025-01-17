package requests

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
	UserIds []uint `json:"user_ids"`
}
