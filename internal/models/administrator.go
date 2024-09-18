package models

// import "gorm.io/gorm"
//
// // Administrator 管理员。后端手动向数据库添加的管理员（要和管理员给用户设置管理权限区分开）
// // 管理员给用户设置管理权限，那么要在Administrator表中添加一条数据（前台和后台的邮箱相同，密码不同，默认密码为abc123）；
// // 如果管理员给另一个管理员设置了一个用户权限，那么要在user表中添加一条数据（前台和后台的邮箱相同，密码不同，默认密码为abc123）
// type Administrator struct {
// 	gorm.Model        // ID CreatedAt UpdatedAt DeletedAt
// 	Email      string `json:"username" gorm:"size:16;not null"`
// 	Password   string `json:"password" gorm:"size:16;not null"`
// 	// todo
// 	// 等待补充管理员个人中心字段
//
// }
