package model

import "github.com/yockii/ruomu-core/database"

type Role struct {
	Id         int64             `json:"id,omitempty" xorm:"pk"`
	RoleName   string            `json:"roleName,omitempty" xorm:"comment('角色名称')"`
	RoleDesc   string            `json:"roleDesc,omitempty" xorm:"comment('角色描述')"`
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
}

func (_ Role) TableComment() string {
	return "角色表"
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
