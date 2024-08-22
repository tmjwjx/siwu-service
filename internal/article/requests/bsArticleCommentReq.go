package requests

// CommentsListReq 评论列表加载请求
type CommentsListReq struct {
	Offset int `json:"offset"` // 分页查询的起始位置
	Limit  int `json:"limit"`  // 分页查询要返回记录的数量
}

// CommentsListRes 评论列表加载响应
type CommentsListRes struct {
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

// AddCommentReq 添加评论请求
type AddCommentReq struct {
	ArticleID    uint   `form:"article_id"`     // 所属文章ID，外键
	UserID       uint   `form:"user_id"`        // 作者ID
	HighestID    *uint  `form:"highest_id"`     // 最上层一级评论
	ParentID     *uint  `form:"parent_id"`      // 上一条评论ID , 允许为 null，表示顶级评论
	ParentUserID *uint  `form:"parent_user_id"` // 上一条评论的发布用户ID
	Content      string `form:"content"`        // 评论内容
	//CommentPath  string `form:"comment_path"`   // 用户发的评论中的图片
	LikesCount int `form:"likes_count"` // 点赞数量
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
	ID      uint   `form:"id"`      // 评论ID
	Content string `form:"content"` // 评论内容
	//CommentPath string `form:"comment_path"` // 评论中的图片
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
