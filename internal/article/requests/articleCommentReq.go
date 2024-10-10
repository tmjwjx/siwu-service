package requests

type ArticleCommentReq struct {
	ArticleID    uint     `json:"article_id"`
	UserID       uint     `json:"user_id"`
	HighestID    *uint    `json:"highest_id"`
	ParentID     *uint    `json:"parent_id"`
	ParentUserID *uint    `json:"parent_user_id"` // 上一条评论的发布用户ID
	Content      string   `json:"content"`
	Path         []string `json:"path" form:"path"` // 标签头像
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
	UserID    uint `json:"user_id"`
}

type TopCommentsRes struct {
	ID           uint   `json:"id"`             // 评论ID
	Nickname     string `json:"nickname"`       // 用户昵称
	CreateAT     string `json:"create_at"`      // 评论创建时间
	ArticleID    uint   `json:"article_id"`     // 文章ID
	UserID       uint   `json:"user_id"`        // 用户ID
	HighestID    *uint  `json:"highest_id"`     // 顶级评论ID
	ParentID     *uint  `json:"parent_id"`      // 上一级评论ID
	ParentUserID *uint  `json:"parent_user_id"` // 上一条评论的发布用户ID
	Content      string `json:"content"`        // 评论内容
	LikesCount   int    `json:"likes_count"`    // 点赞数量
	RepliesCount int64  `json:"replies_count"`  // 回复数量
	Path         string `json:"path"`           // 用户头像
	CommentPath  string `json:"comment_path"`   // 用户发的评论中的图片
	Status       int    `json:"status"`         // 判断用户对评论的点赞情况
}

type RepliesReq2 struct {
	HighestID *uint `json:"highest_id"` // 顶级评论ID
	Offset    int   `json:"offset"`     // 分页查询的起始位置
	Limit     int   `json:"limit"`      // 分页查询要返回记录的数量
	UserID    uint  `json:"user_id"`
}

type RepliesRes struct {
	ID             uint   `json:"id"`              // 评论ID
	Nickname       string `json:"nickname"`        // 用户昵称
	ParentNickname string `json:"parent_nickname"` // 用户回复对象的昵称
	CreateAT       string `json:"create_at"`       // 评论创建时间
	ArticleID      uint   `json:"article_id"`      // 文章ID
	UserID         uint   `json:"user_id"`         // 用户ID
	HighestID      *uint  `json:"highest_id"`      // 顶级评论ID
	ParentID       *uint  `json:"parent_id"`       // 上一级评论ID
	ParentUserID   *uint  `json:"parent_user_id"`  // 上一条评论的发布用户ID
	Content        string `json:"content"`         // 评论内容
	LikesCount     int    `json:"likes_count"`     // 点赞数量
	Path           string `json:"path"`            // 用户头像
	ParentPath     string `json:"parent_path"`     // 用户回复对象的头像
	CommentPath    string `json:"comment_path"`    // 用户发的评论中的图片
	Status         int    `json:"status"`          // 判断用户对评论的点赞情况
}

type DelComment struct {
	ID uint `json:"id"`
}

type PraiseCount struct {
	ID     uint `json:"id"`      // 评论ID
	Status uint `json:"status"`  // 用于判断是增加还是减少点赞数量，(1 : 代表增加， 2 : 代表减少)
	UserID uint `json:"user_id"` // 点赞人ID
}
