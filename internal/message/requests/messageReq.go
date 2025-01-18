package requests

// MessageReq
// @Description: 消息请求
// @Author tianjiajie 2024-10-05 17:28:11
type MessageReq struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

// LikeAndCollectionMessageRes
// @Description: 点赞消息响应
// @Author tianjiajie 2024-10-05 17:32:06
type LikeAndCollectionMessageRes struct {
	UserId    uint   `json:"user_id"`
	Nickname  string `json:"nickname"`
	Path      string `json:"path"`
	ArticleId uint   `json:"article_id"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
}

// FollowMessageRes
// @Description: 关注消息响应
// @Author tianjiajie 2024-10-05 17:32:06
type FollowMessageRes struct {
	FollowerId uint   `json:"follower_id"`
	Nickname   string `json:"nickname"`
	Path       string `json:"path"`
	IsFollowed uint   `json:"is_followed"` // 是否关注
	CreatedAt  string `json:"created_at"`
}

//export interface Response {
//code: number;
//data: Data;
//msg: string;
//[property: string]: any;
//}
//
//export interface Data {
//comment_list: CommentList[];
//"01J93QH4R80A7WAASN2H0BKMG9": any;
//[property: string]: any;
//}
//
//export interface CommentList {
///**
// * 文章ID
// */
//article_id: number;
///**
// * 被评论评论(若评论为二级评论)(可以是对作者文章的评论）
// */
//comment: string;
///**
// * 评论的ID
// */
//comment_id: number;
///**
// * 评论内容
// */
//content: string;
///**
// * 创建时间
// */
//created_at: string;
///**
// * 点赞数量
// */
//likes_count: number;
///**
// * 评论者昵称
// */
//nickname: string;
///**
// * 评论者头像
// */
//path: string;
///**
// * 是否点赞
// */
//status: number;
///**
// * 被评论文章
// */
//title: string;
///**
// * 回复者用户ID
// */
//user_id: number;
//[property: string]: any;
//}

// CommentMessageRes
// @Description: 评论消息响应
// @Author tianjiajie 2025-01-18 11:19:42
type CommentMessageRes struct {
	Nickname   string `json:"nickname"`    // 评论者昵称
	Path       string `json:"path"`        // 评论者头像URL
	Title      string `json:"title"`       // 被评论文章标题
	Content    string `json:"content"`     // 评论内容
	Comment    string `json:"comment"`     // 被评论的评论
	CreatedAt  string `json:"created_at"`  // 评论创建时间
	LikesCount int    `json:"likes_count"` // 点赞数量
	Status     int    `json:"status"`      // 是否点赞
	UserId     uint   `json:"user_id"`     // 回复者用户ID
	ArticleId  uint   `json:"article_id"`  // 文章ID
	CommentId  uint   `json:"comment_id"`  // 评论的ID
}
