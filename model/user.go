package model

import "github.com/yockii/ruomu-core/database"

type User struct {
	Id           int64             `json:"id,omitempty" xorm:"pk"`
	Username     string            `json:"username,omitempty" xorm:"varchar(30) index"`
	Password     string            `json:"password,omitempty"`
	RealName     string            `json:"realName,omitempty"`
	ExternalId   string            `json:"externalId,omitempty" xorm:"varchar(50) index"`
	ExternalType string            `json:"externalType,omitempty"`
	Status       int               `json:"status,omitempty"`
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
