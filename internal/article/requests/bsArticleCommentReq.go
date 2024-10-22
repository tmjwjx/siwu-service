package requests

// BatchReviewReq 批量审核请求
type BatchReviewReq struct {
	IDs []uint // 审核通过的评论id
}

// BatchReviewRes 批量审核响应
type BatchReviewRes struct {
	ComList2 ComList2
}

type ComList2 struct {
}

// CommentsListReq 评论列表加载请求
type CommentsListReq struct {
	Offset      int    `json:"offset"`       // 分页查询的起始位置
	Limit       int    `json:"limit"`        // 分页查询要返回记录的数量
	Type        int    `json:"type"`         // 1:全部 ,2:待审 3:已批准
	Email       string `json:"email"`        // 用户邮箱
	Nickname    string `json:"nickname"`     // 用户昵称
	Title       string `json:"title"`        // 文章标题
	ParentEmail string `json:"parent_email"` // 用户回复对象的昵称
}

// CommentsListRes 评论列表加载响应
type CommentsListRes struct {
	Comlist *[]*ComList `json:"comlist"`
	Total   int         `json:"total"`
}

type ComList struct {
	ID uint `json:"id"` // 评论ID

	Nickname string `json:"nickname"` // 用户昵称
	Email    string `json:"email"`    // 邮箱，唯一
	Path     string `json:"path"`     // 用户头像

	ParentNickname string `json:"parent_nickname"` // 用户回复对象的昵称
	Content        string `json:"content"`         // 评论内容
	CommentPath    string `json:"comment_path"`    // 用户发的评论中的图片

	ArticleID uint   `json:"article_id"` // 文章ID
	Title     string `json:"title"`      // 文章标题
	Summary   string `json:"summary"`    // 文章摘要

	Examine int `json:"examine"` // 是否审核:1审核2:未审核
}

// AddCommentReq 添加评论请求
type AddCommentReq struct {
	ArticleID    uint     `json:"article_id" form:"article_id"`         // 所属文章ID，外键
	UserID       uint     `json:"user_id" form:"user_id"`               // 作者ID
	HighestID    *uint    `json:"highest_id" form:"highest_id"`         // 最上层一级评论(顶级评论)
	ParentID     *uint    `json:"parent_id" form:"parent_id"`           // 被回复评论ID , 允许为 null，表示顶级评论
	ParentUserID *uint    `json:"parent_user_id" form:"parent_user_id"` // 上一条评论的发布用户ID
	Content      string   `json:"content" form:"content"`               // 评论内容
	CommentPath  []string `json:"comment_path" form:"comment_path"`     // 用户发的评论中的图片
	LikesCount   int      `json:"likes_count" form:"likes_count"`       // 点赞数量
}

// DelCommentReq 删除评论请求
type DelCommentReq struct {
	ID uint `json:"id"` // 评论ID
}

// BsBatchDelCommentReq 批量删除请求
type BsBatchDelCommentReq struct {
	ID []uint `json:"id"` // 标签ID
}

// UpdateCommentReq 更新评论请求
type UpdateCommentReq struct {
	ID          uint     `form:"id"`           // 评论ID
	Content     string   `form:"content"`      // 评论内容
	CommentPath []string `form:"comment_path"` // 评论中的图片
}

// QueryCommentReq 查询评论请求
type QueryCommentReq struct {
	UserID uint `json:"user_id"` // 作者ID
	Offset int  `json:"offset"`  // 分页查询的起始位置
	Limit  int  `json:"limit"`   // 分页查询要返回记录的数量
}

// QueryCommentRes 评论列表加载响应
type QueryCommentRes struct {
	ID uint `json:"id"` // 评论ID

	Nickname string `json:"nickname"` // 用户昵称
	Email    string `json:"email"`    // 邮箱，唯一
	Path     string `json:"path"`     // 用户头像

	ParentNickname string `json:"parent_nickname"` // 用户回复对象的昵称
	Content        string `json:"content"`         // 评论内容
	CommentPath    string `json:"comment_path"`    // 用户发的评论中的图片

	ArticleID uint   `json:"article_id"` // 文章ID
	Title     string `json:"title"`      // 文章标题
	Summary   string `json:"summary"`    // 文章摘要
}
