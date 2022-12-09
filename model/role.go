package model

import (
	"github.com/tidwall/gjson"
	"github.com/yockii/ruomu-core/database"
)

type Role struct {
	Id         int64             `json:"id,omitempty" xorm:"pk"`
	RoleName   string            `json:"roleName,omitempty" xorm:"comment('角色名称')"`
	RoleDesc   string            `json:"roleDesc,omitempty" xorm:"comment('角色描述')"`
	RoleType   int               `json:"roleType,omitempty" xorm:"comment('角色类型 1-普通角色 99-超级管理员角色')"`
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
	UpdateTime database.DateTime `json:"updateTime" xorm:"updated"`
}

func (_ Role) TableComment() string {
	return "角色表"
}
func (r *Role) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.Id = j.Get("id").Int()
	r.RoleName = j.Get("roleName").String()
	r.RoleDesc = j.Get("roleDesc").String()
	r.RoleType = int(j.Get("roleType").Int())
	return nil
}

type UserRole struct {
	Id         int64             `json:"id,omitempty" xorm:"pk"`
	UserId     int64             `json:"userId,omitempty"`
	RoleId     int64             `json:"roleId,omitempty"`
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
}

func (_ UserRole) TableComment() string {
	return "用户角色表"
}
func (ur *UserRole) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ur.Id = j.Get("id").Int()
	ur.UserId = j.Get("userId").Int()
	ur.RoleId = j.Get("roleId").Int()
	return nil
}
