package models

import "gorm.io/gorm"

type RequestMethod struct {
	gorm.Model
	Name string `json:"name"` //请求方法名
}
