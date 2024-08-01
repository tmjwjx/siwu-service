package models

// User 用户简略信息
type User struct {
	ID           uint   `json:"id" gorm:"primaryKey"` // 主键
	Nickname     string `json:"nickname"`             // 昵称
	Email        string `json:"email" gorm:"unique"`  // 邮箱，唯一
	Password     string `json:"password"`             // 密码
	PersonalHeat int    `json:"personal_heat"`        // 个人热度
	FansCount    uint   `json:"fans_count"`           // 粉丝数
	// many2many:user_follows_tags 中间表名称
	FollowsTags []Tag `gorm:"many2many:user_follows_tags;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:TagID"` // 关联中间表 user_follows_tags , USerID -> ID
}
