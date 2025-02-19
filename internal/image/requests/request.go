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

// CompressImageReq 请求参数结构体
type CompressImageReq struct {
	Path   string `json:"path" binding:"required"`              // 图片路径（相对于静态文件夹的路径）
	Width  int    `json:"width" binding:"min=0"`                // 目标宽度
	Height int    `json:"height" binding:"min=0"`               // 目标高度
	Level  int    `json:"level" binding:"required,min=0,max=9"` // 压缩级别 (0-9)
}

// WatermarkParam 与图片水印有关的参数
type WatermarkParam struct {
	Watermark string  `yaml:"watermark"` // 图片水印的内容
	Size      float64 `yaml:"size"`      // 水印的大小
	Scale     float64 `yaml:"scale"`     // 图片压缩时的缩放比例
	Margin    float64 `yaml:"margin"`    // 水印距离左下角的距离 ，规则: 从右向左移动margin个位置的距离，从下往上移动3*margin个位置的距离
}
