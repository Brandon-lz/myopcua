package db

import (
	"gorm.io/gorm"
)

func initModels() {
	modelsToMigrate.Add(&WebHookConditions{})
	modelsToMigrate.Add(&WebHooks{})
	modelsToMigrate.Add(&NeedNode{})
}


type WebHookConditions struct {
	gorm.Model
	Condition string `json:"condition" gorm:"not null;comment:条件表达式"`
	WebHooks  []WebHooks `gorm:"foreignKey:WebHookConditionRefer"`
}

// 表名	webhook
// 字段	是否必填	唯一字段	类型	默认值	描述
// id	Y	Y	uint		主键
// name 	Y	Y	string		webhook名称
// url	Y	Y	string		url地址
// activate	N	N	bool	TRUE	是否激活

type WebHooks struct {
	gorm.Model
	Name   string `json:"name" gorm:"unique;not null;index;comment:webhook名称"` // not null in db
	Url    string `json:"url" gorm:"unique;not null;index;comment:url地址"`
	Active bool   `json:"active" gorm:"default:true;comment:是否激活"`
	WebHookConditionRefer *uint `json:"web_hook_condition_refer" gorm:"comment:WebHookCondition foreign key"`   // 外键，在数据库中最好不要not null，在逻辑上去判空
	NeedNodes []NeedNode `gorm:"foreignKey:WebHookRefer"`
}

type NeedNode struct {
	gorm.Model
	NodeName string `json:"node_name" gorm:"not null;index;comment:需要订阅的节点名称"`
	WebHookRefer *uint `json:"web_hook_refer" gorm:"comment:WebHook foreign key"`   // 外键，在数据库中最好不要not null，在逻辑上去判空
}