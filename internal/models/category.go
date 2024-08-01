package models

// Category 类目表
type Category struct {
	ID            uint   `json:"id" gorm:"primaryKey"` // 主键
	CategoryName  string `json:"category_name"`        // 类目名称
	CategoryImage string `json:"category_image"`       // 类目图标
}
