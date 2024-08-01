package models

// UserHeadImage 头像信息
type UserHeadImage struct {
	ID      uint   `json:"id" gorm:"primaryKey"`  // 用户头像图片
	UserID  uint   `json:"user_id" gorm:"unique"` // 用户ID，外键，唯一
	HeadImg []byte `json:"head_img"`              // 头像照片
	ImgName string `json:"img_name"`              // 图像名称
	ImgType string `json:"img_type"`              // 图像类型
	ImgSize uint   `json:"img_size"`              // 图像大小
	// references:ID 表示User表中的ID字段
	User User `gorm:"foreignKey:UserID;references:ID"` // UserHeadImage UserID -> User ID
}
