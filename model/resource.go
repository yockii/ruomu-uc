package model

import "github.com/yockii/ruomu-core/database"

type Resource struct {
	Id              int64             `json:"id,omitempty" xorm:"pk"`
	ResourceName    string            `json:"resourceName,omitempty"`
	ResourceCode    string            `json:"resourceCode,omitempty"`
	ResourceContent string            `json:"resourceContent,omitempty"`
	ResourceType    string            `json:"resourceType,omitempty"`
	Action          string            `json:"action,omitempty"`
	CreateTime      database.DateTime `json:"createTime" xorm:"created"`
}

func (_ Resource) TableComment() string {
	return "资源表"
}
