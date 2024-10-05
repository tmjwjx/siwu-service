package requests

import (
	"time"
)

// CommentMesRes 评论消息(前台)
type CommentMesRes struct {
	CommentList *[]*CommentObj
}

type CommentObj struct {
	Nickname   string `json:"nickname"`    // 评论者昵称
	Path       string `json:"path"`        // 评论者头像
	Title      string `json:"title"`       // 被评论文章
	Content    string `json:"content"`     // 评论内容
	Comment    string `json:"comment"`     // 被评论评论(若评论为二级评论)
	CreatedAt  string `json:"created_at"`  // 创建时间
	LikesCount int    `json:"likes_count"` // 点赞数量
	ID         uint   `json:"id"`          // 被回复的评论的ID
	Status     int    `json:"status"`      // 用于判断该用户对该评论的点赞状态(0: 未点赞 , 1: 点赞)
}

type Result1 struct {
	Nickname   string    `json:"nickname"`    // 评论者昵称
	UserID     uint      `json:"user_id"`     // 评论者ID (用于查询用户用户头像)
	Title      string    `json:"title"`       // 被评论文章
	ID         uint      `json:"id"`          // 回复的id
	Content    string    `json:"content"`     // 评论内容
	ParentID   *uint     `json:"parent_id"`   // 上一条评论ID , 允许为 null，表示顶级评论(用于查询被回复的评论的id和content)
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	LikesCount int       `json:"likes_count"` // 点赞数量
}

type Result2 struct {
	ID      uint   `json:"id"`      // 被回复的评论的id
	Content string `json:"content"` // 被回复的评论的content
}

type reply struct {
}