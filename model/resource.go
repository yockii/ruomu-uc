package model

import (
	"github.com/tidwall/gjson"
)

type Resource struct {
	ID           uint64 `json:"id,omitempty,string" gorm:"primaryKey"`
	ResourceName string `json:"resourceName,omitempty" gorm:"comment:'资源名称'"` // 资源名称
	ResourceCode string `json:"resourceCode,omitempty" gorm:"comment:'资源代码'"` // 资源认证代码
	HttpMethod   string `json:"httpMethod,omitempty" gorm:"comment:'http方法'"` // http方法
	CreateTime   int64  `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime   int64  `json:"updateTime" gorm:"autoUpdateTime"`
}

func (_ Resource) TableComment() string {
	return "资源表"
}
func (r *Resource) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	r.ID = j.Get("id").Uint()
	r.ResourceName = j.Get("resourceName").String()
	r.ResourceCode = j.Get("resourceCode").String()
	r.HttpMethod = j.Get("httpMethod").String()
	return nil
}

type RoleResource struct {
	ID         uint64 `json:"id,omitempty,string" gorm:"primaryKey"`
	RoleID     uint64 `json:"roleId,omitempty,string"`
	ResourceID uint64 `json:"resourceId,omitempty,string"`
	CreateTime int64  `json:"createTime" gorm:"autoCreateTime"`
}

func (_ RoleResource) TableComment() string {
	return "角色资源表"
}
func (rr *RoleResource) UnmarshalJSON(b []byte) error {
	j := gjson.ParseBytes(b)
	rr.ID = j.Get("id").Uint()
	rr.RoleID = j.Get("roleId").Uint()
	rr.ResourceID = j.Get("resourceId").Uint()
	return nil
}
