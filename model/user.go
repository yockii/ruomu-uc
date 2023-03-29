package model

import (
	"github.com/tidwall/gjson"
)

type User struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey"`
	Username     string `json:"username,omitempty" gorm:"size:30;index;comment:'用户名'"`
	Password     string `json:"password,omitempty" gorm:"comment:'密码'"`
	RealName     string `json:"realName,omitempty" gorm:"comment:'真实姓名'"`
	ExternalId   string `json:"externalId,omitempty" gorm:"size:50;index;comment:'外部关联ID'"`
	ExternalType string `json:"externalType,omitempty" gorm:"comment:'关联类型'"`
	Status       int    `json:"status,omitempty" gorm:"comment:'状态 1-正常'"`
	CreateTime   int64  `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime   int64  `json:"updateTime" gorm:"autoUpdateTime"`
}

func (_ *User) TableComment() string {
	return "用户表"
}
func (u *User) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	u.ID = j.Get("id").Uint()
	u.Username = j.Get("username").String()
	u.Password = j.Get("password").String()
	u.RealName = j.Get("realName").String()
	u.ExternalId = j.Get("externalId").String()
	u.ExternalType = j.Get("externalType").String()
	u.Status = int(j.Get("status").Int())

	return nil
}

type UserExtend struct {
	UserID       uint64 `json:"userId,string" gorm:"primaryKey"`
	ExternalInfo string `json:"externalInfo,omitempty" gorm:"type:text"`
}

func (_ UserExtend) TableComment() string {
	return "用户扩展表"
}

func (ue *UserExtend) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ue.UserID = j.Get("userId").Uint()
	ue.ExternalInfo = j.Get("externalInfo").String()
	return nil
}
