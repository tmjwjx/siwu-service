package requests

import "time"

// MessageReq
// @Description: 消息请求
// @Author tianjiajie 2024-10-05 17:28:11
type MessageReq struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

// LikeAndCollectionMessageRes
// @Description: 点赞和收藏消息响应
// @Author tianjiajie 2024-10-05 17:32:06
type LikeAndCollectionMessageRes struct {
	UserId     uint       `json:"user_id"`
	Nickname   string     `json:"nickname"`
	Path       string     `json:"path"`
	ArticleId  uint       `json:"article_id"`
	Title      string     `json:"title"`
	CreatedAt  *time.Time `json:"created_at"`
	FormatTime string     `json:"format_time"` // 格式化时间
	DailyTime  string     `json:"daily_time"`  // 日常时间
}

// FollowMessageRes
// @Description: 关注消息响应
// @Author tianjiajie 2024-10-05 17:32:06
type FollowMessageRes struct {
	FollowerId uint       `json:"follower_id"`
	Nickname   string     `json:"nickname"`
	Path       string     `json:"path"`
	IsFollowed uint       `json:"is_followed"` // 是否关注
	CreatedAt  *time.Time `json:"created_at"`
	FormatTime string     `json:"format_time"` // 格式化时间
	DailyTime  string     `json:"daily_time"`  // 日常时间
}

// CommentMessageRes
// @Description: 评论消息响应
// @Author tianjiajie 2025-01-18 11:19:42
type CommentMessageRes struct {
	Nickname   string     `json:"nickname"`    // 评论者昵称
	Path       string     `json:"path"`        // 评论者头像URL
	Title      string     `json:"title"`       // 被评论文章标题
	Content    string     `json:"content"`     // 评论内容
	ParentId   uint       `json:"parent_id"`   // 被评论的评论id
	Comment    string     `json:"comment"`     // 被评论的评论
	CreatedAt  *time.Time `json:"created_at"`  // 评论创建时间
	FormatTime string     `json:"format_time"` // 格式化时间
	DailyTime  string     `json:"daily_time"`  // 日常时间
	LikesCount int        `json:"likes_count"` // 点赞数量
	Status     int        `json:"status"`      // 是否点赞
	UserId     uint       `json:"user_id"`     // 回复者用户ID
	ArticleId  uint       `json:"article_id"`  // 文章ID
	CommentId  uint       `json:"comment_id"`  // 评论的ID
}

// CommentLikeMessageRes
// @Description: 评论点赞消息响应
// @Author tianjiajie 2025-02-06 09:06:33
type CommentLikeMessageRes struct {
	UserId     uint       `json:"user_id"`     // 点赞者用户ID
	Nickname   string     `json:"nickname"`    // 点赞者昵称
	Path       string     `json:"path"`        // 点赞者头像URL
	ArticleId  uint       `json:"article_id"`  // 文章ID
	Title      string     `json:"title"`       // 文章标题
	CommentId  uint       `json:"comment_id"`  // 被点赞评论的ID
	Content    string     `json:"content"`     // 被点赞评论内容
	CreatedAt  *time.Time `json:"created_at"`  // 评论创建时间
	FormatTime string     `json:"format_time"` // 格式化时间
	DailyTime  string     `json:"daily_time"`  // 日常时间
}
