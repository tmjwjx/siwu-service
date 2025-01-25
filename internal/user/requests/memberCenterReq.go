package requests

import "time"

type ArticleInfo struct {
	LikesCount       int `json:"likes_count"`       // 点赞数量
	ViewsCount       int `json:"views_count"`       // 浏览量
	CollectionsCount int `json:"collections_count"` // 收藏数量
}
type InitUserInfoRes struct {
	HeadShot         string `json:"head_shot"`        // 用户头像
	Nickname         string `json:"nickname"`         // 昵称
	Signature        string `json:"signature"`        // 个人签名
	LikesCount       int    `json:"likes_count"`      // 所有文章总点赞数量
	ViewsCount       int    `json:"reads_count"`      // 所有文章总浏览量
	CollectionsCount int    `json:"attentions_count"` // 所有文章总收藏数量
	ConcernsCount    int    `json:"concerns_count"`   // 关注者
	FansCount        int    `json:"fans_count"`       // 被关注者
	Date             string `json:"date"`             // 注册时间
	Tag              string `json:"tag"`              // 网站的名称
	ConcernStatus    int    `json:"concern_status"`   // 关注状态 0：未关注 1：关注 2：用户本人页面
	ArticleCount     int    `json:"article_count"`    // 用户的总文章数量
}
type MiddleInfo struct {
	Date      time.Time `json:"date"`      // 注册时间
	Nickname  string    `json:"nickname"`  // 昵称
	Signature string    `json:"signature"` // 个人签名
	HeadShot  string    `json:"head_shot"` // 用户头像
}

//// 用户
//type UserFollower struct {
//	FollowedId uint `json:"followed_id"`
//}
//type UserFollowed struct {
//	FollowerId uint `json:"follower_id"`
//}

// EditSignatureReq
// @Description: 编辑个签
// @Author wangyulong 2024-10-11 14:46:26
type EditSignatureReq struct {
	Signature string `json:"signature"`
}
