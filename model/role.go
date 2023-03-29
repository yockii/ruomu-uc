package model

import (
	"github.com/tidwall/gjson"
)

type Role struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey"`
	RoleName   string `json:"roleName,omitempty" gorm:"comment:'角色名称'"`
	RoleDesc   string `json:"roleDesc,omitempty" gorm:"comment:'角色描述'"`
	RoleType   int    `json:"roleType,omitempty" gorm:"comment:'角色类型 1-普通角色 -1-超级管理员角色'"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime int64  `json:"updateTime" gorm:"autoUpdateTime"`
}

func (_ Role) TableComment() string {
	return "角色表"
}
func (r *Role) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.ID = j.Get("id").Uint()
	r.RoleName = j.Get("roleName").String()
	r.RoleDesc = j.Get("roleDesc").String()
	r.RoleType = int(j.Get("roleType").Int())
	return nil
}

type UserRole struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey"`
	UserID     uint64 `json:"userId,omitempty,string"`
	RoleID     uint64 `json:"roleId,omitempty,string"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime"`
}

func (_ UserRole) TableComment() string {
	return "用户角色表"
}
func (ur *UserRole) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	ur.ID = j.Get("id").Uint()
	ur.UserID = j.Get("userId").Uint()
	ur.RoleID = j.Get("roleId").Uint()
	return nil
}
