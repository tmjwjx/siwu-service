package requests

//import (
//	"forum/pkg/globals"
//	"gorm.io/gorm"
//)
//
//// Attachment 图片附件表
//type Attachment struct {
//	gorm.Model              //ID CreatedAt UpdatedAt DeletedAt
//	Home       globals.Home `json:"home"`    // 图片所属单位，即属于文章图片还是用户图片或是资源表图片等
//	HomeID     uint         `json:"home_id"` // 图片对应的具体文章或用户的ID等
//	Name       string       `json:"name"`    // 文件名
//	Type       string       `json:"type"`    // 文件类型，例如 images/png
//	Size       int64        `json:"size"`    // 文件大小（以字节为单位）
//	Path       string       `json:"path"`    // 文件路径
//}

// ImageUrl 响应结构体
type ImageUrl struct {
	Data  []*UrlPath `json:"data"`
	Errno int        `json:"errno"`
}

// UrlPath 图片的Url路径
type UrlPath struct {
	Url string `json:"url"`
}
