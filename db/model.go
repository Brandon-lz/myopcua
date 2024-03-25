package db

import (
	"gorm.io/gorm"
)

// 表名	webhook
// 字段	是否必填	唯一字段	类型	默认值	描述
// id	Y	Y	uint		主键
// name 	Y	Y	string		webhook名称
// url	Y	Y	string		url地址
// activate	N	N	bool	TRUE	是否激活

type WebHook struct {
	gorm.Model
	Name   string `json:"name" gorm:"not null;comment:webhook名称"` // not null in db
	Url    string `json:"url" gorm:"unique;not null;comment:url地址"`
	Active bool   `json:"active" gorm:"default:true;comment:是否激活"`
}
