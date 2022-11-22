package model

import "github.com/yockii/ruomu-core/database"

type Role struct {
	Id         int64             `json:"id,omitempty" xorm:"pk"`
	RoleName   string            `json:"roleName,omitempty"`
	RoleDesc   string            `json:"roleDesc,omitempty"`
	CreateTime database.DateTime `json:"createTime" xorm:"created"`
}

func (_ Role) TableComment() string {
	return "角色表"
}
