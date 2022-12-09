package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
)

type User struct {
	Id           int64             `json:"id,omitempty" xorm:"pk"`
	Username     string            `json:"username,omitempty" xorm:"varchar(30) index comment('用户名')"`
	Password     string            `json:"password,omitempty" xorm:"comment('密码')"`
	RealName     string            `json:"realName,omitempty" xorm:"comment('真实姓名')"`
	ExternalId   string            `json:"externalId,omitempty" xorm:"varchar(50) index comment('外部关联ID')"`
	ExternalType string            `json:"externalType,omitempty" xorm:"comment('关联类型')"`
	Status       int               `json:"status,omitempty" xorm:"comment('状态 1-正常')"`
	CreateTime   database.DateTime `json:"createTime" xorm:"created"`
	UpdateTime   database.DateTime `json:"updateTime" xorm:"updated"`
}

func (_ User) TableComment() string {
	return "用户表"
}
func (u *User) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	u.Id = j.Get("id").Int()
	u.Username = j.Get("username").String()
	u.Password = j.Get("password").String()
	u.RealName = j.Get("realName").String()
	u.ExternalId = j.Get("externalId").String()
	u.ExternalType = j.Get("externalType").String()
	u.Status = int(j.Get("status").Int())

	return nil
}

type UserExtend struct {
	UserId       int64  `json:"userId" xorm:"pk"`
	ExternalInfo string `json:"externalInfo,omitempty" xorm:"text"`
}

func (_ UserExtend) TableComment() string {
	return "用户扩展表"
}

func (ue *UserExtend) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ue.UserId = j.Get("userId").Int()
	ue.ExternalInfo = j.Get("externalInfo").String()
	return nil
}
