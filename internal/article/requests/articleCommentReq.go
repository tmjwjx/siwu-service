package requests

import "time"

type ArticleCommentReq struct {
	ArticleID uint   `json:"article_id"`
	UserID    uint   `json:"user_id"`
	HighestID *uint  `json:"highest_id"`
	ParentID  *uint  `json:"parent_id"`
	Content   string `json:"content"`
}

/*type ArticleCommentRes struct {
	ID        uint                `json:"id"`
	ArticleID uint                `json:"article_id"`
	UserID    uint                `json:"user_id"`
	Content   string              `json:"content"`
	ParentID  *uint               `json:"parent_id"`
	CreateAT  time.Time           `json:"create_at"`
	Replies   []ArticleCommentRes `json:"replies"` // 子评论
}

type RepliesReq struct {
	ParentID string `uri:"parentId"` // 绑定路径参数
	Offset   int    `form:"offset"`  // 绑定查询参数
	Limit    int    `form:"limit"`   // 绑定查询参数
}

type CommentReq struct {
	ArticleID string `json:"article_id" uri:"article_id"`
	Offset    int    `json:"offset" form:"offset"`
}*/

type TopCommentsReq struct {
	ArticleID uint `json:"article_id"`
	Offset    int  `json:"offset"`
	Limit     int  `json:"limit"`
}

type TopCommentsRes struct {
	ID           uint      `json:"id"`
	Nickname     string    `json:"nickname"` // 用户昵称
	CreateAT     time.Time `json:"create_at"`
	ArticleID    uint      `json:"article_id"`
	UserID       uint      `json:"user_id"`
	HighestID    *uint     `json:"highest_id"`
	ParentID     *uint     `json:"parent_id"`
	Content      string    `json:"content"`
	LikesCount   int       `json:"likes_count"`   // 点赞数量
	RepliesCount int64     `json:"replies_count"` // 回复数量
	Path         string    `json:"path"`          // 用户头像
	CommentPath  string    `json:"comment_path"`  // 用户发的评论中的图片
}

type RepliesReq2 struct {
	HighestID *uint `json:"highest_id"`
	Offset    int   `json:"offset"`
	Limit     int   `json:"limit"`
}

type RepliesRes struct {
	ID             uint      `json:"id"`
	Nickname       string    `json:"nickname"`        // 用户昵称
	ParentNickname string    `json:"parent_nickname"` //用户回复对象的昵称
	CreateAT       time.Time `json:"create_at"`
	ArticleID      uint      `json:"article_id"`
	UserID         uint      `json:"user_id"`
	HighestID      *uint     `json:"highest_id"`
	ParentID       *uint     `json:"parent_id"`
	Content        string    `json:"content"`
	LikesCount     int       `json:"likes_count"`  // 点赞数量
	Path           string    `json:"path"`         // 用户头像
	ParentPath     string    `json:"parent_path"`  // 用户回复对象的头像
	CommentPath    string    `json:"comment_path"` // 用户发的评论中的图片
}

type DelComment struct {
	ID uint `json:"id"`
}

type PraiseCount struct {
	ID     uint `json:"id"`
	Status uint `json:"status"` // 用于判断是增加还是减少点赞数量，(1 : 代表增加， 2 : 代表减少)
}
