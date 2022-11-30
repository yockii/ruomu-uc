package model

import "github.com/yockii/ruomu-core/database"

type User struct {
	Id           int64             `json:"id,omitempty" xorm:"pk"`
	Username     string            `json:"username,omitempty" xorm:"varchar(30) index comment('用户名')"`
	Password     string            `json:"password,omitempty" xorm:"comment('密码')"`
	RealName     string            `json:"realName,omitempty" xorm:"comment('真实姓名')"`
	ExternalId   string            `json:"externalId,omitempty" xorm:"varchar(50) index comment('外部关联ID')"`
	ExternalType string            `json:"externalType,omitempty" xorm:"comment('关联类型')"`
	Status       int               `json:"status,omitempty" xorm:"comment('状态 1-正常')"`
	CreateTime   database.DateTime `json:"createTime" xorm:"created"`
}

func (_ User) TableComment() string {
	return "用户表"
}

type UserExtend struct {
	UserId       int64  `json:"userId" xorm:"pk"`
	ExternalInfo string `json:"externalInfo,omitempty" xorm:"text"`
}

func (_ UserExtend) TableComment() string {
	return "用户扩展表"
}
